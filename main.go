package main

import (
	"net"

	"github.com/n0stack/go.proto/kvm/v0"

	pqcow2 "github.com/n0stack/go.proto/qcow2/v0"
	ptap "github.com/n0stack/go.proto/tap/v0"
	"github.com/n0stack/n0core/kvm"
	"github.com/n0stack/n0core/tap"
	"github.com/n0stack/n0core/volume/qcow2"

	_ "github.com/mattn/go-sqlite3"
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

	pkvm.RegisterKVMServiceServer(s, &kvm.Agent{})
	pqcow2.RegisterQcow2ServiceServer(s, &qcow2.Agent{})
	ptap.RegisterTapServiceServer(s, &tap.Agent{InterfaceName: "enp0s25"}) // 環境変数から取る

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
