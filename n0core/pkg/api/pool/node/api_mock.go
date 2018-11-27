package node

import (
	"context"

	"code.cloudfoundry.org/bytefmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"google.golang.org/grpc"
)

type MockNodeAPI struct {
	a *NodeAPI
}

func NewMockNodeAPI(datastore *memory.MemoryDatastore) *MockNodeAPI {
	a := CreateNodeAPI(datastore)
	return &MockNodeAPI{a}
}

func (a MockNodeAPI) SetupMockNode(ctx context.Context) (*ppool.Node, error) {
	return a.ApplyNode(ctx, &ppool.ApplyNodeRequest{
		Name:          "mocked",
		Address:       "127.0.20.180",
		CpuMilliCores: 1000,
		MemoryBytes:   10 * bytefmt.GIGABYTE,
		StorageBytes:  10 * bytefmt.GIGABYTE,
	})
}

func (a MockNodeAPI) ListNodes(ctx context.Context, in *ppool.ListNodesRequest, opts ...grpc.CallOption) (*ppool.ListNodesResponse, error) {
	return a.a.ListNodes(ctx, in)
}
func (a MockNodeAPI) GetNode(ctx context.Context, in *ppool.GetNodeRequest, opts ...grpc.CallOption) (*ppool.Node, error) {
	return a.a.GetNode(ctx, in)
}
func (a MockNodeAPI) ApplyNode(ctx context.Context, in *ppool.ApplyNodeRequest, opts ...grpc.CallOption) (*ppool.Node, error) {
	return a.a.ApplyNode(ctx, in)
}
func (a MockNodeAPI) DeleteNode(ctx context.Context, in *ppool.DeleteNodeRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.a.DeleteNode(ctx, in)
}
func (a MockNodeAPI) ScheduleCompute(ctx context.Context, in *ppool.ScheduleComputeRequest, opts ...grpc.CallOption) (*ppool.Node, error) {
	return a.a.ScheduleCompute(ctx, in)
}
func (a MockNodeAPI) ReserveCompute(ctx context.Context, in *ppool.ReserveComputeRequest, opts ...grpc.CallOption) (*ppool.Node, error) {
	return a.a.ReserveCompute(ctx, in)
}
func (a MockNodeAPI) ReleaseCompute(ctx context.Context, in *ppool.ReleaseComputeRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.a.ReleaseCompute(ctx, in)
}
func (a MockNodeAPI) ScheduleStorage(ctx context.Context, in *ppool.ScheduleStorageRequest, opts ...grpc.CallOption) (*ppool.Node, error) {
	return a.a.ScheduleStorage(ctx, in)
}
func (a MockNodeAPI) ReserveStorage(ctx context.Context, in *ppool.ReserveStorageRequest, opts ...grpc.CallOption) (*ppool.Node, error) {
	return a.a.ReserveStorage(ctx, in)
}
func (a MockNodeAPI) ReleaseStorage(ctx context.Context, in *ppool.ReleaseStorageRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.a.ReleaseStorage(ctx, in)
}
