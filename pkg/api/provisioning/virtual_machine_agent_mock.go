package provisioning

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
)

type MockVirtualMachineAgentAPI struct{}

func (a *MockVirtualMachineAgentAPI) CreateVirtualMachineAgent(ctx context.Context, req *CreateVirtualMachineAgentRequest) (*VirtualMachineAgent, error) {
	u, _ := uuid.FromString("1d5fd196-b6c9-4f58-86f2-3ef227018e47")

	return &VirtualMachineAgent{
		Name:          req.Name,
		Uuid:          u.String(),
		Vcpus:         req.Vcpus,
		MemoryBytes:   req.MemoryBytes,
		State:         VirtualMachineAgentState_RUNNING,
		Blockdev:      req.Blockdev,
		Netdev:        req.Netdev,
		WebsocketPort: 10000,
	}, nil
}
func (a *MockVirtualMachineAgentAPI) DeleteVirtualMachineAgent(ctx context.Context, req *DeleteVirtualMachineAgentRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
func (a *MockVirtualMachineAgentAPI) BootVirtualMachineAgent(ctx context.Context, req *BootVirtualMachineAgentRequest) (*BootVirtualMachineAgentResponse, error) {
	return &BootVirtualMachineAgentResponse{
		State: VirtualMachineAgentState_RUNNING,
	}, nil
}
func (a *MockVirtualMachineAgentAPI) RebootVirtualMachineAgent(ctx context.Context, req *RebootVirtualMachineAgentRequest) (*RebootVirtualMachineAgentResponse, error) {
	return &RebootVirtualMachineAgentResponse{
		State: VirtualMachineAgentState_RUNNING,
	}, nil
}
func (a *MockVirtualMachineAgentAPI) ShutdownVirtualMachineAgent(ctx context.Context, req *ShutdownVirtualMachineAgentRequest) (*ShutdownVirtualMachineAgentResponse, error) {
	return &ShutdownVirtualMachineAgentResponse{
		State: VirtualMachineAgentState_SHUTDOWN,
	}, nil
}
