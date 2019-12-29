package network

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"n0st.ac/n0stack/n0core/pkg/datastore/memory"
	"n0st.ac/n0stack/n0proto.go/pool/v0"
	"google.golang.org/grpc"
)

var factroyIndex = 0

type MockNetworkAPI struct {
	api *NetworkAPI
}

func NewMockNetworkAPI(datastore *memory.MemoryDatastore) *MockNetworkAPI {
	n := CreateNetworkAPI(datastore)
	return &MockNetworkAPI{n}
}

func (a MockNetworkAPI) FactoryNetwork(ctx context.Context) (*ppool.Network, error) {
	factroyIndex++

	return a.api.ApplyNetwork(ctx, &ppool.ApplyNetworkRequest{
		Name:     fmt.Sprintf("factory-network%d", factroyIndex),
		Domain:   fmt.Sprintf("factory-network%d.test", factroyIndex),
		Ipv4Cidr: fmt.Sprintf("10.0.%d.0/24", factroyIndex),
		Ipv6Cidr: fmt.Sprintf("fc00:%x::1/64", factroyIndex),
	})
}

func (a MockNetworkAPI) ListNetworks(ctx context.Context, in *ppool.ListNetworksRequest, opts ...grpc.CallOption) (*ppool.ListNetworksResponse, error) {
	return a.api.ListNetworks(ctx, in)
}
func (a MockNetworkAPI) GetNetwork(ctx context.Context, in *ppool.GetNetworkRequest, opts ...grpc.CallOption) (*ppool.Network, error) {
	return a.api.GetNetwork(ctx, in)
}
func (a MockNetworkAPI) ApplyNetwork(ctx context.Context, in *ppool.ApplyNetworkRequest, opts ...grpc.CallOption) (*ppool.Network, error) {
	return a.api.ApplyNetwork(ctx, in)
}
func (a MockNetworkAPI) DeleteNetwork(ctx context.Context, in *ppool.DeleteNetworkRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.api.DeleteNetwork(ctx, in)
}
func (a MockNetworkAPI) ReserveNetworkInterface(ctx context.Context, in *ppool.ReserveNetworkInterfaceRequest, opts ...grpc.CallOption) (*ppool.Network, error) {
	return a.api.ReserveNetworkInterface(ctx, in)
}
func (a MockNetworkAPI) ReleaseNetworkInterface(ctx context.Context, in *ppool.ReleaseNetworkInterfaceRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.api.ReleaseNetworkInterface(ctx, in)
}
