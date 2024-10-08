// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: director/v1/director.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Director_Open_FullMethodName   = "/director.knita.io.Director/Open"
	Director_Exec_FullMethodName   = "/director.knita.io.Director/Exec"
	Director_Import_FullMethodName = "/director.knita.io.Director/Import"
	Director_Export_FullMethodName = "/director.knita.io.Director/Export"
	Director_Close_FullMethodName  = "/director.knita.io.Director/Close"
)

// DirectorClient is the client API for Director service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DirectorClient interface {
	Open(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenResponse, error)
	Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (Director_ExecClient, error)
	Import(ctx context.Context, in *ImportRequest, opts ...grpc.CallOption) (*ImportResponse, error)
	Export(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error)
	Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error)
}

type directorClient struct {
	cc grpc.ClientConnInterface
}

func NewDirectorClient(cc grpc.ClientConnInterface) DirectorClient {
	return &directorClient{cc}
}

func (c *directorClient) Open(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenResponse, error) {
	out := new(OpenResponse)
	err := c.cc.Invoke(ctx, Director_Open_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *directorClient) Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (Director_ExecClient, error) {
	stream, err := c.cc.NewStream(ctx, &Director_ServiceDesc.Streams[0], Director_Exec_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &directorExecClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Director_ExecClient interface {
	Recv() (*ExecEvent, error)
	grpc.ClientStream
}

type directorExecClient struct {
	grpc.ClientStream
}

func (x *directorExecClient) Recv() (*ExecEvent, error) {
	m := new(ExecEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *directorClient) Import(ctx context.Context, in *ImportRequest, opts ...grpc.CallOption) (*ImportResponse, error) {
	out := new(ImportResponse)
	err := c.cc.Invoke(ctx, Director_Import_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *directorClient) Export(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error) {
	out := new(ExportResponse)
	err := c.cc.Invoke(ctx, Director_Export_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *directorClient) Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error) {
	out := new(CloseResponse)
	err := c.cc.Invoke(ctx, Director_Close_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DirectorServer is the server API for Director service.
// All implementations must embed UnimplementedDirectorServer
// for forward compatibility
type DirectorServer interface {
	Open(context.Context, *OpenRequest) (*OpenResponse, error)
	Exec(*ExecRequest, Director_ExecServer) error
	Import(context.Context, *ImportRequest) (*ImportResponse, error)
	Export(context.Context, *ExportRequest) (*ExportResponse, error)
	Close(context.Context, *CloseRequest) (*CloseResponse, error)
	mustEmbedUnimplementedDirectorServer()
}

// UnimplementedDirectorServer must be embedded to have forward compatible implementations.
type UnimplementedDirectorServer struct {
}

func (UnimplementedDirectorServer) Open(context.Context, *OpenRequest) (*OpenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Open not implemented")
}
func (UnimplementedDirectorServer) Exec(*ExecRequest, Director_ExecServer) error {
	return status.Errorf(codes.Unimplemented, "method Exec not implemented")
}
func (UnimplementedDirectorServer) Import(context.Context, *ImportRequest) (*ImportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Import not implemented")
}
func (UnimplementedDirectorServer) Export(context.Context, *ExportRequest) (*ExportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Export not implemented")
}
func (UnimplementedDirectorServer) Close(context.Context, *CloseRequest) (*CloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedDirectorServer) mustEmbedUnimplementedDirectorServer() {}

// UnsafeDirectorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DirectorServer will
// result in compilation errors.
type UnsafeDirectorServer interface {
	mustEmbedUnimplementedDirectorServer()
}

func RegisterDirectorServer(s grpc.ServiceRegistrar, srv DirectorServer) {
	s.RegisterService(&Director_ServiceDesc, srv)
}

func _Director_Open_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DirectorServer).Open(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Director_Open_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DirectorServer).Open(ctx, req.(*OpenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Director_Exec_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DirectorServer).Exec(m, &directorExecServer{stream})
}

type Director_ExecServer interface {
	Send(*ExecEvent) error
	grpc.ServerStream
}

type directorExecServer struct {
	grpc.ServerStream
}

func (x *directorExecServer) Send(m *ExecEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _Director_Import_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DirectorServer).Import(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Director_Import_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DirectorServer).Import(ctx, req.(*ImportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Director_Export_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DirectorServer).Export(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Director_Export_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DirectorServer).Export(ctx, req.(*ExportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Director_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DirectorServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Director_Close_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DirectorServer).Close(ctx, req.(*CloseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Director_ServiceDesc is the grpc.ServiceDesc for Director service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Director_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "director.knita.io.Director",
	HandlerType: (*DirectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Open",
			Handler:    _Director_Open_Handler,
		},
		{
			MethodName: "Import",
			Handler:    _Director_Import_Handler,
		},
		{
			MethodName: "Export",
			Handler:    _Director_Export_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Director_Close_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Exec",
			Handler:       _Director_Exec_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "director/v1/director.proto",
}
