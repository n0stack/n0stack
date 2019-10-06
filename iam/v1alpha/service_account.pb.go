// Code generated by protoc-gen-go. DO NOT EDIT.
// source: n0stack/iam/v1alpha/service_account.proto

package piam

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ServiceAccount struct {
	// Name is a unique field.
	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Project string `protobuf:"bytes,2,opt,name=project,proto3" json:"project,omitempty"`
	// Annotations can store metadata used by the system for control.
	// In particular, implementation-dependent fields that can not be set as protobuf fields are targeted.
	// The control specified by n0stack may delete metadata specified by the user.
	Annotations map[string]string `protobuf:"bytes,3,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Labels stores user-defined metadata.
	// The n0stack system must not rewrite this value.
	Labels               map[string]string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	PublicKeys           map[string]string `protobuf:"bytes,33,rep,name=public_keys,json=publicKeys,proto3" json:"public_keys,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ServiceAccount) Reset()         { *m = ServiceAccount{} }
func (m *ServiceAccount) String() string { return proto.CompactTextString(m) }
func (*ServiceAccount) ProtoMessage()    {}
func (*ServiceAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ae0a5ab0582aab, []int{0}
}

func (m *ServiceAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceAccount.Unmarshal(m, b)
}
func (m *ServiceAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceAccount.Marshal(b, m, deterministic)
}
func (m *ServiceAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceAccount.Merge(m, src)
}
func (m *ServiceAccount) XXX_Size() int {
	return xxx_messageInfo_ServiceAccount.Size(m)
}
func (m *ServiceAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceAccount.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceAccount proto.InternalMessageInfo

func (m *ServiceAccount) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ServiceAccount) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *ServiceAccount) GetAnnotations() map[string]string {
	if m != nil {
		return m.Annotations
	}
	return nil
}

func (m *ServiceAccount) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *ServiceAccount) GetPublicKeys() map[string]string {
	if m != nil {
		return m.PublicKeys
	}
	return nil
}

type GetServiceAccountRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetServiceAccountRequest) Reset()         { *m = GetServiceAccountRequest{} }
func (m *GetServiceAccountRequest) String() string { return proto.CompactTextString(m) }
func (*GetServiceAccountRequest) ProtoMessage()    {}
func (*GetServiceAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ae0a5ab0582aab, []int{1}
}

func (m *GetServiceAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetServiceAccountRequest.Unmarshal(m, b)
}
func (m *GetServiceAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetServiceAccountRequest.Marshal(b, m, deterministic)
}
func (m *GetServiceAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetServiceAccountRequest.Merge(m, src)
}
func (m *GetServiceAccountRequest) XXX_Size() int {
	return xxx_messageInfo_GetServiceAccountRequest.Size(m)
}
func (m *GetServiceAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetServiceAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetServiceAccountRequest proto.InternalMessageInfo

func (m *GetServiceAccountRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreateServiceAccountRequest struct {
	ServiceAccount       *ServiceAccount `protobuf:"bytes,1,opt,name=service_account,json=serviceAccount,proto3" json:"service_account,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CreateServiceAccountRequest) Reset()         { *m = CreateServiceAccountRequest{} }
func (m *CreateServiceAccountRequest) String() string { return proto.CompactTextString(m) }
func (*CreateServiceAccountRequest) ProtoMessage()    {}
func (*CreateServiceAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ae0a5ab0582aab, []int{2}
}

func (m *CreateServiceAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateServiceAccountRequest.Unmarshal(m, b)
}
func (m *CreateServiceAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateServiceAccountRequest.Marshal(b, m, deterministic)
}
func (m *CreateServiceAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateServiceAccountRequest.Merge(m, src)
}
func (m *CreateServiceAccountRequest) XXX_Size() int {
	return xxx_messageInfo_CreateServiceAccountRequest.Size(m)
}
func (m *CreateServiceAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateServiceAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateServiceAccountRequest proto.InternalMessageInfo

func (m *CreateServiceAccountRequest) GetServiceAccount() *ServiceAccount {
	if m != nil {
		return m.ServiceAccount
	}
	return nil
}

type UpdateServiceAccountRequest struct {
	ServiceAccount       *ServiceAccount       `protobuf:"bytes,1,opt,name=service_account,json=serviceAccount,proto3" json:"service_account,omitempty"`
	UpdateMask           *field_mask.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UpdateServiceAccountRequest) Reset()         { *m = UpdateServiceAccountRequest{} }
func (m *UpdateServiceAccountRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateServiceAccountRequest) ProtoMessage()    {}
func (*UpdateServiceAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ae0a5ab0582aab, []int{3}
}

func (m *UpdateServiceAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateServiceAccountRequest.Unmarshal(m, b)
}
func (m *UpdateServiceAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateServiceAccountRequest.Marshal(b, m, deterministic)
}
func (m *UpdateServiceAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateServiceAccountRequest.Merge(m, src)
}
func (m *UpdateServiceAccountRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateServiceAccountRequest.Size(m)
}
func (m *UpdateServiceAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateServiceAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateServiceAccountRequest proto.InternalMessageInfo

func (m *UpdateServiceAccountRequest) GetServiceAccount() *ServiceAccount {
	if m != nil {
		return m.ServiceAccount
	}
	return nil
}

func (m *UpdateServiceAccountRequest) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type DeleteServiceAccountRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteServiceAccountRequest) Reset()         { *m = DeleteServiceAccountRequest{} }
func (m *DeleteServiceAccountRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteServiceAccountRequest) ProtoMessage()    {}
func (*DeleteServiceAccountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77ae0a5ab0582aab, []int{4}
}

func (m *DeleteServiceAccountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteServiceAccountRequest.Unmarshal(m, b)
}
func (m *DeleteServiceAccountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteServiceAccountRequest.Marshal(b, m, deterministic)
}
func (m *DeleteServiceAccountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteServiceAccountRequest.Merge(m, src)
}
func (m *DeleteServiceAccountRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteServiceAccountRequest.Size(m)
}
func (m *DeleteServiceAccountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteServiceAccountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteServiceAccountRequest proto.InternalMessageInfo

func (m *DeleteServiceAccountRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*ServiceAccount)(nil), "n0stack.iam.v1alpha.ServiceAccount")
	proto.RegisterMapType((map[string]string)(nil), "n0stack.iam.v1alpha.ServiceAccount.AnnotationsEntry")
	proto.RegisterMapType((map[string]string)(nil), "n0stack.iam.v1alpha.ServiceAccount.LabelsEntry")
	proto.RegisterMapType((map[string]string)(nil), "n0stack.iam.v1alpha.ServiceAccount.PublicKeysEntry")
	proto.RegisterType((*GetServiceAccountRequest)(nil), "n0stack.iam.v1alpha.GetServiceAccountRequest")
	proto.RegisterType((*CreateServiceAccountRequest)(nil), "n0stack.iam.v1alpha.CreateServiceAccountRequest")
	proto.RegisterType((*UpdateServiceAccountRequest)(nil), "n0stack.iam.v1alpha.UpdateServiceAccountRequest")
	proto.RegisterType((*DeleteServiceAccountRequest)(nil), "n0stack.iam.v1alpha.DeleteServiceAccountRequest")
}

func init() {
	proto.RegisterFile("n0stack/iam/v1alpha/service_account.proto", fileDescriptor_77ae0a5ab0582aab)
}

var fileDescriptor_77ae0a5ab0582aab = []byte{
	// 638 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xcb, 0x6e, 0x13, 0x31,
	0x14, 0xd5, 0x24, 0x7d, 0x88, 0x3b, 0x52, 0x5b, 0x4c, 0x41, 0xa3, 0x09, 0x8b, 0x30, 0x6c, 0x4a,
	0x45, 0xec, 0x36, 0x65, 0x51, 0x52, 0x81, 0x54, 0xa0, 0x74, 0xd1, 0x22, 0x41, 0x78, 0x2c, 0xd8,
	0x54, 0xce, 0xd4, 0x4d, 0xa7, 0xf3, 0x32, 0x63, 0x4f, 0x51, 0x84, 0xd8, 0xb0, 0xe0, 0x03, 0xca,
	0x8a, 0x3f, 0xe0, 0x33, 0xf8, 0x02, 0x36, 0xfc, 0x02, 0x7b, 0x7e, 0x01, 0x8d, 0xc7, 0x29, 0xcd,
	0x64, 0x12, 0x4d, 0x25, 0xc4, 0xca, 0xbe, 0xf6, 0x7d, 0x9c, 0x73, 0x7d, 0x6c, 0xc3, 0x9d, 0x68,
	0x4d, 0x48, 0xea, 0xfa, 0xc4, 0xa3, 0x21, 0x39, 0x5d, 0xa7, 0x01, 0x3f, 0xa6, 0x44, 0xb0, 0xe4,
	0xd4, 0x73, 0xd9, 0x01, 0x75, 0xdd, 0x38, 0x8d, 0x24, 0xe6, 0x49, 0x2c, 0x63, 0x74, 0x4d, 0xbb,
	0x62, 0x8f, 0x86, 0x58, 0xbb, 0xda, 0x37, 0xfb, 0x71, 0xdc, 0x0f, 0x18, 0xa1, 0xdc, 0x23, 0x34,
	0x8a, 0x62, 0x49, 0xa5, 0x17, 0x47, 0x22, 0x0f, 0xb1, 0x1b, 0x7a, 0x57, 0x59, 0xbd, 0xf4, 0x88,
	0xb0, 0x90, 0xcb, 0x81, 0xde, 0x6c, 0x16, 0x37, 0x8f, 0x3c, 0x16, 0x1c, 0x1e, 0x84, 0x54, 0xf8,
	0xda, 0xe3, 0xae, 0x1a, 0xdc, 0x56, 0x9f, 0x45, 0x2d, 0xf1, 0x9e, 0xf6, 0xfb, 0x2c, 0x21, 0x31,
	0x57, 0x05, 0xc6, 0x8b, 0x39, 0xbf, 0xeb, 0xb0, 0xf0, 0x32, 0x47, 0xbe, 0x9d, 0x03, 0x47, 0x08,
	0x66, 0x22, 0x1a, 0x32, 0xcb, 0x68, 0x1a, 0x2b, 0x57, 0xba, 0x6a, 0x8e, 0x2c, 0x98, 0xe7, 0x49,
	0x7c, 0xc2, 0x5c, 0x69, 0xd5, 0xd4, 0xf2, 0xd0, 0x44, 0x6f, 0xc0, 0xbc, 0x90, 0xd5, 0xaa, 0x37,
	0xeb, 0x2b, 0x66, 0xfb, 0x1e, 0x2e, 0xa1, 0x8d, 0x47, 0xeb, 0xe0, 0xed, 0xbf, 0x61, 0x3b, 0x91,
	0x4c, 0x06, 0xdd, 0x8b, 0x89, 0xd0, 0x2e, 0xcc, 0x05, 0xb4, 0xc7, 0x02, 0x61, 0xcd, 0xa8, 0x94,
	0xa4, 0x4a, 0xca, 0x7d, 0x15, 0x91, 0x67, 0xd3, 0xe1, 0xe8, 0x15, 0x98, 0x3c, 0xed, 0x05, 0x9e,
	0x7b, 0xe0, 0xb3, 0x81, 0xb0, 0x6e, 0xa9, 0x6c, 0x1b, 0x55, 0xb2, 0x3d, 0x57, 0x61, 0x7b, 0x6c,
	0xa0, 0x33, 0x02, 0x3f, 0x5f, 0xb0, 0x1f, 0xc2, 0x52, 0x11, 0x3f, 0x5a, 0x82, 0xba, 0xcf, 0x06,
	0xba, 0x6f, 0xd9, 0x14, 0x2d, 0xc3, 0xec, 0x29, 0x0d, 0x52, 0xa6, 0x9b, 0x96, 0x1b, 0x9d, 0xda,
	0xa6, 0x61, 0xdf, 0x07, 0xf3, 0x02, 0xd8, 0x4b, 0x85, 0x3e, 0x80, 0xc5, 0x02, 0xb2, 0xcb, 0x84,
	0x3b, 0x18, 0xac, 0x5d, 0x26, 0x47, 0xa9, 0x76, 0xd9, 0xbb, 0x94, 0x89, 0xd2, 0xa3, 0x77, 0x7c,
	0x68, 0x3c, 0x4e, 0x18, 0x95, 0xac, 0x3c, 0x64, 0x1f, 0x16, 0x0b, 0xca, 0x57, 0xd1, 0x66, 0xfb,
	0x76, 0x85, 0x16, 0x77, 0x17, 0xc4, 0x88, 0xed, 0x7c, 0x33, 0xa0, 0xf1, 0x9a, 0x1f, 0xfe, 0x9f,
	0x6a, 0x68, 0x0b, 0xcc, 0x54, 0x15, 0x53, 0xf7, 0x47, 0xb5, 0xca, 0x6c, 0xdb, 0x38, 0xbf, 0x62,
	0x78, 0x78, 0xc5, 0xf0, 0xd3, 0xec, 0x8a, 0x3d, 0xa3, 0xc2, 0xef, 0x42, 0xee, 0x9e, 0xcd, 0x9d,
	0x75, 0x68, 0x3c, 0x61, 0x01, 0x9b, 0x84, 0xb4, 0xa4, 0x95, 0xed, 0x1f, 0xb3, 0x70, 0x7d, 0xd4,
	0x5b, 0x5b, 0xe8, 0xab, 0x01, 0x57, 0xc7, 0x4e, 0x05, 0xb5, 0x4a, 0x49, 0x4d, 0x3a, 0x3d, 0xbb,
	0x4a, 0x0f, 0x9c, 0xb5, 0x4f, 0x3f, 0x7f, 0x7d, 0xa9, 0xad, 0xa2, 0x15, 0xf5, 0xfa, 0x4c, 0x79,
	0xc0, 0x04, 0xf9, 0x90, 0x81, 0xfe, 0x88, 0xbe, 0x1b, 0xb0, 0x5c, 0xa6, 0x00, 0xb4, 0x56, 0x5a,
	0x6f, 0x8a, 0x58, 0xaa, 0x21, 0x7c, 0xa1, 0x10, 0xee, 0x39, 0x9b, 0x15, 0x10, 0x16, 0x1f, 0x5d,
	0x85, 0xb8, 0x53, 0xd4, 0x88, 0xa2, 0x50, 0x26, 0xab, 0x09, 0x14, 0xa6, 0x28, 0xf0, 0x52, 0x14,
	0xda, 0xff, 0x90, 0xc2, 0x99, 0x01, 0xcb, 0x65, 0x7a, 0x9b, 0x40, 0x61, 0x8a, 0x34, 0xed, 0x1b,
	0x63, 0x0a, 0xdf, 0xc9, 0x7e, 0x98, 0xa1, 0x34, 0x56, 0x2b, 0x4b, 0xe3, 0xd1, 0x67, 0xe3, 0x6c,
	0x9b, 0xa2, 0x4d, 0x98, 0xd7, 0x10, 0x9c, 0xd6, 0xf9, 0x14, 0x39, 0xc7, 0x52, 0x72, 0xd1, 0x21,
	0xa4, 0xef, 0xc9, 0xe3, 0xb4, 0x87, 0xdd, 0x38, 0x24, 0xc3, 0x7f, 0x53, 0x8f, 0xab, 0x46, 0xad,
	0xbd, 0x44, 0x39, 0x0f, 0x3c, 0x57, 0xbd, 0xa3, 0xe4, 0x44, 0xc4, 0x51, 0x67, 0x6c, 0xe5, 0x6d,
	0x33, 0xf3, 0xc7, 0xd4, 0x25, 0x25, 0xff, 0xee, 0x16, 0xf7, 0x68, 0xd8, 0x9b, 0x53, 0x54, 0x36,
	0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0x0f, 0x8b, 0x02, 0xb2, 0x9a, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServiceAccountServiceClient is the client API for ServiceAccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceAccountServiceClient interface {
	GetServiceAccount(ctx context.Context, in *GetServiceAccountRequest, opts ...grpc.CallOption) (*ServiceAccount, error)
	CreateServiceAccount(ctx context.Context, in *CreateServiceAccountRequest, opts ...grpc.CallOption) (*ServiceAccount, error)
	UpdateServiceAccount(ctx context.Context, in *UpdateServiceAccountRequest, opts ...grpc.CallOption) (*ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, in *DeleteServiceAccountRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type serviceAccountServiceClient struct {
	cc *grpc.ClientConn
}

func NewServiceAccountServiceClient(cc *grpc.ClientConn) ServiceAccountServiceClient {
	return &serviceAccountServiceClient{cc}
}

func (c *serviceAccountServiceClient) GetServiceAccount(ctx context.Context, in *GetServiceAccountRequest, opts ...grpc.CallOption) (*ServiceAccount, error) {
	out := new(ServiceAccount)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ServiceAccountService/GetServiceAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceAccountServiceClient) CreateServiceAccount(ctx context.Context, in *CreateServiceAccountRequest, opts ...grpc.CallOption) (*ServiceAccount, error) {
	out := new(ServiceAccount)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ServiceAccountService/CreateServiceAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceAccountServiceClient) UpdateServiceAccount(ctx context.Context, in *UpdateServiceAccountRequest, opts ...grpc.CallOption) (*ServiceAccount, error) {
	out := new(ServiceAccount)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ServiceAccountService/UpdateServiceAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceAccountServiceClient) DeleteServiceAccount(ctx context.Context, in *DeleteServiceAccountRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ServiceAccountService/DeleteServiceAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceAccountServiceServer is the server API for ServiceAccountService service.
type ServiceAccountServiceServer interface {
	GetServiceAccount(context.Context, *GetServiceAccountRequest) (*ServiceAccount, error)
	CreateServiceAccount(context.Context, *CreateServiceAccountRequest) (*ServiceAccount, error)
	UpdateServiceAccount(context.Context, *UpdateServiceAccountRequest) (*ServiceAccount, error)
	DeleteServiceAccount(context.Context, *DeleteServiceAccountRequest) (*empty.Empty, error)
}

// UnimplementedServiceAccountServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceAccountServiceServer struct {
}

func (*UnimplementedServiceAccountServiceServer) GetServiceAccount(ctx context.Context, req *GetServiceAccountRequest) (*ServiceAccount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServiceAccount not implemented")
}
func (*UnimplementedServiceAccountServiceServer) CreateServiceAccount(ctx context.Context, req *CreateServiceAccountRequest) (*ServiceAccount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateServiceAccount not implemented")
}
func (*UnimplementedServiceAccountServiceServer) UpdateServiceAccount(ctx context.Context, req *UpdateServiceAccountRequest) (*ServiceAccount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateServiceAccount not implemented")
}
func (*UnimplementedServiceAccountServiceServer) DeleteServiceAccount(ctx context.Context, req *DeleteServiceAccountRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteServiceAccount not implemented")
}

func RegisterServiceAccountServiceServer(s *grpc.Server, srv ServiceAccountServiceServer) {
	s.RegisterService(&_ServiceAccountService_serviceDesc, srv)
}

func _ServiceAccountService_GetServiceAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServiceAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceAccountServiceServer).GetServiceAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ServiceAccountService/GetServiceAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceAccountServiceServer).GetServiceAccount(ctx, req.(*GetServiceAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceAccountService_CreateServiceAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServiceAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceAccountServiceServer).CreateServiceAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ServiceAccountService/CreateServiceAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceAccountServiceServer).CreateServiceAccount(ctx, req.(*CreateServiceAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceAccountService_UpdateServiceAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateServiceAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceAccountServiceServer).UpdateServiceAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ServiceAccountService/UpdateServiceAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceAccountServiceServer).UpdateServiceAccount(ctx, req.(*UpdateServiceAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServiceAccountService_DeleteServiceAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteServiceAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceAccountServiceServer).DeleteServiceAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ServiceAccountService/DeleteServiceAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceAccountServiceServer).DeleteServiceAccount(ctx, req.(*DeleteServiceAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServiceAccountService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "n0stack.iam.v1alpha.ServiceAccountService",
	HandlerType: (*ServiceAccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetServiceAccount",
			Handler:    _ServiceAccountService_GetServiceAccount_Handler,
		},
		{
			MethodName: "CreateServiceAccount",
			Handler:    _ServiceAccountService_CreateServiceAccount_Handler,
		},
		{
			MethodName: "UpdateServiceAccount",
			Handler:    _ServiceAccountService_UpdateServiceAccount_Handler,
		},
		{
			MethodName: "DeleteServiceAccount",
			Handler:    _ServiceAccountService_DeleteServiceAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "n0stack/iam/v1alpha/service_account.proto",
}
