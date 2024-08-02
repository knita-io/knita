# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: director/v1/director.proto
# Protobuf Python Version: 5.26.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import executor_pb2 as executor_dot_v1_dot_executor__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1a\x64irector/v1/director.proto\x12\x08\x64irector\x1a\x1a\x65xecutor/v1/executor.proto\"=\n\x0bOpenRequest\x12\x10\n\x08\x62uild_id\x18\x01 \x01(\t\x12\x1c\n\x04opts\x18\x02 \x01(\x0b\x32\x0e.executor.Opts\"b\n\x0cOpenResponse\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x16\n\x0ework_directory\x18\x02 \x01(\t\x12&\n\x08sys_info\x18\x03 \x01(\x0b\x32\x14.executor.SystemInfo\"Y\n\rImportRequest\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x10\n\x08src_path\x18\x02 \x01(\t\x12\"\n\x04opts\x18\x03 \x01(\x0b\x32\x14.director.ImportOpts\"1\n\nImportOpts\x12\x11\n\tdest_path\x18\x01 \x01(\t\x12\x10\n\x08\x65xcludes\x18\x02 \x03(\t\"\x10\n\x0eImportResponse\"Y\n\rExportRequest\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x10\n\x08src_path\x18\x02 \x01(\t\x12\"\n\x04opts\x18\x03 \x01(\x0b\x32\x14.director.ExportOpts\"1\n\nExportOpts\x12\x11\n\tdest_path\x18\x01 \x01(\t\x12\x10\n\x08\x65xcludes\x18\x02 \x03(\t\"\x10\n\x0e\x45xportResponse\"C\n\x0b\x45xecRequest\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12 \n\x04opts\x18\x02 \x01(\x0b\x32\x12.executor.ExecOpts\"\xcc\x01\n\tExecEvent\x12.\n\nexec_start\x18\x01 \x01(\x0b\x32\x18.director.ExecStartEventH\x00\x12+\n\x06stdout\x18\x02 \x01(\x0b\x32\x19.director.ExecStdoutEventH\x00\x12+\n\x06stderr\x18\x03 \x01(\x0b\x32\x19.director.ExecStderrEventH\x00\x12*\n\x08\x65xec_end\x18\x04 \x01(\x0b\x32\x16.director.ExecEndEventH\x00\x42\t\n\x07payload\"\x10\n\x0e\x45xecStartEvent\"\x1f\n\x0f\x45xecStdoutEvent\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\x0c\"\x1f\n\x0f\x45xecStderrEvent\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\x0c\"0\n\x0c\x45xecEndEvent\x12\r\n\x05\x65rror\x18\x01 \x01(\t\x12\x11\n\texit_code\x18\x02 \x01(\x05\x32\xab\x02\n\x08\x44irector\x12\x35\n\x04Open\x12\x15.director.OpenRequest\x1a\x16.director.OpenResponse\x12\x34\n\x04\x45xec\x12\x15.director.ExecRequest\x1a\x13.director.ExecEvent0\x01\x12;\n\x06Import\x12\x17.director.ImportRequest\x1a\x18.director.ImportResponse\x12;\n\x06\x45xport\x12\x17.director.ExportRequest\x1a\x18.director.ExportResponse\x12\x38\n\x05\x43lose\x12\x16.executor.CloseRequest\x1a\x17.executor.CloseResponseB+Z)github.com/knita-io/knita/api/director/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'director.v1.director_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z)github.com/knita-io/knita/api/director/v1'
  _globals['_OPENREQUEST']._serialized_start=68
  _globals['_OPENREQUEST']._serialized_end=129
  _globals['_OPENRESPONSE']._serialized_start=131
  _globals['_OPENRESPONSE']._serialized_end=229
  _globals['_IMPORTREQUEST']._serialized_start=231
  _globals['_IMPORTREQUEST']._serialized_end=320
  _globals['_IMPORTOPTS']._serialized_start=322
  _globals['_IMPORTOPTS']._serialized_end=371
  _globals['_IMPORTRESPONSE']._serialized_start=373
  _globals['_IMPORTRESPONSE']._serialized_end=389
  _globals['_EXPORTREQUEST']._serialized_start=391
  _globals['_EXPORTREQUEST']._serialized_end=480
  _globals['_EXPORTOPTS']._serialized_start=482
  _globals['_EXPORTOPTS']._serialized_end=531
  _globals['_EXPORTRESPONSE']._serialized_start=533
  _globals['_EXPORTRESPONSE']._serialized_end=549
  _globals['_EXECREQUEST']._serialized_start=551
  _globals['_EXECREQUEST']._serialized_end=618
  _globals['_EXECEVENT']._serialized_start=621
  _globals['_EXECEVENT']._serialized_end=825
  _globals['_EXECSTARTEVENT']._serialized_start=827
  _globals['_EXECSTARTEVENT']._serialized_end=843
  _globals['_EXECSTDOUTEVENT']._serialized_start=845
  _globals['_EXECSTDOUTEVENT']._serialized_end=876
  _globals['_EXECSTDERREVENT']._serialized_start=878
  _globals['_EXECSTDERREVENT']._serialized_end=909
  _globals['_EXECENDEVENT']._serialized_start=911
  _globals['_EXECENDEVENT']._serialized_end=959
  _globals['_DIRECTOR']._serialized_start=962
  _globals['_DIRECTOR']._serialized_end=1261
# @@protoc_insertion_point(module_scope)
