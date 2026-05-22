package exec

import (
	"io"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
)

// Opt configures an Exec call.
type Opt func(*Opts)

// Opts holds the options for an Exec invocation.
type Opts struct {
	*executorv1.ExecOpts
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

// WithCommand specifies the command name and args.
func WithCommand(name string, args ...string) Opt {
	return func(o *Opts) {
		if o.ExecOpts == nil {
			o.ExecOpts = &executorv1.ExecOpts{}
		}
		o.Name = name
		o.Args = args
	}
}

// WithStdout directs stdout to the provided writer.
func WithStdout(w io.Writer) Opt {
	return func(o *Opts) {
		o.Stdout = w
	}
}

// WithStderr directs stderr to the provided writer.
func WithStderr(w io.Writer) Opt {
	return func(o *Opts) {
		o.Stderr = w
	}
}

// WithEnv adds environment variables (e.g. "KEY=VALUE").
func WithEnv(env ...string) Opt {
	return func(o *Opts) {
		if o.ExecOpts == nil {
			o.ExecOpts = &executorv1.ExecOpts{}
		}
		o.Env = append(o.Env, env...)
	}
}

// WithDisplayName sets the display name for the exec.
func WithDisplayName(displayName string) Opt {
	return func(o *Opts) {
		o.DisplayName = displayName
	}
}

// WithLabel sets a single label.
func WithLabel(key, value string) Opt {
	return WithLabels(key, value)
}

// WithLabels sets multiple labels from an alternating key/value list.
// Panics if you pass an odd number of args.
func WithLabels(kv ...string) Opt {
	m := runtime.KVMap("WithLabels", kv)
	return func(o *Opts) {
		if o.ExecOpts == nil {
			o.ExecOpts = &executorv1.ExecOpts{}
		}
		o.Meta = runtime.MergeLabels(o.Meta, m)
	}
}

// WithAnnotation sets a single annotation.
func WithAnnotation(key, value string) Opt {
	return WithAnnotations(key, value)
}

// WithAnnotations sets multiple annotations from an alternating key/value list.
// Panics if you pass an odd number of args.
func WithAnnotations(kv ...string) Opt {
	m := runtime.KVMap("WithAnnotations", kv)
	return func(o *Opts) {
		if o.ExecOpts == nil {
			o.ExecOpts = &executorv1.ExecOpts{}
		}
		o.Meta = runtime.MergeAnnotations(o.Meta, m)
	}
}
