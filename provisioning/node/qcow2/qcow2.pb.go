// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/agent/qcow2/qcow2.proto

/*
Package qcow2 is a generated protocol buffer package.

It is generated from these files:
	proto/agent/qcow2/qcow2.proto

It has these top-level messages:
	Qcow2
	ApplyQcow2Request
	DownloadQcow2Request
	BuildQcow2WithPackerRequest
	DeleteQcow2Request
*/
package qcow2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Qcow2 struct {
	// url is id for Qcow2.
	// Location of qcow2 file.
	Url string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	// サイズを指定する
	Bytes uint64 `protobuf:"varint,2,opt,name=bytes" json:"bytes,omitempty"`
}

func (m *Qcow2) Reset()                    { *m = Qcow2{} }
func (m *Qcow2) String() string            { return proto.CompactTextString(m) }
func (*Qcow2) ProtoMessage()               {}
func (*Qcow2) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Qcow2) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Qcow2) GetBytes() uint64 {
	if m != nil {
		return m.Bytes
	}
	return 0
}

type ApplyQcow2Request struct {
	Qcow2 *Qcow2 `protobuf:"bytes,1,opt,name=qcow2" json:"qcow2,omitempty"`
}

func (m *ApplyQcow2Request) Reset()                    { *m = ApplyQcow2Request{} }
func (m *ApplyQcow2Request) String() string            { return proto.CompactTextString(m) }
func (*ApplyQcow2Request) ProtoMessage()               {}
func (*ApplyQcow2Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ApplyQcow2Request) GetQcow2() *Qcow2 {
	if m != nil {
		return m.Qcow2
	}
	return nil
}

type DownloadQcow2Request struct {
	Qcow2 *Qcow2 `protobuf:"bytes,1,opt,name=qcow2" json:"qcow2,omitempty"`
	// URLからイメージファイルをダウンロードすることができる
	SourceUrl string `protobuf:"bytes,2,opt,name=source_url,json=sourceUrl" json:"source_url,omitempty"`
}

func (m *DownloadQcow2Request) Reset()                    { *m = DownloadQcow2Request{} }
func (m *DownloadQcow2Request) String() string            { return proto.CompactTextString(m) }
func (*DownloadQcow2Request) ProtoMessage()               {}
func (*DownloadQcow2Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DownloadQcow2Request) GetQcow2() *Qcow2 {
	if m != nil {
		return m.Qcow2
	}
	return nil
}

func (m *DownloadQcow2Request) GetSourceUrl() string {
	if m != nil {
		return m.SourceUrl
	}
	return ""
}

type BuildQcow2WithPackerRequest struct {
	Qcow2         *Qcow2 `protobuf:"bytes,1,opt,name=qcow2" json:"qcow2,omitempty"`
	Repository    string `protobuf:"bytes,2,opt,name=repository" json:"repository,omitempty"`
	WorkDirectory string `protobuf:"bytes,3,opt,name=work_directory,json=workDirectory" json:"work_directory,omitempty"`
	TemplateFile  string `protobuf:"bytes,4,opt,name=template_file,json=templateFile" json:"template_file,omitempty"`
}

func (m *BuildQcow2WithPackerRequest) Reset()                    { *m = BuildQcow2WithPackerRequest{} }
func (m *BuildQcow2WithPackerRequest) String() string            { return proto.CompactTextString(m) }
func (*BuildQcow2WithPackerRequest) ProtoMessage()               {}
func (*BuildQcow2WithPackerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *BuildQcow2WithPackerRequest) GetQcow2() *Qcow2 {
	if m != nil {
		return m.Qcow2
	}
	return nil
}

func (m *BuildQcow2WithPackerRequest) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *BuildQcow2WithPackerRequest) GetWorkDirectory() string {
	if m != nil {
		return m.WorkDirectory
	}
	return ""
}

func (m *BuildQcow2WithPackerRequest) GetTemplateFile() string {
	if m != nil {
		return m.TemplateFile
	}
	return ""
}

type DeleteQcow2Request struct {
	Qcow2 *Qcow2 `protobuf:"bytes,1,opt,name=qcow2" json:"qcow2,omitempty"`
}

func (m *DeleteQcow2Request) Reset()                    { *m = DeleteQcow2Request{} }
func (m *DeleteQcow2Request) String() string            { return proto.CompactTextString(m) }
func (*DeleteQcow2Request) ProtoMessage()               {}
func (*DeleteQcow2Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DeleteQcow2Request) GetQcow2() *Qcow2 {
	if m != nil {
		return m.Qcow2
	}
	return nil
}

func init() {
	proto.RegisterType((*Qcow2)(nil), "n0stack.n0core.agent.qcow2.Qcow2")
	proto.RegisterType((*ApplyQcow2Request)(nil), "n0stack.n0core.agent.qcow2.ApplyQcow2Request")
	proto.RegisterType((*DownloadQcow2Request)(nil), "n0stack.n0core.agent.qcow2.DownloadQcow2Request")
	proto.RegisterType((*BuildQcow2WithPackerRequest)(nil), "n0stack.n0core.agent.qcow2.BuildQcow2WithPackerRequest")
	proto.RegisterType((*DeleteQcow2Request)(nil), "n0stack.n0core.agent.qcow2.DeleteQcow2Request")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Qcow2Service service

type Qcow2ServiceClient interface {
	ApplyQcow2(ctx context.Context, in *ApplyQcow2Request, opts ...grpc.CallOption) (*Qcow2, error)
	DownloadQcow2(ctx context.Context, in *DownloadQcow2Request, opts ...grpc.CallOption) (*Qcow2, error)
	BuildQcow2WithPacker(ctx context.Context, in *BuildQcow2WithPackerRequest, opts ...grpc.CallOption) (*Qcow2, error)
	DeleteQcow2(ctx context.Context, in *DeleteQcow2Request, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type qcow2ServiceClient struct {
	cc *grpc.ClientConn
}

func NewQcow2ServiceClient(cc *grpc.ClientConn) Qcow2ServiceClient {
	return &qcow2ServiceClient{cc}
}

func (c *qcow2ServiceClient) ApplyQcow2(ctx context.Context, in *ApplyQcow2Request, opts ...grpc.CallOption) (*Qcow2, error) {
	out := new(Qcow2)
	err := grpc.Invoke(ctx, "/n0stack.n0core.agent.qcow2.Qcow2Service/ApplyQcow2", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *qcow2ServiceClient) DownloadQcow2(ctx context.Context, in *DownloadQcow2Request, opts ...grpc.CallOption) (*Qcow2, error) {
	out := new(Qcow2)
	err := grpc.Invoke(ctx, "/n0stack.n0core.agent.qcow2.Qcow2Service/DownloadQcow2", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *qcow2ServiceClient) BuildQcow2WithPacker(ctx context.Context, in *BuildQcow2WithPackerRequest, opts ...grpc.CallOption) (*Qcow2, error) {
	out := new(Qcow2)
	err := grpc.Invoke(ctx, "/n0stack.n0core.agent.qcow2.Qcow2Service/BuildQcow2WithPacker", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *qcow2ServiceClient) DeleteQcow2(ctx context.Context, in *DeleteQcow2Request, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/n0stack.n0core.agent.qcow2.Qcow2Service/DeleteQcow2", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Qcow2Service service

type Qcow2ServiceServer interface {
	ApplyQcow2(context.Context, *ApplyQcow2Request) (*Qcow2, error)
	DownloadQcow2(context.Context, *DownloadQcow2Request) (*Qcow2, error)
	BuildQcow2WithPacker(context.Context, *BuildQcow2WithPackerRequest) (*Qcow2, error)
	DeleteQcow2(context.Context, *DeleteQcow2Request) (*google_protobuf.Empty, error)
}

func RegisterQcow2ServiceServer(s *grpc.Server, srv Qcow2ServiceServer) {
	s.RegisterService(&_Qcow2Service_serviceDesc, srv)
}

func _Qcow2Service_ApplyQcow2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyQcow2Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Qcow2ServiceServer).ApplyQcow2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.n0core.agent.qcow2.Qcow2Service/ApplyQcow2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Qcow2ServiceServer).ApplyQcow2(ctx, req.(*ApplyQcow2Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Qcow2Service_DownloadQcow2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadQcow2Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Qcow2ServiceServer).DownloadQcow2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.n0core.agent.qcow2.Qcow2Service/DownloadQcow2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Qcow2ServiceServer).DownloadQcow2(ctx, req.(*DownloadQcow2Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Qcow2Service_BuildQcow2WithPacker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuildQcow2WithPackerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Qcow2ServiceServer).BuildQcow2WithPacker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.n0core.agent.qcow2.Qcow2Service/BuildQcow2WithPacker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Qcow2ServiceServer).BuildQcow2WithPacker(ctx, req.(*BuildQcow2WithPackerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Qcow2Service_DeleteQcow2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteQcow2Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Qcow2ServiceServer).DeleteQcow2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/n0stack.n0core.agent.qcow2.Qcow2Service/DeleteQcow2",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Qcow2ServiceServer).DeleteQcow2(ctx, req.(*DeleteQcow2Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Qcow2Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "n0stack.n0core.agent.qcow2.Qcow2Service",
	HandlerType: (*Qcow2ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ApplyQcow2",
			Handler:    _Qcow2Service_ApplyQcow2_Handler,
		},
		{
			MethodName: "DownloadQcow2",
			Handler:    _Qcow2Service_DownloadQcow2_Handler,
		},
		{
			MethodName: "BuildQcow2WithPacker",
			Handler:    _Qcow2Service_BuildQcow2WithPacker_Handler,
		},
		{
			MethodName: "DeleteQcow2",
			Handler:    _Qcow2Service_DeleteQcow2_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/agent/qcow2/qcow2.proto",
}

func init() { proto.RegisterFile("proto/agent/qcow2/qcow2.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x54, 0xdd, 0x8a, 0xd3, 0x40,
	0x14, 0xde, 0x6c, 0xb7, 0xc2, 0x9e, 0xdd, 0x8a, 0x0e, 0x45, 0x4a, 0x96, 0x95, 0x35, 0x22, 0xec,
	0x8d, 0x33, 0xb5, 0x7b, 0xb1, 0xd7, 0x2e, 0xd5, 0x2b, 0x05, 0xad, 0x48, 0xc1, 0x9b, 0x9a, 0x4c,
	0x4f, 0xd3, 0xa1, 0xd3, 0x99, 0x74, 0x32, 0x69, 0xc9, 0x13, 0xf9, 0x1e, 0x3e, 0x99, 0x64, 0xa6,
	0xc5, 0x4a, 0x6b, 0x2c, 0xd8, 0x9b, 0x90, 0x7c, 0x7c, 0x3f, 0x99, 0x7c, 0xe7, 0x04, 0xae, 0x33,
	0xa3, 0xad, 0x66, 0x71, 0x8a, 0xca, 0xb2, 0x05, 0xd7, 0xab, 0x9e, 0xbf, 0x52, 0x87, 0x93, 0x50,
	0x75, 0x73, 0x1b, 0xf3, 0x19, 0x55, 0x5d, 0xae, 0x0d, 0x52, 0xc7, 0xa3, 0x8e, 0x11, 0x5e, 0xa5,
	0x5a, 0xa7, 0x12, 0x99, 0x63, 0x26, 0xc5, 0x84, 0xe1, 0x3c, 0xb3, 0xa5, 0x17, 0x46, 0x0c, 0x9a,
	0x9f, 0x2b, 0x16, 0x79, 0x02, 0x8d, 0xc2, 0xc8, 0x4e, 0x70, 0x13, 0xdc, 0x9e, 0x0f, 0xaa, 0x5b,
	0xd2, 0x86, 0x66, 0x52, 0x5a, 0xcc, 0x3b, 0xa7, 0x37, 0xc1, 0xed, 0xd9, 0xc0, 0x3f, 0x44, 0x1f,
	0xe0, 0xe9, 0xdb, 0x2c, 0x93, 0xa5, 0x53, 0x0d, 0x70, 0x51, 0x60, 0x6e, 0xc9, 0x3d, 0x34, 0x5d,
	0x96, 0x93, 0x5f, 0xf4, 0x5e, 0xd0, 0xbf, 0xbf, 0x0e, 0xf5, 0x42, 0xcf, 0x8f, 0x14, 0xb4, 0xfb,
	0x7a, 0xa5, 0xa4, 0x8e, 0xc7, 0x47, 0x31, 0x24, 0xd7, 0x00, 0xb9, 0x2e, 0x0c, 0xc7, 0x51, 0x75,
	0x9a, 0x53, 0x77, 0x9a, 0x73, 0x8f, 0x7c, 0x35, 0x32, 0xfa, 0x19, 0xc0, 0xd5, 0x43, 0x21, 0xa4,
	0x4f, 0x1b, 0x0a, 0x3b, 0xfd, 0x14, 0xf3, 0x19, 0x9a, 0xff, 0xce, 0x7d, 0x0e, 0x60, 0x30, 0xd3,
	0xb9, 0xb0, 0xda, 0x94, 0xeb, 0xdc, 0x2d, 0x84, 0xbc, 0x82, 0xc7, 0x2b, 0x6d, 0x66, 0xa3, 0xb1,
	0x30, 0xc8, 0x1d, 0xa7, 0xe1, 0x38, 0xad, 0x0a, 0xed, 0x6f, 0x40, 0xf2, 0x12, 0x5a, 0x16, 0xe7,
	0x99, 0x8c, 0x2d, 0x8e, 0x26, 0x42, 0x62, 0xe7, 0xcc, 0xb1, 0x2e, 0x37, 0xe0, 0x7b, 0x21, 0x31,
	0xfa, 0x08, 0xa4, 0x8f, 0x12, 0x2d, 0x1e, 0xe5, 0x93, 0xf5, 0x7e, 0x34, 0xe0, 0xd2, 0x01, 0x5f,
	0xd0, 0x2c, 0x05, 0x47, 0xf2, 0x1d, 0xe0, 0x77, 0xc5, 0xe4, 0x75, 0x9d, 0xd1, 0xce, 0x28, 0x84,
	0xff, 0xce, 0x8d, 0x4e, 0xc8, 0x04, 0x5a, 0x7f, 0xd4, 0x4e, 0xba, 0x75, 0xaa, 0x7d, 0x13, 0x72,
	0x58, 0x8e, 0x85, 0xf6, 0xbe, 0xb6, 0xc9, 0x7d, 0x9d, 0xb8, 0x66, 0x3e, 0x0e, 0x4b, 0x1d, 0xc2,
	0xc5, 0x56, 0x3f, 0x84, 0xd6, 0x9e, 0x6d, 0xa7, 0xc8, 0xf0, 0x19, 0xf5, 0x0b, 0x4b, 0x37, 0x0b,
	0x4b, 0xdf, 0x55, 0x0b, 0x1b, 0x9d, 0x3c, 0xdc, 0x7d, 0x7b, 0x93, 0x0a, 0x3b, 0x2d, 0x12, 0xca,
	0xf5, 0x9c, 0xad, 0x5d, 0x99, 0x77, 0xad, 0xd6, 0x7b, 0x29, 0x72, 0xa1, 0x95, 0x50, 0x29, 0x53,
	0x7a, 0x8c, 0xfe, 0x07, 0x91, 0x3c, 0x72, 0x36, 0x77, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x07,
	0x2d, 0x6f, 0x05, 0x42, 0x04, 0x00, 0x00,
}
