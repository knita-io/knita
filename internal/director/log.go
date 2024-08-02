package director

import (
	builtinv1 "github.com/knita-io/knita/api/events/builtin/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/log"
)

type Log struct {
	*log.BuildLog
}

func NewLog(stream event.Stream, buildID string) *Log {
	return &Log{
		BuildLog: log.NewBuildLog(stream, buildID, builtinv1.NewDirectorLogEventSource()),
	}
}
