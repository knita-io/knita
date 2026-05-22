package knita

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	directorv1 "github.com/knita-io/knita/api/director/v1"
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/sdk/go/knita/runtime"
)

// FatalFunc is a function that will be called when a MustXXX function encounters an error.
type FatalFunc func(err error)

// Log is a simple logger for the Knita SDK to use to write logs.
type Log interface {
	Printf(format string, args ...interface{})
}

// Opt configures the Options.
type Opt func(*Opts)

// Opts holds customizable behaviors.
type Opts struct {
	Log       Log
	FatalFunc FatalFunc
}

// WithLog sets a custom Log that Knita will write to.
func WithLog(log Log) Opt {
	return func(o *Opts) {
		o.Log = log
	}
}

// WithFatalFunc sets a custom function that all MustXXX functions will call on errors.
func WithFatalFunc(fn FatalFunc) Opt {
	return func(o *Opts) {
		o.FatalFunc = fn
	}
}

type defaultLog struct{}

func (l *defaultLog) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Client connects back to the Knita CLI process to orchestrate builds.
type Client struct {
	syslog    Log
	fatalFunc FatalFunc
	client    directorv1.DirectorClient
	buildID   string
}

// MustNewClient is like NewClient, but it calls the configured FatalFunc if an error occurs.
func MustNewClient(opts ...Opt) *Client {
	c, err := NewClient(opts...)
	if err != nil {
		o := makeOpts(opts...)
		o.FatalFunc(fmt.Errorf("error creating client: %w", err))
	}
	return c
}

// NewClient returns a Knita client that is configured to connect back to the Knita CLI process.
func NewClient(opts ...Opt) (*Client, error) {
	o := makeOpts(opts...)
	buildID := os.Getenv("KNITA_BUILD_ID")
	if buildID == "" {
		return nil, fmt.Errorf("error expected KNITA_BUILD_ID to be set")
	}
	socket := os.Getenv("KNITA_SOCKET")
	if socket == "" {
		return nil, fmt.Errorf("error expected KNITA_SOCKET to be set")
	}
	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial("unix", addr)
	}
	conn, err := grpc.Dial(socket, grpc.WithInsecure(), grpc.WithDialer(dialer))
	if err != nil {
		return nil, fmt.Errorf("error dialing local Knita socket %s: %w", socket, err)
	}
	return &Client{
		client:    directorv1.NewDirectorClient(conn),
		syslog:    o.Log,
		fatalFunc: o.FatalFunc,
		buildID:   buildID,
	}, nil
}

// Runtime opens a new remote runtime configured based on options.
func (c *Client) Runtime(opts ...runtime.Opt) (*Runtime, error) {
	return c.RuntimeWithContext(context.Background(), opts...)
}

// MustRuntime is like Runtime, but it calls the configured FatalFunc if an error occurs.
func (c *Client) MustRuntime(opts ...runtime.Opt) *Runtime {
	rt, err := c.Runtime(opts...)
	if err != nil {
		c.fatalFunc(fmt.Errorf("error creating runtime: %w", err))
	}
	return rt
}

// RuntimeWithContext is like Runtime, but it allows a context to be set.
func (c *Client) RuntimeWithContext(ctx context.Context, opts ...runtime.Opt) (*Runtime, error) {
	req := &directorv1.OpenRequest{BuildId: c.buildID, Opts: &executorv1.RuntimeOpts{}}
	for _, opt := range opts {
		opt(req)
	}
	res, err := c.client.Open(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error opening runtime: %w", err)
	}
	return &Runtime{
		syslog:              c.syslog,
		fatalFunc:           c.fatalFunc,
		client:              c.client,
		runtimeID:           res.RuntimeId,
		remoteWorkDirectory: res.WorkDirectory,
		remoteSysInfo:       res.SysInfo,
	}, nil
}

func makeOpts(opts ...Opt) *Opts {
	o := &Opts{}
	for _, opt := range opts {
		opt(o)
	}
	if o.Log == nil {
		o.Log = &defaultLog{}
	}
	if o.FatalFunc == nil {
		o.FatalFunc = func(err error) {
			fmt.Fprintf(os.Stderr, "Must function encountered a fatal error: %v\n", err)
			os.Exit(1)
		}
	}
	return o
}
