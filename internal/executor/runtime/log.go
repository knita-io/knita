package runtime

import (
	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/log"
)

type Log struct {
	*log.BuildLog
	runtimeID string
}

func NewLog(stream event.Stream, buildID string, runtimeID string) *Log {
	source := &builtinv1.LogEventSource{Source: &builtinv1.LogEventSource_Runtime{
		Runtime: &builtinv1.LogSourceRuntime{RuntimeId: runtimeID}}}
	return &Log{
		BuildLog:  log.NewBuildLog(stream, buildID, source),
		runtimeID: runtimeID,
	}
}

func (l *Log) Named(name string) *Log {
	return &Log{
		BuildLog:  l.BuildLog.Named(name),
		runtimeID: l.runtimeID,
	}
}

func (l *Log) ExecSource(execID string, system bool) *Log {
	source := &builtinv1.LogEventSource{Source: &builtinv1.LogEventSource_Exec{
		Exec: &builtinv1.LogSourceExec{RuntimeId: l.runtimeID, ExecId: execID, System: system}}}
	return &Log{
		BuildLog:  l.BuildLog.Source(source),
		runtimeID: l.runtimeID,
	}
}
