package director

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	directorv1 "github.com/knita-io/knita/api/director/v1"
	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/file"
)

const heartbeatTimeout = time.Second * 5

type Runtime struct {
	syslog              *zap.SugaredLogger
	log                 *Log
	opts                *executorv1.Opts
	buildID             string
	runtimeID           string
	localWorkFS         file.WriteFS
	client              executorv1.ExecutorClient
	ctx                 context.Context
	cancel              context.CancelFunc
	remoteWorkDirectory string
	remoteSysInfo       *executorv1.SystemInfo
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

// ID returns the unique ID of the runtime.
func (c *Runtime) ID() string {
	return c.runtimeID
}

// WorkDirectory returns the fully qualified remote work directory of the runtime.
// Specify a relative path to have it joined to the work directory.
// This is helpful when exec'ing commands that reference file paths within the runtime.
func (c *Runtime) WorkDirectory(path string) string {
	if path == "" {
		return c.remoteWorkDirectory
	}
	return filepath.Join(c.remoteWorkDirectory, path)
}

// SysInfo returns information about the runtime execution environment.
func (c *Runtime) SysInfo() *executorv1.SystemInfo {
	return c.remoteSysInfo
}

// Import files and directories from the local filesystem to the remote runtime.
// Events associated with the import will be published to the configured event stream.
func (c *Runtime) Import(ctx context.Context, src string, opts *directorv1.ImportOpts) error {
	stream, err := c.client.Import(ctx)
	if err != nil {
		return fmt.Errorf("error opening import stream: %w", err)
	}
	c.syslog.Infow("Import stream opened", "src", src)
	importID := uuid.New().String()
	c.log.MustPublish(builtinv1.NewImportStartEvent(c.runtimeID, importID))
	defer c.log.MustPublish(builtinv1.NewImportEndEvent(c.runtimeID, importID))
	sendCallback := func(header *executorv1.FileTransferHeader) {
		if header.IsDir {
			c.log.Printf("Imported directory src=%s, dest=%s, mode=%s", header.SrcPath, header.DestPath, os.FileMode(header.Mode))
		} else {
			c.log.Printf("Imported file src=%s, dest=%s, mode=%s, size=%d", header.SrcPath, header.DestPath, os.FileMode(header.Mode), header.Size)
		}
	}
	skipCallback := func(path string, isDir bool, excludedBy string) {
		if isDir {
			c.log.Printf("Skipped directory import src=%s, excluded_by=%s", path, excludedBy)
		} else {
			c.log.Printf("Skipped file import src=%s, excluded_by=%s", path, excludedBy)
		}
	}
	sendOpts := []file.SendOpt{
		file.WithSendCallback(sendCallback),
		file.WithSkipCallback(skipCallback)}
	if opts != nil {
		if len(opts.Excludes) > 0 {
			sendOpts = append(sendOpts, file.WithExcludes(opts.Excludes))
		}
		if opts.DestPath != "" {
			sendOpts = append(sendOpts, file.WithDest(opts.DestPath))
		}
	}
	sender := file.NewSender(c.syslog, c.localWorkFS.ReadFS(), stream, c.runtimeID, importID, sendOpts...)
	_, err = sender.Send(src)
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

// Export files and directories from the remote runtime to the local filesystem.
// Events associated with the export will be published to the configured event stream.
func (c *Runtime) Export(ctx context.Context, src string, opts *directorv1.ExportOpts) error {
	req := &executorv1.ExportRequest{
		RuntimeId: c.runtimeID,
		ExportId:  uuid.New().String(),
		SrcPath:   src,
		Opts:      &executorv1.ExportOpts{},
	}
	if opts != nil {
		req.Opts.DestPath = opts.DestPath
		req.Opts.Excludes = opts.Excludes
	}
	c.log.MustPublish(builtinv1.NewExportStartEvent(c.runtimeID, req.ExportId))
	defer c.log.MustPublish(builtinv1.NewExportEndEvent(c.runtimeID, req.ExportId))
	stream, err := c.client.Export(ctx, req)
	if err != nil {
		return fmt.Errorf("error opening export stream: %w", err)
	}
	c.syslog.Infow("Export stream opened", "src", src)
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

// Exec executes a command inside the runtime.
// Events associated with the exec will be published to the configured event stream.
func (c *Runtime) Exec(ctx context.Context, execID string, opts *executorv1.ExecOpts) (res *executorv1.ExecResponse, err error) {
	barrier, cancel := event.NewBarrier(c.log.Stream())
	defer cancel()
	defer func() {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		exitCode := int32(-1)
		if res != nil {
			exitCode = res.ExitCode
		}
		c.log.MustPublish(builtinv1.NewExecEndEvent(c.runtimeID, execID, errMsg, exitCode))
	}()
	c.log.MustPublish(builtinv1.NewExecStartEvent(c.runtimeID, execID, opts))
	res, err = c.client.Exec(ctx, &executorv1.ExecRequest{RuntimeId: c.runtimeID, ExecId: execID, BarrierId: barrier.ID(), Opts: opts})
	if err != nil {
		err = fmt.Errorf("error in exec: %w", err)
		return nil, err
	}
	c.syslog.Debugf("Waiting for exec sync point...")
	if err = barrier.Wait(ctx); err != nil {
		err = fmt.Errorf("error waiting for sync point: %w", err)
		return nil, err
	}
	return res, nil
}

// Start the runtime. A runtime must be started prior to use.
func (c *Runtime) Start(ctx context.Context) error {
	barrier, cancel := event.NewBarrier(c.log.Stream())
	defer cancel()
	backgroundCtx, cancel := context.WithCancel(context.Background())
	stream, err := c.client.Events(backgroundCtx, &executorv1.EventsRequest{BuildId: c.buildID, RuntimeId: c.runtimeID, BarrierId: barrier.ID()})
	if err != nil {
		cancel()
		return fmt.Errorf("error opening event stream: %w", err)
	}
	c.syslog.Info("Started streaming runtime events")
	c.ctx = backgroundCtx
	c.cancel = cancel
	go c.forwardEvents(stream)
	c.syslog.Debugf("Waiting for events sync point...")
	if err := barrier.Wait(ctx); err != nil {
		cancel()
		return fmt.Errorf("error waiting for sync point: %w", err)
	}
	c.syslog.Infow("Opening remote runtime...")
	c.log.MustPublish(builtinv1.NewRuntimeOpenStartEvent(c.runtimeID, c.opts))
	openRes, err := c.client.Open(ctx, &executorv1.OpenRequest{BuildId: c.buildID, RuntimeId: c.runtimeID, Opts: c.opts})
	if err != nil {
		c.cancel = nil
		cancel()
		c.log.MustPublish(builtinv1.NewRuntimeOpenEndEvent(c.runtimeID))
		c.log.MustPublish(builtinv1.NewRuntimeCloseStartEvent(c.runtimeID))
		c.log.MustPublish(builtinv1.NewRuntimeCloseEndEvent(c.runtimeID))
		return fmt.Errorf("error opening remote runtime: %w", err)
	}
	go c.keepalive()
	c.syslog.Infow("Opened runtime")
	c.remoteWorkDirectory = openRes.WorkDirectory
	c.remoteSysInfo = openRes.SysInfo
	c.log.MustPublish(builtinv1.NewRuntimeOpenEndEvent(c.runtimeID))
	return nil
}

// Close the runtime. The runtime cannot be reused after a call to close.
func (c *Runtime) Close(ctx context.Context) error {
	barrier, cancel := event.NewBarrier(c.log.Stream())
	defer cancel()
	c.log.MustPublish(builtinv1.NewRuntimeCloseStartEvent(c.runtimeID))
	defer c.log.MustPublish(builtinv1.NewRuntimeCloseEndEvent(c.runtimeID))
	_, err := c.client.Close(ctx, &executorv1.CloseRequest{RuntimeId: c.runtimeID, BarrierId: barrier.ID()})
	if err != nil {
		return fmt.Errorf("error closing runtime: %w", err)
	}
	c.syslog.Debugf("Waiting for close sync point...")
	if err := barrier.Wait(ctx); err != nil {
		return fmt.Errorf("error waiting for sync point: %w", err)
	}
	return nil
}

// forwardEvents receives runtime events from the executor and republishes them to the local event stream.
// Cancelling c.ctx will exit the forwarding loop.
func (c *Runtime) forwardEvents(stream executorv1.Executor_EventsClient) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				c.syslog.Infow("Event stream closed")
			} else {
				c.syslog.Errorw("Event stream closed with error", "error", err)
			}
			return
		} else {
			event, err := event.Unmarshal(msg)
			if err != nil {
				c.syslog.Errorf("Ignoring error unmarshalling event: %v", err)
			} else {
				c.syslog.Debugw("Event received", "type", fmt.Sprintf("%T", event.Payload), "sequence", event.Meta.Sequence)
				c.log.MustRepublish(event)
			}
		}
	}
}

// keepalive periodically sends heartbeats to the executor to keep the runtime alive.
// Cancelling c.ctx will exit the heartbeat loop.
func (c *Runtime) keepalive() {
	// NOTE: All this is trying to achieve is to keep the runtime alive.
	// If the runtime closes for any reason, even if due to a keepalive failure, we expect the next
	// interaction with the runtime to fail and to cause the director to unfold.
	defer c.syslog.Info("Keepalive finished")
	for c.ctx.Err() == nil {
		func() {
			ctx, cancel := context.WithTimeout(c.ctx, heartbeatTimeout)
			defer cancel()
			start := time.Now()
			res, err := c.client.Heartbeat(ctx, &executorv1.HeartbeatRequest{RuntimeId: c.runtimeID})
			if err != nil {
				if c.ctx.Err() == nil {
					c.syslog.Warnf("Will retry error heartbeating runtime: %v", err)
					time.Sleep(time.Second * 5)
				} else {
					return
				}
			} else {
				c.syslog.Debugf("Extended runtime deadline by: %d seconds", res.GetExtendedBy().Seconds)
				remaining := float64(res.GetExtendedBy().GetSeconds()) - time.Now().Sub(start).Seconds()
				time.Sleep(time.Duration(remaining/3) * time.Second)
			}
		}()
	}
}
