package blockstorage

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/bytefmt"
	"github.com/golang/protobuf/ptypes/empty"
	"n0st.ac/n0stack/n0core/pkg/api/pool/node"
	"n0st.ac/n0stack/n0core/pkg/datastore/memory"
	pprovisioning "n0st.ac/n0stack/n0proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

type MockBlockStorageAPI struct {
	api     *BlockStorageAPI
	NodeAPI *node.MockNodeAPI
}

var factroyIndex = 0

func NewMockBlcokStorageAPI(datastore *memory.MemoryDatastore) *MockBlockStorageAPI {
	na := node.NewMockNodeAPI(datastore)

	a := CreateBlockStorageAPI(datastore, na)
	a.getAgent = func(ctx context.Context, nodeName string) (BlockStorageAgentServiceClient, func() error, error) {
		return NewMockBlockStorageAgentClient(), func() error { return nil }, nil
	}

	return &MockBlockStorageAPI{a, na}
}

func (a MockBlockStorageAPI) FactoryBlockStorage(ctx context.Context, nodeName string) (*pprovisioning.BlockStorage, error) {
	factroyIndex++

	return a.api.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: fmt.Sprintf("factory-blockstorage%d", factroyIndex),
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: nodeName,
		},
		LimitBytes:   10 * bytefmt.GIGABYTE,
		RequestBytes: 10 * bytefmt.GIGABYTE,
	})
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
func (a MockBlockStorageAPI) DeleteBlockStorage(ctx context.Context, in *pprovisioning.DeleteBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.DeleteBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) UndeleteBlockStorage(ctx context.Context, in *pprovisioning.UndeleteBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.BlockStorage, error) {
	return a.api.UndeleteBlockStorage(ctx, in)
}
func (a MockBlockStorageAPI) PurgeBlockStorage(ctx context.Context, in *pprovisioning.PurgeBlockStorageRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.api.PurgeBlockStorage(ctx, in)
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
func (a MockBlockStorageAPI) DownloadBlockStorage(ctx context.Context, in *pprovisioning.DownloadBlockStorageRequest, opts ...grpc.CallOption) (*pprovisioning.DownloadBlockStorageResponse, error) {
	return a.api.DownloadBlockStorage(ctx, in)
}
