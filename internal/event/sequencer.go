package event

import (
	"sync"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

type Sequencer struct {
	stream       Stream
	mu           sync.Mutex
	lastSequence uint64
}

func NewSequencer(stream Stream) *Sequencer {
	return &Sequencer{stream: stream}
}

func (s *Sequencer) Publish(event *executorv1.Event) {
	s.mu.Lock()
	s.lastSequence++
	event.Sequence = s.lastSequence
	s.mu.Unlock()
	s.stream.Publish(event)
}
