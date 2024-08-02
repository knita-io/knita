from events.v1 import event_pb2 as _event_pb2
from google.protobuf import duration_pb2 as _duration_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class RuntimeType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    RUNTIME_UNSPECIFIED: _ClassVar[RuntimeType]
    RUNTIME_HOST: _ClassVar[RuntimeType]
    RUNTIME_DOCKER: _ClassVar[RuntimeType]
RUNTIME_UNSPECIFIED: RuntimeType
RUNTIME_HOST: RuntimeType
RUNTIME_DOCKER: RuntimeType

class ExecutorInfo(_message.Message):
    __slots__ = ("name",)
    NAME_FIELD_NUMBER: _ClassVar[int]
    name: str
    def __init__(self, name: _Optional[str] = ...) -> None: ...

class SystemInfo(_message.Message):
    __slots__ = ("os", "arch", "total_cpu_cores", "total_memory")
    OS_FIELD_NUMBER: _ClassVar[int]
    ARCH_FIELD_NUMBER: _ClassVar[int]
    TOTAL_CPU_CORES_FIELD_NUMBER: _ClassVar[int]
    TOTAL_MEMORY_FIELD_NUMBER: _ClassVar[int]
    os: str
    arch: str
    total_cpu_cores: int
    total_memory: int
    def __init__(self, os: _Optional[str] = ..., arch: _Optional[str] = ..., total_cpu_cores: _Optional[int] = ..., total_memory: _Optional[int] = ...) -> None: ...

class IntrospectRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class IntrospectResponse(_message.Message):
    __slots__ = ("sys_info", "executor_info", "labels")
    SYS_INFO_FIELD_NUMBER: _ClassVar[int]
    EXECUTOR_INFO_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    sys_info: SystemInfo
    executor_info: ExecutorInfo
    labels: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, sys_info: _Optional[_Union[SystemInfo, _Mapping]] = ..., executor_info: _Optional[_Union[ExecutorInfo, _Mapping]] = ..., labels: _Optional[_Iterable[str]] = ...) -> None: ...

class EventsRequest(_message.Message):
    __slots__ = ("build_id", "runtime_id", "barrier_id")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    BARRIER_ID_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    runtime_id: str
    barrier_id: str
    def __init__(self, build_id: _Optional[str] = ..., runtime_id: _Optional[str] = ..., barrier_id: _Optional[str] = ...) -> None: ...

class OpenRequest(_message.Message):
    __slots__ = ("build_id", "runtime_id", "opts")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    runtime_id: str
    opts: Opts
    def __init__(self, build_id: _Optional[str] = ..., runtime_id: _Optional[str] = ..., opts: _Optional[_Union[Opts, _Mapping]] = ...) -> None: ...

class OpenResponse(_message.Message):
    __slots__ = ("work_directory", "sys_info")
    WORK_DIRECTORY_FIELD_NUMBER: _ClassVar[int]
    SYS_INFO_FIELD_NUMBER: _ClassVar[int]
    work_directory: str
    sys_info: SystemInfo
    def __init__(self, work_directory: _Optional[str] = ..., sys_info: _Optional[_Union[SystemInfo, _Mapping]] = ...) -> None: ...

class HeartbeatRequest(_message.Message):
    __slots__ = ("runtime_id",)
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    def __init__(self, runtime_id: _Optional[str] = ...) -> None: ...

class HeartbeatResponse(_message.Message):
    __slots__ = ("extended_by",)
    EXTENDED_BY_FIELD_NUMBER: _ClassVar[int]
    extended_by: _duration_pb2.Duration
    def __init__(self, extended_by: _Optional[_Union[_duration_pb2.Duration, _Mapping]] = ...) -> None: ...

class Opts(_message.Message):
    __slots__ = ("type", "labels", "tags", "host", "docker")
    class TagsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    TYPE_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    HOST_FIELD_NUMBER: _ClassVar[int]
    DOCKER_FIELD_NUMBER: _ClassVar[int]
    type: RuntimeType
    labels: _containers.RepeatedScalarFieldContainer[str]
    tags: _containers.ScalarMap[str, str]
    host: HostOpts
    docker: DockerOpts
    def __init__(self, type: _Optional[_Union[RuntimeType, str]] = ..., labels: _Optional[_Iterable[str]] = ..., tags: _Optional[_Mapping[str, str]] = ..., host: _Optional[_Union[HostOpts, _Mapping]] = ..., docker: _Optional[_Union[DockerOpts, _Mapping]] = ...) -> None: ...

class HostOpts(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class DockerOpts(_message.Message):
    __slots__ = ("image",)
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    image: DockerPullOpts
    def __init__(self, image: _Optional[_Union[DockerPullOpts, _Mapping]] = ...) -> None: ...

class DockerPullOpts(_message.Message):
    __slots__ = ("image_uri", "pull_strategy", "auth")
    class PullStrategy(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
        __slots__ = ()
        PULL_STRATEGY_UNSPECIFIED: _ClassVar[DockerPullOpts.PullStrategy]
        PULL_STRATEGY_NEVER: _ClassVar[DockerPullOpts.PullStrategy]
        PULL_STRATEGY_ALWAYS: _ClassVar[DockerPullOpts.PullStrategy]
        PULL_STRATEGY_NOT_EXISTS: _ClassVar[DockerPullOpts.PullStrategy]
    PULL_STRATEGY_UNSPECIFIED: DockerPullOpts.PullStrategy
    PULL_STRATEGY_NEVER: DockerPullOpts.PullStrategy
    PULL_STRATEGY_ALWAYS: DockerPullOpts.PullStrategy
    PULL_STRATEGY_NOT_EXISTS: DockerPullOpts.PullStrategy
    IMAGE_URI_FIELD_NUMBER: _ClassVar[int]
    PULL_STRATEGY_FIELD_NUMBER: _ClassVar[int]
    AUTH_FIELD_NUMBER: _ClassVar[int]
    image_uri: str
    pull_strategy: DockerPullOpts.PullStrategy
    auth: DockerPullAuth
    def __init__(self, image_uri: _Optional[str] = ..., pull_strategy: _Optional[_Union[DockerPullOpts.PullStrategy, str]] = ..., auth: _Optional[_Union[DockerPullAuth, _Mapping]] = ...) -> None: ...

class DockerPullAuth(_message.Message):
    __slots__ = ("basic", "aws_ecr")
    BASIC_FIELD_NUMBER: _ClassVar[int]
    AWS_ECR_FIELD_NUMBER: _ClassVar[int]
    basic: BasicAuth
    aws_ecr: AWSECRAuth
    def __init__(self, basic: _Optional[_Union[BasicAuth, _Mapping]] = ..., aws_ecr: _Optional[_Union[AWSECRAuth, _Mapping]] = ...) -> None: ...

class BasicAuth(_message.Message):
    __slots__ = ("username", "password")
    USERNAME_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    username: str
    password: str
    def __init__(self, username: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class AWSECRAuth(_message.Message):
    __slots__ = ("region", "aws_access_key_id", "aws_secret_key")
    REGION_FIELD_NUMBER: _ClassVar[int]
    AWS_ACCESS_KEY_ID_FIELD_NUMBER: _ClassVar[int]
    AWS_SECRET_KEY_FIELD_NUMBER: _ClassVar[int]
    region: str
    aws_access_key_id: str
    aws_secret_key: str
    def __init__(self, region: _Optional[str] = ..., aws_access_key_id: _Optional[str] = ..., aws_secret_key: _Optional[str] = ...) -> None: ...

class ExecRequest(_message.Message):
    __slots__ = ("runtime_id", "exec_id", "barrier_id", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    EXEC_ID_FIELD_NUMBER: _ClassVar[int]
    BARRIER_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    exec_id: str
    barrier_id: str
    opts: ExecOpts
    def __init__(self, runtime_id: _Optional[str] = ..., exec_id: _Optional[str] = ..., barrier_id: _Optional[str] = ..., opts: _Optional[_Union[ExecOpts, _Mapping]] = ...) -> None: ...

class ExecOpts(_message.Message):
    __slots__ = ("name", "args", "env", "tags")
    class TagsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: str
        def __init__(self, key: _Optional[str] = ..., value: _Optional[str] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    ARGS_FIELD_NUMBER: _ClassVar[int]
    ENV_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    name: str
    args: _containers.RepeatedScalarFieldContainer[str]
    env: _containers.RepeatedScalarFieldContainer[str]
    tags: _containers.ScalarMap[str, str]
    def __init__(self, name: _Optional[str] = ..., args: _Optional[_Iterable[str]] = ..., env: _Optional[_Iterable[str]] = ..., tags: _Optional[_Mapping[str, str]] = ...) -> None: ...

class ExecResponse(_message.Message):
    __slots__ = ("exit_code",)
    EXIT_CODE_FIELD_NUMBER: _ClassVar[int]
    exit_code: int
    def __init__(self, exit_code: _Optional[int] = ...) -> None: ...

class FileTransfer(_message.Message):
    __slots__ = ("runtime_id", "transfer_id", "file_id", "header", "body", "trailer")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    TRANSFER_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_ID_FIELD_NUMBER: _ClassVar[int]
    HEADER_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    TRAILER_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    transfer_id: str
    file_id: str
    header: FileTransferHeader
    body: FileTransferBody
    trailer: FileTransferTrailer
    def __init__(self, runtime_id: _Optional[str] = ..., transfer_id: _Optional[str] = ..., file_id: _Optional[str] = ..., header: _Optional[_Union[FileTransferHeader, _Mapping]] = ..., body: _Optional[_Union[FileTransferBody, _Mapping]] = ..., trailer: _Optional[_Union[FileTransferTrailer, _Mapping]] = ...) -> None: ...

class FileTransferHeader(_message.Message):
    __slots__ = ("is_dir", "src_path", "dest_path", "mode", "size")
    IS_DIR_FIELD_NUMBER: _ClassVar[int]
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    MODE_FIELD_NUMBER: _ClassVar[int]
    SIZE_FIELD_NUMBER: _ClassVar[int]
    is_dir: bool
    src_path: str
    dest_path: str
    mode: int
    size: int
    def __init__(self, is_dir: bool = ..., src_path: _Optional[str] = ..., dest_path: _Optional[str] = ..., mode: _Optional[int] = ..., size: _Optional[int] = ...) -> None: ...

class FileTransferBody(_message.Message):
    __slots__ = ("offset", "data")
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    offset: int
    data: bytes
    def __init__(self, offset: _Optional[int] = ..., data: _Optional[bytes] = ...) -> None: ...

class FileTransferTrailer(_message.Message):
    __slots__ = ("md5",)
    MD5_FIELD_NUMBER: _ClassVar[int]
    md5: bytes
    def __init__(self, md5: _Optional[bytes] = ...) -> None: ...

class ImportResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ExportRequest(_message.Message):
    __slots__ = ("runtime_id", "export_id", "src_path", "opts")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    EXPORT_ID_FIELD_NUMBER: _ClassVar[int]
    SRC_PATH_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    export_id: str
    src_path: str
    opts: ExportOpts
    def __init__(self, runtime_id: _Optional[str] = ..., export_id: _Optional[str] = ..., src_path: _Optional[str] = ..., opts: _Optional[_Union[ExportOpts, _Mapping]] = ...) -> None: ...

class ExportOpts(_message.Message):
    __slots__ = ("dest_path", "excludes")
    DEST_PATH_FIELD_NUMBER: _ClassVar[int]
    EXCLUDES_FIELD_NUMBER: _ClassVar[int]
    dest_path: str
    excludes: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, dest_path: _Optional[str] = ..., excludes: _Optional[_Iterable[str]] = ...) -> None: ...

class CloseRequest(_message.Message):
    __slots__ = ("runtime_id", "barrier_id")
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    BARRIER_ID_FIELD_NUMBER: _ClassVar[int]
    runtime_id: str
    barrier_id: str
    def __init__(self, runtime_id: _Optional[str] = ..., barrier_id: _Optional[str] = ...) -> None: ...

class CloseResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
