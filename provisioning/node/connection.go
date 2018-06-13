package node

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

type NodeConnections struct {
	NodeAPI pprovisioning.NodeServiceClient

	// conns map[string]*grpc.ClientConn
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
// Nodeがいない場合 (nil, nil)
// Node is NotReady (nil, nil)
func (nc NodeConnections) GetConnection(name string) (*grpc.ClientConn, error) {
	n, err := nc.NodeAPI.GetNode(context.Background(), &pprovisioning.GetNodeRequest{Name: name})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}

		return nil, err
	}

	// nodeがreadyか
	if n.Status.State == pprovisioning.NodeStatus_NotReady {
		return nil, nil
	}

	// Nodeのコネクションをmapでconnsを使って、効率化

	// portはendpointから取る
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Spec.Address, 20181), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Fail to dial to node, err:%v", err.Error())
	}

	return conn, nil
}
