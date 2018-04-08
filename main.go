package main

import (
	"fmt"
	"net"

	"github.com/n0stack/proto.go/node/v0"

	"github.com/kelseyhightower/envconfig"
	"github.com/n0stack/proto.go/kvm/v0"

	"github.com/n0stack/n0core/kvm"
	"github.com/n0stack/n0core/node"
	"github.com/n0stack/n0core/qcow2"
	"github.com/n0stack/n0core/tap"
	pqcow2 "github.com/n0stack/proto.go/qcow2/v0"
	ptap "github.com/n0stack/proto.go/tap/v0"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Environment struct {
	Starter               string `default:""`
	BindAddr              string `default:"" envconfig:"BIND_ADDR"`
	BindPort              int    `default:"20180" envconfig:"BIND_PORT"`
	OutboundInterfaceName string `required:"true" envconfig:"OUTBOUND_INTERFACE_NAME"`
}

var (
	env Environment
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", env.BindAddr, env.BindPort))
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
	ptap.RegisterTapServiceServer(s, &tap.Agent{InterfaceName: env.OutboundInterfaceName})

	n, err := node.GetAgent(env.BindAddr, env.BindPort, env.Starter)
	if err != nil {
		panic(err)
	}
	pnode.RegisterNodeServiceServer(s, n)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func init() {
	// packageがインストールされているか確認する

	// env vars
	err := envconfig.Process("n0core", &env)
	if err != nil {
		panic("Failed to parse environment variables: " + err.Error())
	}
}
