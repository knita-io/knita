package v1

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
