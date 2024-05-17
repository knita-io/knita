# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: executor/v1/executor.proto
# Protobuf Python Version: 5.26.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1a\x65xecutor/v1/executor.proto\x12\x08\x65xecutor\"\x13\n\x11IntrospectRequest\"L\n\x12IntrospectResponse\x12\n\n\x02os\x18\x01 \x01(\t\x12\x0c\n\x04\x61rch\x18\x02 \x01(\t\x12\x0c\n\x04ncpu\x18\x04 \x01(\r\x12\x0e\n\x06labels\x18\x03 \x03(\t\"Q\n\x0bOpenRequest\x12\x10\n\x08\x62uild_id\x18\x01 \x01(\t\x12\x12\n\nruntime_id\x18\x02 \x01(\t\x12\x1c\n\x04opts\x18\x03 \x01(\x0b\x32\x0e.executor.Opts\"&\n\x0cOpenResponse\x12\x16\n\x0ework_directory\x18\x01 \x01(\t\"\xe4\x01\n\x04Opts\x12#\n\x04type\x18\x01 \x01(\x0e\x32\x15.executor.RuntimeType\x12\x0e\n\x06labels\x18\x04 \x03(\t\x12&\n\x04tags\x18\x05 \x03(\x0b\x32\x18.executor.Opts.TagsEntry\x12\"\n\x04host\x18\x02 \x01(\x0b\x32\x12.executor.HostOptsH\x00\x12&\n\x06\x64ocker\x18\x03 \x01(\x0b\x32\x14.executor.DockerOptsH\x00\x1a+\n\tTagsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x42\x06\n\x04opts\"\n\n\x08HostOpts\"5\n\nDockerOpts\x12\'\n\x05image\x18\x01 \x01(\x0b\x32\x18.executor.DockerPullOpts\"\x89\x02\n\x0e\x44ockerPullOpts\x12\x11\n\timage_uri\x18\x01 \x01(\t\x12<\n\rpull_strategy\x18\x02 \x01(\x0e\x32%.executor.DockerPullOpts.PullStrategy\x12&\n\x04\x61uth\x18\x03 \x01(\x0b\x32\x18.executor.DockerPullAuth\"~\n\x0cPullStrategy\x12\x1d\n\x19PULL_STRATEGY_UNSPECIFIED\x10\x00\x12\x17\n\x13PULL_STRATEGY_NEVER\x10\x01\x12\x18\n\x14PULL_STRATEGY_ALWAYS\x10\x02\x12\x1c\n\x18PULL_STRATEGY_NOT_EXISTS\x10\x03\"g\n\x0e\x44ockerPullAuth\x12$\n\x05\x62\x61sic\x18\x01 \x01(\x0b\x32\x13.executor.BasicAuthH\x00\x12\'\n\x07\x61ws_ecr\x18\x02 \x01(\x0b\x32\x14.executor.AWSECRAuthH\x00\x42\x06\n\x04\x61uth\"/\n\tBasicAuth\x12\x10\n\x08username\x18\x01 \x01(\t\x12\x10\n\x08password\x18\x02 \x01(\t\"O\n\nAWSECRAuth\x12\x0e\n\x06region\x18\x01 \x01(\t\x12\x19\n\x11\x61ws_access_key_id\x18\x02 \x01(\t\x12\x16\n\x0e\x61ws_secret_key\x18\x03 \x01(\t\"T\n\x0b\x45xecRequest\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x0f\n\x07\x65xec_id\x18\x02 \x01(\t\x12 \n\x04opts\x18\x03 \x01(\x0b\x32\x12.executor.ExecOpts\"\x8c\x01\n\x08\x45xecOpts\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04\x61rgs\x18\x02 \x03(\t\x12\x0b\n\x03\x65nv\x18\x03 \x03(\t\x12*\n\x04tags\x18\x04 \x03(\x0b\x32\x1c.executor.ExecOpts.TagsEntry\x1a+\n\tTagsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"!\n\x0c\x45xecResponse\x12\x11\n\texit_code\x18\x01 \x01(\x05\"\xce\x01\n\x0c\x46ileTransfer\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x11\n\timport_id\x18\x02 \x01(\t\x12\x0f\n\x07\x66ile_id\x18\x03 \x01(\t\x12,\n\x06header\x18\x04 \x01(\x0b\x32\x1c.executor.FileTransferHeader\x12(\n\x04\x62ody\x18\x05 \x01(\x0b\x32\x1a.executor.FileTransferBody\x12.\n\x07trailer\x18\x06 \x01(\x0b\x32\x1d.executor.FileTransferTrailer\"e\n\x12\x46ileTransferHeader\x12\x0e\n\x06is_dir\x18\x01 \x01(\x08\x12\x10\n\x08src_path\x18\x02 \x01(\t\x12\x11\n\tdest_path\x18\x03 \x01(\t\x12\x0c\n\x04mode\x18\x04 \x01(\r\x12\x0c\n\x04size\x18\x05 \x01(\x04\"0\n\x10\x46ileTransferBody\x12\x0e\n\x06offset\x18\x01 \x01(\x04\x12\x0c\n\x04\x64\x61ta\x18\x02 \x01(\x0c\"\"\n\x13\x46ileTransferTrailer\x12\x0b\n\x03md5\x18\x01 \x01(\x0c\"\x10\n\x0eImportResponse\"[\n\rExportRequest\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x11\n\texport_id\x18\x02 \x01(\t\x12\x10\n\x08src_path\x18\x03 \x01(\t\x12\x11\n\tdest_path\x18\x04 \x01(\t\"\"\n\x0c\x43loseRequest\x12\x12\n\nruntime_id\x18\x01 \x01(\t\"\x0f\n\rCloseResponse\"#\n\rEventsRequest\x12\x12\n\nruntime_id\x18\x01 \x01(\t\"\xb0\x04\n\x05\x45vent\x12\x10\n\x08\x62uild_id\x18\x01 \x01(\t\x12\x12\n\ngroup_name\x18\x03 \x01(\t\x12\x10\n\x08sequence\x18\x04 \x01(\x04\x12\x36\n\x0eruntime_opened\x18\x05 \x01(\x0b\x32\x1c.executor.RuntimeOpenedEventH\x00\x12.\n\nexec_start\x18\x06 \x01(\x0b\x32\x18.executor.ExecStartEventH\x00\x12*\n\x08\x65xec_end\x18\x07 \x01(\x0b\x32\x16.executor.ExecEndEventH\x00\x12\x32\n\x0cimport_start\x18\x08 \x01(\x0b\x32\x1a.executor.ImportStartEventH\x00\x12.\n\nimport_end\x18\t \x01(\x0b\x32\x18.executor.ImportEndEventH\x00\x12\x32\n\x0c\x65xport_start\x18\n \x01(\x0b\x32\x1a.executor.ExportStartEventH\x00\x12.\n\nexport_end\x18\x0b \x01(\x0b\x32\x18.executor.ExportEndEventH\x00\x12\'\n\x06stdout\x18\x0c \x01(\x0b\x32\x15.executor.StdoutEventH\x00\x12\'\n\x06stderr\x18\r \x01(\x0b\x32\x15.executor.StderrEventH\x00\x12\x36\n\x0eruntime_closed\x18\x0e \x01(\x0b\x32\x1c.executor.RuntimeClosedEventH\x00\x42\t\n\x07payload\"F\n\x12RuntimeOpenedEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x1c\n\x04opts\x18\x02 \x01(\x0b\x32\x0e.executor.Opts\"(\n\x12RuntimeClosedEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\"E\n\x0bStdoutEvent\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\x0c\x12(\n\x06source\x18\x02 \x01(\x0b\x32\x18.executor.LogEventSource\"E\n\x0bStderrEvent\x12\x0c\n\x04\x64\x61ta\x18\x01 \x01(\x0c\x12(\n\x06source\x18\x02 \x01(\x0b\x32\x18.executor.LogEventSource\"\xa3\x01\n\x0eLogEventSource\x12-\n\x07runtime\x18\x02 \x01(\x0b\x32\x1a.executor.LogSourceRuntimeH\x00\x12\'\n\x04\x65xec\x18\x03 \x01(\x0b\x32\x17.executor.LogSourceExecH\x00\x12/\n\x08\x64irector\x18\x04 \x01(\x0b\x32\x1b.executor.LogSourceDirectorH\x00\x42\x08\n\x06source\"&\n\x10LogSourceRuntime\x12\x12\n\nruntime_id\x18\x01 \x01(\t\"D\n\rLogSourceExec\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x0f\n\x07\x65xec_id\x18\x02 \x01(\t\x12\x0e\n\x06system\x18\x03 \x01(\x08\"\x13\n\x11LogSourceDirector\"\xb6\x01\n\x0e\x45xecStartEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x0f\n\x07\x65xec_id\x18\x02 \x01(\t\x12 \n\x04opts\x18\x03 \x01(\x0b\x32\x12.executor.ExecOpts\x12\x30\n\x04tags\x18\x04 \x03(\x0b\x32\".executor.ExecStartEvent.TagsEntry\x1a+\n\tTagsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\"9\n\x10ImportStartEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x11\n\timport_id\x18\x02 \x01(\t\"7\n\x0eImportEndEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x11\n\timport_id\x18\x02 \x01(\t\"9\n\x10\x45xportStartEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x11\n\texport_id\x18\x02 \x01(\t\"7\n\x0e\x45xportEndEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x11\n\texport_id\x18\x02 \x01(\t\"U\n\x0c\x45xecEndEvent\x12\x12\n\nruntime_id\x18\x01 \x01(\t\x12\x0f\n\x07\x65xec_id\x18\x02 \x01(\t\x12\r\n\x05\x65rror\x18\x03 \x01(\t\x12\x11\n\texit_code\x18\x04 \x01(\x05*L\n\x0bRuntimeType\x12\x17\n\x13RUNTIME_UNSPECIFIED\x10\x00\x12\x10\n\x0cRUNTIME_HOST\x10\x01\x12\x12\n\x0eRUNTIME_DOCKER\x10\x02\x32\xac\x03\n\x08\x45xecutor\x12G\n\nIntrospect\x12\x1b.executor.IntrospectRequest\x1a\x1c.executor.IntrospectResponse\x12\x35\n\x04Open\x12\x15.executor.OpenRequest\x1a\x16.executor.OpenResponse\x12\x35\n\x04\x45xec\x12\x15.executor.ExecRequest\x1a\x16.executor.ExecResponse\x12<\n\x06Import\x12\x16.executor.FileTransfer\x1a\x18.executor.ImportResponse(\x01\x12;\n\x06\x45xport\x12\x17.executor.ExportRequest\x1a\x16.executor.FileTransfer0\x01\x12\x38\n\x05\x43lose\x12\x16.executor.CloseRequest\x1a\x17.executor.CloseResponse\x12\x34\n\x06\x45vents\x12\x17.executor.EventsRequest\x1a\x0f.executor.Event0\x01\x42+Z)github.com/knita-io/knita/api/executor/v1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'executor.v1.executor_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z)github.com/knita-io/knita/api/executor/v1'
  _globals['_OPTS_TAGSENTRY']._loaded_options = None
  _globals['_OPTS_TAGSENTRY']._serialized_options = b'8\001'
  _globals['_EXECOPTS_TAGSENTRY']._loaded_options = None
  _globals['_EXECOPTS_TAGSENTRY']._serialized_options = b'8\001'
  _globals['_EXECSTARTEVENT_TAGSENTRY']._loaded_options = None
  _globals['_EXECSTARTEVENT_TAGSENTRY']._serialized_options = b'8\001'
  _globals['_RUNTIMETYPE']._serialized_start=3546
  _globals['_RUNTIMETYPE']._serialized_end=3622
  _globals['_INTROSPECTREQUEST']._serialized_start=40
  _globals['_INTROSPECTREQUEST']._serialized_end=59
  _globals['_INTROSPECTRESPONSE']._serialized_start=61
  _globals['_INTROSPECTRESPONSE']._serialized_end=137
  _globals['_OPENREQUEST']._serialized_start=139
  _globals['_OPENREQUEST']._serialized_end=220
  _globals['_OPENRESPONSE']._serialized_start=222
  _globals['_OPENRESPONSE']._serialized_end=260
  _globals['_OPTS']._serialized_start=263
  _globals['_OPTS']._serialized_end=491
  _globals['_OPTS_TAGSENTRY']._serialized_start=440
  _globals['_OPTS_TAGSENTRY']._serialized_end=483
  _globals['_HOSTOPTS']._serialized_start=493
  _globals['_HOSTOPTS']._serialized_end=503
  _globals['_DOCKEROPTS']._serialized_start=505
  _globals['_DOCKEROPTS']._serialized_end=558
  _globals['_DOCKERPULLOPTS']._serialized_start=561
  _globals['_DOCKERPULLOPTS']._serialized_end=826
  _globals['_DOCKERPULLOPTS_PULLSTRATEGY']._serialized_start=700
  _globals['_DOCKERPULLOPTS_PULLSTRATEGY']._serialized_end=826
  _globals['_DOCKERPULLAUTH']._serialized_start=828
  _globals['_DOCKERPULLAUTH']._serialized_end=931
  _globals['_BASICAUTH']._serialized_start=933
  _globals['_BASICAUTH']._serialized_end=980
  _globals['_AWSECRAUTH']._serialized_start=982
  _globals['_AWSECRAUTH']._serialized_end=1061
  _globals['_EXECREQUEST']._serialized_start=1063
  _globals['_EXECREQUEST']._serialized_end=1147
  _globals['_EXECOPTS']._serialized_start=1150
  _globals['_EXECOPTS']._serialized_end=1290
  _globals['_EXECOPTS_TAGSENTRY']._serialized_start=440
  _globals['_EXECOPTS_TAGSENTRY']._serialized_end=483
  _globals['_EXECRESPONSE']._serialized_start=1292
  _globals['_EXECRESPONSE']._serialized_end=1325
  _globals['_FILETRANSFER']._serialized_start=1328
  _globals['_FILETRANSFER']._serialized_end=1534
  _globals['_FILETRANSFERHEADER']._serialized_start=1536
  _globals['_FILETRANSFERHEADER']._serialized_end=1637
  _globals['_FILETRANSFERBODY']._serialized_start=1639
  _globals['_FILETRANSFERBODY']._serialized_end=1687
  _globals['_FILETRANSFERTRAILER']._serialized_start=1689
  _globals['_FILETRANSFERTRAILER']._serialized_end=1723
  _globals['_IMPORTRESPONSE']._serialized_start=1725
  _globals['_IMPORTRESPONSE']._serialized_end=1741
  _globals['_EXPORTREQUEST']._serialized_start=1743
  _globals['_EXPORTREQUEST']._serialized_end=1834
  _globals['_CLOSEREQUEST']._serialized_start=1836
  _globals['_CLOSEREQUEST']._serialized_end=1870
  _globals['_CLOSERESPONSE']._serialized_start=1872
  _globals['_CLOSERESPONSE']._serialized_end=1887
  _globals['_EVENTSREQUEST']._serialized_start=1889
  _globals['_EVENTSREQUEST']._serialized_end=1924
  _globals['_EVENT']._serialized_start=1927
  _globals['_EVENT']._serialized_end=2487
  _globals['_RUNTIMEOPENEDEVENT']._serialized_start=2489
  _globals['_RUNTIMEOPENEDEVENT']._serialized_end=2559
  _globals['_RUNTIMECLOSEDEVENT']._serialized_start=2561
  _globals['_RUNTIMECLOSEDEVENT']._serialized_end=2601
  _globals['_STDOUTEVENT']._serialized_start=2603
  _globals['_STDOUTEVENT']._serialized_end=2672
  _globals['_STDERREVENT']._serialized_start=2674
  _globals['_STDERREVENT']._serialized_end=2743
  _globals['_LOGEVENTSOURCE']._serialized_start=2746
  _globals['_LOGEVENTSOURCE']._serialized_end=2909
  _globals['_LOGSOURCERUNTIME']._serialized_start=2911
  _globals['_LOGSOURCERUNTIME']._serialized_end=2949
  _globals['_LOGSOURCEEXEC']._serialized_start=2951
  _globals['_LOGSOURCEEXEC']._serialized_end=3019
  _globals['_LOGSOURCEDIRECTOR']._serialized_start=3021
  _globals['_LOGSOURCEDIRECTOR']._serialized_end=3040
  _globals['_EXECSTARTEVENT']._serialized_start=3043
  _globals['_EXECSTARTEVENT']._serialized_end=3225
  _globals['_EXECSTARTEVENT_TAGSENTRY']._serialized_start=440
  _globals['_EXECSTARTEVENT_TAGSENTRY']._serialized_end=483
  _globals['_IMPORTSTARTEVENT']._serialized_start=3227
  _globals['_IMPORTSTARTEVENT']._serialized_end=3284
  _globals['_IMPORTENDEVENT']._serialized_start=3286
  _globals['_IMPORTENDEVENT']._serialized_end=3341
  _globals['_EXPORTSTARTEVENT']._serialized_start=3343
  _globals['_EXPORTSTARTEVENT']._serialized_end=3400
  _globals['_EXPORTENDEVENT']._serialized_start=3402
  _globals['_EXPORTENDEVENT']._serialized_end=3457
  _globals['_EXECENDEVENT']._serialized_start=3459
  _globals['_EXECENDEVENT']._serialized_end=3544
  _globals['_EXECUTOR']._serialized_start=3625
  _globals['_EXECUTOR']._serialized_end=4053
# @@protoc_insertion_point(module_scope)