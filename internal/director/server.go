package director

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	directorv1 "github.com/knita-io/knita/api/director/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
)

type Server struct {
	syslog     *zap.SugaredLogger
	stream     event.Stream
	controller *BuildController
	mu         sync.RWMutex
	runtimes   map[string]*Runtime
	directorv1.UnimplementedDirectorServer
}

func NewServer(syslog *zap.SugaredLogger, stream event.Stream, controller *BuildController) *Server {
	return &Server{
		syslog:     syslog,
		stream:     stream,
		controller: controller,
		runtimes:   map[string]*Runtime{},
	}
}

func (s *Server) Open(ctx context.Context, req *directorv1.OpenRequest) (*directorv1.OpenResponse, error) {
	if req.BuildId == "" {
		return nil, fmt.Errorf("error build id must be set")
	}
	if req.BuildId != s.controller.BuildID() {
		return nil, fmt.Errorf("error invalid build id")
	}
	runtime, err := s.controller.Runtime(ctx, req.Opts)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.runtimes[runtime.ID()] = runtime
	s.mu.Unlock()
	return &directorv1.OpenResponse{
		RuntimeId:     runtime.ID(),
		WorkDirectory: runtime.WorkDirectory(""),
		SysInfo:       runtime.SysInfo(),
	}, nil
}

func (s *Server) Exec(req *directorv1.ExecRequest, stream directorv1.Director_ExecServer) error {
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return err
	}
	var (
		closed bool
		execID = uuid.New().String()
	)
	done := s.stream.Subscribe(func(event *executorv1.Event) {
		if closed {
			return
		}
		execEvent := &directorv1.ExecEvent{}
		switch p := event.Payload.(type) {
		case *executorv1.Event_ExecStart:
			if p.ExecStart.RuntimeId != runtime.ID() || p.ExecStart.ExecId != execID {
				return
			}
			execEvent.Payload = &directorv1.ExecEvent_ExecStart{ExecStart: &directorv1.ExecStartEvent{}}
		case *executorv1.Event_Stdout:
			src, ok := p.Stdout.Source.Source.(*executorv1.LogEventSource_Exec)
			if !ok || src.Exec.RuntimeId != runtime.ID() || src.Exec.ExecId != execID || src.Exec.System {
				return
			}
			execEvent.Payload = &directorv1.ExecEvent_Stdout{Stdout: &directorv1.ExecStdoutEvent{
				Data: p.Stdout.Data,
			}}
		case *executorv1.Event_Stderr:
			src, ok := p.Stderr.Source.Source.(*executorv1.LogEventSource_Exec)
			if !ok || src.Exec.RuntimeId != runtime.ID() || src.Exec.ExecId != execID || src.Exec.System {
				return
			}
			execEvent.Payload = &directorv1.ExecEvent_Stderr{Stderr: &directorv1.ExecStderrEvent{
				Data: p.Stderr.Data,
			}}
		case *executorv1.Event_ExecEnd:
			if p.ExecEnd.RuntimeId != runtime.ID() || p.ExecEnd.ExecId != execID {
				return
			}
			execEvent.Payload = &directorv1.ExecEvent_ExecEnd{ExecEnd: &directorv1.ExecEndEvent{
				Error:    p.ExecEnd.Error,
				ExitCode: p.ExecEnd.ExitCode,
			}}
		default:
			return
		}
		s.syslog.Debugf("Forwarded exec event to SDK: %#v", execEvent)
		if err := stream.Send(execEvent); err != nil {
			s.syslog.Errorf("Exec stream closed early: %v", err)
			closed = true
		}
	})
	defer func() {
		done()
		s.syslog.Info("Unsubscribed exec client")
	}()

	_, err = runtime.Exec(stream.Context(), execID, req.Opts)
	return err
}

func (s *Server) Import(ctx context.Context, req *directorv1.ImportRequest) (*directorv1.ImportResponse, error) {
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	err = runtime.Import(ctx, req.SrcPath, req.DestPath)
	if err != nil {
		return nil, err
	}
	return &directorv1.ImportResponse{}, nil
}

func (s *Server) Export(ctx context.Context, req *directorv1.ExportRequest) (*directorv1.ExportResponse, error) {
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	err = runtime.Export(ctx, req.SrcPath, req.DestPath)
	if err != nil {
		return nil, err
	}
	return &directorv1.ExportResponse{}, nil
}

func (s *Server) Close(ctx context.Context, req *executorv1.CloseRequest) (*executorv1.CloseResponse, error) {
	s.mu.Lock()
	runtime, ok := s.runtimes[req.RuntimeId]
	if ok {
		delete(s.runtimes, req.RuntimeId)
	}
	s.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("error runtime not found: %s", req.RuntimeId)
	}
	err := runtime.Close(ctx)
	if err != nil {
		return nil, err
	}
	return &executorv1.CloseResponse{}, nil
}

func (s *Server) getRuntime(runtimeID string) (*Runtime, error) {
	s.mu.RLock()
	runtime, ok := s.runtimes[runtimeID]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("error runtime not found: %s", runtimeID)
	}
	return runtime, nil
}
