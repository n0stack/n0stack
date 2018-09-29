package main

import (
	"fmt"
	"log"
	"net"

	"github.com/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0core/pkg/api/provisioning"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Mock(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	provisioning.RegisterBlockStorageAgentServiceServer(s, &provisioning.MockBlockStorageAgentAPI{})
	provisioning.RegisterVirtualMachineAgentServiceServer(s, &provisioning.MockVirtualMachineAgentAPI{})
	reflection.Register(s)

	if err := node.RegisterNodeToAPI(ctx.String("name"), ctx.String("advertise-address"), ctx.String("node-api-endpoint")); err != nil {
		return err
	}

	log.Printf("[INFO] Started API")
	return s.Serve(lis)
}
