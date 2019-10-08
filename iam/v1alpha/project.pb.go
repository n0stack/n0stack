// Code generated by protoc-gen-go. DO NOT EDIT.
// source: n0stack/iam/v1alpha/project.proto

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

type ProjectMembership int32

const (
	ProjectMembership_PROJECT_MEMBERSHIP_UNSPECIFIED ProjectMembership = 0
	// Owners have all of permissions.
	ProjectMembership_OWNER ProjectMembership = 1
	// Members have only assined permissions by Roles.
	ProjectMembership_MEMBER ProjectMembership = 2
)

var ProjectMembership_name = map[int32]string{
	0: "PROJECT_MEMBERSHIP_UNSPECIFIED",
	1: "OWNER",
	2: "MEMBER",
}

var ProjectMembership_value = map[string]int32{
	"PROJECT_MEMBERSHIP_UNSPECIFIED": 0,
	"OWNER":                          1,
	"MEMBER":                         2,
}

func (x ProjectMembership) String() string {
	return proto.EnumName(ProjectMembership_name, int32(x))
}

func (ProjectMembership) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{0}
}

type Project struct {
	// Name is a unique field.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Annotations can store metadata used by the system for control.
	// In particular, implementation-dependent fields that can not be set as protobuf fields are targeted.
	// The control specified by n0stack may delete metadata specified by the user.
	Annotations map[string]string `protobuf:"bytes,3,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Labels stores user-defined metadata.
	// The n0stack system must not rewrite this value.
	Labels               map[string]string            `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	DisplayName          string                       `protobuf:"bytes,9,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	Membership           map[string]ProjectMembership `protobuf:"bytes,32,rep,name=membership,proto3" json:"membership,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3,enum=n0stack.iam.v1alpha.ProjectMembership"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *Project) Reset()         { *m = Project{} }
func (m *Project) String() string { return proto.CompactTextString(m) }
func (*Project) ProtoMessage()    {}
func (*Project) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{0}
}

func (m *Project) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Project.Unmarshal(m, b)
}
func (m *Project) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Project.Marshal(b, m, deterministic)
}
func (m *Project) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Project.Merge(m, src)
}
func (m *Project) XXX_Size() int {
	return xxx_messageInfo_Project.Size(m)
}
func (m *Project) XXX_DiscardUnknown() {
	xxx_messageInfo_Project.DiscardUnknown(m)
}

var xxx_messageInfo_Project proto.InternalMessageInfo

func (m *Project) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Project) GetAnnotations() map[string]string {
	if m != nil {
		return m.Annotations
	}
	return nil
}

func (m *Project) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Project) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *Project) GetMembership() map[string]ProjectMembership {
	if m != nil {
		return m.Membership
	}
	return nil
}

type ListProjectsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListProjectsRequest) Reset()         { *m = ListProjectsRequest{} }
func (m *ListProjectsRequest) String() string { return proto.CompactTextString(m) }
func (*ListProjectsRequest) ProtoMessage()    {}
func (*ListProjectsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{1}
}

func (m *ListProjectsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListProjectsRequest.Unmarshal(m, b)
}
func (m *ListProjectsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListProjectsRequest.Marshal(b, m, deterministic)
}
func (m *ListProjectsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListProjectsRequest.Merge(m, src)
}
func (m *ListProjectsRequest) XXX_Size() int {
	return xxx_messageInfo_ListProjectsRequest.Size(m)
}
func (m *ListProjectsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListProjectsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListProjectsRequest proto.InternalMessageInfo

type ListProjectsResponse struct {
	Projects             []*Project `protobuf:"bytes,1,rep,name=projects,proto3" json:"projects,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ListProjectsResponse) Reset()         { *m = ListProjectsResponse{} }
func (m *ListProjectsResponse) String() string { return proto.CompactTextString(m) }
func (*ListProjectsResponse) ProtoMessage()    {}
func (*ListProjectsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{2}
}

func (m *ListProjectsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListProjectsResponse.Unmarshal(m, b)
}
func (m *ListProjectsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListProjectsResponse.Marshal(b, m, deterministic)
}
func (m *ListProjectsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListProjectsResponse.Merge(m, src)
}
func (m *ListProjectsResponse) XXX_Size() int {
	return xxx_messageInfo_ListProjectsResponse.Size(m)
}
func (m *ListProjectsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListProjectsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListProjectsResponse proto.InternalMessageInfo

func (m *ListProjectsResponse) GetProjects() []*Project {
	if m != nil {
		return m.Projects
	}
	return nil
}

type GetProjectRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetProjectRequest) Reset()         { *m = GetProjectRequest{} }
func (m *GetProjectRequest) String() string { return proto.CompactTextString(m) }
func (*GetProjectRequest) ProtoMessage()    {}
func (*GetProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{3}
}

func (m *GetProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetProjectRequest.Unmarshal(m, b)
}
func (m *GetProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetProjectRequest.Marshal(b, m, deterministic)
}
func (m *GetProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetProjectRequest.Merge(m, src)
}
func (m *GetProjectRequest) XXX_Size() int {
	return xxx_messageInfo_GetProjectRequest.Size(m)
}
func (m *GetProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetProjectRequest proto.InternalMessageInfo

func (m *GetProjectRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreateProjectRequest struct {
	Project              *Project `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateProjectRequest) Reset()         { *m = CreateProjectRequest{} }
func (m *CreateProjectRequest) String() string { return proto.CompactTextString(m) }
func (*CreateProjectRequest) ProtoMessage()    {}
func (*CreateProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{4}
}

func (m *CreateProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateProjectRequest.Unmarshal(m, b)
}
func (m *CreateProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateProjectRequest.Marshal(b, m, deterministic)
}
func (m *CreateProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateProjectRequest.Merge(m, src)
}
func (m *CreateProjectRequest) XXX_Size() int {
	return xxx_messageInfo_CreateProjectRequest.Size(m)
}
func (m *CreateProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateProjectRequest proto.InternalMessageInfo

func (m *CreateProjectRequest) GetProject() *Project {
	if m != nil {
		return m.Project
	}
	return nil
}

type UpdateProjectRequest struct {
	Project              *Project              `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	UpdateMask           *field_mask.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UpdateProjectRequest) Reset()         { *m = UpdateProjectRequest{} }
func (m *UpdateProjectRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateProjectRequest) ProtoMessage()    {}
func (*UpdateProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{5}
}

func (m *UpdateProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateProjectRequest.Unmarshal(m, b)
}
func (m *UpdateProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateProjectRequest.Marshal(b, m, deterministic)
}
func (m *UpdateProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateProjectRequest.Merge(m, src)
}
func (m *UpdateProjectRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateProjectRequest.Size(m)
}
func (m *UpdateProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateProjectRequest proto.InternalMessageInfo

func (m *UpdateProjectRequest) GetProject() *Project {
	if m != nil {
		return m.Project
	}
	return nil
}

func (m *UpdateProjectRequest) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

type DeleteProjectRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteProjectRequest) Reset()         { *m = DeleteProjectRequest{} }
func (m *DeleteProjectRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteProjectRequest) ProtoMessage()    {}
func (*DeleteProjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{6}
}

func (m *DeleteProjectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteProjectRequest.Unmarshal(m, b)
}
func (m *DeleteProjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteProjectRequest.Marshal(b, m, deterministic)
}
func (m *DeleteProjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteProjectRequest.Merge(m, src)
}
func (m *DeleteProjectRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteProjectRequest.Size(m)
}
func (m *DeleteProjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteProjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteProjectRequest proto.InternalMessageInfo

func (m *DeleteProjectRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type AddProjectMembershipRequest struct {
	ProjectName          string            `protobuf:"bytes,1,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
	UserName             string            `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Membership           ProjectMembership `protobuf:"varint,3,opt,name=membership,proto3,enum=n0stack.iam.v1alpha.ProjectMembership" json:"membership,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AddProjectMembershipRequest) Reset()         { *m = AddProjectMembershipRequest{} }
func (m *AddProjectMembershipRequest) String() string { return proto.CompactTextString(m) }
func (*AddProjectMembershipRequest) ProtoMessage()    {}
func (*AddProjectMembershipRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{7}
}

func (m *AddProjectMembershipRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddProjectMembershipRequest.Unmarshal(m, b)
}
func (m *AddProjectMembershipRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddProjectMembershipRequest.Marshal(b, m, deterministic)
}
func (m *AddProjectMembershipRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddProjectMembershipRequest.Merge(m, src)
}
func (m *AddProjectMembershipRequest) XXX_Size() int {
	return xxx_messageInfo_AddProjectMembershipRequest.Size(m)
}
func (m *AddProjectMembershipRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddProjectMembershipRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddProjectMembershipRequest proto.InternalMessageInfo

func (m *AddProjectMembershipRequest) GetProjectName() string {
	if m != nil {
		return m.ProjectName
	}
	return ""
}

func (m *AddProjectMembershipRequest) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *AddProjectMembershipRequest) GetMembership() ProjectMembership {
	if m != nil {
		return m.Membership
	}
	return ProjectMembership_PROJECT_MEMBERSHIP_UNSPECIFIED
}

type DeleteProjectMembershipRequest struct {
	ProjectName          string   `protobuf:"bytes,1,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
	UserName             string   `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteProjectMembershipRequest) Reset()         { *m = DeleteProjectMembershipRequest{} }
func (m *DeleteProjectMembershipRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteProjectMembershipRequest) ProtoMessage()    {}
func (*DeleteProjectMembershipRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_165db4196db99c73, []int{8}
}

func (m *DeleteProjectMembershipRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteProjectMembershipRequest.Unmarshal(m, b)
}
func (m *DeleteProjectMembershipRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteProjectMembershipRequest.Marshal(b, m, deterministic)
}
func (m *DeleteProjectMembershipRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteProjectMembershipRequest.Merge(m, src)
}
func (m *DeleteProjectMembershipRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteProjectMembershipRequest.Size(m)
}
func (m *DeleteProjectMembershipRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteProjectMembershipRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteProjectMembershipRequest proto.InternalMessageInfo

func (m *DeleteProjectMembershipRequest) GetProjectName() string {
	if m != nil {
		return m.ProjectName
	}
	return ""
}

func (m *DeleteProjectMembershipRequest) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func init() {
	proto.RegisterEnum("n0stack.iam.v1alpha.ProjectMembership", ProjectMembership_name, ProjectMembership_value)
	proto.RegisterType((*Project)(nil), "n0stack.iam.v1alpha.Project")
	proto.RegisterMapType((map[string]string)(nil), "n0stack.iam.v1alpha.Project.AnnotationsEntry")
	proto.RegisterMapType((map[string]string)(nil), "n0stack.iam.v1alpha.Project.LabelsEntry")
	proto.RegisterMapType((map[string]ProjectMembership)(nil), "n0stack.iam.v1alpha.Project.MembershipEntry")
	proto.RegisterType((*ListProjectsRequest)(nil), "n0stack.iam.v1alpha.ListProjectsRequest")
	proto.RegisterType((*ListProjectsResponse)(nil), "n0stack.iam.v1alpha.ListProjectsResponse")
	proto.RegisterType((*GetProjectRequest)(nil), "n0stack.iam.v1alpha.GetProjectRequest")
	proto.RegisterType((*CreateProjectRequest)(nil), "n0stack.iam.v1alpha.CreateProjectRequest")
	proto.RegisterType((*UpdateProjectRequest)(nil), "n0stack.iam.v1alpha.UpdateProjectRequest")
	proto.RegisterType((*DeleteProjectRequest)(nil), "n0stack.iam.v1alpha.DeleteProjectRequest")
	proto.RegisterType((*AddProjectMembershipRequest)(nil), "n0stack.iam.v1alpha.AddProjectMembershipRequest")
	proto.RegisterType((*DeleteProjectMembershipRequest)(nil), "n0stack.iam.v1alpha.DeleteProjectMembershipRequest")
}

func init() { proto.RegisterFile("n0stack/iam/v1alpha/project.proto", fileDescriptor_165db4196db99c73) }

var fileDescriptor_165db4196db99c73 = []byte{
	// 875 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0x4f, 0x6f, 0x1a, 0x47,
	0x14, 0xef, 0x82, 0x63, 0x9b, 0x87, 0x93, 0x92, 0x09, 0x6d, 0xe9, 0x12, 0x45, 0x30, 0x87, 0x96,
	0xa2, 0xb0, 0x9b, 0x62, 0xa9, 0x72, 0x48, 0xd5, 0x96, 0xd8, 0xeb, 0xd4, 0x91, 0xc1, 0x68, 0xdd,
	0xa8, 0x52, 0x2f, 0x74, 0x80, 0x09, 0x6c, 0xd8, 0x7f, 0xdd, 0x59, 0x5c, 0xa1, 0xc4, 0x97, 0x56,
	0xca, 0xa5, 0xb7, 0xe6, 0x33, 0xf4, 0xda, 0x43, 0xa5, 0x7e, 0x92, 0x7e, 0x85, 0x7e, 0x90, 0x6a,
	0x67, 0x07, 0x58, 0x60, 0x81, 0xb8, 0xf1, 0x69, 0x67, 0x67, 0x7e, 0xef, 0xbd, 0xdf, 0xbc, 0x79,
	0xbf, 0x37, 0x03, 0x45, 0xfb, 0x01, 0xf3, 0x49, 0x77, 0xa8, 0x1a, 0xc4, 0x52, 0x2f, 0x3e, 0x27,
	0xa6, 0x3b, 0x20, 0xaa, 0xeb, 0x39, 0x2f, 0x68, 0xd7, 0x57, 0x5c, 0xcf, 0xf1, 0x1d, 0x74, 0x47,
	0x40, 0x14, 0x83, 0x58, 0x8a, 0x80, 0xc8, 0x77, 0xfb, 0x8e, 0xd3, 0x37, 0xa9, 0x4a, 0x5c, 0x43,
	0x25, 0xb6, 0xed, 0xf8, 0xc4, 0x37, 0x1c, 0x9b, 0x85, 0x26, 0x72, 0x5e, 0xac, 0xf2, 0xbf, 0xce,
	0xe8, 0xb9, 0x4a, 0x2d, 0xd7, 0x1f, 0x8b, 0xc5, 0xc2, 0xe2, 0xe2, 0x73, 0x83, 0x9a, 0xbd, 0xb6,
	0x45, 0xd8, 0x50, 0x20, 0xee, 0xf3, 0x4f, 0xb7, 0xd2, 0xa7, 0x76, 0x85, 0xfd, 0x4c, 0xfa, 0x7d,
	0xea, 0xa9, 0x8e, 0xcb, 0x03, 0x2c, 0x07, 0xc3, 0xbf, 0x6e, 0xc1, 0x4e, 0x2b, 0x64, 0x8c, 0x10,
	0x6c, 0xd9, 0xc4, 0xa2, 0x39, 0xa9, 0x20, 0x95, 0x52, 0x3a, 0x1f, 0xa3, 0x33, 0x48, 0x47, 0x8c,
	0x72, 0xc9, 0x42, 0xb2, 0x94, 0xae, 0x56, 0x94, 0x98, 0x5d, 0x29, 0xc2, 0x8d, 0x52, 0x9f, 0xe1,
	0x35, 0xdb, 0xf7, 0xc6, 0x7a, 0xd4, 0x03, 0xfa, 0x06, 0xb6, 0x4d, 0xd2, 0xa1, 0x26, 0xcb, 0x6d,
	0x71, 0x5f, 0xa5, 0xb5, 0xbe, 0x4e, 0x39, 0x34, 0x74, 0x23, 0xec, 0x50, 0x11, 0xf6, 0x7a, 0x06,
	0x73, 0x4d, 0x32, 0x6e, 0x73, 0xba, 0x29, 0x4e, 0x37, 0x2d, 0xe6, 0x9a, 0x01, 0xeb, 0x53, 0x00,
	0x8b, 0x5a, 0x1d, 0xea, 0xb1, 0x81, 0xe1, 0xe6, 0x0a, 0x3c, 0xd0, 0xfd, 0xb5, 0x81, 0x1a, 0x53,
	0x78, 0x18, 0x2c, 0x62, 0x2f, 0x7f, 0x05, 0x99, 0xc5, 0x3d, 0xa1, 0x0c, 0x24, 0x87, 0x74, 0x2c,
	0x52, 0x15, 0x0c, 0x51, 0x16, 0x6e, 0x5c, 0x10, 0x73, 0x44, 0x73, 0x09, 0x3e, 0x17, 0xfe, 0xd4,
	0x12, 0x07, 0x92, 0xfc, 0x10, 0xd2, 0x91, 0x7d, 0x5c, 0xc9, 0x94, 0xc2, 0xfb, 0x0b, 0xcc, 0x62,
	0xcc, 0xbf, 0x8c, 0x9a, 0xdf, 0xaa, 0x7e, 0xb2, 0x6e, 0xa3, 0x33, 0x6f, 0x91, 0x30, 0xf8, 0x03,
	0xb8, 0x73, 0x6a, 0x30, 0x5f, 0x60, 0x98, 0x4e, 0x7f, 0x1a, 0x51, 0xe6, 0xe3, 0x16, 0x64, 0xe7,
	0xa7, 0x99, 0xeb, 0xd8, 0x8c, 0xa2, 0x03, 0xd8, 0x15, 0x55, 0xce, 0x72, 0x12, 0x4f, 0xee, 0xdd,
	0x75, 0x31, 0xf5, 0x29, 0x1a, 0x7f, 0x0a, 0xb7, 0x9f, 0xd0, 0x89, 0x43, 0x11, 0x26, 0xae, 0xee,
	0x70, 0x13, 0xb2, 0x87, 0x1e, 0x25, 0x3e, 0x5d, 0xc0, 0x7e, 0x01, 0x3b, 0xc2, 0x19, 0x87, 0x6f,
	0x8a, 0x3c, 0x01, 0xe3, 0xdf, 0x24, 0xc8, 0x3e, 0x73, 0x7b, 0xd7, 0xe6, 0x10, 0x3d, 0x82, 0xf4,
	0x88, 0xfb, 0xe3, 0xda, 0xe3, 0xa9, 0x4f, 0x57, 0x65, 0x25, 0x94, 0xa7, 0x32, 0x91, 0xa7, 0x72,
	0x1c, 0xc8, 0xb3, 0x41, 0xd8, 0x50, 0x87, 0x10, 0x1e, 0x8c, 0x71, 0x19, 0xb2, 0x47, 0xd4, 0xa4,
	0x4b, 0x64, 0xe2, 0x32, 0xf1, 0x87, 0x04, 0xf9, 0x7a, 0xaf, 0xb7, 0x7c, 0x7e, 0xc2, 0xa6, 0x08,
	0x7b, 0x82, 0x53, 0x3b, 0x62, 0x9b, 0x16, 0x73, 0x5c, 0x0e, 0x79, 0x48, 0x8d, 0x18, 0xf5, 0xc2,
	0xf5, 0xb0, 0xc6, 0x76, 0x83, 0x09, 0xbe, 0x78, 0x3c, 0xa7, 0x95, 0xe4, 0x95, 0x4a, 0x28, 0x62,
	0x89, 0x7f, 0x84, 0x7b, 0x73, 0x7b, 0xba, 0x76, 0xa6, 0xe5, 0x26, 0xdc, 0x5e, 0xf2, 0x8d, 0x30,
	0xdc, 0x6b, 0xe9, 0x67, 0x4f, 0xb5, 0xc3, 0xef, 0xda, 0x0d, 0xad, 0xf1, 0x58, 0xd3, 0xcf, 0xbf,
	0x3d, 0x69, 0xb5, 0x9f, 0x35, 0xcf, 0x5b, 0xda, 0xe1, 0xc9, 0xf1, 0x89, 0x76, 0x94, 0x79, 0x0f,
	0xa5, 0xe0, 0xc6, 0xd9, 0xf7, 0x4d, 0x4d, 0xcf, 0x48, 0x08, 0x60, 0x3b, 0x84, 0x65, 0x12, 0xd5,
	0x3f, 0x77, 0xe1, 0x96, 0x70, 0x78, 0x4e, 0xbd, 0x0b, 0xa3, 0x4b, 0xd1, 0x6b, 0x09, 0xf6, 0xa2,
	0x25, 0x8f, 0xe2, 0xdb, 0x53, 0x8c, 0x58, 0xe4, 0xcf, 0xde, 0x02, 0x19, 0xea, 0x07, 0x17, 0x7f,
	0xf9, 0xe7, 0xdf, 0x37, 0x89, 0x3c, 0xfa, 0x98, 0xdf, 0x00, 0x31, 0x97, 0x07, 0x43, 0xaf, 0x00,
	0x66, 0x42, 0x41, 0xf1, 0xe7, 0xb1, 0xa4, 0x24, 0x79, 0x6d, 0xed, 0xe2, 0x12, 0x0f, 0x8b, 0x51,
	0x61, 0x65, 0x58, 0xf5, 0x65, 0x90, 0xf8, 0x4b, 0xf4, 0x46, 0x82, 0x9b, 0x73, 0xf2, 0x43, 0xf1,
	0xbb, 0x8b, 0x93, 0xe8, 0x06, 0x12, 0x0f, 0x39, 0x89, 0x7d, 0x5c, 0x5a, 0x43, 0x62, 0x72, 0x85,
	0x72, 0x32, 0xb5, 0xa9, 0xe4, 0x02, 0x56, 0x73, 0x1a, 0x5e, 0xc1, 0x2a, 0x4e, 0xe7, 0x6f, 0xc7,
	0x4a, 0xfe, 0x1f, 0xac, 0x5e, 0xc1, 0xcd, 0xb9, 0xba, 0x5f, 0x41, 0x2a, 0x4e, 0xef, 0xf2, 0x87,
	0x4b, 0xfd, 0x42, 0x0b, 0xee, 0xfa, 0xc9, 0x49, 0x95, 0x37, 0x9f, 0xd4, 0x5f, 0x12, 0x64, 0xe3,
	0xba, 0x03, 0x7a, 0x10, 0xcb, 0x62, 0x4d, 0x23, 0xd9, 0x90, 0xa1, 0xa7, 0x9c, 0xd2, 0x11, 0xfe,
	0x7a, 0x73, 0x86, 0xb8, 0x7a, 0x2f, 0xd5, 0x59, 0x6f, 0x50, 0x5f, 0x4e, 0x45, 0x7d, 0x59, 0x93,
	0xca, 0xe8, 0x6f, 0x09, 0x3e, 0x5a, 0xd1, 0x2b, 0xd0, 0xfe, 0xe6, 0xec, 0x5d, 0x95, 0xfa, 0x13,
	0x4e, 0xbd, 0x5e, 0x7e, 0x57, 0xea, 0x8f, 0x5f, 0x4b, 0xbf, 0xd7, 0x09, 0x3a, 0x80, 0x1d, 0x11,
	0x0d, 0x57, 0xa6, 0x43, 0x84, 0x07, 0xbe, 0xef, 0xb2, 0x9a, 0xaa, 0xf6, 0x0d, 0x7f, 0x30, 0xea,
	0x28, 0x5d, 0xc7, 0x52, 0x27, 0xcf, 0x43, 0xf1, 0x2d, 0x4b, 0x89, 0x6a, 0x86, 0xb8, 0xae, 0x69,
	0x74, 0xf9, 0x53, 0x42, 0x7d, 0xc1, 0x1c, 0xbb, 0xb6, 0x34, 0xf3, 0x43, 0x21, 0xc0, 0x2b, 0xa4,
	0xab, 0xc6, 0x3c, 0x2f, 0x1f, 0xb9, 0x06, 0xb1, 0x3a, 0xdb, 0xbc, 0x5e, 0xf6, 0xff, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0xca, 0x2d, 0xf3, 0x6a, 0x81, 0x0a, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProjectServiceClient is the client API for ProjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProjectServiceClient interface {
	// あるユーザーがメンバーであるプロジェクトすべてを返す
	// errors:
	//   NotFound: どのプロジェクトにも属していない場合
	ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error)
	// Summary: プロジェクトの詳細を取得する
	// errors:
	//   NotFound: memberじゃない場合
	//   Unauthorized: ログインしていない場合
	GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*Project, error)
	// ログインしているユーザーがオーナーとなるプロジェクトを作成する
	// errors:
	//   Unauthorized: ログインしていない場合
	CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error)
	UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*Project, error)
	DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	AddProjectMembership(ctx context.Context, in *AddProjectMembershipRequest, opts ...grpc.CallOption) (*Project, error)
	DeleteProjectMembership(ctx context.Context, in *DeleteProjectMembershipRequest, opts ...grpc.CallOption) (*Project, error)
}

type projectServiceClient struct {
	cc *grpc.ClientConn
}

func NewProjectServiceClient(cc *grpc.ClientConn) ProjectServiceClient {
	return &projectServiceClient{cc}
}

func (c *projectServiceClient) ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error) {
	out := new(ListProjectsResponse)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ProjectService/ListProjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ProjectService/GetProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ProjectService/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ProjectService/UpdateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ProjectService/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) AddProjectMembership(ctx context.Context, in *AddProjectMembershipRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ProjectService/AddProjectMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) DeleteProjectMembership(ctx context.Context, in *DeleteProjectMembershipRequest, opts ...grpc.CallOption) (*Project, error) {
	out := new(Project)
	err := c.cc.Invoke(ctx, "/n0stack.iam.v1alpha.ProjectService/DeleteProjectMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectServiceServer is the server API for ProjectService service.
type ProjectServiceServer interface {
	// あるユーザーがメンバーであるプロジェクトすべてを返す
	// errors:
	//   NotFound: どのプロジェクトにも属していない場合
	ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error)
	// Summary: プロジェクトの詳細を取得する
	// errors:
	//   NotFound: memberじゃない場合
	//   Unauthorized: ログインしていない場合
	GetProject(context.Context, *GetProjectRequest) (*Project, error)
	// ログインしているユーザーがオーナーとなるプロジェクトを作成する
	// errors:
	//   Unauthorized: ログインしていない場合
	CreateProject(context.Context, *CreateProjectRequest) (*Project, error)
	UpdateProject(context.Context, *UpdateProjectRequest) (*Project, error)
	DeleteProject(context.Context, *DeleteProjectRequest) (*empty.Empty, error)
	AddProjectMembership(context.Context, *AddProjectMembershipRequest) (*Project, error)
	DeleteProjectMembership(context.Context, *DeleteProjectMembershipRequest) (*Project, error)
}

// UnimplementedProjectServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProjectServiceServer struct {
}

func (*UnimplementedProjectServiceServer) ListProjects(ctx context.Context, req *ListProjectsRequest) (*ListProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjects not implemented")
}
func (*UnimplementedProjectServiceServer) GetProject(ctx context.Context, req *GetProjectRequest) (*Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}
func (*UnimplementedProjectServiceServer) CreateProject(ctx context.Context, req *CreateProjectRequest) (*Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (*UnimplementedProjectServiceServer) UpdateProject(ctx context.Context, req *UpdateProjectRequest) (*Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProject not implemented")
}
func (*UnimplementedProjectServiceServer) DeleteProject(ctx context.Context, req *DeleteProjectRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
func (*UnimplementedProjectServiceServer) AddProjectMembership(ctx context.Context, req *AddProjectMembershipRequest) (*Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProjectMembership not implemented")
}
func (*UnimplementedProjectServiceServer) DeleteProjectMembership(ctx context.Context, req *DeleteProjectMembershipRequest) (*Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProjectMembership not implemented")
}

func RegisterProjectServiceServer(s *grpc.Server, srv ProjectServiceServer) {
	s.RegisterService(&_ProjectService_serviceDesc, srv)
}

func _ProjectService_ListProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).ListProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ProjectService/ListProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).ListProjects(ctx, req.(*ListProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ProjectService/GetProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetProject(ctx, req.(*GetProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ProjectService/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).CreateProject(ctx, req.(*CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).UpdateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ProjectService/UpdateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).UpdateProject(ctx, req.(*UpdateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ProjectService/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).DeleteProject(ctx, req.(*DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_AddProjectMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddProjectMembershipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).AddProjectMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ProjectService/AddProjectMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).AddProjectMembership(ctx, req.(*AddProjectMembershipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_DeleteProjectMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProjectMembershipRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).DeleteProjectMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.iam.v1alpha.ProjectService/DeleteProjectMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).DeleteProjectMembership(ctx, req.(*DeleteProjectMembershipRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProjectService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "n0stack.iam.v1alpha.ProjectService",
	HandlerType: (*ProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListProjects",
			Handler:    _ProjectService_ListProjects_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _ProjectService_GetProject_Handler,
		},
		{
			MethodName: "CreateProject",
			Handler:    _ProjectService_CreateProject_Handler,
		},
		{
			MethodName: "UpdateProject",
			Handler:    _ProjectService_UpdateProject_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _ProjectService_DeleteProject_Handler,
		},
		{
			MethodName: "AddProjectMembership",
			Handler:    _ProjectService_AddProjectMembership_Handler,
		},
		{
			MethodName: "DeleteProjectMembership",
			Handler:    _ProjectService_DeleteProjectMembership_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "n0stack/iam/v1alpha/project.proto",
}