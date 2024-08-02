from . import executor_pb2 as _executor_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class OpenRequest(_message.Message):
    __slots__ = ("build_id", "opts")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    opts: _executor_pb2.Opts
    def __init__(self, build_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.Opts, _Mapping]] = ...) -> None: ...

class OpenResponse(_message.Message):
    __slots__ = ("runtime_id", "work_directory", "sys_info")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    WORK_DIRECTORY_FIELD_NUMBER: _ClassVar[int]
    SYS_INFO_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    work_directory: str
    sys_info: _executor_pb2.SystemInfo
    def __init__(self, runtime_id: _Optional[str] = ..., work_directory: _Optional[str] = ..., sys_info: _Optional[_Union[_executor_pb2.SystemInfo, _Mapping]] = ...) -> None: ...

class ImportRequest(_message.Message):
    __slots__ = ("runtime_id", "src_path", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    src_path: str
    opts: ImportOpts
    def __init__(self, runtime_id: _Optional[str] = ..., src_path: _Optional[str] = ..., opts: _Optional[_Union[ImportOpts, _Mapping]] = ...) -> None: ...

class ImportOpts(_message.Message):
    __slots__ = ("dest_path", "excludes")
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    EXCLUDES_FIELD_NUMBER: _ClassVar[int]
    dest_path: str
    excludes: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, dest_path: _Optional[str] = ..., excludes: _Optional[_Iterable[str]] = ...) -> None: ...

class ImportResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ExportRequest(_message.Message):
    __slots__ = ("runtime_id", "src_path", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    src_path: str
    opts: ExportOpts
    def __init__(self, runtime_id: _Optional[str] = ..., src_path: _Optional[str] = ..., opts: _Optional[_Union[ExportOpts, _Mapping]] = ...) -> None: ...

class ExportOpts(_message.Message):
    __slots__ = ("dest_path", "excludes")
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    EXCLUDES_FIELD_NUMBER: _ClassVar[int]
    dest_path: str
    excludes: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, dest_path: _Optional[str] = ..., excludes: _Optional[_Iterable[str]] = ...) -> None: ...

class ExportResponse(_message.Message):
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
    __slots__ = ("exec_start", "stdout", "stderr", "exec_end")
    EXEC_START_FIELD_NUMBER: _ClassVar[int]
    STDOUT_FIELD_NUMBER: _ClassVar[int]
    STDERR_FIELD_NUMBER: _ClassVar[int]
    EXEC_END_FIELD_NUMBER: _ClassVar[int]
    exec_start: ExecStartEvent
    stdout: ExecStdoutEvent
    stderr: ExecStderrEvent
    exec_end: ExecEndEvent
    def __init__(self, exec_start: _Optional[_Union[ExecStartEvent, _Mapping]] = ..., stdout: _Optional[_Union[ExecStdoutEvent, _Mapping]] = ..., stderr: _Optional[_Union[ExecStderrEvent, _Mapping]] = ..., exec_end: _Optional[_Union[ExecEndEvent, _Mapping]] = ...) -> None: ...

class ExecStartEvent(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

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

class ExecEndEvent(_message.Message):
    __slots__ = ("error", "exit_code")
    ERROR_FIELD_NUMBER: _ClassVar[int]
    EXIT_CODE_FIELD_NUMBER: _ClassVar[int]
    error: str
    exit_code: int
    def __init__(self, error: _Optional[str] = ..., exit_code: _Optional[int] = ...) -> None: ...
