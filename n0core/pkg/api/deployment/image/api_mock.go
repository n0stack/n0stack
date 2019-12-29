package image

// type MockImageAPI struct {
// 	api             *ImageAPI
// 	NodeAPI         *node.MockNodeAPI
// 	BlockStorageAPI *blockstorage.MockBlockStorageAPI
// }

// func NewMockVirtualMachineAPI(datastore *memory.MemoryDatastore) *MockImageAPI {
// 	noa := node.NewMockNodeAPI(datastore)
// 	nea := network.NewMockNetworkAPI(datastore)
// 	bsa := blockstorage.NewMockBlcokStorageAPI(datastore)

// 	a := CreateVirtualMachineAPI(datastore, noa, nea, bsa)
// 	a.getAgent = func(ctx context.Context, nodeName string) (VirtualMachineAgentServiceClient, func() error, error) {
// 		return NewMockVirtualMachineAgentClientMock(), func() error { return nil }, nil
// 	}

// 	return &MockImageAPI{a, noa, nea, bsa}
// }

// func (a MockVirtualMachineAPI) CreateVirtualMachine(ctx context.Context, in *pprovisioning.CreateVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
// 	return a.api.CreateVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) ListVirtualMachines(ctx context.Context, in *pprovisioning.ListVirtualMachinesRequest, opts ...grpc.CallOption) (*pprovisioning.ListVirtualMachinesResponse, error) {
// 	return a.api.ListVirtualMachines(ctx, in)
// }
// func (a MockVirtualMachineAPI) GetVirtualMachine(ctx context.Context, in *pprovisioning.GetVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
// 	return a.api.GetVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) UpdateVirtualMachine(ctx context.Context, in *pprovisioning.UpdateVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
// 	return a.api.UpdateVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) DeleteVirtualMachine(ctx context.Context, in *pprovisioning.DeleteVirtualMachineRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
// 	return a.api.DeleteVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) BootVirtualMachine(ctx context.Context, in *pprovisioning.BootVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
// 	return a.api.BootVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) RebootVirtualMachine(ctx context.Context, in *pprovisioning.RebootVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
// 	return a.api.RebootVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) ShutdownVirtualMachine(ctx context.Context, in *pprovisioning.ShutdownVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
// 	return a.api.ShutdownVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) SaveVirtualMachine(ctx context.Context, in *pprovisioning.SaveVirtualMachineRequest, opts ...grpc.CallOption) (*pprovisioning.VirtualMachine, error) {
// 	return a.api.SaveVirtualMachine(ctx, in)
// }
// func (a MockVirtualMachineAPI) OpenConsole(ctx context.Context, in *pprovisioning.OpenConsoleRequest, opts ...grpc.CallOption) (*pprovisioning.OpenConsoleResponse, error) {
// 	return a.api.OpenConsole(ctx, in)
// }
