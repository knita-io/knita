package v1

import (
	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

func NewRuntimeTenderStartEvent(buildID string, tenderID string, opts *executorv1.Opts) *RuntimeTenderStartEvent {
	return &RuntimeTenderStartEvent{BuildId: buildID, TenderId: tenderID, Opts: opts}
}

func NewRuntimeTenderEndEvent(tenderID string) *RuntimeTenderEndEvent {
	return &RuntimeTenderEndEvent{TenderId: tenderID}
}

func NewRuntimeSettlementStartEvent(tenderID string, contractID string, runtimeID string) *RuntimeSettlementStartEvent {
	return &RuntimeSettlementStartEvent{TenderId: tenderID, ContractId: contractID, RuntimeId: runtimeID}
}

func NewRuntimeSettlementEndEvent(tenderID string, contractID string, runtimeID string) *RuntimeSettlementEndEvent {
	return &RuntimeSettlementEndEvent{TenderId: tenderID, ContractId: contractID, RuntimeId: runtimeID}
}

func NewRuntimeOpenStartEvent(runtimeID string, opts *executorv1.Opts) *RuntimeOpenStartEvent {
	return &RuntimeOpenStartEvent{RuntimeId: runtimeID, Opts: opts}
}

func NewRuntimeOpenEndEvent(runtimeID string) *RuntimeOpenEndEvent {
	return &RuntimeOpenEndEvent{RuntimeId: runtimeID}
}

func NewRuntimeCloseStartEvent(runtimeID string) *RuntimeCloseStartEvent {
	return &RuntimeCloseStartEvent{RuntimeId: runtimeID}
}

func NewRuntimeCloseEndEvent(runtimeID string) *RuntimeCloseEndEvent {
	return &RuntimeCloseEndEvent{RuntimeId: runtimeID}
}

func NewExecStartEvent(runtimeID string, execID string, opts *executorv1.ExecOpts) *ExecStartEvent {
	return &ExecStartEvent{RuntimeId: runtimeID, ExecId: execID, Opts: opts}
}

func NewExecEndEvent(runtimeID string, execID string, error string, exitCode int32) *ExecEndEvent {
	return &ExecEndEvent{RuntimeId: runtimeID, ExecId: execID, Error: error, ExitCode: exitCode}
}

func NewImportStartEvent(runtimeID string, importID string) *ImportStartEvent {
	return &ImportStartEvent{RuntimeId: runtimeID, ImportId: importID}
}

func NewImportEndEvent(runtimeID string, importID string) *ImportEndEvent {
	return &ImportEndEvent{RuntimeId: runtimeID, ImportId: importID}
}

func NewExportStartEvent(runtimeID string, exportID string) *ExportStartEvent {
	return &ExportStartEvent{RuntimeId: runtimeID, ExportId: exportID}
}

func NewExportEndEvent(runtimeID string, exportID string) *ExportEndEvent {
	return &ExportEndEvent{RuntimeId: runtimeID, ExportId: exportID}
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

func NewStdoutEvent(data []byte, source *LogEventSource) *StdoutEvent {
	return &StdoutEvent{Data: data, Source: source}
}

func NewStderrEvent(data []byte, source *LogEventSource) *StderrEvent {
	return &StderrEvent{Data: data, Source: source}
}
