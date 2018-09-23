package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/n0stack/n0core/pkg/api/pool/node"

	"github.com/n0stack/proto.go/pool/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/n0stack/n0core/pkg/datastore/etcd"

	"github.com/urfave/cli"
)

func ServeNodeAPI(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}

	ne, err := etcd.NewEtcdDatastore("node", strings.Split(ctx.String("etcd-endpoints"), ","))
	if err != nil {
		return err
	}
	defer ne.Close()

	na, err := node.CreateNodeAPI(ne)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	ppool.RegisterNodeServiceServer(s, na)
	reflection.Register(s)

	log.Printf("[INFO] Started API")
	return s.Serve(lis)
}
