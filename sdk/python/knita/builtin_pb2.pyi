from . import executor_pb2 as _executor_pb2
from . import broker_pb2 as _broker_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Error(_message.Message):
    __slots__ = ("message",)
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...

class DirectorInfo(_message.Message):
    __slots__ = ("version", "sys_info")
    VERSION_FIELD_NUMBER: _ClassVar[int]
    SYS_INFO_FIELD_NUMBER: _ClassVar[int]
    version: str
    sys_info: _executor_pb2.SystemInfo
    def __init__(self, version: _Optional[str] = ..., sys_info: _Optional[_Union[_executor_pb2.SystemInfo, _Mapping]] = ...) -> None: ...

class BuildStartEvent(_message.Message):
    __slots__ = ("build_id", "director_info")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    DIRECTOR_INFO_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    director_info: DirectorInfo
    def __init__(self, build_id: _Optional[str] = ..., director_info: _Optional[_Union[DirectorInfo, _Mapping]] = ...) -> None: ...

class BuildResult(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class BuildEndEvent(_message.Message):
    __slots__ = ("build_id", "error", "result")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    error: Error
    result: BuildResult
    def __init__(self, build_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[BuildResult, _Mapping]] = ...) -> None: ...

class RuntimeTenderStartEvent(_message.Message):
    __slots__ = ("build_id", "tender_id", "opts")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    TENDER_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    tender_id: str
    opts: _executor_pb2.RuntimeOpts
    def __init__(self, build_id: _Optional[str] = ..., tender_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.RuntimeOpts, _Mapping]] = ...) -> None: ...

class RuntimeTenderResult(_message.Message):
    __slots__ = ("contracts",)
    CONTRACTS_FIELD_NUMBER: _ClassVar[int]
    contracts: _containers.RepeatedCompositeFieldContainer[_broker_pb2.RuntimeContract]
    def __init__(self, contracts: _Optional[_Iterable[_Union[_broker_pb2.RuntimeContract, _Mapping]]] = ...) -> None: ...

class RuntimeTenderEndEvent(_message.Message):
    __slots__ = ("tender_id", "error", "result")
    TENDER_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    tender_id: str
    error: Error
    result: RuntimeTenderResult
    def __init__(self, tender_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[RuntimeTenderResult, _Mapping]] = ...) -> None: ...

class RuntimeSettlementStartEvent(_message.Message):
    __slots__ = ("tender_id", "contract_id", "runtime_id")
    TENDER_ID_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_ID_FIELD_NUMBER: _ClassVar[int]
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    tender_id: str
    contract_id: str
    runtime_id: str
    def __init__(self, tender_id: _Optional[str] = ..., contract_id: _Optional[str] = ..., runtime_id: _Optional[str] = ...) -> None: ...

class RuntimeSettlementResult(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RuntimeSettlementEndEvent(_message.Message):
    __slots__ = ("tender_id", "contract_id", "runtime_id", "error", "result")
    TENDER_ID_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_ID_FIELD_NUMBER: _ClassVar[int]
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    tender_id: str
    contract_id: str
    runtime_id: str
    error: Error
    result: RuntimeSettlementResult
    def __init__(self, tender_id: _Optional[str] = ..., contract_id: _Optional[str] = ..., runtime_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[RuntimeSettlementResult, _Mapping]] = ...) -> None: ...

class RuntimeOpenStartEvent(_message.Message):
    __slots__ = ("runtime_id", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    opts: _executor_pb2.RuntimeOpts
    def __init__(self, runtime_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.RuntimeOpts, _Mapping]] = ...) -> None: ...

class RuntimeOpenResult(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RuntimeOpenEndEvent(_message.Message):
    __slots__ = ("runtime_id", "error", "result")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    error: Error
    result: RuntimeOpenResult
    def __init__(self, runtime_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[RuntimeOpenResult, _Mapping]] = ...) -> None: ...

class RuntimeCloseStartEvent(_message.Message):
    __slots__ = ("runtime_id",)
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    def __init__(self, runtime_id: _Optional[str] = ...) -> None: ...

class RuntimeCloseResult(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class RuntimeCloseEndEvent(_message.Message):
    __slots__ = ("runtime_id", "error", "result")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    error: Error
    result: RuntimeCloseResult
    def __init__(self, runtime_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[RuntimeCloseResult, _Mapping]] = ...) -> None: ...

class StdoutEvent(_message.Message):
    __slots__ = ("data", "source")
    DATA_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    source: LogEventSource
    def __init__(self, data: _Optional[bytes] = ..., source: _Optional[_Union[LogEventSource, _Mapping]] = ...) -> None: ...

class StderrEvent(_message.Message):
    __slots__ = ("data", "source")
    DATA_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    source: LogEventSource
    def __init__(self, data: _Optional[bytes] = ..., source: _Optional[_Union[LogEventSource, _Mapping]] = ...) -> None: ...

class LogEventSource(_message.Message):
    __slots__ = ("runtime", "exec", "director")
    RUNTIME_FIELD_NUMBER: _ClassVar[int]
    EXEC_FIELD_NUMBER: _ClassVar[int]
    DIRECTOR_FIELD_NUMBER: _ClassVar[int]
    runtime: LogSourceRuntime
    exec: LogSourceExec
    director: LogSourceDirector
    def __init__(self, runtime: _Optional[_Union[LogSourceRuntime, _Mapping]] = ..., exec: _Optional[_Union[LogSourceExec, _Mapping]] = ..., director: _Optional[_Union[LogSourceDirector, _Mapping]] = ...) -> None: ...

class LogSourceRuntime(_message.Message):
    __slots__ = ("runtime_id",)
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    def __init__(self, runtime_id: _Optional[str] = ...) -> None: ...

class LogSourceExec(_message.Message):
    __slots__ = ("runtime_id", "exec_id", "system")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    EXEC_ID_FIELD_NUMBER: _ClassVar[int]
    SYSTEM_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    exec_id: str
    system: bool
    def __init__(self, runtime_id: _Optional[str] = ..., exec_id: _Optional[str] = ..., system: bool = ...) -> None: ...

class LogSourceDirector(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ExecStartEvent(_message.Message):
    __slots__ = ("runtime_id", "exec_id", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    EXEC_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    exec_id: str
    opts: _executor_pb2.ExecOpts
    def __init__(self, runtime_id: _Optional[str] = ..., exec_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.ExecOpts, _Mapping]] = ...) -> None: ...

class ExecResult(_message.Message):
    __slots__ = ("exit_code",)
    EXIT_CODE_FIELD_NUMBER: _ClassVar[int]
    exit_code: int
    def __init__(self, exit_code: _Optional[int] = ...) -> None: ...

class ExecEndEvent(_message.Message):
    __slots__ = ("runtime_id", "exec_id", "error", "result")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    EXEC_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    exec_id: str
    error: Error
    result: ExecResult
    def __init__(self, runtime_id: _Optional[str] = ..., exec_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[ExecResult, _Mapping]] = ...) -> None: ...

class ImportStartEvent(_message.Message):
    __slots__ = ("runtime_id", "import_id")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    IMPORT_ID_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    import_id: str
    def __init__(self, runtime_id: _Optional[str] = ..., import_id: _Optional[str] = ...) -> None: ...

class ImportResult(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ImportEndEvent(_message.Message):
    __slots__ = ("runtime_id", "import_id", "error", "result")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    IMPORT_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    import_id: str
    error: Error
    result: ImportResult
    def __init__(self, runtime_id: _Optional[str] = ..., import_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[ImportResult, _Mapping]] = ...) -> None: ...

class ExportStartEvent(_message.Message):
    __slots__ = ("runtime_id", "export_id")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    EXPORT_ID_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    export_id: str
    def __init__(self, runtime_id: _Optional[str] = ..., export_id: _Optional[str] = ...) -> None: ...

class ExportResult(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ExportEndEvent(_message.Message):
    __slots__ = ("runtime_id", "export_id", "error", "result")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    EXPORT_ID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    export_id: str
    error: Error
    result: ExportResult
    def __init__(self, runtime_id: _Optional[str] = ..., export_id: _Optional[str] = ..., error: _Optional[_Union[Error, _Mapping]] = ..., result: _Optional[_Union[ExportResult, _Mapping]] = ...) -> None: ...

class SyncPointReachedEvent(_message.Message):
    __slots__ = ("barrier_id",)
    BARRIER_ID_FIELD_NUMBER: _ClassVar[int]
    barrier_id: str
    def __init__(self, barrier_id: _Optional[str] = ...) -> None: ...
