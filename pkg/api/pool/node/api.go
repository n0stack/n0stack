package node

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/proto.go/budget/v0"
	"github.com/n0stack/proto.go/pool/v0"
)

type NodeAPI struct {
	dataStore datastore.Datastore
	// list      *memberlist.Memberlist
}

func CreateNodeAPI(ds datastore.Datastore) (*NodeAPI, error) {
	a := &NodeAPI{
		dataStore: ds,
	}

	// c := memberlist.DefaultLANConfig()
	// c.Events = &NodeAPIEventDelegate{ds: ds}
	// // c.Name = a.id.String()

	// var err error
	// a.list, err = memberlist.Create(c)
	// if err != nil {
	// 	return nil, err
	// }

	// if starter != "" {
	// 	_, err := a.list.Join([]string{starter})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return a, nil
}

func (a NodeAPI) ListNodes(ctx context.Context, req *ppool.ListNodesRequest) (*ppool.ListNodesResponse, error) {
	res := &ppool.ListNodesResponse{}
	f := func(s int) []proto.Message {
		res.Nodes = make([]*ppool.Node, s)
		for i := range res.Nodes {
			res.Nodes[i] = &ppool.Node{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Nodes {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		log.Printf("[WARNING] Failed to list data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.Nodes) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NodeAPI) GetNode(ctx context.Context, req *ppool.GetNodeRequest) (*ppool.Node, error) {
	res := &ppool.Node{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res == nil {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NodeAPI) ApplyNode(ctx context.Context, req *ppool.ApplyNodeRequest) (*ppool.Node, error) {
	res := &ppool.Node{
		Metadata: req.Metadata,
		Spec:     req.Spec,
	}

	prev := &ppool.Node{}
	if err := a.dataStore.Get(req.Metadata.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	}

	var err error
	res.Metadata.Version, err = datastore.CheckVersion(prev, req)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to check version, detail='%s'", err.Error())
	}

	res.Status = prev.Status
	if res.Status == nil { // version == 0 でもあるはず
		res.Status = &ppool.NodeStatus{}

		// TODO: 何らかの仕組みで死活監視
		res.Status.State = ppool.NodeStatus_Ready
		// res.Status.State = ppool.NodeStatus_NotReady
		// for _, m := range a.list.Members() {
		// 	if m.Name == res.Metadata.Name {
		// 		res.Status.State = ppool.NodeStatus_Ready
		// 	}
		// }
	}

	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	}

	return res, nil
}

func (a NodeAPI) DeleteNode(ctx context.Context, req *ppool.DeleteNodeRequest) (*empty.Empty, error) {
	if err := a.dataStore.Delete(req.Name); err != nil {
		log.Printf("[WARNING] Failed to delete data from db: err='%s'", err.Error())
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "Failed to delete '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return &empty.Empty{}, nil
}

func (a NodeAPI) ScheduleCompute(ctx context.Context, req *ppool.ScheduleComputeRequest) (*ppool.ReserveComputeResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a NodeAPI) ReserveCompute(ctx context.Context, req *ppool.ReserveComputeRequest) (*ppool.ReserveComputeResponse, error) {
	if req.ComputeName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Do not set field 'compute_name' as blank")
	}

	n := &ppool.Node{}
	if err := a.dataStore.Get(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if n == nil {
		return nil, grpc.Errorf(codes.NotFound, "Node '%s' is not found", req.Name)
	}
	if _, ok := n.Status.ReservedComputes[req.ComputeName]; ok {
		return nil, grpc.Errorf(codes.AlreadyExists, "Compute '%s' is already exists on node '%s'", req.ComputeName, req.Name)
	}

	if err := CheckCompute(req.Compute.RequestCpuMilliCore, n.Spec.CpuMilliCores, req.Compute.RequestMemoryBytes, n.Spec.MemoryBytes, n.Status.ReservedComputes); err != nil {
		grpc.Errorf(codes.ResourceExhausted, "Compute resource is exhausted on node '%s': %s", req.Name, err.Error())
	}

	res := &ppool.ReserveComputeResponse{
		Name:        req.Name,
		ComputeName: req.ComputeName,
		Compute: &pbudget.Compute{
			Annotations:         req.Compute.Annotations,
			RequestCpuMilliCore: req.Compute.RequestCpuMilliCore,
			LimitCpuMilliCore:   req.Compute.LimitCpuMilliCore,
			RequestMemoryBytes:  req.Compute.RequestMemoryBytes,
			LimitMemoryBytes:    req.Compute.LimitMemoryBytes,
		},
	}
	n.Status.ReservedComputes[req.ComputeName] = res.Compute
	if err := a.dataStore.Apply(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a NodeAPI) ReleaseCompute(ctx context.Context, req *ppool.ReleaseComputeRequest) (*empty.Empty, error) {
	n := &ppool.Node{}
	if err := a.dataStore.Get(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if n == nil {
		return nil, grpc.Errorf(codes.NotFound, "Node '%s' is not found", req.Name)
	}

	if _, ok := n.Status.ReservedComputes[req.ComputeName]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Compute '%s' is not found on node '%s'", req.ComputeName, req.Name)
	}

	delete(n.Status.ReservedComputes, req.ComputeName)
	if err := a.dataStore.Apply(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return &empty.Empty{}, nil
}

func (a NodeAPI) ScheduleStorage(ctx context.Context, req *ppool.ScheduleStorageRequest) (*ppool.ReserveStorageResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a NodeAPI) ReserveStorage(ctx context.Context, req *ppool.ReserveStorageRequest) (*ppool.ReserveStorageResponse, error) {
	if req.StorageName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Do not set field 'storage_name' as blank")
	}

	n := &ppool.Node{}
	if err := a.dataStore.Get(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if n == nil {
		return nil, grpc.Errorf(codes.NotFound, "Node '%s' is not found", req.Name)
	}
	if _, ok := n.Status.ReservedStorages[req.StorageName]; ok {
		return nil, grpc.Errorf(codes.AlreadyExists, "Storage '%s' is already exists on node '%s'", req.StorageName, req.Name)
	}

	if err := CheckStorage(req.Storage.RequestBytes, n.Spec.StorageBytes, n.Status.ReservedStorages); err != nil {
		grpc.Errorf(codes.ResourceExhausted, "Storage resource is exhausted on node '%s': %s", req.Name, err.Error())
	}

	res := &ppool.ReserveStorageResponse{
		Name:        req.Name,
		StorageName: req.StorageName,
		Storage: &pbudget.Storage{
			Annotations:  req.Storage.Annotations,
			RequestBytes: req.Storage.RequestBytes,
			LimitBytes:   req.Storage.LimitBytes,
		},
	}
	n.Status.ReservedStorages[req.StorageName] = res.Storage
	if err := a.dataStore.Apply(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a NodeAPI) ReleaseStorage(ctx context.Context, req *ppool.ReleaseStorageRequest) (*empty.Empty, error) {
	n := &ppool.Node{}
	if err := a.dataStore.Get(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if n == nil {
		return nil, grpc.Errorf(codes.NotFound, "Node '%s' is not found", req.Name)
	}

	if _, ok := n.Status.ReservedStorages[req.StorageName]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Storage '%s' is not found on node '%s'", req.StorageName, req.Name)
	}

	delete(n.Status.ReservedStorages, req.StorageName)
	if err := a.dataStore.Apply(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return &empty.Empty{}, nil
}
