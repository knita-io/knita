package executor

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/file"
)

type RuntimeSupervisor struct {
	log       *zap.SugaredLogger
	stream    event.Stream
	runtime   runtime.Runtime
	mu        sync.RWMutex
	importers map[string]*file.Receiver
}

func NewRuntimeSupervisor(log *zap.SugaredLogger, stream event.Stream, runtime runtime.Runtime) *RuntimeSupervisor {
	return &RuntimeSupervisor{
		log:       log.Named("mediator").With("runtime_id", runtime.ID()),
		stream:    stream,
		runtime:   runtime,
		importers: map[string]*file.Receiver{},
	}
}

func (s *RuntimeSupervisor) Exec(ctx context.Context, req *v1.ExecRequest) (*v1.ExecResponse, error) {
	s.runtime.Log().Publish(v1.NewExecStartEvent(req.RuntimeId, req.ExecId, req.Opts))
	res, err := s.runtime.Exec(ctx, req.ExecId, req.Opts)
	if err != nil {
		s.runtime.Log().Publish(v1.NewExecEndEvent(req.RuntimeId, req.ExecId, err.Error(), -1))
		return nil, err
	}
	s.runtime.Log().Publish(v1.NewExecEndEvent(req.RuntimeId, req.ExecId, "", res.ExitCode))
	return &v1.ExecResponse{ExitCode: res.ExitCode}, nil
}

func (s *RuntimeSupervisor) Import(req *v1.FileTransfer) error {
	s.mu.Lock()
	imp, ok := s.importers[req.FileId]
	if !ok {
		// TODO garbage collect imports if the next request never comes?
		imp = file.NewReceiver(s.log, s.runtime, file.WithRecvCallback(func(header *v1.FileTransferHeader) {
			if header.IsDir {
				s.runtime.Log().Printf("Imported directory src=%s, dest=%s, mode=%s", header.SrcPath, header.DestPath, os.FileMode(header.Mode))
			} else {
				s.runtime.Log().Printf("Imported file src=%s, dest=%s, mode=%s, size=%d", header.SrcPath, header.DestPath, os.FileMode(header.Mode), header.Size)
			}
		}))
		s.importers[req.FileId] = imp
	}
	s.mu.Unlock()
	err := imp.Next(req)
	if err != nil {
		s.mu.Lock()
		delete(s.importers, req.FileId)
		s.mu.Unlock()
		return err
	}
	return nil
}

func (s *RuntimeSupervisor) Export(req *v1.ExportRequest, stream v1.Executor_ExportServer) error {
	sender := file.NewSender(s.log, s.runtime.ReadFS(), s.runtime.ID(), file.WithSendCallback(func(header *v1.FileTransferHeader) {
		if header.IsDir {
			s.runtime.Log().Printf("Exported directory src=%s, dest=%s, mode=%s", header.SrcPath, header.DestPath, os.FileMode(header.Mode))
		} else {
			s.runtime.Log().Printf("Exported file src=%s, dest=%s, mode=%s, size=%d", header.SrcPath, header.DestPath, os.FileMode(header.Mode), header.Size)
		}
	}))
	_, err := sender.Send(stream, req.SrcPath, req.DestPath)
	return err
}

func (s *RuntimeSupervisor) Close(ctx context.Context, req *v1.CloseRequest) (*v1.CloseResponse, error) {
	err := s.runtime.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing runtime: %w", err)
	}
	return &v1.CloseResponse{}, nil
}
