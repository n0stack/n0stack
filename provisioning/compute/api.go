package compute

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/datastore"
	"github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ComputeAPI struct {
	dataStore datastore.Datastore
}

func CreateComputeAPI(ds datastore.Datastore) (*ComputeAPI, error) {
	return &ComputeAPI{
		dataStore: ds,
	}, nil
}

func (a ComputeAPI) ListComputes(ctx context.Context, req *pprovisioning.ListComputesRequest) (*pprovisioning.ListComputesResponse, error) {
	res := &pprovisioning.ListComputesResponse{}
	f := func(s int) []proto.Message {
		res.Computes = make([]*pprovisioning.Compute, s)
		for i := range res.Computes {
			res.Computes[i] = &pprovisioning.Compute{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Computes {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}
	if len(res.Computes) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a ComputeAPI) GetCompute(ctx context.Context, req *pprovisioning.GetComputeRequest) (*pprovisioning.Compute, error) {
	res := &pprovisioning.Compute{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}

	if res.Metadata == nil {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a ComputeAPI) ApplyCompute(ctx context.Context, req *pprovisioning.ApplyComputeRequest) (*pprovisioning.Compute, error) {
	prev := &pprovisioning.Compute{}
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

	res := &pprovisioning.Compute{
		Metadata: req.Metadata,
		Spec:     req.Spec,
	}
	if prev.Status == nil {
		res.Status = &pprovisioning.ComputeStatus{}
	} else {
		res.Status = prev.Status
	}

	res.Metadata.Version++

	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to apply for db, got:%v.", err.Error())
	}
	log.Printf("[INFO] On Applly, applied Network:%v", res)

	return res, nil
}

func (a ComputeAPI) DeleteCompute(ctx context.Context, req *pprovisioning.DeleteComputeRequest) (*empty.Empty, error) {
	d, err := a.dataStore.Delete(req.Name)
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}
	if d < 1 {
		return &empty.Empty{}, grpc.Errorf(codes.NotFound, "")
	}

	return &empty.Empty{}, nil
}

func (a ComputeAPI) WatchCompute(req *pprovisioning.WatchComputesRequest, res pprovisioning.ComputeService_WatchComputeServer) error {
	return nil
}

func (a ComputeAPI) Boot(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, nil
}

func (a ComputeAPI) Reboot(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, nil
}

func (a ComputeAPI) HardReboot(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, nil
}

func (a ComputeAPI) Shutdown(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, nil
}

func (a ComputeAPI) HardShutdown(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, nil
}

func (a ComputeAPI) Save(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, nil
}
