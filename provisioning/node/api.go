package node

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/memberlist"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/datastore"
	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"
)

type NodeAPI struct {
	ds   datastore.Datastore
	list *memberlist.Memberlist
}

func CreateNodeAPI(ds datastore.Datastore, starter string) (*NodeAPI, error) {
	a := &NodeAPI{
		ds: ds,
	}

	c := memberlist.DefaultLANConfig()
	c.Events = &NodeAPIEventDelegate{ds: ds}
	// c.Name = a.id.String()

	var err error
	a.list, err = memberlist.Create(c)
	if err != nil {
		return nil, err
	}

	if starter != "" {
		_, err := a.list.Join([]string{starter})
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *NodeAPI) ListNodes(ctx context.Context, req *pprovisioning.ListNodesRequest) (res *pprovisioning.ListNodesResponse, errRes error) {
	res = &pprovisioning.ListNodesResponse{}
	f := func(s int) []proto.Message {
		res.Nodes = make([]*pprovisioning.Node, s)
		for i := range res.Nodes {
			res.Nodes[i] = &pprovisioning.Node{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Nodes {
			m[i] = v
		}

		return m
	}

	if err := a.ds.List(f); err != nil {
		errRes = grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
		return
	}
	if len(res.Nodes) == 0 {
		errRes = grpc.Errorf(codes.NotFound, "")
		return
	}

	return
}

func (a *NodeAPI) GetNode(ctx context.Context, req *pprovisioning.GetNodeRequest) (res *pprovisioning.Node, errRes error) {
	res = &pprovisioning.Node{}
	if err := a.ds.Get(req.Name, res); err != nil {
		errRes = grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}

	if res == nil {
		errRes = grpc.Errorf(codes.NotFound, "")
		return
	}

	return
}

func (a *NodeAPI) ApplyNode(ctx context.Context, req *pprovisioning.ApplyNodeRequest) (res *pprovisioning.Node, errRes error) {
	res = &pprovisioning.Node{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &pprovisioning.NodeStatus{},
	}

	res.Status.State = pprovisioning.NodeStatus_NotReady
	for _, m := range a.list.Members() {
		if m.Name == res.Metadata.Name {
			res.Status.State = pprovisioning.NodeStatus_Ready
		}
	}

	if err := a.ds.Apply(req.Metadata.Name, res); err != nil {
		errRes = grpc.Errorf(codes.Internal, "message:Failed to apply for db.\tgot:%v", err.Error())
		return
	}

	return
}

func (a *NodeAPI) DeleteNode(ctx context.Context, req *pprovisioning.DeleteNodeRequest) (res *empty.Empty, errRes error) {
	d, err := a.ds.Delete(req.Name)
	if err != nil {
		errRes = grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
		return
	}
	if d > 0 {
		errRes = grpc.Errorf(codes.NotFound, "")
		return
	}

	return
}
