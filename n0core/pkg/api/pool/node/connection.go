package node

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	ppool "n0st.ac/n0stack/n0proto.go/pool/v0"
	"github.com/pkg/errors"
)

func GetConnection(ctx context.Context, api ppool.NodeServiceClient, nodeName string) (*grpc.ClientConn, error) {
	n, err := api.GetNode(ctx, &ppool.GetNodeRequest{Name: nodeName})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get node from API")
	}

	if n.State != ppool.Node_READY {
		return nil, nil
	}

	// port を何かから取れるようにする
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Address, 20181), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "Fail to dial to node")
	}

	return conn, nil
}
