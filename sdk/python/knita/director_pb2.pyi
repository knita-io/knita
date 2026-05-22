from . import executor_pb2 as _executor_pb2
from . import event_pb2 as _event_pb2
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
    opts: _executor_pb2.RuntimeOpts
    def __init__(self, build_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.RuntimeOpts, _Mapping]] = ...) -> None: ...

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
    __slots__ = ("runtime_id", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    opts: ImportOpts
    def __init__(self, runtime_id: _Optional[str] = ..., opts: _Optional[_Union[ImportOpts, _Mapping]] = ...) -> None: ...

class ImportOpts(_message.Message):
    __slots__ = ("src_path", "dest_path", "excludes", "meta", "display_name")
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    EXCLUDES_FIELD_NUMBER: _ClassVar[int]
    META_FIELD_NUMBER: _ClassVar[int]
    DISPLAY_NAME_FIELD_NUMBER: _ClassVar[int]
    src_path: str
    dest_path: str
    excludes: _containers.RepeatedScalarFieldContainer[str]
    meta: _executor_pb2.OptsMeta
    display_name: str
    def __init__(self, src_path: _Optional[str] = ..., dest_path: _Optional[str] = ..., excludes: _Optional[_Iterable[str]] = ..., meta: _Optional[_Union[_executor_pb2.OptsMeta, _Mapping]] = ..., display_name: _Optional[str] = ...) -> None: ...

class ImportResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ExportRequest(_message.Message):
    __slots__ = ("runtime_id", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    opts: ExportOpts
    def __init__(self, runtime_id: _Optional[str] = ..., opts: _Optional[_Union[ExportOpts, _Mapping]] = ...) -> None: ...

class ExportOpts(_message.Message):
    __slots__ = ("src_path", "dest_path", "excludes", "meta", "display_name")
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    EXCLUDES_FIELD_NUMBER: _ClassVar[int]
    META_FIELD_NUMBER: _ClassVar[int]
    DISPLAY_NAME_FIELD_NUMBER: _ClassVar[int]
    src_path: str
    dest_path: str
    excludes: _containers.RepeatedScalarFieldContainer[str]
    meta: _executor_pb2.OptsMeta
    display_name: str
    def __init__(self, src_path: _Optional[str] = ..., dest_path: _Optional[str] = ..., excludes: _Optional[_Iterable[str]] = ..., meta: _Optional[_Union[_executor_pb2.OptsMeta, _Mapping]] = ..., display_name: _Optional[str] = ...) -> None: ...

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

class CloseRequest(_message.Message):
    __slots__ = ("runtime_id",)
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    def __init__(self, runtime_id: _Optional[str] = ...) -> None: ...

class CloseResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
