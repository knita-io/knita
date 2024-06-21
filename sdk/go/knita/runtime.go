package knita

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	directorv1 "github.com/knita-io/knita/api/director/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/sdk/go/knita/runtime/exec"
)

// Runtime represents a local handle to a remote runtime hosted by an executor.
type Runtime struct {
	syslog              Log
	fatalFunc           FatalFunc
	runtimeID           string
	remoteWorkDirectory string
	remoteSysInfo       *executorv1.SystemInfo
	client              directorv1.DirectorClient
}

// ID returns the unique ID of the runtime.
func (c *Runtime) ID() string {
	return c.runtimeID
}

// SysInfo returns information about the configuration of the runtime.
func (c *Runtime) SysInfo() *executorv1.SystemInfo {
	return c.remoteSysInfo
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

// Import files from the local work directory into the runtime's remote work directory.
// src and dest must be relative paths. src may be a glob (doublestar syntax supported).
// If dest is empty, all files identified by src will be copied to their original location in dest.
func (c *Runtime) Import(src string, dest string) error {
	return c.ImportWithContext(context.Background(), src, dest)
}

// MustImport is like Import, but it calls the configured FatalFunc if an error occurs.
func (c *Runtime) MustImport(src string, dest string) {
	err := c.Import(src, dest)
	if err != nil {
		c.fatalFunc(fmt.Errorf("error importing: %w", err))
	}
}

// ImportWithContext is like Import, but it allows a context to be set.
func (c *Runtime) ImportWithContext(ctx context.Context, src string, dest string) error {
	_, err := c.client.Import(ctx, &directorv1.ImportRequest{
		RuntimeId: c.runtimeID,
		SrcPath:   src,
		DestPath:  dest,
	})
	return err
}

// Export files from the runtime's remote work directory into the local work directory.
// src and dest must be relative paths. src may be a glob (doublestar syntax supported).
// If dest is empty, all files identified by src will be copied to their original location in dest.
func (c *Runtime) Export(src string, dest string) error {
	return c.ExportWithContext(context.Background(), src, dest)
}

// MustExport is like Export, but it calls the configured FatalFunc if an error occurs.
func (c *Runtime) MustExport(src string, dest string) {
	err := c.Export(src, dest)
	if err != nil {
		c.fatalFunc(fmt.Errorf("error exporting: %w", err))
	}
}

// ExportWithContext is like Export, but it allows a context to be set.
func (c *Runtime) ExportWithContext(ctx context.Context, src string, dest string) error {
	_, err := c.client.Export(ctx, &directorv1.ExportRequest{
		RuntimeId: c.runtimeID,
		SrcPath:   src,
		DestPath:  dest,
	})
	return err
}

// Exec executes a command inside the remote runtime.
// Check the returned ExecResponse to see the command's exit code (a non-zero code is not an error).
func (c *Runtime) Exec(opts ...exec.Opt) (*executorv1.ExecResponse, error) {
	return c.ExecWithContext(context.Background(), opts...)
}

// MustExec is like Exec, but it calls the configured FatalFunc if an error occurs or the command exits with a non-zero exit code.
func (c *Runtime) MustExec(opts ...exec.Opt) *executorv1.ExecResponse {
	res, err := c.Exec(opts...)
	if err != nil {
		c.fatalFunc(fmt.Errorf("error execing: %w", err))
	}
	if res.ExitCode != 0 {
		c.fatalFunc(fmt.Errorf("error non-zero exit code"))
	}
	return res
}

// ExecWithContext is like Exec, but it allows a context to be set.
func (c *Runtime) ExecWithContext(ctx context.Context, opts ...exec.Opt) (*executorv1.ExecResponse, error) {
	o := &exec.Opts{ExecOpts: &executorv1.ExecOpts{}}
	for _, opt := range opts {
		opt.Apply(o)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := c.client.Exec(ctx, &directorv1.ExecRequest{RuntimeId: c.runtimeID, Opts: o.ExecOpts})
	if err != nil {
		return nil, fmt.Errorf("error in exec: %w", err)
	}
	var execEnd *executorv1.ExecResponse
	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if execEnd != nil {
					return execEnd, nil
				}
				return nil, fmt.Errorf("error stream closed before exec end event was observed")
			}
			return nil, err
		}
		switch p := msg.Payload.(type) {
		case *directorv1.ExecEvent_Stdout:
			if o.Stdout != nil {
				o.Stdout.Write(p.Stdout.Data)
			}
		case *directorv1.ExecEvent_Stderr:
			if o.Stderr != nil {
				o.Stderr.Write(p.Stderr.Data)
			}
		case *directorv1.ExecEvent_ExecEnd:
			execEnd = &executorv1.ExecResponse{ExitCode: p.ExecEnd.ExitCode}
		}
	}
}

// Close the runtime. After a call to close the runtime can no longer be used.
func (c *Runtime) Close() error {
	return c.CloseWithContext(context.Background())
}

// MustClose is like Close, but it calls the configured FatalFunc if an error occurs.
func (c *Runtime) MustClose() {
	err := c.Close()
	if err != nil {
		c.fatalFunc(fmt.Errorf("error closing runtime: %w", err))
	}
}

// CloseWithContext is like Close, but it allows a context to be set.
func (c *Runtime) CloseWithContext(ctx context.Context) error {
	_, err := c.client.Close(ctx, &executorv1.CloseRequest{
		RuntimeId: c.runtimeID,
	})
	return err
}
