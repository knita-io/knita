package runtime

import (
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/log"
)

type Log struct {
	*log.BuildLog
	runtimeID string
}

func NewLog(stream event.Stream, buildID string, runtimeID string) *Log {
	return &Log{
		BuildLog:  log.NewBuildLog(stream, buildID, executorv1.NewRuntimeLogEventSource(runtimeID)),
		runtimeID: runtimeID,
	}
}

func (l *Log) Named(name string) *Log {
	return &Log{
		BuildLog:  l.BuildLog.Named(name),
		runtimeID: l.runtimeID,
	}
}

func (l *Log) ExecSource(execID string) *Log {
	return &Log{
		BuildLog:  l.BuildLog.Source(executorv1.NewExecLogEventSource(l.runtimeID, execID)),
		runtimeID: l.runtimeID,
	}
}
