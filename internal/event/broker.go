package event

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Stream interface {
	Publish(event *Event)
	Subscribe(handler Handler, opts ...Opt) func()
}

type Predicate func(event *Event) bool

type Handler func(event *Event)

type Broker struct {
	syslog        *zap.SugaredLogger
	mu            sync.RWMutex
	subscriptions map[*subscription]struct{}
}

type subscription struct {
	handler Handler
	opts    *Opts
}

func NewBroker(syslog *zap.SugaredLogger) *Broker {
	return &Broker{
		syslog:        syslog.Named("event_broker"),
		subscriptions: map[*subscription]struct{}{},
	}
}

func (b *Broker) Publish(event *Event) {
	b.mu.RLock()
	subscriptions := make(map[*subscription]struct{})
	for k, v := range b.subscriptions {
		subscriptions[k] = v
	}
	b.mu.RUnlock()
	delivered := 0
	filtered := 0
	if event.Payload == nil {
		b.syslog.Panicf("event payload is nil: %v", event)
	}
	for sub := range subscriptions {
		send := true
		for _, pred := range sub.opts.Predicates {
			if !pred(event) {
				send = false
				filtered++
				break
			}
		}
		if send {
			sub.handler(event)
			delivered++
		}
	}
	b.syslog.Debugw("Published event",
		"type", fmt.Sprintf("%T", event.Payload),
		"sequence", event.Meta.Sequence,
		"delivered", delivered,
		"filtered", filtered)
}

func (b *Broker) Subscribe(handler Handler, opts ...Opt) func() {
	o := &Opts{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	sub := &subscription{
		handler: handler,
		opts:    o,
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscriptions[sub] = struct{}{}
	b.syslog.Debugw("Registered subscriber")
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if _, ok := b.subscriptions[sub]; ok {
			delete(b.subscriptions, sub)
			b.syslog.Debugw("Unregistered subscriber")
		}
	}
}

type Opts struct {
	Predicates []Predicate
}

type Opt interface {
	Apply(opts *Opts)
}

type withPredicate struct {
	predicate Predicate
}

func (o *withPredicate) Apply(opts *Opts) {
	opts.Predicates = append(opts.Predicates, o.predicate)
}

func WithPredicate(predicate Predicate) Opt {
	return &withPredicate{predicate: predicate}
}
