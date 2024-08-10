package director

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	directorv1 "github.com/knita-io/knita/api/director/v1"
	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	"github.com/knita-io/knita/internal/event"
)

type Server struct {
	syslog   *zap.SugaredLogger
	build    *Build
	mu       sync.RWMutex
	runtimes map[string]*Runtime
	directorv1.UnimplementedDirectorServer
}

func NewServer(syslog *zap.SugaredLogger, build *Build) *Server {
	return &Server{
		syslog:   syslog,
		build:    build,
		runtimes: map[string]*Runtime{},
	}
}

// Open opens a new runtime.
func (s *Server) Open(ctx context.Context, req *directorv1.OpenRequest) (*directorv1.OpenResponse, error) {
	if err := s.validateOpenRequest(req); err != nil {
		return nil, err
	}
	runtime, err := s.build.OpenRuntime(ctx, req.Opts)
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

// Exec executes a command inside the specified runtime and streams the output back to the client.
func (s *Server) Exec(req *directorv1.ExecRequest, stream directorv1.Director_ExecServer) error {
	if err := validateExecRequest(req); err != nil {
		return err
	}
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return err
	}
	var (
		closed bool
		execID = uuid.New().String()
	)
	done := s.build.Log().Stream().Subscribe(func(event *event.Event) {
		if closed {
			return
		}
		execEvent := &directorv1.ExecEvent{}
		switch p := event.Payload.(type) {
		case *builtinv1.ExecStartEvent:
			if p.RuntimeId != runtime.ID() || p.ExecId != execID {
				return
			}
			execEvent.Payload = &directorv1.ExecEvent_ExecStart{ExecStart: &directorv1.ExecStartEvent{}}
		case *builtinv1.StdoutEvent:
			src, ok := p.Source.Source.(*builtinv1.LogEventSource_Exec)
			if !ok || src.Exec.RuntimeId != runtime.ID() || src.Exec.ExecId != execID || src.Exec.System {
				return
			}
			execEvent.Payload = &directorv1.ExecEvent_Stdout{Stdout: &directorv1.ExecStdoutEvent{Data: p.Data}}
		case *builtinv1.StderrEvent:
			src, ok := p.Source.Source.(*builtinv1.LogEventSource_Exec)
			if !ok || src.Exec.RuntimeId != runtime.ID() || src.Exec.ExecId != execID || src.Exec.System {
				return
			}
			execEvent.Payload = &directorv1.ExecEvent_Stderr{Stderr: &directorv1.ExecStderrEvent{Data: p.Data}}
		case *builtinv1.ExecEndEvent:
			if p.RuntimeId != runtime.ID() || p.ExecId != execID {
				return
			}
			switch s := p.Status.(type) {
			case *builtinv1.ExecEndEvent_Result:
				execEvent.Payload = &directorv1.ExecEvent_ExecEnd{
					ExecEnd: &directorv1.ExecEndEvent{ExitCode: s.Result.ExitCode}}
			case *builtinv1.ExecEndEvent_Error:
				execEvent.Payload = &directorv1.ExecEvent_ExecEnd{
					ExecEnd: &directorv1.ExecEndEvent{Error: s.Error.Message}}
			}
		default:
			return
		}
		s.syslog.Debugf("Forwarded exec event to SDK: %T", execEvent.Payload)
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

// Import files and directories from the local filesystem to the remote runtime.
func (s *Server) Import(ctx context.Context, req *directorv1.ImportRequest) (*directorv1.ImportResponse, error) {
	if err := validateImportRequest(req); err != nil {
		return nil, err
	}
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	err = runtime.Import(ctx, req.SrcPath, req.Opts)
	if err != nil {
		return nil, err
	}
	return &directorv1.ImportResponse{}, nil
}

// Export files and directories from the remote runtime to the local filesystem.
func (s *Server) Export(ctx context.Context, req *directorv1.ExportRequest) (*directorv1.ExportResponse, error) {
	if err := validateExportRequest(req); err != nil {
		return nil, err
	}
	runtime, err := s.getRuntime(req.RuntimeId)
	if err != nil {
		return nil, err
	}
	err = runtime.Export(ctx, req.SrcPath, req.Opts)
	if err != nil {
		return nil, err
	}
	return &directorv1.ExportResponse{}, nil
}

// Close the runtime. The runtime cannot be reused after a call to close.
func (s *Server) Close(ctx context.Context, req *directorv1.CloseRequest) (*directorv1.CloseResponse, error) {
	if err := validateCloseRequest(req); err != nil {
		return nil, err
	}
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
	return &directorv1.CloseResponse{}, nil
}

// getRuntime returns the runtime with the specified ID.
func (s *Server) getRuntime(runtimeID string) (*Runtime, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runtime, ok := s.runtimes[runtimeID]
	if !ok {
		return nil, fmt.Errorf("error runtime not found: %s", runtimeID)
	}
	return runtime, nil
}

// validateOpenRequest validates an OpenRequest.
func (s *Server) validateOpenRequest(req *directorv1.OpenRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.BuildId == "" {
		return fmt.Errorf("empty build_id")
	}
	if req.BuildId != s.build.BuildID() {
		return fmt.Errorf("invalid build_id")
	}
	if req.Opts == nil {
		return fmt.Errorf("empty opts")
	}
	return nil
}

// validateExecRequest validates an ExecRequest.
func validateExecRequest(req *directorv1.ExecRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	if req.Opts == nil {
		return fmt.Errorf("empty opts")
	}
	if req.Opts.Name == "" {
		return fmt.Errorf("empty opts name")
	}
	return nil
}

// validateImportRequest validates an ImportRequest.
func validateImportRequest(req *directorv1.ImportRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	if req.SrcPath == "" {
		return fmt.Errorf("empty source path")
	}
	// NOTE: An empty dest path is valid
	return nil
}

// validateExportRequest validates an ExportRequest.
func validateExportRequest(req *directorv1.ExportRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	if req.SrcPath == "" {
		return fmt.Errorf("empty src_path")
	}
	// NOTE: An empty dest path is valid
	return nil
}

// validateCloseRequest validates a CloseRequest.
func validateCloseRequest(req *directorv1.CloseRequest) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if req.RuntimeId == "" {
		return fmt.Errorf("empty runtime_id")
	}
	return nil
}
