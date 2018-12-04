package blockstorage

import (
	"context"
	"path/filepath"

	"github.com/golang/protobuf/ptypes/empty"
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
