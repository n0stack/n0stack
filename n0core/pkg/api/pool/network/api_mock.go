package network

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"google.golang.org/grpc"
)

type MockNetworkAPI struct {
	a *NetworkAPI
}

func NewMockNetworkAPI(datastore *memory.MemoryDatastore) ppool.NetworkServiceClient {
	n := CreateNetworkAPI(datastore)
	return &MockNetworkAPI{n}
}
func (a MockNetworkAPI) ListNetworks(ctx context.Context, in *ppool.ListNetworksRequest, opts ...grpc.CallOption) (*ppool.ListNetworksResponse, error) {
	return a.a.ListNetworks(ctx, in)
}
func (a MockNetworkAPI) GetNetwork(ctx context.Context, in *ppool.GetNetworkRequest, opts ...grpc.CallOption) (*ppool.Network, error) {
	return a.a.GetNetwork(ctx, in)
}
func (a MockNetworkAPI) ApplyNetwork(ctx context.Context, in *ppool.ApplyNetworkRequest, opts ...grpc.CallOption) (*ppool.Network, error) {
	return a.a.ApplyNetwork(ctx, in)
}
func (a MockNetworkAPI) DeleteNetwork(ctx context.Context, in *ppool.DeleteNetworkRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.a.DeleteNetwork(ctx, in)
}
func (a MockNetworkAPI) ReserveNetworkInterface(ctx context.Context, in *ppool.ReserveNetworkInterfaceRequest, opts ...grpc.CallOption) (*ppool.Network, error) {
	return a.a.ReserveNetworkInterface(ctx, in)
}
func (a MockNetworkAPI) ReleaseNetworkInterface(ctx context.Context, in *ppool.ReleaseNetworkInterfaceRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.a.ReleaseNetworkInterface(ctx, in)
}
