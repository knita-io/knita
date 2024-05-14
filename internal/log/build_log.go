package log

import (
	"fmt"
	"io"
	"strings"

	"github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
)

const writerBufferSize = 4 * 1024

type closer interface {
	Close() error
}

type BuildLog struct {
	name       string
	dependents []closer
	stream     event.Stream
	sequencer  *event.Sequencer
	source     *v1.LogOutEventSource
	buildID    string
}

func NewBuildLog(stream event.Stream, buildID string, source *v1.LogOutEventSource) *BuildLog {
	return &BuildLog{
		name:      "",
		source:    source,
		stream:    stream,
		sequencer: event.NewSequencer(stream),
		buildID:   buildID,
	}
}

func (l *BuildLog) Stream() event.Stream {
	return l.stream
}

func (l *BuildLog) Publish(event *v1.Event) {
	event.BuildId = l.buildID
	event.GroupName = l.name
	l.sequencer.Publish(event)
}

func (l *BuildLog) Republish(event *v1.Event) {
	if event.BuildId != l.buildID {
		panic("build id mismatch")
	}
	l.sequencer.Publish(event)
}

func (l *BuildLog) Named(name string) *BuildLog {
	if l.name != "" {
		name = fmt.Sprintf("%s.%s", l.name, name)
	}
	log := l.clone()
	log.name = name
	return log
}

func (l *BuildLog) Source(source *v1.LogOutEventSource) *BuildLog {
	log := l.clone()
	log.source = source
	return log
}

func (l *BuildLog) Stdout() io.WriteCloser {
	r, w := io.Pipe()
	buf := make([]byte, writerBufferSize)
	go func() {
		for {
			n, err := r.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				event := v1.NewStdoutEvent(buf[:n], l.source)
				l.Publish(event)
			}
		}
	}()
	l.dependents = append(l.dependents, r)
	return w
}

func (l *BuildLog) Stderr() io.WriteCloser {
	r, w := io.Pipe()
	buf := make([]byte, writerBufferSize)
	go func() {
		for {
			n, err := r.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				event := v1.NewStderrEvent(buf[:n], l.source)
				l.Publish(event)
			}
		}
	}()
	l.dependents = append(l.dependents, r)
	return w
}

func (l *BuildLog) Printf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(str, "\n") {
		str += "\n"
	}
	event := v1.NewStdoutEvent([]byte(str), l.source)
	l.Publish(event)
}

func (l *BuildLog) Close() error {
	for _, d := range l.dependents {
		d.Close()
	}
	return nil
}

func (l *BuildLog) clone() *BuildLog {
	log := &BuildLog{
		name:      l.name,
		source:    l.source,
		stream:    l.stream,
		sequencer: l.sequencer,
		buildID:   l.buildID,
	}
	l.dependents = append(l.dependents, log)
	return log
}
