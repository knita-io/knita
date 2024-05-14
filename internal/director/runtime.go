package director

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/uuid"
	"go.uber.org/zap"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/file"
)

type Runtime struct {
	syslog              *zap.SugaredLogger
	log                 *Log
	opts                *executorv1.Opts
	buildID             string
	runtimeID           string
	localWorkFS         file.WriteFS
	client              executorv1.ExecutorClient
	eventCancel         context.CancelFunc
	remoteWorkDirectory string
}

func newRuntime(
	syslog *zap.SugaredLogger,
	log *Log,
	buildID string,
	runtimeID string,
	client executorv1.ExecutorClient,
	localWorkFS file.WriteFS,
	opts *executorv1.Opts) *Runtime {

	return &Runtime{
		syslog:      syslog.Named("runtime").With("runtime_id", runtimeID),
		log:         log,
		buildID:     buildID,
		runtimeID:   runtimeID,
		client:      client,
		localWorkFS: localWorkFS,
		opts:        opts,
	}
}

func (c *Runtime) start(ctx context.Context) error {
	eCtx, cancel := context.WithCancel(context.Background())
	stream, err := c.client.Events(eCtx, &executorv1.EventsRequest{RuntimeId: c.runtimeID})
	if err != nil {
		cancel()
		return fmt.Errorf("error opening event stream: %w", err)
	}
	initial, err := stream.Recv()
	if err != nil {
		cancel()
		return fmt.Errorf("error waiting for initial event: %w", err)
	}
	if initial.Sequence != 0 {
		cancel()
		return fmt.Errorf("error unexpected initial event")
	}
	c.syslog.Info("Started streaming runtime events")
	c.eventCancel = cancel
	go func() {
		for {
			evt, err := stream.Recv()
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					c.syslog.Infow("Event stream closed")
				} else {
					c.syslog.Errorw("Event stream closed", "error", err)
				}
				return
			} else {
				c.syslog.Debugw("Event received", "type", fmt.Sprintf("%T", evt.Payload), "sequence", evt.Sequence)
				c.log.Republish(evt)
			}
		}
	}()
	c.syslog.Infow("Opening remote runtime...")
	openRes, err := c.client.Open(ctx, &executorv1.OpenRequest{
		BuildId:   c.buildID,
		RuntimeId: c.runtimeID,
		Opts:      c.opts,
	})
	if err != nil {
		c.eventCancel = nil
		cancel()
		return fmt.Errorf("error opening remote runtime: %w", err)
	}
	c.syslog.Infow("Opened runtime")
	c.remoteWorkDirectory = openRes.WorkDirectory
	return nil
}

func (c *Runtime) ID() string {
	return c.runtimeID
}

func (c *Runtime) WorkDirectory(path string) string {
	if path == "" {
		return c.remoteWorkDirectory
	}
	return filepath.Join(c.remoteWorkDirectory, path)
}

func (c *Runtime) Import(ctx context.Context, src string, dest string) error {
	if filepath.IsAbs(src) {
		return fmt.Errorf("error src dir must be relative")
	}
	if filepath.IsAbs(dest) {
		return fmt.Errorf("error dest dir must be relative")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := c.client.Import(ctx)
	if err != nil {
		return fmt.Errorf("error opening import stream: %w", err)
	}
	c.syslog.Infow("Import stream opened", "src", src, "dest", dest)
	sender := file.NewSender(c.syslog, c.localWorkFS.ReadFS(), c.runtimeID)
	_, err = sender.Send(stream, src, dest)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_, err = stream.CloseAndRecv()
		}
		return err
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}
	return err
}

func (c *Runtime) Export(ctx context.Context, src string, dest string) error {
	if filepath.IsAbs(src) {
		return fmt.Errorf("error src dir must be relative")
	}
	if filepath.IsAbs(dest) {
		return fmt.Errorf("error dest dir must be relative")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	exportReq := &executorv1.ExportRequest{
		RuntimeId: c.runtimeID,
		ExportId:  uuid.New().String(),
		SrcPath:   src,
		DestPath:  dest,
	}
	stream, err := c.client.Export(ctx, exportReq)
	if err != nil {
		return fmt.Errorf("error opening export stream: %w", err)
	}
	c.syslog.Infow("Export stream opened", "src", src, "dest", dest)
	receivers := make(map[string]*file.Receiver)
	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// TODO verify all success?
				return nil
			}
			return fmt.Errorf("error receiving from stream: %w", err)
		}
		recv, ok := receivers[msg.FileId]
		if !ok {
			recv = file.NewReceiver(c.syslog, c.localWorkFS)
			receivers[msg.FileId] = recv
		}
		err = recv.Next(msg)
		if err != nil {
			delete(receivers, msg.FileId)
			return err
		}
	}
}

func (c *Runtime) Exec(ctx context.Context, execID string, opts *executorv1.ExecOpts) (*executorv1.ExecResponse, error) {
	doneC := make(chan struct{})
	done := c.log.Stream().Subscribe(func(event *executorv1.Event) {
		close(doneC)
	}, event.WithPredicate(func(event *executorv1.Event) bool {
		closed, ok := event.Payload.(*executorv1.Event_ExecEnd)
		return ok && closed.ExecEnd.RuntimeId == c.runtimeID && closed.ExecEnd.ExecId == execID
	}))
	defer done()
	res, err := c.client.Exec(ctx, &executorv1.ExecRequest{
		RuntimeId: c.runtimeID,
		ExecId:    execID,
		Opts:      opts,
	})
	if err != nil {
		return nil, fmt.Errorf("error in exec: %w", err)
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-doneC:
	}
	return res, nil
}

func (c *Runtime) Close(ctx context.Context) error {
	doneC := make(chan struct{})
	if c.eventCancel != nil {
		done := c.log.Stream().Subscribe(func(event *executorv1.Event) {
			c.eventCancel()
			close(doneC)
		}, event.WithPredicate(func(event *executorv1.Event) bool {
			closed, ok := event.Payload.(*executorv1.Event_RuntimeClosed)
			return ok && closed.RuntimeClosed.RuntimeId == c.runtimeID
		}))
		defer done()
	} else {
		close(doneC)
	}
	_, err := c.client.Close(ctx, &executorv1.CloseRequest{
		RuntimeId: c.runtimeID,
	})
	if err != nil {
		return fmt.Errorf("error closing runtime: %w", err)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-doneC:
	}
	return nil
}
