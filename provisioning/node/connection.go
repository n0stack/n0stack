package node

import (
	"context"
	"fmt"

	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

type NodeConnections struct {
	NodeAPI pprovisioning.NodeServiceClient
}

func NewNodeConnections(api string) (*NodeConnections, error) {
	conn, err := grpc.Dial(api, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Fail to dial to node, err:%v", err.Error())
	}

	nc := &NodeConnections{
		NodeAPI: pprovisioning.NewNodeServiceClient(conn),
	}

	return nc, nil
}

// GetConnection return a connection to Node having name of arguments.
// TODO: Node is not existとNode is NotReadyの場合を考える必要がある
func (nc NodeConnections) GetConnection(name string) (*grpc.ClientConn, error) {
	n, err := nc.NodeAPI.GetNode(context.Background(), &pprovisioning.GetNodeRequest{Name: name})
	if err != nil {
		// if status.Code(err) == codes.NotFound {
		// 	return nil, nil
		// }

		return nil, err
	}

	if n.Status.State == pprovisioning.NodeStatus_NotReady {
		return nil, nil
	}

	// portはendpointから取る
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Spec.Address, 20181), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Fail to dial to node, err:%v", err.Error())
	}

	return conn, nil
}
