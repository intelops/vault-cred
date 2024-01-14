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
	VaultCred_GetCredential_FullMethodName                 = "/vaultcredpb.VaultCred/GetCredential"
	VaultCred_PutCredential_FullMethodName                 = "/vaultcredpb.VaultCred/PutCredential"
	VaultCred_DeleteCredential_FullMethodName              = "/vaultcredpb.VaultCred/DeleteCredential"
	VaultCred_CreateAppRoleToken_FullMethodName            = "/vaultcredpb.VaultCred/CreateAppRoleToken"
	VaultCred_DeleteAppRole_FullMethodName                 = "/vaultcredpb.VaultCred/DeleteAppRole"
	VaultCred_GetCredentialWithAppRoleToken_FullMethodName = "/vaultcredpb.VaultCred/GetCredentialWithAppRoleToken"
	VaultCred_ConfigureClusterK8SAuth_FullMethodName       = "/vaultcredpb.VaultCred/ConfigureClusterK8SAuth"
	VaultCred_CreateK8SAuthRole_FullMethodName             = "/vaultcredpb.VaultCred/CreateK8SAuthRole"
	VaultCred_UpdateK8SAuthRole_FullMethodName             = "/vaultcredpb.VaultCred/UpdateK8SAuthRole"
	VaultCred_DeleteK8SAuthRole_FullMethodName             = "/vaultcredpb.VaultCred/DeleteK8SAuthRole"
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
	GetCredential(ctx context.Context, in *GetCredentialRequest, opts ...grpc.CallOption) (*GetCredentialResponse, error)
	PutCredential(ctx context.Context, in *PutCredentialRequest, opts ...grpc.CallOption) (*PutCredentialResponse, error)
	DeleteCredential(ctx context.Context, in *DeleteCredentialRequest, opts ...grpc.CallOption) (*DeleteCredentialResponse, error)
	CreateAppRoleToken(ctx context.Context, in *CreateAppRoleTokenRequest, opts ...grpc.CallOption) (*CreateAppRoleTokenResponse, error)
	DeleteAppRole(ctx context.Context, in *DeleteAppRoleRequest, opts ...grpc.CallOption) (*DeleteAppRoleResponse, error)
	GetCredentialWithAppRoleToken(ctx context.Context, in *GetCredentialWithAppRoleTokenRequest, opts ...grpc.CallOption) (*GetCredentialWithAppRoleTokenResponse, error)
	ConfigureClusterK8SAuth(ctx context.Context, in *ConfigureClusterK8SAuthRequest, opts ...grpc.CallOption) (*ConfigureClusterK8SAuthResponse, error)
	CreateK8SAuthRole(ctx context.Context, in *CreateK8SAuthRoleRequest, opts ...grpc.CallOption) (*CreateK8SAuthRoleResponse, error)
	UpdateK8SAuthRole(ctx context.Context, in *UpdateK8SAuthRoleRequest, opts ...grpc.CallOption) (*UpdateK8SAuthRoleResponse, error)
	DeleteK8SAuthRole(ctx context.Context, in *DeleteK8SAuthRoleRequest, opts ...grpc.CallOption) (*DeleteK8SAuthRoleResponse, error)
}

type vaultCredClient struct {
	cc grpc.ClientConnInterface
}

func NewVaultCredClient(cc grpc.ClientConnInterface) VaultCredClient {
	return &vaultCredClient{cc}
}

func (c *vaultCredClient) GetCredential(ctx context.Context, in *GetCredentialRequest, opts ...grpc.CallOption) (*GetCredentialResponse, error) {
	out := new(GetCredentialResponse)
	err := c.cc.Invoke(ctx, VaultCred_GetCredential_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) PutCredential(ctx context.Context, in *PutCredentialRequest, opts ...grpc.CallOption) (*PutCredentialResponse, error) {
	out := new(PutCredentialResponse)
	err := c.cc.Invoke(ctx, VaultCred_PutCredential_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) DeleteCredential(ctx context.Context, in *DeleteCredentialRequest, opts ...grpc.CallOption) (*DeleteCredentialResponse, error) {
	out := new(DeleteCredentialResponse)
	err := c.cc.Invoke(ctx, VaultCred_DeleteCredential_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) CreateAppRoleToken(ctx context.Context, in *CreateAppRoleTokenRequest, opts ...grpc.CallOption) (*CreateAppRoleTokenResponse, error) {
	out := new(CreateAppRoleTokenResponse)
	err := c.cc.Invoke(ctx, VaultCred_CreateAppRoleToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) DeleteAppRole(ctx context.Context, in *DeleteAppRoleRequest, opts ...grpc.CallOption) (*DeleteAppRoleResponse, error) {
	out := new(DeleteAppRoleResponse)
	err := c.cc.Invoke(ctx, VaultCred_DeleteAppRole_FullMethodName, in, out, opts...)
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

func (c *vaultCredClient) ConfigureClusterK8SAuth(ctx context.Context, in *ConfigureClusterK8SAuthRequest, opts ...grpc.CallOption) (*ConfigureClusterK8SAuthResponse, error) {
	out := new(ConfigureClusterK8SAuthResponse)
	err := c.cc.Invoke(ctx, VaultCred_ConfigureClusterK8SAuth_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) CreateK8SAuthRole(ctx context.Context, in *CreateK8SAuthRoleRequest, opts ...grpc.CallOption) (*CreateK8SAuthRoleResponse, error) {
	out := new(CreateK8SAuthRoleResponse)
	err := c.cc.Invoke(ctx, VaultCred_CreateK8SAuthRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) UpdateK8SAuthRole(ctx context.Context, in *UpdateK8SAuthRoleRequest, opts ...grpc.CallOption) (*UpdateK8SAuthRoleResponse, error) {
	out := new(UpdateK8SAuthRoleResponse)
	err := c.cc.Invoke(ctx, VaultCred_UpdateK8SAuthRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultCredClient) DeleteK8SAuthRole(ctx context.Context, in *DeleteK8SAuthRoleRequest, opts ...grpc.CallOption) (*DeleteK8SAuthRoleResponse, error) {
	out := new(DeleteK8SAuthRoleResponse)
	err := c.cc.Invoke(ctx, VaultCred_DeleteK8SAuthRole_FullMethodName, in, out, opts...)
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
	GetCredential(context.Context, *GetCredentialRequest) (*GetCredentialResponse, error)
	PutCredential(context.Context, *PutCredentialRequest) (*PutCredentialResponse, error)
	DeleteCredential(context.Context, *DeleteCredentialRequest) (*DeleteCredentialResponse, error)
	CreateAppRoleToken(context.Context, *CreateAppRoleTokenRequest) (*CreateAppRoleTokenResponse, error)
	DeleteAppRole(context.Context, *DeleteAppRoleRequest) (*DeleteAppRoleResponse, error)
	GetCredentialWithAppRoleToken(context.Context, *GetCredentialWithAppRoleTokenRequest) (*GetCredentialWithAppRoleTokenResponse, error)
	ConfigureClusterK8SAuth(context.Context, *ConfigureClusterK8SAuthRequest) (*ConfigureClusterK8SAuthResponse, error)
	CreateK8SAuthRole(context.Context, *CreateK8SAuthRoleRequest) (*CreateK8SAuthRoleResponse, error)
	UpdateK8SAuthRole(context.Context, *UpdateK8SAuthRoleRequest) (*UpdateK8SAuthRoleResponse, error)
	DeleteK8SAuthRole(context.Context, *DeleteK8SAuthRoleRequest) (*DeleteK8SAuthRoleResponse, error)
	mustEmbedUnimplementedVaultCredServer()
}

// UnimplementedVaultCredServer must be embedded to have forward compatible implementations.
type UnimplementedVaultCredServer struct {
}

func (UnimplementedVaultCredServer) GetCredential(context.Context, *GetCredentialRequest) (*GetCredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCredential not implemented")
}
func (UnimplementedVaultCredServer) PutCredential(context.Context, *PutCredentialRequest) (*PutCredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutCredential not implemented")
}
func (UnimplementedVaultCredServer) DeleteCredential(context.Context, *DeleteCredentialRequest) (*DeleteCredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCredential not implemented")
}
func (UnimplementedVaultCredServer) CreateAppRoleToken(context.Context, *CreateAppRoleTokenRequest) (*CreateAppRoleTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAppRoleToken not implemented")
}
func (UnimplementedVaultCredServer) DeleteAppRole(context.Context, *DeleteAppRoleRequest) (*DeleteAppRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAppRole not implemented")
}
func (UnimplementedVaultCredServer) GetCredentialWithAppRoleToken(context.Context, *GetCredentialWithAppRoleTokenRequest) (*GetCredentialWithAppRoleTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCredentialWithAppRoleToken not implemented")
}
func (UnimplementedVaultCredServer) ConfigureClusterK8SAuth(context.Context, *ConfigureClusterK8SAuthRequest) (*ConfigureClusterK8SAuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureClusterK8SAuth not implemented")
}
func (UnimplementedVaultCredServer) CreateK8SAuthRole(context.Context, *CreateK8SAuthRoleRequest) (*CreateK8SAuthRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateK8SAuthRole not implemented")
}
func (UnimplementedVaultCredServer) UpdateK8SAuthRole(context.Context, *UpdateK8SAuthRoleRequest) (*UpdateK8SAuthRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateK8SAuthRole not implemented")
}
func (UnimplementedVaultCredServer) DeleteK8SAuthRole(context.Context, *DeleteK8SAuthRoleRequest) (*DeleteK8SAuthRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteK8SAuthRole not implemented")
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

func _VaultCred_GetCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).GetCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_GetCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).GetCredential(ctx, req.(*GetCredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_PutCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutCredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).PutCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_PutCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).PutCredential(ctx, req.(*PutCredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_DeleteCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).DeleteCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_DeleteCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).DeleteCredential(ctx, req.(*DeleteCredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_CreateAppRoleToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAppRoleTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).CreateAppRoleToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_CreateAppRoleToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).CreateAppRoleToken(ctx, req.(*CreateAppRoleTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_DeleteAppRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAppRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).DeleteAppRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_DeleteAppRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).DeleteAppRole(ctx, req.(*DeleteAppRoleRequest))
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

func _VaultCred_ConfigureClusterK8SAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureClusterK8SAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).ConfigureClusterK8SAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_ConfigureClusterK8SAuth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).ConfigureClusterK8SAuth(ctx, req.(*ConfigureClusterK8SAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_CreateK8SAuthRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateK8SAuthRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).CreateK8SAuthRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_CreateK8SAuthRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).CreateK8SAuthRole(ctx, req.(*CreateK8SAuthRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_UpdateK8SAuthRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateK8SAuthRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).UpdateK8SAuthRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_UpdateK8SAuthRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).UpdateK8SAuthRole(ctx, req.(*UpdateK8SAuthRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultCred_DeleteK8SAuthRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteK8SAuthRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultCredServer).DeleteK8SAuthRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultCred_DeleteK8SAuthRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultCredServer).DeleteK8SAuthRole(ctx, req.(*DeleteK8SAuthRoleRequest))
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
			MethodName: "GetCredential",
			Handler:    _VaultCred_GetCredential_Handler,
		},
		{
			MethodName: "PutCredential",
			Handler:    _VaultCred_PutCredential_Handler,
		},
		{
			MethodName: "DeleteCredential",
			Handler:    _VaultCred_DeleteCredential_Handler,
		},
		{
			MethodName: "CreateAppRoleToken",
			Handler:    _VaultCred_CreateAppRoleToken_Handler,
		},
		{
			MethodName: "DeleteAppRole",
			Handler:    _VaultCred_DeleteAppRole_Handler,
		},
		{
			MethodName: "GetCredentialWithAppRoleToken",
			Handler:    _VaultCred_GetCredentialWithAppRoleToken_Handler,
		},
		{
			MethodName: "ConfigureClusterK8SAuth",
			Handler:    _VaultCred_ConfigureClusterK8SAuth_Handler,
		},
		{
			MethodName: "CreateK8SAuthRole",
			Handler:    _VaultCred_CreateK8SAuthRole_Handler,
		},
		{
			MethodName: "UpdateK8SAuthRole",
			Handler:    _VaultCred_UpdateK8SAuthRole_Handler,
		},
		{
			MethodName: "DeleteK8SAuthRole",
			Handler:    _VaultCred_DeleteK8SAuthRole_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vault-cred.proto",
}
