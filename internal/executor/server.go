package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	stdruntime "runtime"
	"sync"
	"time"

	"github.com/moby/moby/client"
	"github.com/pbnjay/memory"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"

	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/executor/runtime/docker"
	"github.com/knita-io/knita/internal/executor/runtime/host"
	"github.com/knita-io/knita/internal/file"
)

const deadlineExtensionPeriod = time.Minute * 2

type runtimeFactory func(ctx context.Context, log *runtime.Log, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error)

type Config struct {
	// Name is the name of the executor.
	Name string
	// Labels the executor will advertise to the broker.
	Labels []string
}

type Server struct {
	executorv1.UnimplementedExecutorServer
	syslog         *zap.SugaredLogger
	config         Config
	runtimeFactory runtimeFactory
	ctx            context.Context
	ctxCancel      context.CancelFunc
	mu             sync.RWMutex
	runtimes       map[string]*runtimeState
}

func NewServer(syslog *zap.SugaredLogger, config Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	syslog = syslog.Named("executor")
	exec := &Server{
		syslog:    syslog,
		config:    config,
		ctx:       ctx,
		ctxCancel: cancel,
		runtimes:  map[string]*runtimeState{}}
	exec.runtimeFactory = exec.defaultRuntimeFactory()
	go exec.watchdog()
	return exec
}

func (s *Server) Events(req *executorv1.EventsRequest, stream executorv1.Executor_EventsServer) error {
	if err := validateEventsRequest(req); err != nil {
		return err
	}
	syslog := s.syslog.With("runtime_id", req.RuntimeId)
	runtime, created := s.findOrCreateRuntime(req.BuildId, req.RuntimeId)
	if !created {
		return fmt.Errorf("event stream already opened")
	}
	defer s.Close(context.Background(), &executorv1.CloseRequest{RuntimeId: req.RuntimeId})
	var (
		closed   bool
		sendErrC = make(chan error)
	)
	fail := func(err error) {
		if !closed {
			closed = true
			sendErrC <- err
		}
	}
	done := runtime.log.Stream().Subscribe(func(event *event.Event) {
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
	runtime.log.MustPublish(&builtinv1.SyncPointReachedEvent{BarrierId: req.BarrierId})
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
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	if runtime.IsOpen() {
		return nil, fmt.Errorf("error runtime already open")
	}
	inner, err := s.runtimeFactory(ctx, runtime.log, req.BuildId, req.RuntimeId, req.Opts)
	if err != nil {
		return nil, fmt.Errorf("error creating runtime: %w", err)
	}
	err = inner.Start(ctx)
	if err != nil {
		inner.Close()
		return nil, fmt.Errorf("error starting runtime: %w", err)
	}
	runtime.SetDeadline(time.Now().Add(deadlineExtensionPeriod))
	runtime.Open(inner)
	s.syslog.Infow("Opened runtime", "runtime_id", req.RuntimeId)
	return &executorv1.OpenResponse{WorkDirectory: runtime.Directory(), SysInfo: s.getSysInfo()}, nil
}

func (s *Server) Exec(ctx context.Context, req *executorv1.ExecRequest) (*executorv1.ExecResponse, error) {
	if err := validateExecRequest(req); err != nil {
		return nil, err
	}
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	if !runtime.IsOpen() {
		return nil, fmt.Errorf("runtime not found")
	}
	res, err := runtime.Exec(ctx, req.ExecId, req.Opts)
	if err != nil {
		return nil, err
	}
	runtime.log.MustPublish(&builtinv1.SyncPointReachedEvent{BarrierId: req.BarrierId})
	return &executorv1.ExecResponse{ExitCode: res.ExitCode}, nil
}

func (s *Server) Import(stream executorv1.Executor_ImportServer) error {
	var (
		runtime   *runtimeState
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
			runtime, err = s.getRuntime(req.RuntimeId)
			if err != nil {
				return err
			}
			if !runtime.IsOpen() {
				return fmt.Errorf("runtime not found")
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
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return err
	}
	if !runtime.IsOpen() {
		return fmt.Errorf("runtime not found")
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
	sendOpts := []file.SendOpt{
		file.WithSendCallback(sendCallback),
		file.WithSkipCallback(skipCallback)}
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

func (s *Server) Heartbeat(ctx context.Context, req *executorv1.HeartbeatRequest) (*executorv1.HeartbeatResponse, error) {
	if err := validateHeartbeatRequest(req); err != nil {
		return nil, err
	}
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(deadlineExtensionPeriod)
	runtime.SetDeadline(deadline)
	s.syslog.Debugf("Extended runtime %s deadline to: %s", runtime.ID(), deadline)
	return &executorv1.HeartbeatResponse{ExtendedBy: durationpb.New(deadlineExtensionPeriod)}, nil
}

func (s *Server) Close(ctx context.Context, req *executorv1.CloseRequest) (*executorv1.CloseResponse, error) {
	if err := validateCloseRequest(req); err != nil {
		return nil, err
	}
	runtime, removed := s.findAndRemoveRuntime(req.RuntimeId)
	if !removed {
		return nil, fmt.Errorf("error runtime not found")
	}
	s.syslog.Infow("Closing runtime", "runtime_id", req.RuntimeId)
	err := runtime.Close()
	if err != nil {
		s.syslog.Errorw("Ignoring error closing runtime: %v", err)
	}
	runtime.Log().MustPublish(&builtinv1.SyncPointReachedEvent{BarrierId: req.BarrierId})
	return &executorv1.CloseResponse{}, nil
}

func (s *Server) Stop() {
	s.ctxCancel()
	var runtimeIDs []string
	s.mu.Lock()
	for id, _ := range s.runtimes {
		runtimeIDs = append(runtimeIDs, id)
	}
	s.mu.Unlock()
	for _, id := range runtimeIDs {
		s.Close(context.Background(), &executorv1.CloseRequest{RuntimeId: id})
	}
}

// getRuntime returns the runtime with the specified ID.
// If the runtime is not found, it returns an error.
func (s *Server) getRuntime(id string) (*runtimeState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runtime, ok := s.runtimes[id]
	if !ok {
		return nil, fmt.Errorf("error runtime not found")
	}
	return runtime, nil
}

// findOrCreateRuntime tries to find the runtime with the specified ID.
// If the runtime is found, it returns the runtime and "false" indicating that it already exists.
// If the runtime is not found, it creates a new one, adds it to the list of runtimes, and then returns the new runtime and "true".
func (s *Server) findOrCreateRuntime(buildID string, runtimeID string) (*runtimeState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	runtime, ok := s.runtimes[runtimeID]
	if ok {
		return runtime, false
	}
	runtime = newRuntimeState(s.syslog, buildID, runtimeID)
	s.runtimes[runtimeID] = runtime
	return runtime, true
}

// findAndRemoveRuntime finds and removes the runtime with the specified runtime ID from the list of runtimes if it is open.
// If the runtime is found, it returns the runtime and "true" indicating that it was removed.
// If the runtime is not found, or it is found but not opened, it returns nil and "false".
func (s *Server) findAndRemoveRuntime(runtimeID string) (*runtimeState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	runtime, ok := s.runtimes[runtimeID]
	if !ok {
		return nil, false
	}
	if !runtime.IsOpen() {
		return nil, false
	}
	delete(s.runtimes, runtimeID)
	return runtime, true
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

// watchdog continuously monitors the deadlines of runtime terminates any runtime that exceeds its deadline.
func (s *Server) watchdog() {
	for {
		s.mu.Lock()
		var (
			deadlinedRuntimes []string
			wakeupIn          = deadlineExtensionPeriod
		)
		for id, runtime := range s.runtimes {
			deadline := runtime.Deadline()
			if deadline.Before(time.Now()) {
				deadlinedRuntimes = append(deadlinedRuntimes, id)
				continue
			}
			deadlineIn := (time.Now().Sub(deadline)) + 1
			if deadlineIn < wakeupIn {
				wakeupIn = deadlineIn
			}
		}
		s.mu.Unlock()
		for _, id := range deadlinedRuntimes {
			s.syslog.Warnf("Runtime %s has deadlined", id)
			s.Close(context.Background(), &executorv1.CloseRequest{RuntimeId: id})
		}
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(wakeupIn):
		}
	}
}

func (s *Server) defaultRuntimeFactory() runtimeFactory {
	return func(ctx context.Context, log *runtime.Log, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error) {
		switch opts.Type {
		case executorv1.RuntimeType_RUNTIME_HOST:
			return host.NewRuntime(s.syslog, log, runtimeID)
		case executorv1.RuntimeType_RUNTIME_DOCKER:
			dOpts := opts.GetDocker()
			if dOpts == nil {
				return nil, fmt.Errorf("error no docker options provided")
			}
			dClient, err := client.NewClientWithOpts(client.FromEnv)
			if err != nil {
				return nil, fmt.Errorf("error making Docker API client: %w", err)
			}
			dRuntime, err := docker.NewRuntime(s.syslog, log, runtimeID, dOpts, dClient)
			if err != nil {
				dClient.Close()
				return nil, fmt.Errorf("error creating Docker runtime: %w", err)
			}
			return dRuntime, nil
		default:
			return nil, fmt.Errorf("error unsupported runtime: %T", opts.Type)
		}
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
