package director

import (
	executorv1 "github.com/knita-io/knita/api/executor/v1"
	"github.com/knita-io/knita/internal/event"
	"github.com/knita-io/knita/internal/log"
)

type Log struct {
	*log.BuildLog
}

func NewLog(stream event.Stream, buildID string) *Log {
	return &Log{
		BuildLog: log.NewBuildLog(stream, buildID, executorv1.NewDirectorLogEventSource()),
	}
}
