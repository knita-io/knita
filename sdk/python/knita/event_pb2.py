# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: events/v1/event.proto
# Protobuf Python Version: 5.26.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import any_pb2 as google_dot_protobuf_dot_any__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x15\x65vents/v1/event.proto\x12\x0f\x65vents.knita.io\x1a\x19google/protobuf/any.proto\"\x95\x02\n\x04Meta\x12\x10\n\x08\x62uild_id\x18\x01 \x01(\t\x12\x16\n\x0e\x63orrelation_id\x18\x02 \x01(\t\x12\x10\n\x08sequence\x18\x03 \x01(\x04\x12\x31\n\x06labels\x18\x04 \x03(\x0b\x32!.events.knita.io.Meta.LabelsEntry\x12;\n\x0b\x61nnotations\x18\x05 \x03(\x0b\x32&.events.knita.io.Meta.AnnotationsEntry\x1a-\n\x0bLabelsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x1a\x32\n\x10\x41nnotationsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"S\n\x05\x45vent\x12#\n\x04meta\x18\x01 \x01(\x0b\x32\x15.events.knita.io.Meta\x12%\n\x07payload\x18\x02 \x01(\x0b\x32\x14.google.protobuf.AnyB)Z\'github.com/knita-io/knita/api/events/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'events.v1.event_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\'github.com/knita-io/knita/api/events/v1'
  _globals['_META_LABELSENTRY']._loaded_options = None
  _globals['_META_LABELSENTRY']._serialized_options = b'8\001'
  _globals['_META_ANNOTATIONSENTRY']._loaded_options = None
  _globals['_META_ANNOTATIONSENTRY']._serialized_options = b'8\001'
  _globals['_META']._serialized_start=70
  _globals['_META']._serialized_end=347
  _globals['_META_LABELSENTRY']._serialized_start=250
  _globals['_META_LABELSENTRY']._serialized_end=295
  _globals['_META_ANNOTATIONSENTRY']._serialized_start=297
  _globals['_META_ANNOTATIONSENTRY']._serialized_end=347
  _globals['_EVENT']._serialized_start=349
  _globals['_EVENT']._serialized_end=432
# @@protoc_insertion_point(module_scope)