// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: secretsmanager.proto

package secretsmanagerplugin

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

// RemoteSecretsManagerClient is the client API for RemoteSecretsManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemoteSecretsManagerClient interface {
	Get(ctx context.Context, in *SecretsGetRequest, opts ...grpc.CallOption) (*SecretsGetResponse, error)
	Set(ctx context.Context, in *SecretsSetRequest, opts ...grpc.CallOption) (*SecretsErrorResponse, error)
	Del(ctx context.Context, in *SecretsDelRequest, opts ...grpc.CallOption) (*SecretsErrorResponse, error)
	Keys(ctx context.Context, in *SecretsKeysRequest, opts ...grpc.CallOption) (*SecretsKeysResponse, error)
	Rename(ctx context.Context, in *SecretsRenameRequest, opts ...grpc.CallOption) (*SecretsErrorResponse, error)
}

type remoteSecretsManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoteSecretsManagerClient(cc grpc.ClientConnInterface) RemoteSecretsManagerClient {
	return &remoteSecretsManagerClient{cc}
}

func (c *remoteSecretsManagerClient) Get(ctx context.Context, in *SecretsGetRequest, opts ...grpc.CallOption) (*SecretsGetResponse, error) {
	out := new(SecretsGetResponse)
	err := c.cc.Invoke(ctx, "/secretsmanagerplugin.RemoteSecretsManager/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteSecretsManagerClient) Set(ctx context.Context, in *SecretsSetRequest, opts ...grpc.CallOption) (*SecretsErrorResponse, error) {
	out := new(SecretsErrorResponse)
	err := c.cc.Invoke(ctx, "/secretsmanagerplugin.RemoteSecretsManager/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteSecretsManagerClient) Del(ctx context.Context, in *SecretsDelRequest, opts ...grpc.CallOption) (*SecretsErrorResponse, error) {
	out := new(SecretsErrorResponse)
	err := c.cc.Invoke(ctx, "/secretsmanagerplugin.RemoteSecretsManager/Del", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteSecretsManagerClient) Keys(ctx context.Context, in *SecretsKeysRequest, opts ...grpc.CallOption) (*SecretsKeysResponse, error) {
	out := new(SecretsKeysResponse)
	err := c.cc.Invoke(ctx, "/secretsmanagerplugin.RemoteSecretsManager/Keys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remoteSecretsManagerClient) Rename(ctx context.Context, in *SecretsRenameRequest, opts ...grpc.CallOption) (*SecretsErrorResponse, error) {
	out := new(SecretsErrorResponse)
	err := c.cc.Invoke(ctx, "/secretsmanagerplugin.RemoteSecretsManager/Rename", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemoteSecretsManagerServer is the server API for RemoteSecretsManager service.
// All implementations must embed UnimplementedRemoteSecretsManagerServer
// for forward compatibility
type RemoteSecretsManagerServer interface {
	Get(context.Context, *SecretsGetRequest) (*SecretsGetResponse, error)
	Set(context.Context, *SecretsSetRequest) (*SecretsErrorResponse, error)
	Del(context.Context, *SecretsDelRequest) (*SecretsErrorResponse, error)
	Keys(context.Context, *SecretsKeysRequest) (*SecretsKeysResponse, error)
	Rename(context.Context, *SecretsRenameRequest) (*SecretsErrorResponse, error)
	mustEmbedUnimplementedRemoteSecretsManagerServer()
}

// UnimplementedRemoteSecretsManagerServer must be embedded to have forward compatible implementations.
type UnimplementedRemoteSecretsManagerServer struct {
}

func (UnimplementedRemoteSecretsManagerServer) Get(context.Context, *SecretsGetRequest) (*SecretsGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedRemoteSecretsManagerServer) Set(context.Context, *SecretsSetRequest) (*SecretsErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedRemoteSecretsManagerServer) Del(context.Context, *SecretsDelRequest) (*SecretsErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedRemoteSecretsManagerServer) Keys(context.Context, *SecretsKeysRequest) (*SecretsKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Keys not implemented")
}
func (UnimplementedRemoteSecretsManagerServer) Rename(context.Context, *SecretsRenameRequest) (*SecretsErrorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rename not implemented")
}
func (UnimplementedRemoteSecretsManagerServer) mustEmbedUnimplementedRemoteSecretsManagerServer() {}

// UnsafeRemoteSecretsManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemoteSecretsManagerServer will
// result in compilation errors.
type UnsafeRemoteSecretsManagerServer interface {
	mustEmbedUnimplementedRemoteSecretsManagerServer()
}

func RegisterRemoteSecretsManagerServer(s grpc.ServiceRegistrar, srv RemoteSecretsManagerServer) {
	s.RegisterService(&RemoteSecretsManager_ServiceDesc, srv)
}

func _RemoteSecretsManager_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretsGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteSecretsManagerServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/secretsmanagerplugin.RemoteSecretsManager/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteSecretsManagerServer).Get(ctx, req.(*SecretsGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoteSecretsManager_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretsSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteSecretsManagerServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/secretsmanagerplugin.RemoteSecretsManager/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteSecretsManagerServer).Set(ctx, req.(*SecretsSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoteSecretsManager_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretsDelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteSecretsManagerServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/secretsmanagerplugin.RemoteSecretsManager/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteSecretsManagerServer).Del(ctx, req.(*SecretsDelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoteSecretsManager_Keys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretsKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteSecretsManagerServer).Keys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/secretsmanagerplugin.RemoteSecretsManager/Keys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteSecretsManagerServer).Keys(ctx, req.(*SecretsKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoteSecretsManager_Rename_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretsRenameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoteSecretsManagerServer).Rename(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/secretsmanagerplugin.RemoteSecretsManager/Rename",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoteSecretsManagerServer).Rename(ctx, req.(*SecretsRenameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RemoteSecretsManager_ServiceDesc is the grpc.ServiceDesc for RemoteSecretsManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RemoteSecretsManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "secretsmanagerplugin.RemoteSecretsManager",
	HandlerType: (*RemoteSecretsManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _RemoteSecretsManager_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _RemoteSecretsManager_Set_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _RemoteSecretsManager_Del_Handler,
		},
		{
			MethodName: "Keys",
			Handler:    _RemoteSecretsManager_Keys_Handler,
		},
		{
			MethodName: "Rename",
			Handler:    _RemoteSecretsManager_Rename_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "secretsmanager.proto",
}