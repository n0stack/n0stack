package main

import (
	"fmt"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ServeMockAgent(ctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}

	// とりあえず log を表示するため利用する
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zapLogger),
			// grpc_auth.StreamServerInterceptor(auth),
			// grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapLogger),
			// grpc_auth.UnaryServerInterceptor(auth),
			// grpc_prometheus.UnaryServerInterceptor,
		)))
	provisioning.RegisterBlockStorageAgentServiceServer(grpcServer, &provisioning.MockBlockStorageAgentAPI{})
	provisioning.RegisterVirtualMachineAgentServiceServer(grpcServer, &provisioning.MockVirtualMachineAgentAPI{})
	reflection.Register(grpcServer)

	if err := node.RegisterNodeToAPI(ctx.String("name"), ctx.String("advertise-address"), ctx.String("node-api-endpoint")); err != nil {
		return err
	}

	log.Printf("[INFO] Started API")
	return grpcServer.Serve(lis)
}
