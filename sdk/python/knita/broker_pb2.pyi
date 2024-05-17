from . import executor_pb2 as _executor_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class RuntimeTender(_message.Message):
    __slots__ = ("tender_id", "opts")
    TENDER_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    tender_id: str
    opts: _executor_pb2.Opts
    def __init__(self, tender_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.Opts, _Mapping]] = ...) -> None: ...

class RuntimeContracts(_message.Message):
    __slots__ = ("contracts",)
    CONTRACTS_FIELD_NUMBER: _ClassVar[int]
    contracts: _containers.RepeatedCompositeFieldContainer[RuntimeContract]
    def __init__(self, contracts: _Optional[_Iterable[_Union[RuntimeContract, _Mapping]]] = ...) -> None: ...

class RuntimeContract(_message.Message):
    __slots__ = ("contract_id", "runtime_id", "opts")
    CONTRACT_ID_FIELD_NUMBER: _ClassVar[int]
    RUNTIME_ID_FIELD_NUMBER: _ClassVar[int]
    OPTS_FIELD_NUMBER: _ClassVar[int]
    contract_id: str
    runtime_id: str
    opts: _executor_pb2.Opts
    def __init__(self, contract_id: _Optional[str] = ..., runtime_id: _Optional[str] = ..., opts: _Optional[_Union[_executor_pb2.Opts, _Mapping]] = ...) -> None: ...

class RuntimeSettlement(_message.Message):
    __slots__ = ("connection_info",)
    CONNECTION_INFO_FIELD_NUMBER: _ClassVar[int]
    connection_info: RuntimeConnectionInfo
    def __init__(self, connection_info: _Optional[_Union[RuntimeConnectionInfo, _Mapping]] = ...) -> None: ...

class RuntimeConnectionInfo(_message.Message):
    __slots__ = ("unix",)
    UNIX_FIELD_NUMBER: _ClassVar[int]
    unix: RuntimeTransportUnix
    def __init__(self, unix: _Optional[_Union[RuntimeTransportUnix, _Mapping]] = ...) -> None: ...

class RuntimeTransportUnix(_message.Message):
    __slots__ = ("socket_path",)
    SOCKET_PATH_FIELD_NUMBER: _ClassVar[int]
    socket_path: str
    def __init__(self, socket_path: _Optional[str] = ...) -> None: ...
