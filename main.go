package main

import (
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/n0stack/n0core/kvm"
	"github.com/n0stack/n0core/node"
	"github.com/n0stack/n0core/qcow2"
	"github.com/n0stack/n0core/tap"
	pkvm "github.com/n0stack/proto.go/kvm/v0"
	pnode "github.com/n0stack/proto.go/node/v0"
	pqcow2 "github.com/n0stack/proto.go/qcow2/v0"
	ptap "github.com/n0stack/proto.go/tap/v0"
)

type Environment struct {
	Starter               string `default:""`
	BindAddr              string `default:"" envconfig:"BIND_ADDR"`
	BindPort              int    `default:"20180" envconfig:"BIND_PORT"`
	VNCPortMin            uint   `default:"5900" envconfig:"VNC_PORT_MIN"`
	VNCPortMax            uint   `default:"5999" envconfig:"VNC_PORT_MAX"`
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

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(),
		),
	)

	pkvm.RegisterKVMServiceServer(s, &kvm.Agent{VNCPortMin: env.VNCPortMin, VNCPortMax: env.VNCPortMax})
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
