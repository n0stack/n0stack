package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/n0stack/proto.go/pool/v0"

	"github.com/n0stack/n0core/pkg/api/provisioning"
	"github.com/n0stack/proto.go/provisioning/v0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/n0stack/n0core/pkg/datastore/etcd"

	"github.com/urfave/cli"
)

func ServeBlockStorageAPI(ctx *cli.Context) error {
	b := fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port"))
	lis, err := net.Listen("tcp", b)
	if err != nil {
		return err
	}

	ve, err := etcd.NewEtcdDatastore("block_storage", strings.Split(ctx.String("etcd-endpoints"), ","))
	if err != nil {
		return err
	}
	defer ve.Close()

	conn, err := grpc.Dial(ctx.String("node-api-endpoint"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dial:", err)
	}
	defer conn.Close()
	noc := ppool.NewNodeServiceClient(conn)

	va, err := provisioning.CreateBlockStorageAPI(ve, noc)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pprovisioning.RegisterBlockStorageServiceServer(s, va)
	reflection.Register(s)

	log.Printf("[INFO] Started API")
	return s.Serve(lis)
}
