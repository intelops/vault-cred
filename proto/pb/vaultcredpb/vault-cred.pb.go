// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: vault-cred.proto

package vaultcredpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetCredRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CredentialType string `protobuf:"bytes,1,opt,name=credentialType,proto3" json:"credentialType,omitempty"`
	CredEntityName string `protobuf:"bytes,2,opt,name=credEntityName,proto3" json:"credEntityName,omitempty"`
	CredIdentifier string `protobuf:"bytes,3,opt,name=credIdentifier,proto3" json:"credIdentifier,omitempty"`
}

func (x *GetCredRequest) Reset() {
	*x = GetCredRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCredRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCredRequest) ProtoMessage() {}

func (x *GetCredRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCredRequest.ProtoReflect.Descriptor instead.
func (*GetCredRequest) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{0}
}

func (x *GetCredRequest) GetCredentialType() string {
	if x != nil {
		return x.CredentialType
	}
	return ""
}

func (x *GetCredRequest) GetCredEntityName() string {
	if x != nil {
		return x.CredEntityName
	}
	return ""
}

func (x *GetCredRequest) GetCredIdentifier() string {
	if x != nil {
		return x.CredIdentifier
	}
	return ""
}

type GetCredResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// service-cred credential, for example: "userName": "iam-root", "password:: "hello"
	// client-cert credential, for example: "clientId": "intelops-user", "ca.crt": "...", "client.crt": "...", "client.key": "..."
	Credential map[string]string `protobuf:"bytes,1,rep,name=credential,proto3" json:"credential,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GetCredResponse) Reset() {
	*x = GetCredResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCredResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCredResponse) ProtoMessage() {}

func (x *GetCredResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCredResponse.ProtoReflect.Descriptor instead.
func (*GetCredResponse) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{1}
}

func (x *GetCredResponse) GetCredential() map[string]string {
	if x != nil {
		return x.Credential
	}
	return nil
}

type PutCredRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CredentialType string `protobuf:"bytes,1,opt,name=credentialType,proto3" json:"credentialType,omitempty"`
	CredEntityName string `protobuf:"bytes,2,opt,name=credEntityName,proto3" json:"credEntityName,omitempty"`
	CredIdentifier string `protobuf:"bytes,3,opt,name=credIdentifier,proto3" json:"credIdentifier,omitempty"`
	// service-cred credential, for example: "userName": "iam-root", "password:: "hello"
	// client-cert credential, for example: "clientId": "intelops-user", "ca.crt": "...", "client.crt": "...", "client.key": "..."
	Credential map[string]string `protobuf:"bytes,6,rep,name=credential,proto3" json:"credential,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *PutCredRequest) Reset() {
	*x = PutCredRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutCredRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutCredRequest) ProtoMessage() {}

func (x *PutCredRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutCredRequest.ProtoReflect.Descriptor instead.
func (*PutCredRequest) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{2}
}

func (x *PutCredRequest) GetCredentialType() string {
	if x != nil {
		return x.CredentialType
	}
	return ""
}

func (x *PutCredRequest) GetCredEntityName() string {
	if x != nil {
		return x.CredEntityName
	}
	return ""
}

func (x *PutCredRequest) GetCredIdentifier() string {
	if x != nil {
		return x.CredIdentifier
	}
	return ""
}

func (x *PutCredRequest) GetCredential() map[string]string {
	if x != nil {
		return x.Credential
	}
	return nil
}

type PutCredResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PutCredResponse) Reset() {
	*x = PutCredResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutCredResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutCredResponse) ProtoMessage() {}

func (x *PutCredResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutCredResponse.ProtoReflect.Descriptor instead.
func (*PutCredResponse) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{3}
}

type DeleteCredRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CredentialType string `protobuf:"bytes,1,opt,name=credentialType,proto3" json:"credentialType,omitempty"`
	CredEntityName string `protobuf:"bytes,2,opt,name=credEntityName,proto3" json:"credEntityName,omitempty"`
	CredIdentifier string `protobuf:"bytes,3,opt,name=credIdentifier,proto3" json:"credIdentifier,omitempty"`
}

func (x *DeleteCredRequest) Reset() {
	*x = DeleteCredRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCredRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCredRequest) ProtoMessage() {}

func (x *DeleteCredRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCredRequest.ProtoReflect.Descriptor instead.
func (*DeleteCredRequest) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteCredRequest) GetCredentialType() string {
	if x != nil {
		return x.CredentialType
	}
	return ""
}

func (x *DeleteCredRequest) GetCredEntityName() string {
	if x != nil {
		return x.CredEntityName
	}
	return ""
}

func (x *DeleteCredRequest) GetCredIdentifier() string {
	if x != nil {
		return x.CredIdentifier
	}
	return ""
}

type DeleteCredResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteCredResponse) Reset() {
	*x = DeleteCredResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCredResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCredResponse) ProtoMessage() {}

func (x *DeleteCredResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCredResponse.ProtoReflect.Descriptor instead.
func (*DeleteCredResponse) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{5}
}

type GetAppRoleTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppRoleName    string `protobuf:"bytes,1,opt,name=appRoleName,proto3" json:"appRoleName,omitempty"`
	CredentialPath string `protobuf:"bytes,2,opt,name=credentialPath,proto3" json:"credentialPath,omitempty"`
}

func (x *GetAppRoleTokenRequest) Reset() {
	*x = GetAppRoleTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAppRoleTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAppRoleTokenRequest) ProtoMessage() {}

func (x *GetAppRoleTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAppRoleTokenRequest.ProtoReflect.Descriptor instead.
func (*GetAppRoleTokenRequest) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{6}
}

func (x *GetAppRoleTokenRequest) GetAppRoleName() string {
	if x != nil {
		return x.AppRoleName
	}
	return ""
}

func (x *GetAppRoleTokenRequest) GetCredentialPath() string {
	if x != nil {
		return x.CredentialPath
	}
	return ""
}

type GetAppRoleTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *GetAppRoleTokenResponse) Reset() {
	*x = GetAppRoleTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAppRoleTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAppRoleTokenResponse) ProtoMessage() {}

func (x *GetAppRoleTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAppRoleTokenResponse.ProtoReflect.Descriptor instead.
func (*GetAppRoleTokenResponse) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{7}
}

func (x *GetAppRoleTokenResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type GetCredentialWithAppRoleTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token          string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	CredentialPath string `protobuf:"bytes,2,opt,name=credentialPath,proto3" json:"credentialPath,omitempty"`
}

func (x *GetCredentialWithAppRoleTokenRequest) Reset() {
	*x = GetCredentialWithAppRoleTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCredentialWithAppRoleTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCredentialWithAppRoleTokenRequest) ProtoMessage() {}

func (x *GetCredentialWithAppRoleTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCredentialWithAppRoleTokenRequest.ProtoReflect.Descriptor instead.
func (*GetCredentialWithAppRoleTokenRequest) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{8}
}

func (x *GetCredentialWithAppRoleTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *GetCredentialWithAppRoleTokenRequest) GetCredentialPath() string {
	if x != nil {
		return x.CredentialPath
	}
	return ""
}

type GetCredentialWithAppRoleTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Credential map[string]string `protobuf:"bytes,1,rep,name=credential,proto3" json:"credential,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GetCredentialWithAppRoleTokenResponse) Reset() {
	*x = GetCredentialWithAppRoleTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vault_cred_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCredentialWithAppRoleTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCredentialWithAppRoleTokenResponse) ProtoMessage() {}

func (x *GetCredentialWithAppRoleTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_vault_cred_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCredentialWithAppRoleTokenResponse.ProtoReflect.Descriptor instead.
func (*GetCredentialWithAppRoleTokenResponse) Descriptor() ([]byte, []int) {
	return file_vault_cred_proto_rawDescGZIP(), []int{9}
}

func (x *GetCredentialWithAppRoleTokenResponse) GetCredential() map[string]string {
	if x != nil {
		return x.Credential
	}
	return nil
}

var File_vault_cred_proto protoreflect.FileDescriptor

var file_vault_cred_proto_rawDesc = []byte{
	0x0a, 0x10, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x2d, 0x63, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x22,
	0x88, 0x01, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x72,
	0x65, 0x64, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22, 0x9e, 0x01, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x1a, 0x3d, 0x0a, 0x0f,
	0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x94, 0x02, 0x0a, 0x0e,
	0x50, 0x75, 0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26,
	0x0a, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x63, 0x72, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x26,
	0x0a, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x4b, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x61, 0x6c, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x76, 0x61, 0x75,
	0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x74, 0x43, 0x72, 0x65, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x61, 0x6c, 0x1a, 0x3d, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61,
	0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x11, 0x0a, 0x0f, 0x50, 0x75, 0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x8b, 0x01, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x63,
	0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65,
	0x64, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63,
	0x72, 0x65, 0x64, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x22, 0x14, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x72, 0x65,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x62, 0x0a, 0x16, 0x47, 0x65, 0x74,
	0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x70, 0x70, 0x52, 0x6f, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x61, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63,
	0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x22, 0x2f, 0x0a,
	0x17, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x64,
	0x0a, 0x24, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x57,
	0x69, 0x74, 0x68, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x26, 0x0a, 0x0e,
	0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x50, 0x61, 0x74, 0x68, 0x22, 0xca, 0x01, 0x0a, 0x25, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x57, 0x69, 0x74, 0x68, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x42, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x57, 0x69,
	0x74, 0x68, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61,
	0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x1a, 0x3d, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x32, 0xd7, 0x03, 0x0a, 0x09, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x43, 0x72, 0x65, 0x64, 0x12,
	0x46, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x12, 0x1b, 0x2e, 0x76, 0x61, 0x75,
	0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63,
	0x72, 0x65, 0x64, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x07, 0x50, 0x75, 0x74, 0x43, 0x72,
	0x65, 0x64, 0x12, 0x1b, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62,
	0x2e, 0x50, 0x75, 0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x2e, 0x50, 0x75,
	0x74, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4f, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x72, 0x65, 0x64, 0x12, 0x1e, 0x2e,
	0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x72, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x5e, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x23, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70,
	0x62, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74,
	0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x88, 0x01, 0x0a, 0x1d, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x61, 0x6c, 0x57, 0x69, 0x74, 0x68, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x31, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x57, 0x69,
	0x74, 0x68, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65,
	0x64, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61,
	0x6c, 0x57, 0x69, 0x74, 0x68, 0x41, 0x70, 0x70, 0x52, 0x6f, 0x6c, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0e, 0x5a, 0x0c, 0x2f,
	0x76, 0x61, 0x75, 0x6c, 0x74, 0x63, 0x72, 0x65, 0x64, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_vault_cred_proto_rawDescOnce sync.Once
	file_vault_cred_proto_rawDescData = file_vault_cred_proto_rawDesc
)

func file_vault_cred_proto_rawDescGZIP() []byte {
	file_vault_cred_proto_rawDescOnce.Do(func() {
		file_vault_cred_proto_rawDescData = protoimpl.X.CompressGZIP(file_vault_cred_proto_rawDescData)
	})
	return file_vault_cred_proto_rawDescData
}

var file_vault_cred_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_vault_cred_proto_goTypes = []interface{}{
	(*GetCredRequest)(nil),                        // 0: vaultcredpb.GetCredRequest
	(*GetCredResponse)(nil),                       // 1: vaultcredpb.GetCredResponse
	(*PutCredRequest)(nil),                        // 2: vaultcredpb.PutCredRequest
	(*PutCredResponse)(nil),                       // 3: vaultcredpb.PutCredResponse
	(*DeleteCredRequest)(nil),                     // 4: vaultcredpb.DeleteCredRequest
	(*DeleteCredResponse)(nil),                    // 5: vaultcredpb.DeleteCredResponse
	(*GetAppRoleTokenRequest)(nil),                // 6: vaultcredpb.GetAppRoleTokenRequest
	(*GetAppRoleTokenResponse)(nil),               // 7: vaultcredpb.GetAppRoleTokenResponse
	(*GetCredentialWithAppRoleTokenRequest)(nil),  // 8: vaultcredpb.GetCredentialWithAppRoleTokenRequest
	(*GetCredentialWithAppRoleTokenResponse)(nil), // 9: vaultcredpb.GetCredentialWithAppRoleTokenResponse
	nil, // 10: vaultcredpb.GetCredResponse.CredentialEntry
	nil, // 11: vaultcredpb.PutCredRequest.CredentialEntry
	nil, // 12: vaultcredpb.GetCredentialWithAppRoleTokenResponse.CredentialEntry
}
var file_vault_cred_proto_depIdxs = []int32{
	10, // 0: vaultcredpb.GetCredResponse.credential:type_name -> vaultcredpb.GetCredResponse.CredentialEntry
	11, // 1: vaultcredpb.PutCredRequest.credential:type_name -> vaultcredpb.PutCredRequest.CredentialEntry
	12, // 2: vaultcredpb.GetCredentialWithAppRoleTokenResponse.credential:type_name -> vaultcredpb.GetCredentialWithAppRoleTokenResponse.CredentialEntry
	0,  // 3: vaultcredpb.VaultCred.GetCred:input_type -> vaultcredpb.GetCredRequest
	2,  // 4: vaultcredpb.VaultCred.PutCred:input_type -> vaultcredpb.PutCredRequest
	4,  // 5: vaultcredpb.VaultCred.DeleteCred:input_type -> vaultcredpb.DeleteCredRequest
	6,  // 6: vaultcredpb.VaultCred.GetAppRoleToken:input_type -> vaultcredpb.GetAppRoleTokenRequest
	8,  // 7: vaultcredpb.VaultCred.GetCredentialWithAppRoleToken:input_type -> vaultcredpb.GetCredentialWithAppRoleTokenRequest
	1,  // 8: vaultcredpb.VaultCred.GetCred:output_type -> vaultcredpb.GetCredResponse
	3,  // 9: vaultcredpb.VaultCred.PutCred:output_type -> vaultcredpb.PutCredResponse
	5,  // 10: vaultcredpb.VaultCred.DeleteCred:output_type -> vaultcredpb.DeleteCredResponse
	7,  // 11: vaultcredpb.VaultCred.GetAppRoleToken:output_type -> vaultcredpb.GetAppRoleTokenResponse
	9,  // 12: vaultcredpb.VaultCred.GetCredentialWithAppRoleToken:output_type -> vaultcredpb.GetCredentialWithAppRoleTokenResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_vault_cred_proto_init() }
func file_vault_cred_proto_init() {
	if File_vault_cred_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vault_cred_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCredRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCredResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutCredRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutCredResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCredRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCredResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAppRoleTokenRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAppRoleTokenResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCredentialWithAppRoleTokenRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vault_cred_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCredentialWithAppRoleTokenResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vault_cred_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vault_cred_proto_goTypes,
		DependencyIndexes: file_vault_cred_proto_depIdxs,
		MessageInfos:      file_vault_cred_proto_msgTypes,
	}.Build()
	File_vault_cred_proto = out.File
	file_vault_cred_proto_rawDesc = nil
	file_vault_cred_proto_goTypes = nil
	file_vault_cred_proto_depIdxs = nil
}
