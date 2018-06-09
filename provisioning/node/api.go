package node

import (
	"context"
	"log"

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

func (a NodeAPI) ListNodes(ctx context.Context, req *pprovisioning.ListNodesRequest) (*pprovisioning.ListNodesResponse, error) {
	res := &pprovisioning.ListNodesResponse{}
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
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}
	if len(res.Nodes) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NodeAPI) GetNode(ctx context.Context, req *pprovisioning.GetNodeRequest) (*pprovisioning.Node, error) {
	res := &pprovisioning.Node{}
	if err := a.ds.Get(req.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}

	if res.Metadata == nil {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NodeAPI) ApplyNode(ctx context.Context, req *pprovisioning.ApplyNodeRequest) (*pprovisioning.Node, error) {
	res := &pprovisioning.Node{
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

	prev := &pprovisioning.Node{}
	err := a.ds.Get(req.Metadata.Name, prev)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get db, got:%v.", err.Error())
	}
	if prev.Metadata == nil && req.Metadata.Version != 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set the same version as GetVolume result, have:%d, want:0.", req.Metadata.Version)
	}
	if prev.Metadata != nil && req.Metadata.Version != prev.Metadata.Version {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set the same version as GetVolume result, have:%d, want:%d.", req.Metadata.Version, prev.Metadata.Version)
	}

	res.Metadata.Version++

	if err := a.ds.Apply(req.Metadata.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to apply for db, got:%v.", err.Error())
	}
	log.Printf("[INFO] On Applly, applied Node:%v", res)

	return res, nil
}

func (a NodeAPI) DeleteNode(ctx context.Context, req *pprovisioning.DeleteNodeRequest) (*empty.Empty, error) {
	d, err := a.ds.Delete(req.Name)
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}
	if d < 1 {
		return &empty.Empty{}, grpc.Errorf(codes.NotFound, "")
	}

	return &empty.Empty{}, nil
}
