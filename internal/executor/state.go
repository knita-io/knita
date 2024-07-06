package executor

import (
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"go.uber.org/zap"
	"sync"
	"time"
)

type runtimeState struct {
	runtime.Runtime
	mu       sync.RWMutex
	open     bool
	stream   event.Stream
	deadline time.Time
}

func newRuntimeState(syslog *zap.SugaredLogger) *runtimeState {
	return &runtimeState{
		stream: event.NewBroker(syslog),
	}
}

func (r *runtimeState) IsOpen() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.open
}

func (r *runtimeState) Open(runtime runtime.Runtime) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.open {
		panic("already open")
	}
	r.Runtime = runtime
	r.open = true
}

func (r *runtimeState) Deadline() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.deadline
}

func (r *runtimeState) SetDeadline(deadline time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deadline = deadline
}
