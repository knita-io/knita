# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: broker/v1/broker.proto
# Protobuf Python Version: 5.26.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import executor_pb2 as executor_dot_v1_dot_executor__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x16\x62roker/v1/broker.proto\x12\x06\x62roker\x1a\x1a\x65xecutor/v1/executor.proto\"@\n\rRuntimeTender\x12\x11\n\ttender_id\x18\x01 \x01(\t\x12\x1c\n\x04opts\x18\x02 \x01(\x0b\x32\x0e.executor.Opts\">\n\x10RuntimeContracts\x12*\n\tcontracts\x18\x01 \x03(\x0b\x32\x17.broker.RuntimeContract\"X\n\x0fRuntimeContract\x12\x13\n\x0b\x63ontract_id\x18\x01 \x01(\t\x12\x12\n\nruntime_id\x18\x02 \x01(\t\x12\x1c\n\x04opts\x18\x03 \x01(\x0b\x32\x0e.executor.Opts\"K\n\x11RuntimeSettlement\x12\x36\n\x0f\x63onnection_info\x18\x01 \x01(\x0b\x32\x1d.broker.RuntimeConnectionInfo\"R\n\x15RuntimeConnectionInfo\x12,\n\x04unix\x18\x01 \x01(\x0b\x32\x1c.broker.RuntimeTransportUnixH\x00\x42\x0b\n\ttransport\"+\n\x14RuntimeTransportUnix\x12\x13\n\x0bsocket_path\x18\x01 \x01(\t2\x88\x01\n\rRuntimeBroker\x12\x39\n\x06Tender\x12\x15.broker.RuntimeTender\x1a\x18.broker.RuntimeContracts\x12<\n\x06Settle\x12\x17.broker.RuntimeContract\x1a\x19.broker.RuntimeSettlementB)Z\'github.com/knita-io/knita/api/broker/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'broker.v1.broker_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\'github.com/knita-io/knita/api/broker/v1'
  _globals['_RUNTIMETENDER']._serialized_start=62
  _globals['_RUNTIMETENDER']._serialized_end=126
  _globals['_RUNTIMECONTRACTS']._serialized_start=128
  _globals['_RUNTIMECONTRACTS']._serialized_end=190
  _globals['_RUNTIMECONTRACT']._serialized_start=192
  _globals['_RUNTIMECONTRACT']._serialized_end=280
  _globals['_RUNTIMESETTLEMENT']._serialized_start=282
  _globals['_RUNTIMESETTLEMENT']._serialized_end=357
  _globals['_RUNTIMECONNECTIONINFO']._serialized_start=359
  _globals['_RUNTIMECONNECTIONINFO']._serialized_end=441
  _globals['_RUNTIMETRANSPORTUNIX']._serialized_start=443
  _globals['_RUNTIMETRANSPORTUNIX']._serialized_end=486
  _globals['_RUNTIMEBROKER']._serialized_start=489
  _globals['_RUNTIMEBROKER']._serialized_end=625
# @@protoc_insertion_point(module_scope)
