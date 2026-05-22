from . import executor_pb2 as _executor_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class TenderRequest(_message.Message):
    __slots__ = ("build_id", "tender_id", "opts")
    BUILD_ID_FIELD_NUMBER: _ClassVar[int]
    TENDER_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    build_id: str
    tender_id: str
    opts: _executor_pb2.RuntimeOpts
    def __init__(self, build_id: _Optional[str] = ..., tender_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.RuntimeOpts, _Mapping]] = ...) -> None: ...

class RuntimeContract(_message.Message):
    __slots__ = ("tender_id", "contract_id", "runtime_id", "opts", "sys_info", "executor_info")
    TENDER_ID_FIELD_NUMBER: _ClassVar[int]
    CONTRACT_ID_FIELD_NUMBER: _ClassVar[int]
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    SYS_INFO_FIELD_NUMBER: _ClassVar[int]
    EXECUTOR_INFO_FIELD_NUMBER: _ClassVar[int]
    tender_id: str
    contract_id: str
    runtime_id: str
    opts: _executor_pb2.RuntimeOpts
    sys_info: _executor_pb2.SystemInfo
    executor_info: _executor_pb2.ExecutorInfo
    def __init__(self, tender_id: _Optional[str] = ..., contract_id: _Optional[str] = ..., runtime_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.RuntimeOpts, _Mapping]] = ..., sys_info: _Optional[_Union[_executor_pb2.SystemInfo, _Mapping]] = ..., executor_info: _Optional[_Union[_executor_pb2.ExecutorInfo, _Mapping]] = ...) -> None: ...

class TenderResponse(_message.Message):
    __slots__ = ("contracts",)
    CONTRACTS_FIELD_NUMBER: _ClassVar[int]
    contracts: _containers.RepeatedCompositeFieldContainer[RuntimeContract]
    def __init__(self, contracts: _Optional[_Iterable[_Union[RuntimeContract, _Mapping]]] = ...) -> None: ...

class SettlementRequest(_message.Message):
    __slots__ = ("contract",)
    CONTRACT_FIELD_NUMBER: _ClassVar[int]
    contract: RuntimeContract
    def __init__(self, contract: _Optional[_Union[RuntimeContract, _Mapping]] = ...) -> None: ...

class SettlementResponse(_message.Message):
    __slots__ = ("connection_info",)
    CONNECTION_INFO_FIELD_NUMBER: _ClassVar[int]
    connection_info: RuntimeConnectionInfo
    def __init__(self, connection_info: _Optional[_Union[RuntimeConnectionInfo, _Mapping]] = ...) -> None: ...

class RuntimeConnectionInfo(_message.Message):
    __slots__ = ("unix", "tcp")
    UNIX_FIELD_NUMBER: _ClassVar[int]
    TCP_FIELD_NUMBER: _ClassVar[int]
    unix: RuntimeTransportUnix
    tcp: RuntimeTransportTCP
    def __init__(self, unix: _Optional[_Union[RuntimeTransportUnix, _Mapping]] = ..., tcp: _Optional[_Union[RuntimeTransportTCP, _Mapping]] = ...) -> None: ...

class RuntimeTransportUnix(_message.Message):
    __slots__ = ("socket_path",)
    SOCKET_PATH_FIELD_NUMBER: _ClassVar[int]
    socket_path: str
    def __init__(self, socket_path: _Optional[str] = ...) -> None: ...

class RuntimeTransportTCP(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...
