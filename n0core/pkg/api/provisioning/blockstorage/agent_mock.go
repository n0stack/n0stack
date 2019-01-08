package blockstorage

import (
	"context"
	"path/filepath"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type MockBlockStorageAgentAPI struct{}

func (a MockBlockStorageAgentAPI) CreateEmptyBlockStorage(ctx context.Context, req *CreateEmptyBlockStorageRequest) (*CreateEmptyBlockStorageResponse, error) {
	return &CreateEmptyBlockStorageResponse{
		Path: filepath.Join("/tmp", req.Name),
	}, nil
}
func (a MockBlockStorageAgentAPI) FetchBlockStorage(ctx context.Context, req *FetchBlockStorageRequest) (*FetchBlockStorageResponse, error) {
	return &FetchBlockStorageResponse{
		Path: filepath.Join("/tmp", req.Name),
	}, nil
}
func (a MockBlockStorageAgentAPI) DeleteBlockStorage(ctx context.Context, req *DeleteBlockStorageRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

type MockBlockStorageAgentClient struct {
	api BlockStorageAgentServiceServer
}

func NewMockBlockStorageAgentClient() *MockBlockStorageAgentClient {
	return &MockBlockStorageAgentClient{
		api: &MockBlockStorageAgentAPI{},
	}
}
func (a MockBlockStorageAgentClient) CreateEmptyBlockStorage(ctx context.Context, in *CreateEmptyBlockStorageRequest, opts ...grpc.CallOption) (*CreateEmptyBlockStorageResponse, error) {
	return a.api.CreateEmptyBlockStorage(ctx, in)
}
func (a MockBlockStorageAgentClient) FetchBlockStorage(ctx context.Context, in *FetchBlockStorageRequest, opts ...grpc.CallOption) (*FetchBlockStorageResponse, error) {
	return a.api.FetchBlockStorage(ctx, in)
}
func (a MockBlockStorageAgentClient) DeleteBlockStorage(ctx context.Context, in *DeleteBlockStorageRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.api.DeleteBlockStorage(ctx, in)
}
