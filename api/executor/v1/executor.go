package v1

type IsEvent_Payload interface {
	isEvent_Payload
}

type IsOutEvent_Source interface {
	isLogOutEventSource_Source
}

type RuntimeEvent interface {
	// GetRuntimeId is implemented by events that *can* be emitted in the context of a runtime.
	// If the returned value is non-empty then this event *was* emitted in the context of a runtime.
	GetRuntimeId() string
}

func (e *Event_RuntimeOpened) GetRuntimeId() string {
	return e.RuntimeOpened.GetRuntimeId()
}

func (e *Event_ExecStart) GetRuntimeId() string {
	return e.ExecStart.GetRuntimeId()
}

func (e *Event_ExecEnd) GetRuntimeId() string {
	return e.ExecEnd.GetRuntimeId()
}

func (e *Event_ImportStart) GetRuntimeId() string {
	return e.ImportStart.GetRuntimeId()
}

func (e *Event_ImportEnd) GetRuntimeId() string {
	return e.ImportEnd.GetRuntimeId()
}

func (e *Event_ExportStart) GetRuntimeId() string {
	return e.ExportStart.GetRuntimeId()
}

func (e *Event_ExportEnd) GetRuntimeId() string {
	return e.ExportEnd.GetRuntimeId()
}

func (e *Event_Stdout) GetRuntimeId() string {
	switch s := e.Stdout.Source.Source.(type) {
	case *LogOutEventSource_Runtime:
		return s.Runtime.RuntimeId
	case *LogOutEventSource_Exec:
		return s.Exec.RuntimeId
	}
	return ""
}

func (e *Event_Stderr) GetRuntimeId() string {
	switch s := e.Stderr.Source.Source.(type) {
	case *LogOutEventSource_Runtime:
		return s.Runtime.RuntimeId
	case *LogOutEventSource_Exec:
		return s.Exec.RuntimeId
	}
	return ""
}

func (e *Event_RuntimeClosed) GetRuntimeId() string {
	return e.RuntimeClosed.GetRuntimeId()
}

func NewRuntimeOpenedEvent(runtimeID string, opts *Opts) *Event {
	return &Event{Payload: &Event_RuntimeOpened{RuntimeOpened: &RuntimeOpenedEvent{RuntimeId: runtimeID, Opts: opts}}}
}

func NewRuntimeClosedEvent(runtimeID string) *Event {
	return &Event{Payload: &Event_RuntimeClosed{RuntimeClosed: &RuntimeClosedEvent{RuntimeId: runtimeID}}}
}

func NewExecStartEvent(runtimeID string, execID string, opts *ExecOpts) *Event {
	return &Event{Payload: &Event_ExecStart{ExecStart: &ExecStartEvent{RuntimeId: runtimeID, ExecId: execID, Opts: opts}}}
}

func NewExecEndEvent(runtimeID string, execID string, error string, exitCode int32) *Event {
	return &Event{Payload: &Event_ExecEnd{ExecEnd: &ExecEndEvent{RuntimeId: runtimeID, ExecId: execID, Error: error, ExitCode: exitCode}}}
}

func NewImportStartEvent(runtimeID string, importID string) *Event {
	return &Event{Payload: &Event_ImportStart{ImportStart: &ImportStartEvent{RuntimeId: runtimeID, ImportId: importID}}}
}

func NewImportEndEvent(runtimeID string, importID string) *Event {
	return &Event{Payload: &Event_ImportEnd{ImportEnd: &ImportEndEvent{RuntimeId: runtimeID, ImportId: importID}}}
}

func NewExportStartEvent(runtimeID string, exportID string) *Event {
	return &Event{Payload: &Event_ExportStart{ExportStart: &ExportStartEvent{RuntimeId: runtimeID, ExportId: exportID}}}
}

func NewExportEndEvent(runtimeID string, exportID string) *Event {
	return &Event{Payload: &Event_ExportEnd{ExportEnd: &ExportEndEvent{RuntimeId: runtimeID, ExportId: exportID}}}
}

func NewRuntimeLogEventSource(runtimeID string) *LogOutEventSource {
	return &LogOutEventSource{Source: &LogOutEventSource_Runtime{Runtime: &LogSourceRuntime{RuntimeId: runtimeID}}}
}

func NewExecLogEventSource(runtimeID string, execID string) *LogOutEventSource {
	return &LogOutEventSource{Source: &LogOutEventSource_Exec{Exec: &LogSourceExec{RuntimeId: runtimeID, ExecId: execID}}}
}

func NewDirectorLogEventSource() *LogOutEventSource {
	return &LogOutEventSource{Source: &LogOutEventSource_Director{Director: &LogSourceDirector{}}}
}

func NewStdoutEvent(data []byte, source *LogOutEventSource) *Event {
	return &Event{
		Payload: &Event_Stdout{
			Stdout: &StdoutEvent{Data: data, Source: source},
		},
	}
}

func NewStderrEvent(data []byte, source *LogOutEventSource) *Event {
	return &Event{
		Payload: &Event_Stderr{
			Stderr: &StderrEvent{Data: data, Source: source},
		},
	}
}
