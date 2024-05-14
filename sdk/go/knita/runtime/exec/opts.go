package exec

import (
	"io"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

type Opts struct {
	*executorv1.ExecOpts
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

type Opt interface {
	Apply(opts *Opts)
}

type withCommand struct {
	name string
	args []string
}

func (o *withCommand) Apply(opts *Opts) {
	if opts.ExecOpts == nil {
		opts.ExecOpts = &executorv1.ExecOpts{}
	}
	opts.Name = o.name
	opts.Args = o.args
}

// WithCommand specifies the command to execute.
func WithCommand(name string, args ...string) Opt {
	return &withCommand{name: name, args: args}
}

type withStdout struct {
	stdout io.Writer
}

func (o *withStdout) Apply(opts *Opts) {
	if opts.ExecOpts == nil {
		opts.ExecOpts = &executorv1.ExecOpts{}
	}
	opts.Stdout = o.stdout
}

// WithStdout wires up the command's stdout to the specified writer.
func WithStdout(stdout io.Writer) Opt {
	return &withStdout{stdout: stdout}
}

type withStderr struct {
	stderr io.Writer
}

func (o *withStderr) Apply(opts *Opts) {
	if opts.ExecOpts == nil {
		opts.ExecOpts = &executorv1.ExecOpts{}
	}
	opts.Stderr = o.stderr
}

// WithStderr wires up the command's stderr to the specified writer.
func WithStderr(stderr io.Writer) Opt {
	return &withStderr{stderr: stderr}
}

type withEnv struct {
	env []string
}

func (o *withEnv) Apply(opts *Opts) {
	if opts.ExecOpts == nil {
		opts.ExecOpts = &executorv1.ExecOpts{}
	}
	opts.Env = append(opts.Env, o.env...)
}

// WithEnv provides one or more environment variables to the command.
// Variables are specified as key value pairs separated by an = sign e.g. 'FOO=bar'.
func WithEnv(env ...string) Opt {
	return &withEnv{env: env}
}

type withTag struct {
	key   string
	value string
}

func (o *withTag) Apply(opts *Opts) {
	if opts.ExecOpts == nil {
		opts.ExecOpts = &executorv1.ExecOpts{}
	}
	if opts.Tags == nil {
		opts.Tags = map[string]string{}
	}
	opts.Tags[o.key] = o.value
}

// WithTag tags the command with a key value pair.
func WithTag(key string, value string) Opt {
	return &withTag{key: key, value: value}
}

type withTags struct {
	tags map[string]string
}

func (o *withTags) Apply(opts *Opts) {
	if opts.ExecOpts == nil {
		opts.ExecOpts = &executorv1.ExecOpts{}
	}
	if opts.Tags == nil {
		opts.Tags = map[string]string{}
	}
	for k, v := range o.tags {
		opts.Tags[k] = v
	}
}

// WithTags tags the command with one or more key value pairs.
// The expected format of tags is an array of alternating key value pairs.
// If tags are mismatched this function will panic.
func WithTags(tags ...string) Opt {
	m := map[string]string{}
	if len(tags)%2 != 0 {
		panic("expected alternating kv paris")
	}
	return &withTags{tags: m}
}
