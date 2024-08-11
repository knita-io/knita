# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc
import warnings

from . import event_pb2 as events_dot_v1_dot_event__pb2
from . import executor_pb2 as executor_dot_v1_dot_executor__pb2

GRPC_GENERATED_VERSION = '1.63.0'
GRPC_VERSION = grpc.__version__
EXPECTED_ERROR_RELEASE = '1.65.0'
SCHEDULED_RELEASE_DATE = 'June 25, 2024'
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower
    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    warnings.warn(
        f'The grpc package installed is at version {GRPC_VERSION},'
        + f' but the generated code in executor/v1/executor_pb2_grpc.py depends on'
        + f' grpcio>={GRPC_GENERATED_VERSION}.'
        + f' Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}'
        + f' or downgrade your generated code using grpcio-tools<={GRPC_VERSION}.'
        + f' This warning will become an error in {EXPECTED_ERROR_RELEASE},'
        + f' scheduled for release on {SCHEDULED_RELEASE_DATE}.',
        RuntimeWarning
    )


class ExecutorStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Introspect = channel.unary_unary(
                '/executor.knita.io.Executor/Introspect',
                request_serializer=executor_dot_v1_dot_executor__pb2.IntrospectRequest.SerializeToString,
                response_deserializer=executor_dot_v1_dot_executor__pb2.IntrospectResponse.FromString,
                _registered_method=True)
        self.Events = channel.unary_stream(
                '/executor.knita.io.Executor/Events',
                request_serializer=executor_dot_v1_dot_executor__pb2.EventsRequest.SerializeToString,
                response_deserializer=events_dot_v1_dot_event__pb2.Event.FromString,
                _registered_method=True)
        self.Open = channel.unary_unary(
                '/executor.knita.io.Executor/Open',
                request_serializer=executor_dot_v1_dot_executor__pb2.OpenRequest.SerializeToString,
                response_deserializer=executor_dot_v1_dot_executor__pb2.OpenResponse.FromString,
                _registered_method=True)
        self.Heartbeat = channel.unary_unary(
                '/executor.knita.io.Executor/Heartbeat',
                request_serializer=executor_dot_v1_dot_executor__pb2.HeartbeatRequest.SerializeToString,
                response_deserializer=executor_dot_v1_dot_executor__pb2.HeartbeatResponse.FromString,
                _registered_method=True)
        self.Exec = channel.unary_unary(
                '/executor.knita.io.Executor/Exec',
                request_serializer=executor_dot_v1_dot_executor__pb2.ExecRequest.SerializeToString,
                response_deserializer=executor_dot_v1_dot_executor__pb2.ExecResponse.FromString,
                _registered_method=True)
        self.Import = channel.stream_unary(
                '/executor.knita.io.Executor/Import',
                request_serializer=executor_dot_v1_dot_executor__pb2.FileTransfer.SerializeToString,
                response_deserializer=executor_dot_v1_dot_executor__pb2.ImportResponse.FromString,
                _registered_method=True)
        self.Export = channel.unary_stream(
                '/executor.knita.io.Executor/Export',
                request_serializer=executor_dot_v1_dot_executor__pb2.ExportRequest.SerializeToString,
                response_deserializer=executor_dot_v1_dot_executor__pb2.FileTransfer.FromString,
                _registered_method=True)
        self.Close = channel.unary_unary(
                '/executor.knita.io.Executor/Close',
                request_serializer=executor_dot_v1_dot_executor__pb2.CloseRequest.SerializeToString,
                response_deserializer=executor_dot_v1_dot_executor__pb2.CloseResponse.FromString,
                _registered_method=True)


class ExecutorServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Introspect(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Events(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Open(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Heartbeat(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Exec(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Import(self, request_iterator, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Export(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Close(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ExecutorServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Introspect': grpc.unary_unary_rpc_method_handler(
                    servicer.Introspect,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.IntrospectRequest.FromString,
                    response_serializer=executor_dot_v1_dot_executor__pb2.IntrospectResponse.SerializeToString,
            ),
            'Events': grpc.unary_stream_rpc_method_handler(
                    servicer.Events,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.EventsRequest.FromString,
                    response_serializer=events_dot_v1_dot_event__pb2.Event.SerializeToString,
            ),
            'Open': grpc.unary_unary_rpc_method_handler(
                    servicer.Open,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.OpenRequest.FromString,
                    response_serializer=executor_dot_v1_dot_executor__pb2.OpenResponse.SerializeToString,
            ),
            'Heartbeat': grpc.unary_unary_rpc_method_handler(
                    servicer.Heartbeat,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.HeartbeatRequest.FromString,
                    response_serializer=executor_dot_v1_dot_executor__pb2.HeartbeatResponse.SerializeToString,
            ),
            'Exec': grpc.unary_unary_rpc_method_handler(
                    servicer.Exec,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.ExecRequest.FromString,
                    response_serializer=executor_dot_v1_dot_executor__pb2.ExecResponse.SerializeToString,
            ),
            'Import': grpc.stream_unary_rpc_method_handler(
                    servicer.Import,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.FileTransfer.FromString,
                    response_serializer=executor_dot_v1_dot_executor__pb2.ImportResponse.SerializeToString,
            ),
            'Export': grpc.unary_stream_rpc_method_handler(
                    servicer.Export,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.ExportRequest.FromString,
                    response_serializer=executor_dot_v1_dot_executor__pb2.FileTransfer.SerializeToString,
            ),
            'Close': grpc.unary_unary_rpc_method_handler(
                    servicer.Close,
                    request_deserializer=executor_dot_v1_dot_executor__pb2.CloseRequest.FromString,
                    response_serializer=executor_dot_v1_dot_executor__pb2.CloseResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'executor.knita.io.Executor', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Executor(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Introspect(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/executor.knita.io.Executor/Introspect',
            executor_dot_v1_dot_executor__pb2.IntrospectRequest.SerializeToString,
            executor_dot_v1_dot_executor__pb2.IntrospectResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Events(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(
            request,
            target,
            '/executor.knita.io.Executor/Events',
            executor_dot_v1_dot_executor__pb2.EventsRequest.SerializeToString,
            events_dot_v1_dot_event__pb2.Event.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Open(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/executor.knita.io.Executor/Open',
            executor_dot_v1_dot_executor__pb2.OpenRequest.SerializeToString,
            executor_dot_v1_dot_executor__pb2.OpenResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Heartbeat(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/executor.knita.io.Executor/Heartbeat',
            executor_dot_v1_dot_executor__pb2.HeartbeatRequest.SerializeToString,
            executor_dot_v1_dot_executor__pb2.HeartbeatResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Exec(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/executor.knita.io.Executor/Exec',
            executor_dot_v1_dot_executor__pb2.ExecRequest.SerializeToString,
            executor_dot_v1_dot_executor__pb2.ExecResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Import(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_unary(
            request_iterator,
            target,
            '/executor.knita.io.Executor/Import',
            executor_dot_v1_dot_executor__pb2.FileTransfer.SerializeToString,
            executor_dot_v1_dot_executor__pb2.ImportResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Export(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(
            request,
            target,
            '/executor.knita.io.Executor/Export',
            executor_dot_v1_dot_executor__pb2.ExportRequest.SerializeToString,
            executor_dot_v1_dot_executor__pb2.FileTransfer.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)

    @staticmethod
    def Close(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(
            request,
            target,
            '/executor.knita.io.Executor/Close',
            executor_dot_v1_dot_executor__pb2.CloseRequest.SerializeToString,
            executor_dot_v1_dot_executor__pb2.CloseResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True)
