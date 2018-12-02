package virtualmachine

import (
	"context"
	"fmt"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/network"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

type MockVirtualMachineAPI struct {
	api             *VirtualMachineAPI
	NodeAPI         *node.MockNodeAPI
	NetworkAPI      *network.MockNetworkAPI
	BlockStorageAPI *provisioning.MockBlockStorageAPI
}

// どうもサービスが始まるまでのタイムラグがあるせいで、性能の悪いデバイスでは安定性が悪い
// TODO: 上位層が使いにくくなっているので変える
func UpMockAgent(address string) error {
	addr := fmt.Sprintf("%s:%d", address, 20181)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	provisioning.RegisterBlockStorageAgentServiceServer(grpcServer, &provisioning.MockBlockStorageAgentAPI{})
	RegisterVirtualMachineAgentServiceServer(grpcServer, &VirtualMachineAgentMock{})
	return grpcServer.Serve(lis)
}

func NewMockVirtualMachineAPI(datastore *memory.MemoryDatastore) *MockVirtualMachineAPI {
	noa := node.NewMockNodeAPI(datastore)
	nea := network.NewMockNetworkAPI(datastore)
	bsa := provisioning.NewMockBlcokStorageAPI(datastore)

	a := CreateVirtualMachineAPI(datastore, noa, nea, bsa)
	return &MockVirtualMachineAPI{a, noa, nea, bsa}
}

func (a MockVirtualMachineAPI) CreateVirtualMachine(ctx context.Context, in *pprovisioning.CreateVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
	return a.api.CreateVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) ListVirtualMachines(ctx context.Context, in *pprovisioning.ListVirtualMachinesRequest, opts ...grpc.CallOption) (*pprovisioning.ListVirtualMachinesResponse, error) {
	return a.api.ListVirtualMachines(ctx, in)
}
func (a MockVirtualMachineAPI) GetVirtualMachine(ctx context.Context, in *pprovisioning.GetVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
	return a.api.GetVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) UpdateVirtualMachine(ctx context.Context, in *pprovisioning.UpdateVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
	return a.api.UpdateVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) DeleteVirtualMachine(ctx context.Context, in *pprovisioning.DeleteVirtualMachineRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.api.DeleteVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) BootVirtualMachine(ctx context.Context, in *pprovisioning.BootVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
	return a.api.BootVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) RebootVirtualMachine(ctx context.Context, in *pprovisioning.RebootVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
	return a.api.RebootVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) ShutdownVirtualMachine(ctx context.Context, in *pprovisioning.ShutdownVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
	return a.api.ShutdownVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) SaveVirtualMachine(ctx context.Context, in *pprovisioning.SaveVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
	return a.api.SaveVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAPI) OpenConsole(ctx context.Context, in *pprovisioning.OpenConsoleRequest, opts ...grpc.CallOption) (*pprovisioning.OpenConsoleResponse, error) {
	return a.api.OpenConsole(ctx, in)
}
