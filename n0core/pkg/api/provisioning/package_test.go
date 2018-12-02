package provisioning

import (
	"fmt"
	"net"

	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"google.golang.org/grpc"
)

func UpMockAgent(address string) error {
	addr := fmt.Sprintf("%s:%d", address, 20181)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	RegisterBlockStorageAgentServiceServer(grpcServer, &MockBlockStorageAgentAPI{})
	return grpcServer.Serve(lis)
}

func init() {
	go UpMockAgent(node.MockNodeIP)
}
