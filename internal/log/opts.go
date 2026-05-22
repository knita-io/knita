package log

import (
	eventsv1 "github.com/knita-io/knita/api/events/v1"
	"github.com/knita-io/knita/internal/event"
)

// Opt configures an Event.
type Opt func(*event.Event)

// WithLabel sets a single label key→value on the Event.
func WithLabel(key, value string) Opt {
	return WithLabels(map[string]string{key: value})
}

// WithLabels sets or merges the given labels into the Event.
func WithLabels(labels map[string]string) Opt {
	return func(r *event.Event) {
		if len(labels) == 0 {
			return
		}
		if r.Meta == nil {
			r.Meta = &eventsv1.Meta{}
		}
		if r.Meta.Labels == nil {
			r.Meta.Labels = make(map[string]string)
		}
		for k, v := range labels {
			r.Meta.Labels[k] = v
		}
	}
}

// WithAnnotation sets a single annotation key→value on the Event.
func WithAnnotation(key, value string) Opt {
	return WithAnnotations(map[string]string{key: value})
}

// WithAnnotations sets or merges the given annotations into the Event.
func WithAnnotations(annotations map[string]string) Opt {
	return func(r *event.Event) {
		if len(annotations) == 0 {
			return
		}
		if r.Meta == nil {
			r.Meta = &eventsv1.Meta{}
		}
		if r.Meta.Annotations == nil {
			r.Meta.Annotations = make(map[string]string)
		}
		for k, v := range annotations {
			r.Meta.Annotations[k] = v
		}
	}
}
