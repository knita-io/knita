package runtime

import (
	"context"
	"time"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/file"
)

type Runtime interface {
	file.WriteFS
	ID() string
	Deadline() time.Time
	SetDeadline(deadline time.Time)
	Log() *Log
	Start(ctx context.Context) error
	Exec(ctx context.Context, execID string, opts *executorv1.ExecOpts) (*ExecResult, error)
	Close() error
}

type ExecResult struct {
	ExitCode int32
}
