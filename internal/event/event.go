package event

import (
	"errors"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"

	eventsv1 "github.com/knita-io/knita/api/events/v1"
)

type Event struct {
	Meta    *eventsv1.Meta
	Payload proto.Message
}

func (e *Event) Marshal() (*eventsv1.Event, error) {
	out := &eventsv1.Event{Meta: e.Meta}
	if p, ok := e.Payload.(*anypb.Any); ok {
		out.Payload = p
	} else {
		p, err := anypb.New(e.Payload)
		if err != nil {
			return nil, err
		}
		out.Payload = p
	}
	return out, nil
}

func Unmarshal(event *eventsv1.Event) (*Event, error) {
	var (
		p   proto.Message
		err error
	)
	if event.Payload != nil {
		p, err = anypb.UnmarshalNew(event.Payload, proto.UnmarshalOptions{})
		if err != nil {
			if errors.Is(err, protoregistry.NotFound) {
				p = event.Payload
			} else {
				return nil, err
			}
		}
	}
	return &Event{Meta: event.Meta, Payload: p}, nil
}
