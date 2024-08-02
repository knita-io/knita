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

func (s *Sequencer) MustPublish(event *Event) {
	s.mu.Lock()
	s.lastSequence++
	event.Meta.Sequence = s.lastSequence
	s.mu.Unlock()
	s.stream.MustPublish(event)
}
