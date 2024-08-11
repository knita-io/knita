package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/moby/moby/client"
	"go.uber.org/zap"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/executor/runtime"
	"github.com/knita-io/knita/internal/executor/runtime/docker"
	"github.com/knita-io/knita/internal/executor/runtime/host"
)

const deadlineExtensionPeriod = time.Minute * 2

type runtimeFactory func(ctx context.Context, log *runtime.Log, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error)

type pendingRuntime struct {
	mu  sync.Mutex
	log *runtime.Log
}

func newPendingRuntime(syslog *zap.SugaredLogger, buildID string, runtimeID string) *pendingRuntime {
	stream := event.NewBroker(syslog)
	return &pendingRuntime{log: runtime.NewLog(stream, buildID, runtimeID)}
}

// supervisor manages the lifecycle of runtimes inside an executor.
type supervisor struct {
	syslog          *zap.SugaredLogger
	runtimeFactory  runtimeFactory
	ctx             context.Context
	ctxCancel       context.CancelFunc
	mu              sync.RWMutex
	pendingRuntimes map[string]*pendingRuntime
	openRuntimes    map[string]runtime.Runtime
}

func newSupervisor(syslog *zap.SugaredLogger) *supervisor {
	ctx, cancel := context.WithCancel(context.Background())
	sup := &supervisor{
		syslog:          syslog.Named("supervisor"),
		ctx:             ctx,
		ctxCancel:       cancel,
		pendingRuntimes: map[string]*pendingRuntime{},
		openRuntimes:    map[string]runtime.Runtime{},
	}
	sup.runtimeFactory = defaultRuntimeFactory(syslog)
	go sup.watchdog()
	return sup
}

// PrepareRuntime records the intent to open a new runtime.
// Returns the runtime log, or an error if the runtime has already been prepared.
func (s *supervisor) PrepareRuntime(buildID string, runtimeID string) (*runtime.Log, func(), error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pending, ok := s.pendingRuntimes[runtimeID]
	if ok {
		return nil, nil, fmt.Errorf("runtime %s log is already initialized", runtimeID)
	}
	pending = newPendingRuntime(s.syslog, buildID, runtimeID)
	s.pendingRuntimes[runtimeID] = pending
	return pending.log, func() {
		s.CloseRuntime(runtimeID)
	}, nil
}

// OpenRuntime opens a new runtime. A call to PrepareRuntime must have been made previously.
// Returns an error if called in parallel.
func (s *supervisor) OpenRuntime(ctx context.Context, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error) {
	s.mu.RLock()
	pending, ok := s.pendingRuntimes[runtimeID]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("error pending runtime not found")
	}
	locked := pending.mu.TryLock()
	if !locked {
		return nil, fmt.Errorf("error locking pending runtime")
	}
	defer pending.mu.Unlock()
	runtime, err := s.runtimeFactory(ctx, pending.log, buildID, runtimeID, opts)
	if err != nil {
		return nil, fmt.Errorf("error creating runtime: %w", err)
	}
	runtime.SetDeadline(time.Now().Add(deadlineExtensionPeriod))
	err = runtime.Start(ctx)
	if err != nil {
		runtime.Close()
		return nil, fmt.Errorf("error starting runtime: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pendingRuntimes, runtimeID)
	s.openRuntimes[runtimeID] = runtime
	return runtime, nil
}

// GetRuntime returns the runtime with the specified ID.
// If the runtime is not found, it returns an error.
func (s *supervisor) GetRuntime(id string) (runtime.Runtime, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runtime, ok := s.openRuntimes[id]
	if !ok {
		return nil, fmt.Errorf("error runtime not found")
	}
	return runtime, nil
}

// ExtendRuntime pushes out an open runtimes deadline.
// Returns the amount of time the deadline was extended by.
func (s *supervisor) ExtendRuntime(runtimeID string) (time.Duration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	runtime, ok := s.openRuntimes[runtimeID]
	if !ok {
		return -1, fmt.Errorf("error runtime not found")
	}
	deadline := time.Now().Add(deadlineExtensionPeriod)
	runtime.SetDeadline(deadline)
	s.syslog.Debugf("Extended runtime %s deadline to: %s", runtime.ID(), deadline)
	return deadlineExtensionPeriod, nil
}

// CloseRuntime idempotently closes a runtime (prepared or open).
func (s *supervisor) CloseRuntime(runtimeID string) {
	s.mu.Lock()
	delete(s.pendingRuntimes, runtimeID)
	runtime, ok := s.openRuntimes[runtimeID]
	delete(s.openRuntimes, runtimeID)
	s.mu.Unlock()
	if ok {
		if err := runtime.Close(); err != nil {
			s.syslog.Warnf("Ignoring error closing runtime %s: %v", runtimeID, err)
		}
	}
}

// Stop the supervisor and close all runtimes.
// The supervisor cannot be used again after being closed.
func (s *supervisor) Stop() {
	s.ctxCancel()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, runtime := range s.openRuntimes {
		s.syslog.Infow("Closing runtime", "runtime_id", runtime.ID())
		err := runtime.Close()
		if err != nil {
			s.syslog.Errorf("Ignoring error closing runtime: %v", err)
		}
	}
	s.pendingRuntimes = make(map[string]*pendingRuntime)
	s.openRuntimes = make(map[string]runtime.Runtime)
}

// watchdog continuously monitors runtimes and terminates any runtime that exceed their deadlines.
func (s *supervisor) watchdog() {
	for {
		s.mu.Lock()
		var (
			deadlinedRuntimes []string
			wakeupIn          = deadlineExtensionPeriod
		)
		for id, runtime := range s.openRuntimes {
			deadline := runtime.Deadline()
			if deadline.Before(time.Now()) {
				deadlinedRuntimes = append(deadlinedRuntimes, id)
				continue
			}
			deadlineIn := (time.Now().Sub(deadline)) + 1
			if deadlineIn < wakeupIn {
				wakeupIn = deadlineIn
			}
		}
		s.mu.Unlock()
		for _, id := range deadlinedRuntimes {
			s.syslog.Warnf("Runtime %s has deadlined", id)
			s.CloseRuntime(id)
		}
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(wakeupIn):
		}
	}
}

func defaultRuntimeFactory(syslog *zap.SugaredLogger) runtimeFactory {
	return func(ctx context.Context, log *runtime.Log, buildID string, runtimeID string, opts *executorv1.Opts) (runtime.Runtime, error) {
		switch opts.Type {
		case executorv1.RuntimeType_RUNTIME_HOST:
			return host.NewRuntime(syslog, log, runtimeID)
		case executorv1.RuntimeType_RUNTIME_DOCKER:
			dOpts := opts.GetDocker()
			if dOpts == nil {
				return nil, fmt.Errorf("error no docker options provided")
			}
			dClient, err := client.NewClientWithOpts(client.FromEnv)
			if err != nil {
				return nil, fmt.Errorf("error making Docker API client: %w", err)
			}
			dRuntime, err := docker.NewRuntime(syslog, log, runtimeID, dOpts, dClient)
			if err != nil {
				dClient.Close()
				return nil, fmt.Errorf("error creating Docker runtime: %w", err)
			}
			return dRuntime, nil
		default:
			return nil, fmt.Errorf("error unsupported runtime: %T", opts.Type)
		}
	}
}
