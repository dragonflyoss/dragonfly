// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package manager

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ManagerClient is the client API for Manager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ManagerClient interface {
	// get scheduler server list, using scene as follows:
	// 1. scheduler servers are not exist in local config
	//
	// 2. connection is fail for all servers from config,
	// so need retry one times to get latest servers
	//
	// 3. manager actively triggers fresh
	GetSchedulers(ctx context.Context, in *GetSchedulersRequest, opts ...grpc.CallOption) (*SchedulerNodes, error)
	// keep alive for cdn or scheduler
	KeepAlive(ctx context.Context, in *KeepAliveRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// get cluster config for cdn or scheduler(client) from manager
	GetClusterConfig(ctx context.Context, in *GetClusterConfigRequest, opts ...grpc.CallOption) (*ClusterConfig, error)
}

type managerClient struct {
	cc grpc.ClientConnInterface
}

func NewManagerClient(cc grpc.ClientConnInterface) ManagerClient {
	return &managerClient{cc}
}

func (c *managerClient) GetSchedulers(ctx context.Context, in *GetSchedulersRequest, opts ...grpc.CallOption) (*SchedulerNodes, error) {
	out := new(SchedulerNodes)
	err := c.cc.Invoke(ctx, "/manager.Manager/GetSchedulers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerClient) KeepAlive(ctx context.Context, in *KeepAliveRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/manager.Manager/KeepAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managerClient) GetClusterConfig(ctx context.Context, in *GetClusterConfigRequest, opts ...grpc.CallOption) (*ClusterConfig, error) {
	out := new(ClusterConfig)
	err := c.cc.Invoke(ctx, "/manager.Manager/GetClusterConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManagerServer is the server API for Manager service.
// All implementations must embed UnimplementedManagerServer
// for forward compatibility
type ManagerServer interface {
	// get scheduler server list, using scene as follows:
	// 1. scheduler servers are not exist in local config
	//
	// 2. connection is fail for all servers from config,
	// so need retry one times to get latest servers
	//
	// 3. manager actively triggers fresh
	GetSchedulers(context.Context, *GetSchedulersRequest) (*SchedulerNodes, error)
	// keep alive for cdn or scheduler
	KeepAlive(context.Context, *KeepAliveRequest) (*emptypb.Empty, error)
	// get cluster config for cdn or scheduler(client) from manager
	GetClusterConfig(context.Context, *GetClusterConfigRequest) (*ClusterConfig, error)
	mustEmbedUnimplementedManagerServer()
}

// UnimplementedManagerServer must be embedded to have forward compatible implementations.
type UnimplementedManagerServer struct {
}

func (UnimplementedManagerServer) GetSchedulers(context.Context, *GetSchedulersRequest) (*SchedulerNodes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSchedulers not implemented")
}
func (UnimplementedManagerServer) KeepAlive(context.Context, *KeepAliveRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}
func (UnimplementedManagerServer) GetClusterConfig(context.Context, *GetClusterConfigRequest) (*ClusterConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterConfig not implemented")
}
func (UnimplementedManagerServer) mustEmbedUnimplementedManagerServer() {}

// UnsafeManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ManagerServer will
// result in compilation errors.
type UnsafeManagerServer interface {
	mustEmbedUnimplementedManagerServer()
}

func RegisterManagerServer(s grpc.ServiceRegistrar, srv ManagerServer) {
	s.RegisterService(&_Manager_serviceDesc, srv)
}

func _Manager_GetSchedulers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchedulersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServer).GetSchedulers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.Manager/GetSchedulers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServer).GetSchedulers(ctx, req.(*GetSchedulersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Manager_KeepAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeepAliveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServer).KeepAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.Manager/KeepAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServer).KeepAlive(ctx, req.(*KeepAliveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Manager_GetClusterConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagerServer).GetClusterConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.Manager/GetClusterConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagerServer).GetClusterConfig(ctx, req.(*GetClusterConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Manager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "manager.Manager",
	HandlerType: (*ManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSchedulers",
			Handler:    _Manager_GetSchedulers_Handler,
		},
		{
			MethodName: "KeepAlive",
			Handler:    _Manager_KeepAlive_Handler,
		},
		{
			MethodName: "GetClusterConfig",
			Handler:    _Manager_GetClusterConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/rpc/manager/manager.proto",
}
