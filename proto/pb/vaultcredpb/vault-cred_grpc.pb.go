// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: vault-cred.proto

package vaultcredpb

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
	VaultCred_GetCred_FullMethodName                       = "/vaultcredpb.VaultCred/GetCred"
	VaultCred_PutCred_FullMethodName                       = "/vaultcredpb.VaultCred/PutCred"
	VaultCred_DeleteCred_FullMethodName                    = "/vaultcredpb.VaultCred/DeleteCred"
	VaultCred_GetAppRoleToken_FullMethodName               = "/vaultcredpb.VaultCred/GetAppRoleToken"
	VaultCred_GetCredentialWithAppRoleToken_FullMethodName = "/vaultcredpb.VaultCred/GetCredentialWithAppRoleToken"
)

// VaultCredClient is the client API for VaultCred service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VaultCredClient interface {
	// single RPC for all credential types like user/password, certificates etc.
	// vault secret path  prepated based on <credentialType>/<credEntityName>/<credIdentifier>
	// for example, client-certs path will be certs/client/clientA
	// for example, cassandra root user path will be service-cred/cassandra/root
	// pass authentication token with service account token to authenticate & authorize the request with vault
	GetCred(ctx context.Context, in *GetCredRequest, opts ...grpc.CallOption) (*GetCredResponse, error)
	PutCred(ctx context.Context, in *PutCredRequest, opts ...grpc.CallOption) (*PutCredResponse, error)
	DeleteCred(ctx context.Context, in *DeleteCredRequest, opts ...grpc.CallOption) (*DeleteCredResponse, error)
	GetAppRoleToken(ctx context.Context, in *GetAppRoleTokenRequest, opts ...grpc.CallOption) (*GetAppRoleTokenResponse, error)
	GetCredentialWithAppRoleToken(ctx context.Context, in *GetCredentialWithAppRoleTokenRequest, opts ...grpc.CallOption) (*GetCredentialWithAppRoleTokenResponse, error)
}

type vaultCredClient struct {
	cc grpc.ClientConnInterface
}

func NewVaultCredClient(cc grpc.ClientConnInterface) VaultCredClient {
	return &vaultCredClient{cc}
}

func (c *vaultCredClient) GetCred(ctx context.Context, in *GetCredRequest, opts ...grpc.CallOption) (*GetCredResponse, error) {
	out := new(GetCredResponse)
	err := c.cc.Invoke(ctx, VaultCred_GetCred_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) PutCred(ctx context.Context, in *PutCredRequest, opts ...grpc.CallOption) (*PutCredResponse, error) {
	out := new(PutCredResponse)
	err := c.cc.Invoke(ctx, VaultCred_PutCred_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) DeleteCred(ctx context.Context, in *DeleteCredRequest, opts ...grpc.CallOption) (*DeleteCredResponse, error) {
	out := new(DeleteCredResponse)
	err := c.cc.Invoke(ctx, VaultCred_DeleteCred_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) GetAppRoleToken(ctx context.Context, in *GetAppRoleTokenRequest, opts ...grpc.CallOption) (*GetAppRoleTokenResponse, error) {
	out := new(GetAppRoleTokenResponse)
	err := c.cc.Invoke(ctx, VaultCred_GetAppRoleToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) GetCredentialWithAppRoleToken(ctx context.Context, in *GetCredentialWithAppRoleTokenRequest, opts ...grpc.CallOption) (*GetCredentialWithAppRoleTokenResponse, error) {
	out := new(GetCredentialWithAppRoleTokenResponse)
	err := c.cc.Invoke(ctx, VaultCred_GetCredentialWithAppRoleToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VaultCredServer is the server API for VaultCred service.
// All implementations must embed UnimplementedVaultCredServer
// for forward compatibility
type VaultCredServer interface {
	// single RPC for all credential types like user/password, certificates etc.
	// vault secret path  prepated based on <credentialType>/<credEntityName>/<credIdentifier>
	// for example, client-certs path will be certs/client/clientA
	// for example, cassandra root user path will be service-cred/cassandra/root
	// pass authentication token with service account token to authenticate & authorize the request with vault
	GetCred(context.Context, *GetCredRequest) (*GetCredResponse, error)
	PutCred(context.Context, *PutCredRequest) (*PutCredResponse, error)
	DeleteCred(context.Context, *DeleteCredRequest) (*DeleteCredResponse, error)
	GetAppRoleToken(context.Context, *GetAppRoleTokenRequest) (*GetAppRoleTokenResponse, error)
	GetCredentialWithAppRoleToken(context.Context, *GetCredentialWithAppRoleTokenRequest) (*GetCredentialWithAppRoleTokenResponse, error)
	mustEmbedUnimplementedVaultCredServer()
}

// UnimplementedVaultCredServer must be embedded to have forward compatible implementations.
type UnimplementedVaultCredServer struct {
}

func (UnimplementedVaultCredServer) GetCred(context.Context, *GetCredRequest) (*GetCredResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCred not implemented")
}
func (UnimplementedVaultCredServer) PutCred(context.Context, *PutCredRequest) (*PutCredResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutCred not implemented")
}
func (UnimplementedVaultCredServer) DeleteCred(context.Context, *DeleteCredRequest) (*DeleteCredResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCred not implemented")
}
func (UnimplementedVaultCredServer) GetAppRoleToken(context.Context, *GetAppRoleTokenRequest) (*GetAppRoleTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAppRoleToken not implemented")
}
func (UnimplementedVaultCredServer) GetCredentialWithAppRoleToken(context.Context, *GetCredentialWithAppRoleTokenRequest) (*GetCredentialWithAppRoleTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCredentialWithAppRoleToken not implemented")
}
func (UnimplementedVaultCredServer) mustEmbedUnimplementedVaultCredServer() {}

// UnsafeVaultCredServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VaultCredServer will
// result in compilation errors.
type UnsafeVaultCredServer interface {
	mustEmbedUnimplementedVaultCredServer()
}

func RegisterVaultCredServer(s grpc.ServiceRegistrar, srv VaultCredServer) {
	s.RegisterService(&VaultCred_ServiceDesc, srv)
}

func _VaultCred_GetCred_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCredRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).GetCred(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_GetCred_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).GetCred(ctx, req.(*GetCredRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_PutCred_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutCredRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).PutCred(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_PutCred_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).PutCred(ctx, req.(*PutCredRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_DeleteCred_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCredRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).DeleteCred(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_DeleteCred_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).DeleteCred(ctx, req.(*DeleteCredRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_GetAppRoleToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAppRoleTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).GetAppRoleToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_GetAppRoleToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).GetAppRoleToken(ctx, req.(*GetAppRoleTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_GetCredentialWithAppRoleToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCredentialWithAppRoleTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).GetCredentialWithAppRoleToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_GetCredentialWithAppRoleToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).GetCredentialWithAppRoleToken(ctx, req.(*GetCredentialWithAppRoleTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VaultCred_ServiceDesc is the grpc.ServiceDesc for VaultCred service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VaultCred_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vaultcredpb.VaultCred",
	HandlerType: (*VaultCredServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCred",
			Handler:    _VaultCred_GetCred_Handler,
		},
		{
			MethodName: "PutCred",
			Handler:    _VaultCred_PutCred_Handler,
		},
		{
			MethodName: "DeleteCred",
			Handler:    _VaultCred_DeleteCred_Handler,
		},
		{
			MethodName: "GetAppRoleToken",
			Handler:    _VaultCred_GetAppRoleToken_Handler,
		},
		{
			MethodName: "GetCredentialWithAppRoleToken",
			Handler:    _VaultCred_GetCredentialWithAppRoleToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vault-cred.proto",
}
