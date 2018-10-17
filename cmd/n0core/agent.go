package main

import (
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0core/pkg/api/provisioning"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ServeAgent(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}

	bvm := filepath.Join(ctx.String("base-directory"), "virtual_machine")
	vma, err := provisioning.CreateVirtualMachineAgentAPI(bvm)
	if err != nil {
		return err
	}

	bv := filepath.Join(ctx.String("base-directory"), "block_storage")
	va, err := provisioning.CreateBlockStorageAgentAPI(bv)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	provisioning.RegisterBlockStorageAgentServiceServer(grpcServer, va)
	provisioning.RegisterVirtualMachineAgentServiceServer(grpcServer, vma)
	reflection.Register(grpcServer)

	if err := node.RegisterNodeToAPI(ctx.String("name"), ctx.String("advertise-address"), ctx.String("node-api-endpoint")); err != nil {
		return err
	}

	log.Printf("[INFO] Started API")
	return grpcServer.Serve(lis)
}
