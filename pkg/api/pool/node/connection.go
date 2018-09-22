package node

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	"github.com/n0stack/proto.go/pool/v0"
	"github.com/pkg/errors"
)

// TODO: APIを叩く回数を減らす
type NodeConnections struct {
	nodeAPI ppool.NodeServiceClient
}

func NewNodeConnections(api string) (*NodeConnections, error) {
	conn, err := grpc.Dial(api, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Fail to dial to node, err:%v", err.Error())
	}

	nc := &NodeConnections{
		nodeAPI: ppool.NewNodeServiceClient(conn),
	}

	return nc, nil
}

func (nc NodeConnections) IsExisting(nodeName string) (bool, error) {
	_, err := nc.nodeAPI.GetNode(context.Background(), &ppool.GetNodeRequest{Name: nodeName})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}

		return false, errors.Wrap(err, "Failed to get node from API")
	}

	return true, nil
}

// GetConnection return a connection to Node having name of arguments.
func (nc NodeConnections) GetConnection(nodeName string) (*grpc.ClientConn, error) {
	n, err := nc.nodeAPI.GetNode(context.Background(), &ppool.GetNodeRequest{Name: nodeName})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get node from API")
	}

	if n.Status.State == ppool.NodeStatus_NotReady {
		return nil, nil
	}

	// port を何かから取れるようにする
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Spec.Address, 20181), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "Fail to dial to node")
	}

	return conn, nil
}