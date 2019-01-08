package virtualmachine

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type VirtualMachineAgentMock struct{}

func (a VirtualMachineAgentMock) BootVirtualMachine(ctx context.Context, req *BootVirtualMachineRequest) (*BootVirtualMachineResponse, error) {
	return &BootVirtualMachineResponse{
		State:         VirtualMachineState_RUNNING,
		WebsocketPort: 6900,
	}, nil
}
func (a VirtualMachineAgentMock) RebootVirtualMachine(ctx context.Context, req *RebootVirtualMachineRequest) (*RebootVirtualMachineResponse, error) {
	return &RebootVirtualMachineResponse{
		State: VirtualMachineState_RUNNING,
	}, nil
}
func (a VirtualMachineAgentMock) ShutdownVirtualMachine(ctx context.Context, req *ShutdownVirtualMachineRequest) (*ShutdownVirtualMachineResponse, error) {
	return &ShutdownVirtualMachineResponse{
		State: VirtualMachineState_SHUTDOWN,
	}, nil
}
func (a VirtualMachineAgentMock) DeleteVirtualMachine(ctx context.Context, req *DeleteVirtualMachineRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

type MockVirtualMachineAgentClient struct {
	api VirtualMachineAgentServiceServer
}

func NewMockVirtualMachineAgentClientMock() *MockVirtualMachineAgentClient {
	return &MockVirtualMachineAgentClient{
		api: VirtualMachineAgentMock{},
	}
}
func (a MockVirtualMachineAgentClient) BootVirtualMachine(ctx context.Context, in *BootVirtualMachineRequest, opts ...grpc.CallOption) (*BootVirtualMachineResponse, error) {
	return a.api.BootVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAgentClient) RebootVirtualMachine(ctx context.Context, in *RebootVirtualMachineRequest, opts ...grpc.CallOption) (*RebootVirtualMachineResponse, error) {
	return a.api.RebootVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAgentClient) ShutdownVirtualMachine(ctx context.Context, in *ShutdownVirtualMachineRequest, opts ...grpc.CallOption) (*ShutdownVirtualMachineResponse, error) {
	return a.api.ShutdownVirtualMachine(ctx, in)
}
func (a MockVirtualMachineAgentClient) DeleteVirtualMachine(ctx context.Context, in *DeleteVirtualMachineRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return a.api.DeleteVirtualMachine(ctx, in)
}
