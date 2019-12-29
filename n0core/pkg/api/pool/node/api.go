package node

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"

	"github.com/golang/protobuf/ptypes/empty"
	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/datastore"
	"n0st.ac/n0stack/n0core/pkg/datastore/lock"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	pbudget "n0st.ac/n0stack/n0proto.go/budget/v0"
	ppool "n0st.ac/n0stack/n0proto.go/pool/v0"
)

type NodeAPI struct {
	dataStore datastore.Datastore
	// list      *memberlist.Memberlist
}

func CreateNodeAPI(ds datastore.Datastore) *NodeAPI {
	a := &NodeAPI{
		dataStore: ds.AddPrefix("node"),
	}

	return a
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
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	return res, nil
}

func (a NodeAPI) ApplyNode(ctx context.Context, req *ppool.ApplyNodeRequest) (*ppool.Node, error) {
	if req.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Name is required")
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	res := &ppool.Node{}
	if err := a.dataStore.Get(req.Name, res); err != nil && !datastore.IsNotFound(err) {
		return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	res.Name = req.Name
	res.Annotations = req.Annotations
	res.Labels = req.Labels

	if req.Address == "" {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, grpcutil.Errorf(codes.Internal, "Failed to get gRPC peer information from request")
		}

		switch addr := p.Addr.(type) {
		case *net.TCPAddr:
			res.Address = addr.IP.String()
		default:
			return nil, grpc.Errorf(codes.Internal, "Expected peer address is net.TCPAddr, but received %T", addr)
		}
	} else {
		res.Address = req.Address
	}

	res.IpmiAddress = req.IpmiAddress
	res.Serial = req.Serial
	res.CpuMilliCores = req.CpuMilliCores
	res.MemoryBytes = req.MemoryBytes
	res.StorageBytes = req.StorageBytes
	res.Datacenter = req.Datacenter
	res.AvailabilityZone = req.AvailabilityZone
	res.Cell = req.Cell
	res.Rack = req.Rack
	res.Unit = req.Unit

	// res.State = prev.State
	res.State = ppool.Node_READY

	// TODO: 何らかの仕組みで死活監視
	// res.Status.State = ppool.NodeStatus_NotReady
	// for _, m := range a.list.Members() {
	// 	if m.Name == res.Metadata.Name {
	// 		res.Status.State = ppool.NodeStatus_Ready
	// 	}
	// }

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a NodeAPI) DeleteNode(ctx context.Context, req *ppool.DeleteNodeRequest) (*empty.Empty, error) {
	if req.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set name")
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	node := &ppool.Node{}
	if err := a.dataStore.Get(req.Name, node); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	if IsLockedForDeletion(node) {
		return nil, grpcutil.Errorf(codes.FailedPrecondition, "Node has some computes or storages, so is locked for deletion")
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to delete data from db: name='%s', err='%s'", req.Name, err.Error())
	}

	return &empty.Empty{}, nil
}

func (a NodeAPI) ScheduleCompute(ctx context.Context, req *ppool.ScheduleComputeRequest) (*ppool.Node, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a NodeAPI) ReserveCompute(ctx context.Context, req *ppool.ReserveComputeRequest) (*ppool.Node, error) {
	if req.ComputeName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set field 'compute_name'")
	}
	if req.RequestCpuMilliCore == 0 || req.RequestMemoryBytes == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set field 'request_*'")
	}

	if !lock.WaitUntilLock(a.dataStore, req.NodeName, 5*time.Second, 10*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.NodeName)

	res := &ppool.Node{}
	if err := a.dataStore.Get(req.NodeName, res); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, "Failed to get data from db: err='%s'", err.Error())
	}

	if res.ReservedComputes == nil {
		res.ReservedComputes = make(map[string]*pbudget.Compute)
	}
	if _, ok := res.ReservedComputes[req.ComputeName]; ok {
		return nil, grpc.Errorf(codes.AlreadyExists, "Compute '%s' is already exists on node '%s'", req.ComputeName, req.NodeName)
	}

	if err := CheckCompute(req.RequestCpuMilliCore, res.CpuMilliCores, req.RequestMemoryBytes, res.MemoryBytes, res.ReservedComputes); err != nil {
		return nil, grpc.Errorf(codes.ResourceExhausted, "Compute resource is exhausted on node '%s': %s", req.NodeName, err.Error())
	}

	res.ReservedComputes[req.ComputeName] = &pbudget.Compute{
		Annotations:         req.Annotations,
		RequestCpuMilliCore: req.RequestCpuMilliCore,
		LimitCpuMilliCore:   req.LimitCpuMilliCore,
		RequestMemoryBytes:  req.RequestMemoryBytes,
		LimitMemoryBytes:    req.LimitMemoryBytes,
	}
	if err := a.dataStore.Apply(req.NodeName, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.NodeName)
	}

	return res, nil
}

func (a NodeAPI) ReleaseCompute(ctx context.Context, req *ppool.ReleaseComputeRequest) (*empty.Empty, error) {
	if !lock.WaitUntilLock(a.dataStore, req.NodeName, 5*time.Second, 10*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.NodeName)

	n := &ppool.Node{}
	if err := a.dataStore.Get(req.NodeName, n); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, "Failed to get data from db: err='%s'", err.Error())
	}

	if _, ok := n.ReservedComputes[req.ComputeName]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Compute '%s' is not found on node '%s'", req.ComputeName, req.NodeName)
	}

	delete(n.ReservedComputes, req.ComputeName)
	if err := a.dataStore.Apply(req.NodeName, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.NodeName)
	}

	return &empty.Empty{}, nil
}

func (a NodeAPI) ScheduleStorage(ctx context.Context, req *ppool.ScheduleStorageRequest) (*ppool.Node, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a NodeAPI) ReserveStorage(ctx context.Context, req *ppool.ReserveStorageRequest) (*ppool.Node, error) {
	if req.StorageName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set field 'storage_name'")
	}
	if req.RequestBytes == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set field 'request_*'")
	}

	if !lock.WaitUntilLock(a.dataStore, req.NodeName, 5*time.Second, 50*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.NodeName)

	res := &ppool.Node{}
	if err := a.dataStore.Get(req.NodeName, res); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, "Failed to get data from db: err='%s'", err.Error())
	}

	if res.ReservedStorages == nil {
		res.ReservedStorages = make(map[string]*pbudget.Storage)
	}
	if _, ok := res.ReservedStorages[req.StorageName]; ok {
		return nil, grpc.Errorf(codes.AlreadyExists, "Storage '%s' is already exists on node '%s'", req.StorageName, req.NodeName)
	}

	if err := CheckStorage(req.RequestBytes, res.StorageBytes, res.ReservedStorages); err != nil {
		return nil, grpc.Errorf(codes.ResourceExhausted, "Storage resource is exhausted on node '%s': %s", req.NodeName, err.Error())
	}

	res.ReservedStorages[req.StorageName] = &pbudget.Storage{
		Annotations:  req.Annotations,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}
	if err := a.dataStore.Apply(req.NodeName, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.NodeName)
	}

	return res, nil
}

func (a NodeAPI) ReleaseStorage(ctx context.Context, req *ppool.ReleaseStorageRequest) (*empty.Empty, error) {
	if !lock.WaitUntilLock(a.dataStore, req.NodeName, 5*time.Second, 50*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.NodeName)

	n := &ppool.Node{}
	if err := a.dataStore.Get(req.NodeName, n); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, "Failed to get data from db: err='%s'", err.Error())
	}

	if _, ok := n.ReservedStorages[req.StorageName]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Storage '%s' is not found on node '%s'", req.StorageName, req.NodeName)
	}

	delete(n.ReservedStorages, req.StorageName)
	if err := a.dataStore.Apply(req.NodeName, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.NodeName)
	}

	return &empty.Empty{}, nil
}
