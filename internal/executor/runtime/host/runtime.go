package host

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"go.uber.org/zap"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/file"
)

type Runtime struct {
	file.WriteFS
	syslog    *zap.SugaredLogger
	runtimeID string
	baseDir   string
	log       *runtime.Log
	deadline  time.Time
}

func NewRuntime(syslog *zap.SugaredLogger, log *runtime.Log, runtimeID string) (*Runtime, error) {
	baseDir, err := os.MkdirTemp("", "knita-host-*")
	if err != nil {
		return nil, fmt.Errorf("error creating runtime base dir: %w", err)
	}
	return &Runtime{
		syslog:    syslog.Named("local_runtime"),
		runtimeID: runtimeID,
		baseDir:   baseDir,
		WriteFS:   file.WriteDirFS(baseDir),
		log:       log,
	}, nil
}

func (r *Runtime) ID() string {
	return r.runtimeID
}

func (r *Runtime) Log() *runtime.Log {
	return r.log
}

func (r *Runtime) Deadline() time.Time {
	return r.deadline
}

func (r *Runtime) SetDeadline(deadline time.Time) {
	r.deadline = deadline
}

func (r *Runtime) Start(ctx context.Context) error {
	return nil
}

func (r *Runtime) Exec(ctx context.Context, execID string, opts *executorv1.ExecOpts) (*runtime.ExecResult, error) {
	r.syslog.Infow("Executing command", "name", opts.Name, "args", opts.Args)
	r.Log().ExecSource(execID, true).Printf("Executing command: %s %v", opts.Name, opts.Args)
	execLog := r.Log().ExecSource(execID, false)

	env := os.Environ()
	env = append(env, opts.Env...)

	cmd := exec.CommandContext(ctx, opts.Name, opts.Args...)
	cmd.Dir = r.baseDir
	cmd.Env = env

	w := execLog.Stdout()
	defer w.Close()
	cmd.Stdout = w

	w = execLog.Stderr()
	defer w.Close()
	cmd.Stderr = w

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &runtime.ExecResult{ExitCode: int32(exitErr.ExitCode())}, nil
		}
		return nil, fmt.Errorf("error running command: %w", err)
	}
	return &runtime.ExecResult{ExitCode: 0}, nil
}

func (r *Runtime) Close() error {
	var res error
	// TODO killall
	if r.baseDir != "" {
		if err := os.RemoveAll(r.baseDir); err != nil {
			res = errors.Join(res, err)
		}
	}
	r.log.Close()
	return res
}
