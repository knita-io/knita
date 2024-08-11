package event

import (
	"context"

	"github.com/google/uuid"

	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
)

type Barrier struct {
	id string
	c  chan struct{}
}

func NewBarrier(stream Stream) (*Barrier, func()) {
	b := &Barrier{id: uuid.NewString(), c: make(chan struct{})}
	var closed bool
	done := stream.Subscribe(func(event *Event) {
		if !closed {
			close(b.c)
			closed = true
		}
	}, WithPredicate(func(event *Event) bool {
		sync, ok := event.Payload.(*builtinv1.SyncPointReachedEvent)
		return ok && sync.BarrierId == b.id
	}))
	return b, done
}

func (b *Barrier) ID() string {
	return b.id
}

func (b *Barrier) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.c:
		return nil
	}
}
