// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: observer/v1/observer.proto

package v1

import (
	context "context"
	v1 "github.com/knita-io/knita/api/events/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Observer_Watch_FullMethodName = "/observer.knita.io.Observer/Watch"
)

// ObserverClient is the client API for Observer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObserverClient interface {
	Watch(ctx context.Context, opts ...grpc.CallOption) (Observer_WatchClient, error)
}

type observerClient struct {
	cc grpc.ClientConnInterface
}

func NewObserverClient(cc grpc.ClientConnInterface) ObserverClient {
	return &observerClient{cc}
}

func (c *observerClient) Watch(ctx context.Context, opts ...grpc.CallOption) (Observer_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Observer_ServiceDesc.Streams[0], Observer_Watch_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &observerWatchClient{stream}
	return x, nil
}

type Observer_WatchClient interface {
	Send(*v1.Event) error
	CloseAndRecv() (*WatchResponse, error)
	grpc.ClientStream
}

type observerWatchClient struct {
	grpc.ClientStream
}

func (x *observerWatchClient) Send(m *v1.Event) error {
	return x.ClientStream.SendMsg(m)
}

func (x *observerWatchClient) CloseAndRecv() (*WatchResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ObserverServer is the server API for Observer service.
// All implementations must embed UnimplementedObserverServer
// for forward compatibility
type ObserverServer interface {
	Watch(Observer_WatchServer) error
	mustEmbedUnimplementedObserverServer()
}

// UnimplementedObserverServer must be embedded to have forward compatible implementations.
type UnimplementedObserverServer struct {
}

func (UnimplementedObserverServer) Watch(Observer_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedObserverServer) mustEmbedUnimplementedObserverServer() {}

// UnsafeObserverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObserverServer will
// result in compilation errors.
type UnsafeObserverServer interface {
	mustEmbedUnimplementedObserverServer()
}

func RegisterObserverServer(s grpc.ServiceRegistrar, srv ObserverServer) {
	s.RegisterService(&Observer_ServiceDesc, srv)
}

func _Observer_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ObserverServer).Watch(&observerWatchServer{stream})
}

type Observer_WatchServer interface {
	SendAndClose(*WatchResponse) error
	Recv() (*v1.Event, error)
	grpc.ServerStream
}

type observerWatchServer struct {
	grpc.ServerStream
}

func (x *observerWatchServer) SendAndClose(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *observerWatchServer) Recv() (*v1.Event, error) {
	m := new(v1.Event)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Observer_ServiceDesc is the grpc.ServiceDesc for Observer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Observer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "observer.knita.io.Observer",
	HandlerType: (*ObserverServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Observer_Watch_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "observer/v1/observer.proto",
}
