package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0core/pkg/api/provisioning"
	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockVirtualMachineAgentAPI struct{}

func (a *MockVirtualMachineAgentAPI) CreateVirtualMachineAgent(ctx context.Context, req *provisioning.CreateVirtualMachineAgentRequest) (*provisioning.VirtualMachineAgent, error) {
	return &provisioning.VirtualMachineAgent{
		Name:          req.Name,
		Uuid:          uuid.NewV4().String(),
		Vcpus:         req.Vcpus,
		MemoryBytes:   req.MemoryBytes,
		State:         provisioning.VirtualMachineAgentState_RUNNING,
		Blockdev:      req.Blockdev,
		Netdev:        req.Netdev,
		WebsocketPort: 10000,
	}, nil
}
func (a *MockVirtualMachineAgentAPI) DeleteVirtualMachineAgent(ctx context.Context, req *provisioning.DeleteVirtualMachineAgentRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
func (a *MockVirtualMachineAgentAPI) BootVirtualMachineAgent(ctx context.Context, req *provisioning.BootVirtualMachineAgentRequest) (*provisioning.BootVirtualMachineAgentResponse, error) {
	return &provisioning.BootVirtualMachineAgentResponse{
		State: provisioning.VirtualMachineAgentState_RUNNING,
	}, nil
}
func (a *MockVirtualMachineAgentAPI) RebootVirtualMachineAgent(ctx context.Context, req *provisioning.RebootVirtualMachineAgentRequest) (*provisioning.RebootVirtualMachineAgentResponse, error) {
	return &provisioning.RebootVirtualMachineAgentResponse{
		State: provisioning.VirtualMachineAgentState_RUNNING,
	}, nil
}
func (a *MockVirtualMachineAgentAPI) ShutdownVirtualMachineAgent(ctx context.Context, req *provisioning.ShutdownVirtualMachineAgentRequest) (*provisioning.ShutdownVirtualMachineAgentResponse, error) {
	return &provisioning.ShutdownVirtualMachineAgentResponse{
		State: provisioning.VirtualMachineAgentState_SHUTDOWN,
	}, nil
}

type MockBlockStorageAgentAPI struct{}

func (a MockBlockStorageAgentAPI) CreateEmptyBlockStorageAgent(ctx context.Context, req *provisioning.CreateEmptyBlockStorageAgentRequest) (*provisioning.BlockStorageAgent, error) {
	return &provisioning.BlockStorageAgent{
		Name:  req.Name,
		Bytes: req.Bytes,
		Path:  "/tmp/test",
	}, nil
}
func (a MockBlockStorageAgentAPI) CreateBlockStorageAgentWithDownloading(ctx context.Context, req *provisioning.CreateBlockStorageAgentWithDownloadingRequest) (*provisioning.BlockStorageAgent, error) {
	return &provisioning.BlockStorageAgent{
		Name:  req.Name,
		Bytes: req.Bytes,
		Path:  "/tmp/test",
	}, nil
}
func (a MockBlockStorageAgentAPI) DeleteBlockStorageAgent(ctx context.Context, req *provisioning.DeleteBlockStorageAgentRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func Mock(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	provisioning.RegisterBlockStorageAgentServiceServer(s, &MockBlockStorageAgentAPI{})
	provisioning.RegisterVirtualMachineAgentServiceServer(s, &MockVirtualMachineAgentAPI{})
	reflection.Register(s)

	if err := node.RegisterNodeToAPI(ctx.String("name"), ctx.String("advertise-address"), ctx.String("node-api-endpoint")); err != nil {
		return err
	}

	log.Printf("[INFO] Started API")
	return s.Serve(lis)
}
