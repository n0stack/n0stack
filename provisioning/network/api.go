package network

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/datastore"
	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"
)

type NetworkAPI struct {
	dataStore datastore.Datastore
}

func CreateNetworkAPI(ds datastore.Datastore) (*NetworkAPI, error) {
	a := &NetworkAPI{
		dataStore: ds,
	}

	return a, nil
}

func (a NetworkAPI) ListNetworks(ctx context.Context, req *pprovisioning.ListNetworksRequest) (*pprovisioning.ListNetworksResponse, error) {
	res := &pprovisioning.ListNetworksResponse{}
	f := func(s int) []proto.Message {
		res.Networks = make([]*pprovisioning.Network, s)
		for i := range res.Networks {
			res.Networks[i] = &pprovisioning.Network{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Networks {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}
	if len(res.Networks) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NetworkAPI) GetNetwork(ctx context.Context, req *pprovisioning.GetNetworkRequest) (*pprovisioning.Network, error) {
	res := &pprovisioning.Network{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}

	if res.Metadata == nil {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NetworkAPI) ApplyNetwork(ctx context.Context, req *pprovisioning.ApplyNetworkRequest) (*pprovisioning.Network, error) {
	res := &pprovisioning.Network{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &pprovisioning.NetworkStatus{},
	}

	prev := &pprovisioning.Network{}
	err := a.dataStore.Get(req.Metadata.Name, prev)
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

	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to apply for db, got:%v.", err.Error())
	}
	log.Printf("[INFO] On Applly, applied Network:%v", res)

	return res, nil
}

func (a NetworkAPI) DeleteNetwork(ctx context.Context, req *pprovisioning.DeleteNetworkRequest) (*empty.Empty, error) {
	d, err := a.dataStore.Delete(req.Name)
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}
	if d < 1 {
		return &empty.Empty{}, grpc.Errorf(codes.NotFound, "")
	}

	return &empty.Empty{}, nil
}
