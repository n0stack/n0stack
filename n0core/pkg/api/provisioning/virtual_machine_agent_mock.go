package provisioning

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

type MockVirtualMachineAgentAPI struct{}

func (a *MockVirtualMachineAgentAPI) CreateVirtualMachineAgent(ctx context.Context, req *CreateVirtualMachineAgentRequest) (*VirtualMachineAgent, error) {
	return &VirtualMachineAgent{
		Name:              req.Name,
		Uuid:              req.Uuid,
		Vcpus:             req.Vcpus,
		MemoryBytes:       req.MemoryBytes,
		State:             VirtualMachineAgentState_RUNNING,
		Blockdev:          req.Blockdev,
		Netdev:            req.Netdev,
		LoginUsername:     req.LoginUsername,
		SshAuthorizedKeys: req.SshAuthorizedKeys,
		WebsocketPort:     10000,
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
