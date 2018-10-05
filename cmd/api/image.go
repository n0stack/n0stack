package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/n0stack/proto.go/deployment/v0"
	"github.com/n0stack/proto.go/provisioning/v0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/n0stack/n0core/pkg/api/deployment/image"
	"github.com/n0stack/n0core/pkg/datastore/etcd"

	"github.com/urfave/cli"
)

func ServeImageAPI(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}

	ve, err := etcd.NewEtcdDatastore("virtual_machine", strings.Split(ctx.String("etcd-endpoints"), ","))
	if err != nil {
		return err
	}
	defer ve.Close()

	bsConn, err := grpc.Dial(ctx.String("block-storage-api-endpoint"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dial:", err)
	}
	defer bsConn.Close()
	bsc := pprovisioning.NewBlockStorageServiceClient(bsConn)

	ia, err := image.CreateImageAPI(ve, bsc)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pdeployment.RegisterImageServiceServer(s, ia)
	reflection.Register(s)

	log.Printf("[INFO] Started API")
	return s.Serve(lis)
}
