package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	stdruntime "runtime"

	"github.com/pbnjay/memory"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"

	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/file"
)

type Config struct {
	// Name is the name of the executor.
	Name string
	// Labels the executor will advertise to the broker.
	Labels []string
}

type Server struct {
	executorv1.UnimplementedExecutorServer
	syslog     *zap.SugaredLogger
	config     Config
	supervisor *supervisor
}

func NewServer(syslog *zap.SugaredLogger, config Config) *Server {
	syslog = syslog.Named("executor")
	exec := &Server{
		syslog:     syslog,
		config:     config,
		supervisor: newSupervisor(syslog),
	}
	return exec
}

func (s *Server) Events(req *executorv1.EventsRequest, stream executorv1.Executor_EventsServer) error {
	if err := validateEventsRequest(req); err != nil {
		return err
	}
	log, done, err := s.supervisor.PrepareRuntime(req.BuildId, req.RuntimeId)
	if err != nil {
		return err
	}
	defer done()
	var closed bool
	var sendErrC = make(chan error)
	fail := func(err error) {
		if !closed {
			closed = true
			sendErrC <- err
		}
	}
	done = log.Stream().Subscribe(func(event *event.Event) {
		out, err := event.Marshal()
		if err != nil {
			fail(err)
			return
		}
		if err := stream.Send(out); err != nil {
			fail(err)
			return
		}
	})
	defer done()
	log.Publish(&builtinv1.SyncPointReachedEvent{BarrierId: req.BarrierId})
	syslog := s.syslog.With("runtime_id", req.RuntimeId)
	syslog.Infow("Event stream opened")
	select {
	case <-stream.Context().Done():
		syslog.Infow("Event stream closed")
	case err := <-sendErrC:
		syslog.Infow("Event stream closed with send error: %v", err)
		return err
	}
	return nil
}

func (s *Server) Open(ctx context.Context, req *executorv1.OpenRequest) (*executorv1.OpenResponse, error) {
	if err := s.validateOpenRequest(req); err != nil {
		return nil, err
	}
	s.syslog.Infow("Opening runtime", "runtime_id", req.RuntimeId)
	runtime, err := s.supervisor.OpenRuntime(ctx, req.BuildId, req.RuntimeId, req.Opts)
	if err != nil {
		return nil, err
	}
	s.syslog.Infow("Opened runtime", "runtime_id", req.RuntimeId)
	return &executorv1.OpenResponse{WorkDirectory: runtime.Directory(), SysInfo: s.getSysInfo()}, nil
}

func (s *Server) Exec(ctx context.Context, req *executorv1.ExecRequest) (*executorv1.ExecResponse, error) {
	if err := validateExecRequest(req); err != nil {
		return nil, err
	}
	runtime, err := s.supervisor.GetRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	res, err := runtime.Exec(ctx, req.ExecId, req.Opts)
	if err != nil {
		return nil, err
	}
	runtime.Log().Publish(&builtinv1.SyncPointReachedEvent{BarrierId: req.BarrierId})
	return &executorv1.ExecResponse{ExitCode: res.ExitCode}, nil
}

func (s *Server) Import(stream executorv1.Executor_ImportServer) error {
	var (
		runtime   runtime.Runtime
		importID  string
		receivers = make(map[string]*file.Receiver)
	)
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return stream.SendAndClose(&executorv1.ImportResponse{})
			}
			s.syslog.Errorw("recv error", "error", err)
			return fmt.Errorf("error in receive: %w", err)
		}
		if err = validateFileTransfer(req); err != nil {
			return err
		}
		if runtime == nil {
			runtime, err = s.supervisor.GetRuntime(req.RuntimeId)
			if err != nil {
				return err
			}
			importID = req.TransferId
		}
		if runtime.ID() != req.RuntimeId {
			return fmt.Errorf("invalid runtime id")
		}
		if importID != req.TransferId {
			return fmt.Errorf("invalid import id")
		}
		receiver, ok := receivers[req.FileId]
		if !ok {
			receiver = file.NewReceiver(s.syslog, runtime)
			receivers[req.FileId] = receiver
		}
		err = receiver.Next(req)
		if receiver.State() == file.ReceiveStateDone {
			delete(receivers, req.FileId)
		}
		if err != nil {
			return err
		}
	}
}

func (s *Server) Export(req *executorv1.ExportRequest, stream executorv1.Executor_ExportServer) error {
	if err := validateExportRequest(req); err != nil {
		return err
	}
	runtime, err := s.supervisor.GetRuntime(req.RuntimeId)
	if err != nil {
		return err
	}
	sendCallback := func(header *executorv1.FileTransferHeader) {
		if header.IsDir {
			runtime.Log().Printf("Exported directory src=%s, dest=%s, mode=%s", header.SrcPath, header.DestPath, os.FileMode(header.Mode))
		} else {
			runtime.Log().Printf("Exported file src=%s, dest=%s, mode=%s, size=%d", header.SrcPath, header.DestPath, os.FileMode(header.Mode), header.Size)
		}
	}
	skipCallback := func(path string, isDir bool, excludedBy string) {
		if isDir {
			runtime.Log().Printf("Skipped directory export src=%s, excluded_by=%s", path, excludedBy)
		} else {
			runtime.Log().Printf("Skipped file export src=%s, excluded_by=%s", path, excludedBy)
		}
	}
	sendOpts := []file.SendOpt{file.WithSendCallback(sendCallback), file.WithSkipCallback(skipCallback)}
	if req.Opts != nil {
		if len(req.Opts.Excludes) > 0 {
			sendOpts = append(sendOpts, file.WithExcludes(req.Opts.Excludes))
		}
		if req.Opts.DestPath != "" {
			sendOpts = append(sendOpts, file.WithDest(req.Opts.DestPath))
		}
	}
	sender := file.NewSender(s.syslog, runtime.ReadFS(), stream, runtime.ID(), req.ExportId, sendOpts...)
	_, err = sender.Send(req.SrcPath)
	return err
}

func (s *Server) Heartbeat(ctx context.Context, req *executorv1.HeartbeatRequest) (*executorv1.HeartbeatResponse, error) {
	if err := validateHeartbeatRequest(req); err != nil {
		return nil, err
	}
	extendedBy, err := s.supervisor.ExtendRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	return &executorv1.HeartbeatResponse{ExtendedBy: durationpb.New(extendedBy)}, nil
}

func (s *Server) Close(ctx context.Context, req *executorv1.CloseRequest) (*executorv1.CloseResponse, error) {
	if err := validateCloseRequest(req); err != nil {
		return nil, err
	}
	runtime, err := s.supervisor.GetRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	s.syslog.Infow("Closing runtime", "runtime_id", req.RuntimeId)
	s.supervisor.CloseRuntime(req.RuntimeId)
	s.syslog.Infow("Closed runtime", "runtime_id", req.RuntimeId)
	runtime.Log().Publish(&builtinv1.SyncPointReachedEvent{BarrierId: req.BarrierId})
	return &executorv1.CloseResponse{}, nil
}

func (s *Server) Introspect(ctx context.Context, req *executorv1.IntrospectRequest) (*executorv1.IntrospectResponse, error) {
	if err := validateIntrospectRequest(req); err != nil {
		return nil, err
	}
	return &executorv1.IntrospectResponse{
		SysInfo:      s.getSysInfo(),
		ExecutorInfo: &executorv1.ExecutorInfo{Name: s.config.Name},
		Labels: append([]string{
			stdruntime.GOOS,
			stdruntime.GOARCH,
		}, s.config.Labels...),
	}, nil
}

func (s *Server) Stop() {
	s.supervisor.Stop()
}

// getSysInfo returns system information about the executor server host.
func (s *Server) getSysInfo() *executorv1.SystemInfo {
	return &executorv1.SystemInfo{
		Os:            stdruntime.GOOS,
		Arch:          stdruntime.GOARCH,
		TotalCpuCores: uint32(stdruntime.NumCPU()),
		TotalMemory:   memory.TotalMemory(),
	}
}

// validateOpenRequest validates an OpenRequest.
func (s *Server) validateOpenRequest(req *executorv1.OpenRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.BuildId == "" {
		return fmt.Errorf("empty build_id")
	}
	if req.Opts == nil {
		return fmt.Errorf("empty opts")
	}
	switch req.Opts.Type {
	case executorv1.RuntimeType_RUNTIME_HOST:
		// NOTE: HostOpts is currently an empty struct
	case executorv1.RuntimeType_RUNTIME_DOCKER:
		if req.Opts.Opts == nil {
			return fmt.Errorf("empty opts")
		}
		dOpts, ok := req.Opts.Opts.(*executorv1.Opts_Docker)
		if !ok {
			return fmt.Errorf("expected Docker opts for runtime type Docker")
		}
		if dOpts.Docker.Image == nil {
			return fmt.Errorf("missing Docker image")
		}
		if dOpts.Docker.Image.ImageUri == "" {
			return fmt.Errorf("missing Docker image uri")
		}
		switch dOpts.Docker.Image.PullStrategy {
		case executorv1.DockerPullOpts_PULL_STRATEGY_NEVER:
		case executorv1.DockerPullOpts_PULL_STRATEGY_ALWAYS:
		case executorv1.DockerPullOpts_PULL_STRATEGY_NOT_EXISTS:
		case executorv1.DockerPullOpts_PULL_STRATEGY_UNSPECIFIED:
		default:
			return fmt.Errorf("unknown Docker pull strategy: %s", dOpts.Docker.Image.PullStrategy)
		}
		if dOpts.Docker.Image.Auth != nil {
			if dOpts.Docker.Image.Auth.Auth == nil {
				return fmt.Errorf("missing Docker image auth")
			}
			switch auth := dOpts.Docker.Image.Auth.Auth.(type) {
			case *executorv1.DockerPullAuth_Basic:
				if auth.Basic == nil {
					return fmt.Errorf("missing Docker image basic auth")
				}
				if auth.Basic.Username == "" {
					return fmt.Errorf("missing Docker image basic username")
				}
				if auth.Basic.Password == "" {
					return fmt.Errorf("missing Docker image basic password")
				}
			case *executorv1.DockerPullAuth_AwsEcr:
				if auth.AwsEcr == nil {
					return fmt.Errorf("missing Docker image aws ecr")
				}
				if auth.AwsEcr.Region == "" {
					return fmt.Errorf("missing Docker image aws ecr region")
				}
				if auth.AwsEcr.AwsAccessKeyId == "" {
					return fmt.Errorf("missing Docker image aws ecr access key id")
				}
				if auth.AwsEcr.AwsSecretKey == "" {
					return fmt.Errorf("missing Docker image aws ecr secret key")
				}
			default:
				return fmt.Errorf("unknown auth type")
			}
		}
	default:
		return fmt.Errorf("unknown type: %v", req.Opts.Type)
	}
	return nil
}

// validateCloseRequest validates a CloseRequest.
func validateCloseRequest(req *executorv1.CloseRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	return nil
}

// validateEventsRequest validates an EventsRequest.
func validateEventsRequest(req *executorv1.EventsRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	return nil
}

// validateExecRequest validates an ExecRequest.
func validateExecRequest(req *executorv1.ExecRequest) error {
	if req == nil {
		return errors.New("nil request")
	}
	if req.RuntimeId == "" {
		return errors.New("empty runtime_id")
	}
	if req.ExecId == "" {
		return errors.New("empty exec_id")
	}
	if req.Opts == nil {
		return errors.New("nil opts")
	}
	if req.Opts.Name == "" {
		return errors.New("empty name")
	}
	return nil
}

// validateFileTransfer validates a FileTransfer.
func validateFileTransfer(req *executorv1.FileTransfer) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	if req.TransferId == "" {
		return fmt.Errorf("empty transfer_id")
	}
	if req.FileId == "" {
		return fmt.Errorf("empty file_id")
	}
	if req.Header != nil {
		if req.Header.DestPath == "" {
			return fmt.Errorf("empty dest_path")
		}
		if req.Header.SrcPath == "" {
			return fmt.Errorf("empty src_path")
		}
	}
	if req.Header == nil && req.Body == nil && req.Trailer == nil {
		return fmt.Errorf("empty header, body, and trailer")
	}
	return nil
}

// validateExportRequest validates an ExportRequest.
func validateExportRequest(req *executorv1.ExportRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	if req.ExportId == "" {
		return fmt.Errorf("empty export_id")
	}
	if req.SrcPath == "" {
		return fmt.Errorf("empty src_path")
	}
	// NOTE: An empty dest path is valid
	return nil
}

// validateIntrospectRequest validates an IntrospectRequest.
func validateIntrospectRequest(req *executorv1.IntrospectRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	return nil
}

// validateHeartbeatRequest validates a HeartbeatRequest.
func validateHeartbeatRequest(req *executorv1.HeartbeatRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	return nil
}
