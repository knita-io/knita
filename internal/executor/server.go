package executor

import (
	"context"
	"errors"
	"fmt"
	"github.com/moby/moby/client"
	"github.com/pbnjay/memory"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
	"io"
	stdruntime "runtime"
	"sync"
	"time"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/executor/runtime/docker"
	"github.com/knita-io/knita/internal/executor/runtime/host"
)

const deadlineExtensionPeriod = time.Minute * 2

type runtimeFactory func(ctx context.Context, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error)

type Config struct {
	// Name is the name of the executor.
	Name string
	// Labels the executor will advertise to the broker.
	Labels []string
}

type Server struct {
	executorv1.UnimplementedExecutorServer
	syslog         *zap.SugaredLogger
	stream         event.Stream
	runtimeFactory runtimeFactory
	config         Config
	ctx            context.Context
	ctxCancel      context.CancelFunc

	mu          sync.RWMutex
	supervisors map[string]*RuntimeSupervisor
}

func NewServer(syslog *zap.SugaredLogger, config Config, stream event.Stream) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	syslog = syslog.Named("executor")
	exec := &Server{
		syslog:      syslog,
		stream:      stream,
		config:      config,
		ctx:         ctx,
		ctxCancel:   cancel,
		supervisors: map[string]*RuntimeSupervisor{},
		runtimeFactory: func(ctx context.Context, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error) {
			switch opts.Type {
			case executorv1.RuntimeType_RUNTIME_HOST:
				return host.NewRuntime(syslog, buildID, runtimeID, stream)
			case executorv1.RuntimeType_RUNTIME_DOCKER:
				dOpts := opts.GetDocker()
				if dOpts == nil {
					return nil, fmt.Errorf("error no docker options provided")
				}
				dClient, err := client.NewClientWithOpts(client.FromEnv)
				if err != nil {
					return nil, fmt.Errorf("error making Docker API client: %w", err)
				}
				dRuntime, err := docker.NewRuntime(syslog, buildID, runtimeID, stream, dOpts, dClient)
				if err != nil {
					dClient.Close()
					return nil, fmt.Errorf("error creating Docker runtime: %w", err)
				}
				return dRuntime, nil
			default:
				return nil, fmt.Errorf("error unsupported runtime: %T", opts.Type)
			}
		},
	}
	go exec.watchdog()
	return exec
}

func (s *Server) Events(req *executorv1.EventsRequest, stream executorv1.Executor_EventsServer) error {
	if err := validateEventsRequest(req); err != nil {
		return err
	}
	s.syslog.Infow("Event stream opened")
	var (
		closed bool
		errC   = make(chan error)
	)
	done := s.stream.Subscribe(func(event *executorv1.Event) {
		if err := stream.Send(event); err != nil {
			if !closed {
				closed = true
				errC <- err
			}
		}
	}, event.WithPredicate(func(event *executorv1.Event) bool {
		rEvent, ok := event.Payload.(executorv1.RuntimeEvent)
		return ok && rEvent.GetRuntimeId() == req.RuntimeId
	}))
	defer done()
	if err := stream.Send(&executorv1.Event{Sequence: 0}); err != nil {
		return fmt.Errorf("error sending initial event")
	}
	select {
	case <-stream.Context().Done():
		s.syslog.Infow("Event stream closed")
	case err := <-errC:
		s.syslog.Infow("Event stream closed: %v", err)
	}
	return nil
}

func (s *Server) Exec(ctx context.Context, req *executorv1.ExecRequest) (*executorv1.ExecResponse, error) {
	if err := validateExecRequest(req); err != nil {
		return nil, err
	}
	supervisor, err := s.getSupervisor(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	return supervisor.Exec(ctx, req)
}

func (s *Server) Import(stream executorv1.Executor_ImportServer) error {
	var supervisor *RuntimeSupervisor
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
		if supervisor == nil {
			supervisor, err = s.getSupervisor(req.RuntimeId)
			if err != nil {
				return err
			}
			supervisor.runtime.Log().Publish(executorv1.NewImportStartEvent(req.RuntimeId, req.ImportId))
			defer supervisor.runtime.Log().Publish(executorv1.NewImportEndEvent(req.RuntimeId, req.ImportId))
		}
		if supervisor.runtime.ID() != req.RuntimeId {
			return fmt.Errorf("error runtime switcharoo detected")
		}
		err = supervisor.Import(req)
		if err != nil {
			s.syslog.Errorw("import error", "error", err)
			return err
		}
	}
}

func (s *Server) Export(req *executorv1.ExportRequest, stream executorv1.Executor_ExportServer) error {
	if err := validateExportRequest(req); err != nil {
		return err
	}
	supervisor, err := s.getSupervisor(req.RuntimeId)
	if err != nil {
		return err
	}
	supervisor.runtime.Log().Publish(executorv1.NewExportStartEvent(req.RuntimeId, req.ExportId))
	defer supervisor.runtime.Log().Publish(executorv1.NewExportEndEvent(req.RuntimeId, req.ExportId))
	return supervisor.Export(req, stream)
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	supervisor, ok := s.supervisors[req.RuntimeId]
	if !ok {
		return nil, fmt.Errorf("error supervisor not found")
	}
	deadline := time.Now().Add(deadlineExtensionPeriod)
	supervisor.SetDeadline(deadline)
	s.syslog.Debugf("Extended runtime %s deadline to: %s", supervisor.runtime.ID(), deadline)
	return &executorv1.HeartbeatResponse{ExtendedBy: durationpb.New(deadlineExtensionPeriod)}, nil
}

func (s *Server) Open(ctx context.Context, req *executorv1.OpenRequest) (*executorv1.OpenResponse, error) {
	if err := s.validateOpenRequest(req); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.supervisors[req.RuntimeId]
	if ok {
		return nil, fmt.Errorf("error runtime already initialized")
	}
	rt, err := s.runtimeFactory(ctx, req.BuildId, req.RuntimeId, req.Opts)
	if err != nil {
		return nil, fmt.Errorf("error initializing runtime: %w", err)
	}
	rt.Log().Publish(executorv1.NewRuntimeOpenStartEvent(req.RuntimeId, req.Opts))
	err = rt.Start(ctx)
	if err != nil {
		rt.Log().Publish(executorv1.NewRuntimeCloseStartEvent(req.RuntimeId))
		rt.Close()
		rt.Log().Publish(executorv1.NewRuntimeCloseEndEvent(req.RuntimeId))
		return nil, fmt.Errorf("error starting runtime: %w", err)
	}
	rt.Log().Publish(executorv1.NewRuntimeOpenEndEvent(req.RuntimeId))
	supervisor := NewRuntimeSupervisor(s.syslog, s.stream, rt)
	supervisor.SetDeadline(time.Now().Add(deadlineExtensionPeriod))
	s.supervisors[req.RuntimeId] = supervisor
	s.syslog.Infow("Initialized runtime", "runtime_id", req.RuntimeId)
	return &executorv1.OpenResponse{WorkDirectory: rt.Directory(), SysInfo: s.getSysInfo()}, nil
}

func (s *Server) Close(ctx context.Context, req *executorv1.CloseRequest) (*executorv1.CloseResponse, error) {
	if err := validateCloseRequest(req); err != nil {
		return nil, err
	}
	s.mu.Lock()
	rt, ok := s.supervisors[req.RuntimeId]
	if !ok {
		return nil, fmt.Errorf("error runtime not found")
	}
	delete(s.supervisors, req.RuntimeId)
	s.mu.Unlock()
	s.syslog.Infow("Closing runtime", "runtime_id", req.RuntimeId)
	rt.runtime.Log().Publish(executorv1.NewRuntimeCloseStartEvent(req.RuntimeId))
	defer rt.runtime.Log().Publish(executorv1.NewRuntimeCloseEndEvent(req.RuntimeId))
	return rt.Close(ctx, req)
}

func (s *Server) Stop() {
	s.ctxCancel()
	var runtimeIDs []string
	s.mu.Lock()
	for id, _ := range s.supervisors {
		runtimeIDs = append(runtimeIDs, id)
	}
	s.mu.Unlock()
	for _, id := range runtimeIDs {
		s.Close(context.Background(), &executorv1.CloseRequest{RuntimeId: id})
	}
}

// getSupervisor returns the RuntimeSupervisor with the specified ID from the Server's supervisors map.
// If the supervisor is not found, it returns an error.
func (s *Server) getSupervisor(id string) (*RuntimeSupervisor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	supervisor, ok := s.supervisors[id]
	if !ok {
		return nil, fmt.Errorf("error supervisor not found")
	}
	return supervisor, nil
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
		for id, supervisor := range s.supervisors {
			deadline := supervisor.GetDeadline()
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
	if req.ImportId == "" {
		return fmt.Errorf("empty import_id")
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
