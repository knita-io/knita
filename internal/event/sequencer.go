package event

import (
	"sync"
)

type Sequencer struct {
	stream       Stream
	mu           sync.Mutex
	lastSequence uint64
}

func NewSequencer(stream Stream) *Sequencer {
	return &Sequencer{stream: stream}
}

func (s *Sequencer) Publish(event *Event) {
	s.mu.Lock()
	s.lastSequence++
	event.Meta.Sequence = s.lastSequence
	s.stream.Publish(event)
	s.mu.Unlock() // hold the lock while we publish to prevent re-ordering
}
