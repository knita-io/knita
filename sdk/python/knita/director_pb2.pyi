from . import executor_pb2 as _executor_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class WorkflowUpdate(_message.Message):
    __slots__ = ("add_job", "add_input", "start_job", "complete_job")
    ADD_JOB_FIELD_NUMBER: _ClassVar[int]
    ADD_INPUT_FIELD_NUMBER: _ClassVar[int]
    START_JOB_FIELD_NUMBER: _ClassVar[int]
    COMPLETE_JOB_FIELD_NUMBER: _ClassVar[int]
    add_job: WorkflowAddJob
    add_input: WorkflowAddInput
    start_job: WorkflowStartJob
    complete_job: WorkflowCompleteJob
    def __init__(self, add_job: _Optional[_Union[WorkflowAddJob, _Mapping]] = ..., add_input: _Optional[_Union[WorkflowAddInput, _Mapping]] = ..., start_job: _Optional[_Union[WorkflowStartJob, _Mapping]] = ..., complete_job: _Optional[_Union[WorkflowCompleteJob, _Mapping]] = ...) -> None: ...

class WorkflowAddJob(_message.Message):
    __slots__ = ("job_id", "needs", "provides")
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    NEEDS_FIELD_NUMBER: _ClassVar[int]
    PROVIDES_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    needs: _containers.RepeatedScalarFieldContainer[str]
    provides: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, job_id: _Optional[str] = ..., needs: _Optional[_Iterable[str]] = ..., provides: _Optional[_Iterable[str]] = ...) -> None: ...

class WorkflowAddInput(_message.Message):
    __slots__ = ("input_id",)
    INPUT_ID_FIELD_NUMBER: _ClassVar[int]
    input_id: str
    def __init__(self, input_id: _Optional[str] = ...) -> None: ...

class WorkflowStartJob(_message.Message):
    __slots__ = ("job_id", "input_data")
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    INPUT_DATA_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    input_data: bytes
    def __init__(self, job_id: _Optional[str] = ..., input_data: _Optional[bytes] = ...) -> None: ...

class WorkflowCompleteJob(_message.Message):
    __slots__ = ("job_id", "duration", "output_data")
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    DURATION_FIELD_NUMBER: _ClassVar[int]
    OUTPUT_DATA_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    duration: _duration_pb2.Duration
    output_data: bytes
    def __init__(self, job_id: _Optional[str] = ..., duration: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ..., output_data: _Optional[bytes] = ...) -> None: ...

class WorkflowSignal(_message.Message):
    __slots__ = ("job_ready",)
    JOB_READY_FIELD_NUMBER: _ClassVar[int]
    job_ready: WorkflowJobReady
    def __init__(self, job_ready: _Optional[_Union[WorkflowJobReady, _Mapping]] = ...) -> None: ...

class WorkflowJobReady(_message.Message):
    __slots__ = ("job_id",)
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    def __init__(self, job_id: _Optional[str] = ...) -> None: ...

class OpenRequest(_message.Message):
    __slots__ = ("build_id", "opts")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    opts: _executor_pb2.Opts
    def __init__(self, build_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.Opts, _Mapping]] = ...) -> None: ...

class OpenResponse(_message.Message):
    __slots__ = ("runtime_id", "work_directory")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    WORK_DIRECTORY_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    work_directory: str
    def __init__(self, runtime_id: _Optional[str] = ..., work_directory: _Optional[str] = ...) -> None: ...

class ImportRequest(_message.Message):
    __slots__ = ("runtime_id", "src_path", "dest_path")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    src_path: str
    dest_path: str
    def __init__(self, runtime_id: _Optional[str] = ..., src_path: _Optional[str] = ..., dest_path: _Optional[str] = ...) -> None: ...

class ImportResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ExportRequest(_message.Message):
    __slots__ = ("runtime_id", "src_path", "dest_path")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    src_path: str
    dest_path: str
    def __init__(self, runtime_id: _Optional[str] = ..., src_path: _Optional[str] = ..., dest_path: _Optional[str] = ...) -> None: ...

class ExportResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class EventsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ExecRequest(_message.Message):
    __slots__ = ("runtime_id", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    opts: _executor_pb2.ExecOpts
    def __init__(self, runtime_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.ExecOpts, _Mapping]] = ...) -> None: ...

class ExecEvent(_message.Message):
    __slots__ = ("exec_end", "stdout", "stderr")
    EXEC_END_FIELD_NUMBER: _ClassVar[int]
    STDOUT_FIELD_NUMBER: _ClassVar[int]
    STDERR_FIELD_NUMBER: _ClassVar[int]
    exec_end: ExecEndEvent
    stdout: ExecStdoutEvent
    stderr: ExecStderrEvent
    def __init__(self, exec_end: _Optional[_Union[ExecEndEvent, _Mapping]] = ..., stdout: _Optional[_Union[ExecStdoutEvent, _Mapping]] = ..., stderr: _Optional[_Union[ExecStderrEvent, _Mapping]] = ...) -> None: ...

class ExecEndEvent(_message.Message):
    __slots__ = ("error", "exit_code")
    ERROR_FIELD_NUMBER: _ClassVar[int]
    EXIT_CODE_FIELD_NUMBER: _ClassVar[int]
    error: str
    exit_code: int
    def __init__(self, error: _Optional[str] = ..., exit_code: _Optional[int] = ...) -> None: ...

class ExecStdoutEvent(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...

class ExecStderrEvent(_message.Message):
    __slots__ = ("data",)
    DATA_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    def __init__(self, data: _Optional[bytes] = ...) -> None: ...
