package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdruntime "runtime"
	"sync"

	"github.com/moby/moby/client"
	"go.uber.org/zap"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/executor/runtime/docker"
	"github.com/knita-io/knita/internal/executor/runtime/host"
)

type runtimeFactory func(ctx context.Context, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error)

type Executor struct {
	executorv1.UnimplementedExecutorServer
	log            *zap.SugaredLogger
	stream         event.Stream
	runtimeFactory runtimeFactory

	mu          sync.RWMutex
	supervisors map[string]*RuntimeSupervisor
}

func NewExecutor(log *zap.SugaredLogger, stream event.Stream) *Executor {
	log = log.Named("executor")
	return &Executor{
		log:         log,
		stream:      stream,
		supervisors: map[string]*RuntimeSupervisor{},
		runtimeFactory: func(ctx context.Context, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error) {
			switch opts.Type {
			case executorv1.RuntimeType_RUNTIME_HOST:
				return host.NewRuntime(log, buildID, runtimeID, stream)
			case executorv1.RuntimeType_RUNTIME_DOCKER:
				dOpts := opts.GetDocker()
				if dOpts == nil {
					return nil, fmt.Errorf("error no docker options provided")
				}
				dClient, err := client.NewClientWithOpts(client.FromEnv)
				if err != nil {
					return nil, fmt.Errorf("error making Docker API client: %w", err)
				}
				dRuntime, err := docker.NewRuntime(log, buildID, runtimeID, stream, dOpts, dClient)
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
}

func (s *Executor) Events(req *executorv1.EventsRequest, stream executorv1.Executor_EventsServer) error {
	s.log.Infow("Event stream opened")
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
		s.log.Infow("Event stream closed")
	case err := <-errC:
		s.log.Infow("Event stream closed: %v", err)
	}
	return nil
}

func (s *Executor) Exec(ctx context.Context, req *executorv1.ExecRequest) (*executorv1.ExecResponse, error) {
	supervisor, err := s.getSupervisor(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	return supervisor.Exec(ctx, req)
}

func (s *Executor) Import(stream executorv1.Executor_ImportServer) error {
	var supervisor *RuntimeSupervisor
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return stream.SendAndClose(&executorv1.ImportResponse{})
			}
			s.log.Errorw("recv error", "error", err)
			return fmt.Errorf("error in receive: %w", err)
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
			s.log.Errorw("import error", "error", err)
			return err
		}
	}
}

func (s *Executor) Export(req *executorv1.ExportRequest, stream executorv1.Executor_ExportServer) error {
	supervisor, err := s.getSupervisor(req.RuntimeId)
	if err != nil {
		return err
	}
	supervisor.runtime.Log().Publish(executorv1.NewExportStartEvent(req.RuntimeId, req.ExportId))
	defer supervisor.runtime.Log().Publish(executorv1.NewExportEndEvent(req.RuntimeId, req.ExportId))
	return supervisor.Export(req, stream)
}

func (s *Executor) Introspect(ctx context.Context, req *executorv1.IntrospectRequest) (*executorv1.IntrospectResponse, error) {
	return &executorv1.IntrospectResponse{
		Os:   stdruntime.GOOS,
		Arch: stdruntime.GOARCH,
		Ncpu: uint32(stdruntime.NumCPU()),
		Labels: []string{
			stdruntime.GOOS,
			stdruntime.GOARCH,
		},
	}, nil
}

func (s *Executor) Open(ctx context.Context, req *executorv1.OpenRequest) (*executorv1.OpenResponse, error) {
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
	rt.Log().Publish(executorv1.NewRuntimeOpenedEvent(req.RuntimeId, req.Opts))
	err = rt.Start(ctx)
	if err != nil {
		rt.Close()
		rt.Log().Publish(executorv1.NewRuntimeClosedEvent(req.RuntimeId))
		return nil, fmt.Errorf("error starting runtime: %w", err)
	}
	s.supervisors[req.RuntimeId] = NewRuntimeSupervisor(s.log, s.stream, rt)
	s.log.Infow("Initialized runtime", "runtime_id", req.RuntimeId)
	return &executorv1.OpenResponse{WorkDirectory: rt.Directory()}, nil
}

func (s *Executor) Close(ctx context.Context, req *executorv1.CloseRequest) (*executorv1.CloseResponse, error) {
	s.mu.Lock()
	rt, ok := s.supervisors[req.RuntimeId]
	if !ok {
		return nil, fmt.Errorf("error runtime not found")
	}
	delete(s.supervisors, req.RuntimeId)
	s.mu.Unlock()
	s.log.Infow("Closing runtime", "runtime_id", req.RuntimeId)
	defer rt.runtime.Log().Publish(executorv1.NewRuntimeClosedEvent(req.RuntimeId))
	return rt.Close(ctx, req)
}

func (s *Executor) getSupervisor(id string) (*RuntimeSupervisor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	supervisor, ok := s.supervisors[id]
	if !ok {
		return nil, fmt.Errorf("error supervisor not found")
	}
	return supervisor, nil
}
