package v1

type IsEvent_Payload interface {
	isEvent_Payload
}

type IsOutEvent_Source interface {
	isLogEventSource_Source
}

type RuntimeEvent interface {
	// GetRuntimeId is implemented by events that *can* be emitted in the context of a runtime.
	// If the returned value is non-empty then this event *was* emitted in the context of a runtime.
	GetRuntimeId() string
}

func (e *Event_RuntimeSettlementStart) GetRuntimeId() string {
	return e.RuntimeSettlementStart.GetRuntimeId()
}

func (e *Event_RuntimeSettlementEnd) GetRuntimeId() string {
	return e.RuntimeSettlementEnd.GetRuntimeId()
}

func (e *Event_RuntimeOpenStart) GetRuntimeId() string {
	return e.RuntimeOpenStart.GetRuntimeId()
}

func (e *Event_RuntimeOpenEnd) GetRuntimeId() string {
	return e.RuntimeOpenEnd.GetRuntimeId()
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
	case *LogEventSource_Runtime:
		return s.Runtime.RuntimeId
	case *LogEventSource_Exec:
		return s.Exec.RuntimeId
	}
	return ""
}

func (e *Event_Stderr) GetRuntimeId() string {
	switch s := e.Stderr.Source.Source.(type) {
	case *LogEventSource_Runtime:
		return s.Runtime.RuntimeId
	case *LogEventSource_Exec:
		return s.Exec.RuntimeId
	}
	return ""
}

func (e *Event_RuntimeCloseStart) GetRuntimeId() string {
	return e.RuntimeCloseStart.GetRuntimeId()
}

func (e *Event_RuntimeCloseEnd) GetRuntimeId() string {
	return e.RuntimeCloseEnd.GetRuntimeId()
}

func NewRuntimeTenderStartEvent(buildID string, tenderID string, opts *Opts) *Event {
	return &Event{Payload: &Event_RuntimeTenderStart{RuntimeTenderStart: &RuntimeTenderStartEvent{BuildId: buildID, TenderId: tenderID, Opts: opts}}}
}

func NewRuntimeTenderEndEvent(tenderID string) *Event {
	return &Event{Payload: &Event_RuntimeTenderEnd{RuntimeTenderEnd: &RuntimeTenderEndEvent{TenderId: tenderID}}}
}

func NewRuntimeSettlementStartEvent(tenderID string, contractID string, runtimeID string) *Event {
	return &Event{Payload: &Event_RuntimeSettlementStart{RuntimeSettlementStart: &RuntimeSettlementStartEvent{TenderId: tenderID, ContractId: contractID, RuntimeId: runtimeID}}}
}

func NewRuntimeSettlementEndEvent(tenderID string, contractID string, runtimeID string) *Event {
	return &Event{Payload: &Event_RuntimeSettlementEnd{RuntimeSettlementEnd: &RuntimeSettlementEndEvent{TenderId: tenderID, ContractId: contractID, RuntimeId: runtimeID}}}
}

func NewRuntimeOpenStartEvent(runtimeID string, opts *Opts) *Event {
	return &Event{Payload: &Event_RuntimeOpenStart{RuntimeOpenStart: &RuntimeOpenStartEvent{RuntimeId: runtimeID, Opts: opts}}}
}

func NewRuntimeOpenEndEvent(runtimeID string) *Event {
	return &Event{Payload: &Event_RuntimeOpenEnd{RuntimeOpenEnd: &RuntimeOpenEndEvent{RuntimeId: runtimeID}}}
}

func NewRuntimeCloseStartEvent(runtimeID string) *Event {
	return &Event{Payload: &Event_RuntimeCloseStart{RuntimeCloseStart: &RuntimeCloseStartEvent{RuntimeId: runtimeID}}}
}

func NewRuntimeCloseEndEvent(runtimeID string) *Event {
	return &Event{Payload: &Event_RuntimeCloseEnd{RuntimeCloseEnd: &RuntimeCloseEndEvent{RuntimeId: runtimeID}}}
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

func NewRuntimeLogEventSource(runtimeID string) *LogEventSource {
	return &LogEventSource{Source: &LogEventSource_Runtime{Runtime: &LogSourceRuntime{RuntimeId: runtimeID}}}
}

func NewExecLogEventSource(runtimeID string, execID string, system bool) *LogEventSource {
	return &LogEventSource{Source: &LogEventSource_Exec{Exec: &LogSourceExec{RuntimeId: runtimeID, ExecId: execID, System: system}}}
}

func NewDirectorLogEventSource() *LogEventSource {
	return &LogEventSource{Source: &LogEventSource_Director{Director: &LogSourceDirector{}}}
}

func NewStdoutEvent(data []byte, source *LogEventSource) *Event {
	return &Event{
		Payload: &Event_Stdout{
			Stdout: &StdoutEvent{Data: data, Source: source},
		},
	}
}

func NewStderrEvent(data []byte, source *LogEventSource) *Event {
	return &Event{
		Payload: &Event_Stderr{
			Stderr: &StderrEvent{Data: data, Source: source},
		},
	}
}
