package main

import (
	"net"

	"github.com/kelseyhightower/envconfig"
	"github.com/n0stack/go.proto/kvm/v0"

	pqcow2 "github.com/n0stack/go.proto/qcow2/v0"
	ptap "github.com/n0stack/go.proto/tap/v0"
	"github.com/n0stack/n0core/kvm"
	"github.com/n0stack/n0core/qcow2"
	"github.com/n0stack/n0core/tap"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Environment struct {
	Starter               string `default:""`
	BindPort              string `default:":20180"`
	OutboundInterfaceName string `required:"true"`
}

var (
	Env Environment
)

func main() {
	lis, err := net.Listen("tcp", Env.BindPort)
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
	ptap.RegisterTapServiceServer(s, &tap.Agent{InterfaceName: Env.OutboundInterfaceName})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func init() {
	// packageがインストールされているか確認する

	// env vars
	err := envconfig.Process("n0core", &Env)
	if err != nil {
		panic("Failed to parse environment variables: " + err.Error())
	}
}
