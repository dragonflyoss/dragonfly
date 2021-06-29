// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package dfdaemon

import (
	context "context"
	base "d7y.io/dragonfly/v2/internal/rpc/base"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// DaemonClient is the client API for Daemon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DaemonClient interface {
	// trigger client to download file
	Download(ctx context.Context, in *DownRequest, opts ...grpc.CallOption) (Daemon_DownloadClient, error)
	// get piece tasks from other peers
	GetPieceTasks(ctx context.Context, in *base.PieceTaskRequest, opts ...grpc.CallOption) (*base.PiecePacket, error)
	// check daemon health
	CheckHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type daemonClient struct {
	cc grpc.ClientConnInterface
}

func NewDaemonClient(cc grpc.ClientConnInterface) DaemonClient {
	return &daemonClient{cc}
}

func (c *daemonClient) Download(ctx context.Context, in *DownRequest, opts ...grpc.CallOption) (Daemon_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Daemon_serviceDesc.Streams[0], "/dfdaemon.Daemon/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &daemonDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Daemon_DownloadClient interface {
	Recv() (*DownResult, error)
	grpc.ClientStream
}

type daemonDownloadClient struct {
	grpc.ClientStream
}

func (x *daemonDownloadClient) Recv() (*DownResult, error) {
	m := new(DownResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *daemonClient) GetPieceTasks(ctx context.Context, in *base.PieceTaskRequest, opts ...grpc.CallOption) (*base.PiecePacket, error) {
	out := new(base.PiecePacket)
	err := c.cc.Invoke(ctx, "/dfdaemon.Daemon/GetPieceTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonClient) CheckHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/dfdaemon.Daemon/CheckHealth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaemonServer is the server API for Daemon service.
// All implementations must embed UnimplementedDaemonServer
// for forward compatibility
type DaemonServer interface {
	// trigger client to download file
	Download(*DownRequest, Daemon_DownloadServer) error
	// get piece tasks from other peers
	GetPieceTasks(context.Context, *base.PieceTaskRequest) (*base.PiecePacket, error)
	// check daemon health
	CheckHealth(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedDaemonServer()
}

// UnimplementedDaemonServer must be embedded to have forward compatible implementations.
type UnimplementedDaemonServer struct {
}

func (UnimplementedDaemonServer) Download(*DownRequest, Daemon_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedDaemonServer) GetPieceTasks(context.Context, *base.PieceTaskRequest) (*base.PiecePacket, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPieceTasks not implemented")
}
func (UnimplementedDaemonServer) CheckHealth(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHealth not implemented")
}
func (UnimplementedDaemonServer) mustEmbedUnimplementedDaemonServer() {}

// UnsafeDaemonServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DaemonServer will
// result in compilation errors.
type UnsafeDaemonServer interface {
	mustEmbedUnimplementedDaemonServer()
}

func RegisterDaemonServer(s grpc.ServiceRegistrar, srv DaemonServer) {
	s.RegisterService(&_Daemon_serviceDesc, srv)
}

func _Daemon_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DaemonServer).Download(m, &daemonDownloadServer{stream})
}

type Daemon_DownloadServer interface {
	Send(*DownResult) error
	grpc.ServerStream
}

type daemonDownloadServer struct {
	grpc.ServerStream
}

func (x *daemonDownloadServer) Send(m *DownResult) error {
	return x.ServerStream.SendMsg(m)
}

func _Daemon_GetPieceTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(base.PieceTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).GetPieceTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dfdaemon.Daemon/GetPieceTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).GetPieceTasks(ctx, req.(*base.PieceTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Daemon_CheckHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServer).CheckHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dfdaemon.Daemon/CheckHealth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServer).CheckHealth(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Daemon_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dfdaemon.Daemon",
	HandlerType: (*DaemonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPieceTasks",
			Handler:    _Daemon_GetPieceTasks_Handler,
		},
		{
			MethodName: "CheckHealth",
			Handler:    _Daemon_CheckHealth_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Download",
			Handler:       _Daemon_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/rpc/dfdaemon/dfdaemon.proto",
}
