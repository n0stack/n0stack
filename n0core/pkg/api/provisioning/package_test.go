package provisioning

import (
	"fmt"
	"net"

	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"google.golang.org/grpc"
)

// どうもサービスが始まるまでのタイムラグがあるせいで、性能の悪いデバイスでは安定性が悪い
// TODO: 上位層が使いにくくなっているので変える
func UpMockAgent(address string) error {
	addr := fmt.Sprintf("%s:%d", address, 20181)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	RegisterBlockStorageAgentServiceServer(grpcServer, &MockBlockStorageAgentAPI{})
	RegisterVirtualMachineAgentServiceServer(grpcServer, &MockVirtualMachineAgentAPI{})
	return grpcServer.Serve(lis)
}

func init() {
	go UpMockAgent(node.MockNodeIP)
}
