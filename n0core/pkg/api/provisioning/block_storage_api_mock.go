package provisioning

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

type MockBlockStorageAPI struct {
	api     *BlockStorageAPI
	NodeAPI *node.MockNodeAPI
}

func NewMockBlcokStorageAPI(datastore *memory.MemoryDatastore) *MockBlockStorageAPI {
	na := node.NewMockNodeAPI(datastore)

	a := CreateBlockStorageAPI(datastore, na)
	return &MockBlockStorageAPI{a, na}
}

func (a MockBlockStorageAPI) CreateBlockStorage(ctx context.Context, in *pprovisioning.CreateBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.CreateBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) FetchBlockStorage(ctx context.Context, in *pprovisioning.FetchBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.FetchBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) CopyBlockStorage(ctx context.Context, in *pprovisioning.CopyBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.CopyBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) ListBlockStorages(ctx context.Context, in *pprovisioning.ListBlockStoragesRequest, opts ...grpc.CallOption) (*pprovisioning.ListBlockStoragesResponse, error) {
	return a.api.ListBlockStorages(ctx, in)
}
func (a MockBlockStorageAPI) GetBlockStorage(ctx context.Context, in *pprovisioning.GetBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.GetBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) UpdateBlockStorage(ctx context.Context, in *pprovisioning.UpdateBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.UpdateBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) DeleteBlockStorage(ctx context.Context, in *pprovisioning.DeleteBlockStorageRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.api.DeleteBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) SetInuseBlockStorage(ctx context.Context, in *pprovisioning.SetInuseBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.SetInuseBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) SetAvailableBlockStorage(ctx context.Context, in *pprovisioning.SetAvailableBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.SetAvailableBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) SetProtectedBlockStorage(ctx context.Context, in *pprovisioning.SetProtectedBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.SetProtectedBlockStorage(ctx, in)
}
