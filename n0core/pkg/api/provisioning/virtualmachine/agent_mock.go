package virtualmachine

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
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
