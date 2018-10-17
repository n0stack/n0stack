package provisioning

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type MockBlockStorageAgentAPI struct{}

func (a MockBlockStorageAgentAPI) CreateEmptyBlockStorageAgent(ctx context.Context, req *CreateEmptyBlockStorageAgentRequest) (*BlockStorageAgent, error) {
	return &BlockStorageAgent{
		Name:  req.Name,
		Bytes: req.Bytes,
		Path:  "/tmp/test",
	}, nil
}
func (a MockBlockStorageAgentAPI) CreateBlockStorageAgentWithDownloading(ctx context.Context, req *CreateBlockStorageAgentWithDownloadingRequest) (*BlockStorageAgent, error) {
	return &BlockStorageAgent{
		Name:  req.Name,
		Bytes: req.Bytes,
		Path:  "/tmp/test",
	}, nil
}
func (a MockBlockStorageAgentAPI) DeleteBlockStorageAgent(ctx context.Context, req *DeleteBlockStorageAgentRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
