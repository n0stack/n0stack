package flavor

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0proto.go/deployment/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type FlavorAPI struct {
	dataStore datastore.Datastore
	vmAPI     pprovisioning.VirtualMachineServiceClient
	imageAPI  pdeployment.ImageServiceClient
}

func CreateFlavorAPI(ds datastore.Datastore, vma pprovisioning.VirtualMachineServiceClient, ia pdeployment.ImageServiceClient) *FlavorAPI {
	a := &FlavorAPI{
		dataStore: ds.AddPrefix("flavor"),
		vmAPI:     vma,
		imageAPI:  ia,
	}

	return a
}

func (a FlavorAPI) ListFlavors(ctx context.Context, req *pdeployment.ListFlavorsRequest) (*pdeployment.ListFlavorsResponse, error) {
	res := &pdeployment.ListFlavorsResponse{}
	f := func(s int) []proto.Message {
		res.Flavors = make([]*pdeployment.Flavor, s)
		for i := range res.Flavors {
			res.Flavors[i] = &pdeployment.Flavor{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Flavors {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		log.Printf("[WARNING] Failed to list data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.Flavors) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a FlavorAPI) GetFlavor(ctx context.Context, req *pdeployment.GetFlavorRequest) (*pdeployment.Flavor, error) {
	res := &pdeployment.Flavor{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a FlavorAPI) ApplyFlavor(ctx context.Context, req *pdeployment.ApplyFlavorRequest) (*pdeployment.Flavor, error) {
	res := &pdeployment.Flavor{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	res.Name = req.Name
	res.Annotations = req.Annotations
	res.LimitCpuMilliCore = req.LimitCpuMilliCore
	res.LimitMemoryBytes = req.LimitMemoryBytes
	res.NetworkName = req.NetworkName

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a FlavorAPI) DeleteFlavor(ctx context.Context, req *pdeployment.DeleteFlavorRequest) (*empty.Empty, error) {
	prev := &pdeployment.Flavor{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if prev.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Flavor '%s' is not found", req.Name)
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		log.Printf("[WARNING] Failed to delete data from db: err='%s'", err.Error())
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "Failed to delete '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return &empty.Empty{}, nil
}

// TODO: 異常系の処理
func (a FlavorAPI) GenerateVirtualMachine(ctx context.Context, req *pdeployment.GenerateVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	prev := &pdeployment.Flavor{}
	if err := a.dataStore.Get(req.FlavorName, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.VirtualMachineName)
	}
	if prev.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Flavor '%s' is not found", req.VirtualMachineName)
	}

	bsName := fmt.Sprintf("%s-%s", req.VirtualMachineName, req.ImageName)
	_, err := a.imageAPI.GenerateBlockStorage(context.Background(), &pdeployment.GenerateBlockStorageRequest{
		ImageName:        req.ImageName,
		Tag:              req.ImageTag,
		BlockStorageName: bsName,
		LimitBytes:       prev.LimitStorageBytes,
		RequestBytes:     req.RequestStorageBytes,
	})
	if err != nil {
		log.Printf("Failed to generate blocksotrage: err='%s'", err.Error())
		return nil, err
		// return nil, grpc.Errorf(codes.Internal, "")
	}

	res, err := a.vmAPI.CreateVirtualMachine(context.Background(), &pprovisioning.CreateVirtualMachineRequest{
		Name:                req.VirtualMachineName,
		Annotations:         req.Annotations,
		LimitCpuMilliCore:   prev.LimitCpuMilliCore,
		RequestCpuMilliCore: req.RequestCpuMilliCore,
		LimitMemoryBytes:    prev.LimitMemoryBytes,
		RequestMemoryBytes:  req.RequestMemoryBytes,
		BlockStorageNames:   []string{bsName},
		Nics: []*pprovisioning.VirtualMachineNIC{
			{
				NetworkName: prev.NetworkName,
			},
		},
	})
	if err != nil {
		log.Printf("Failed to create virtual machinesotrage: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "")
	}

	return res, nil
}
