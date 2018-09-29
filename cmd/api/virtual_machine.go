package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/n0stack/n0core/pkg/api/provisioning"
	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/provisioning/v0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/n0stack/n0core/pkg/datastore/etcd"

	"github.com/urfave/cli"
)

func ServeVirtualMachineAPI(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}

	ve, err := etcd.NewEtcdDatastore("virtual_machine", strings.Split(ctx.String("etcd-endpoints"), ","))
	if err != nil {
		return err
	}
	defer ve.Close()

	nodeConn, err := grpc.Dial(ctx.String("node-api-endpoint"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dial:", err)
	}
	defer nodeConn.Close()
	noc := ppool.NewNodeServiceClient(nodeConn)

	networkConn, err := grpc.Dial(ctx.String("network-api-endpoint"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dial:", err)
	}
	defer networkConn.Close()
	nec := ppool.NewNetworkServiceClient(networkConn)

	bsConn, err := grpc.Dial(ctx.String("block-storage-api-endpoint"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dial:", err)
	}
	defer bsConn.Close()
	vc := pprovisioning.NewBlockStorageServiceClient(bsConn)

	va, err := provisioning.CreateVirtualMachineAPI(ve, noc, nec, vc)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pprovisioning.RegisterVirtualMachineServiceServer(s, va)
	reflection.Register(s)

	log.Printf("[INFO] Started API")
	return s.Serve(lis)
}
