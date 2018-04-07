package main

import (
	"net"

	tap "github.com/n0stack/go-proto/device/tap/v0"
	"github.com/n0stack/go-proto/device/volume"

	_ "github.com/mattn/go-sqlite3"
	"github.com/n0stack/go-proto/device/vm"
	"github.com/n0stack/n0core/device/tap/flat"
	"github.com/n0stack/n0core/device/vm/kvm"
	"github.com/n0stack/n0core/device/volume/qcow2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":20180")
	if err != nil {
		panic(err)
	}

	// s := grpc.NewServer(
	// 	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	// 		grpc_recovery.StreamServerInterceptor(),
	// 	)),
	// 	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	// 		grpc_recovery.UnaryServerInterceptor(),
	// 	)),
	// )
	s := grpc.NewServer()

	vm.RegisterStandardServer(s, &kvm.Agent{})
	volume.RegisterStandardServer(s, &qcow2.Agent{})
	tap.RegisterStandardServer(s, &flat.Agent{InterfaceName: "enp0s25"}) // 環境変数から取る

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
