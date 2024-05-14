package runtime

import (
	"context"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/file"
)

type Runtime interface {
	file.WriteFS
	Start(ctx context.Context) error
	ID() string
	Log() *Log
	Exec(ctx context.Context, execID string, opts *executorv1.ExecOpts) (*ExecResult, error)
	Close() error
}

type ExecResult struct {
	ExitCode int32
}
