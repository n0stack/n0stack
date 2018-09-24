package provisioning

import (
	"context"
	"log"
	"reflect"
	"strconv"

	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/provisioning/v0"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0core/pkg/datastore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const AnnotationVNCWebSocketPort = "n0core/provisioning/virtual_machine_vnc_websocket_port"

type VirtualMachineAPI struct {
	dataStore datastore.Datastore

	// dependency APIs
	nodeAPI    ppool.NodeServiceClient
	networkAPI ppool.NetworkServiceClient
	volumeAPI  pprovisioning.VolumeServiceClient

	nodeConnections *node.NodeConnections
}

func CreateVirtualMachineAPI(ds datastore.Datastore, noa ppool.NodeServiceClient, nea ppool.NetworkServiceClient, va pprovisioning.VolumeServiceClient) (*VirtualMachineAPI, error) {
	nc := &node.NodeConnections{
		NodeAPI: noa,
	}

	a := &VirtualMachineAPI{
		dataStore:       ds,
		nodeAPI:         noa,
		networkAPI:      nea,
		volumeAPI:       va,
		nodeConnections: nc,
	}

	return a, nil
}

func (a *VirtualMachineAPI) CreateVirtualMachine(ctx context.Context, req *pprovisioning.ApplyVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	prev := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Metadata.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	} else if !reflect.ValueOf(prev.Metadata).IsNil() {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Metadata.Name)
	}

	if req.Spec.LimitCpuMilliCore%1000 != 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Make limit_cpu_milli_core '%d' a multiple of 1000", req.Spec.LimitCpuMilliCore)
	}

	res := &pprovisioning.VirtualMachine{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &pprovisioning.VirtualMachineStatus{},
	}
	var blockdev []*BlockDev
	if res.Metadata.Annotations == nil {
		res.Metadata.Annotations = make(map[string]string)
	}

	var err error
	res.Status.ComputeNodeName, res.Status.ComputeName, err = a.reserveCompute(
		req.Metadata.Name,
		req.Metadata.Annotations,
		req.Spec.RequestCpuMilliCore,
		req.Spec.LimitCpuMilliCore,
		req.Spec.RequestMemoryBytes,
		req.Spec.LimitMemoryBytes,
	)
	if err != nil {
		log.Printf("Failed to reserve compute: err=%v.", err.Error())
		return nil, err
	}

	// errorについて考える
	conn, err := a.nodeConnections.GetConnection(res.Status.ComputeNodeName)
	cli := NewVirtualMachineAgentServiceClient(conn)
	if err != nil {
		log.Printf("Failed to dial to node: err=%v.", err.Error())
		goto ReleaseCompute
	}
	if conn == nil {
		// TODO: goto ReleaseCompute
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.Status.ComputeNodeName)
	}
	defer conn.Close()

	if blockdev, err = a.reserveVolume(req.Spec.VolumeNames); err != nil {
		log.Printf("Failed to reserve volume: err=%v.", err.Error())
		goto ReleaseCompute
	}

	res.Spec.Nics, res.Status.NetworkInterfaceNames, err = a.reserveNics(req.Metadata.Name, req.Spec.Nics)
	if err != nil {
		log.Printf("Failed to reserve nics: err=%v.", err.Error())
		goto ReleaseVolume
	}

	if vm, err := cli.CreateVirtualMachineAgent(context.Background(), &CreateVirtualMachineAgentRequest{
		Name:        req.Metadata.Name,
		Vcpus:       req.Spec.LimitCpuMilliCore / 1000,
		MemoryBytes: req.Spec.LimitMemoryBytes,
		Netdev:      StructNetDev(req.Spec.Nics, res.Status.NetworkInterfaceNames),
		Blockdev:    blockdev,
	}); err != nil && status.Code(err) != codes.AlreadyExists {
		log.Printf("Failed to create volume on node '%s': err='%s'", res.Status.ComputeNodeName, err.Error()) // TODO: #89
		goto ReleaseNetworkInterface
	} else {
		res.Metadata.Annotations[AnnotationVNCWebSocketPort] = strconv.Itoa(int(vm.WebsocketPort))
		res.Status.State = GetAPIStateFromAgentState(vm.State)
		res.Status.Uuid = vm.Uuid
	}

	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteVirtualMachine
	}

	return res, nil

DeleteVirtualMachine:
	_, err = cli.DeleteVirtualMachineAgent(context.Background(), &DeleteVirtualMachineAgentRequest{
		Name:   req.Metadata.Name,
		Netdev: StructNetDev(res.Spec.Nics, res.Status.NetworkInterfaceNames),
	})
	if err != nil {
		log.Printf("Fail to delete virtual machine on node: err=%s.", err.Error())
	}

ReleaseNetworkInterface:
	if err := a.releaseNics(res.Spec.Nics, res.Status.NetworkInterfaceNames); err != nil {
		log.Printf("Fail to release network interfaces on API: err=%s.", err.Error())
	}

ReleaseVolume:
	if err := a.relaseVolumes(res.Spec.VolumeNames); err != nil {
		log.Printf("Fail to release volume on API: err=%s.", err.Error())
	}

ReleaseCompute:
	if err := a.releaseCompute(res.Status.ComputeNodeName, res.Status.ComputeName); err != nil {
		log.Printf("Fail to release compute on API: err=%s.", err.Error())
	}

	return nil, grpc.Errorf(codes.Internal, "")
}

func (a *VirtualMachineAPI) ListVirtualMachines(ctx context.Context, req *pprovisioning.ListVirtualMachinesRequest) (*pprovisioning.ListVirtualMachinesResponse, error) {
	res := &pprovisioning.ListVirtualMachinesResponse{}
	f := func(s int) []proto.Message {
		res.VirtualMachines = make([]*pprovisioning.VirtualMachine, s)
		for i := range res.VirtualMachines {
			res.VirtualMachines[i] = &pprovisioning.VirtualMachine{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.VirtualMachines {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		log.Printf("[WARNING] Failed to list data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.VirtualMachines) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *VirtualMachineAPI) GetVirtualMachine(ctx context.Context, req *pprovisioning.GetVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	res := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if reflect.ValueOf(res.Metadata).IsNil() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *VirtualMachineAPI) UpdateVirtualMachine(ctx context.Context, req *pprovisioning.UpdateVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a *VirtualMachineAPI) DeleteVirtualMachine(ctx context.Context, req *pprovisioning.DeleteVirtualMachineRequest) (*empty.Empty, error) {
	prev := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, prev); err == nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Name)
	} else if status.Code(err) != codes.NotFound {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	conn, err := a.nodeConnections.GetConnection(prev.Status.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.Status.ComputeNodeName)
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	_, err = cli.DeleteVirtualMachineAgent(context.Background(), &DeleteVirtualMachineAgentRequest{
		Name:   req.Name,
		Netdev: StructNetDev(prev.Spec.Nics, prev.Status.NetworkInterfaceNames),
	})
	if err != nil {
		log.Printf("Fail to delete virtual machine on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete virtual machine on node") // TODO #89
	}

	if err := a.releaseCompute(prev.Status.ComputeNodeName, prev.Status.ComputeName); err != nil {
		return nil, err
	}

	if err := a.relaseVolumes(prev.Spec.VolumeNames); err != nil {
		return nil, err
	}

	if err := a.releaseNics(prev.Spec.Nics, prev.Status.NetworkInterfaceNames); err != nil {
		return nil, err
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	return &empty.Empty{}, nil
}

func (a *VirtualMachineAPI) BootVirtualMachine(ctx context.Context, req *pprovisioning.BootVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	prev := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, prev); err == nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Name)
	} else if status.Code(err) != codes.NotFound {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	res := &pprovisioning.VirtualMachine{
		Metadata: prev.Metadata,
		Spec:     prev.Spec,
		Status:   prev.Status,
	}

	conn, err := a.nodeConnections.GetConnection(prev.Status.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.Status.ComputeNodeName)
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	vm, err := cli.BootVirtualMachineAgent(context.Background(), &BootVirtualMachineAgentRequest{Name: req.Name})
	if err != nil {
		log.Printf("Fail to boot on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to boot volume on node") // TODO #89
	}
	res.Status.State = GetAPIStateFromAgentState(vm.State)

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *VirtualMachineAPI) RebootVirtualMachine(ctx context.Context, req *pprovisioning.RebootVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	prev := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, prev); err == nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Name)
	} else if status.Code(err) != codes.NotFound {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	res := &pprovisioning.VirtualMachine{
		Metadata: prev.Metadata,
		Spec:     prev.Spec,
		Status:   prev.Status,
	}

	conn, err := a.nodeConnections.GetConnection(prev.Status.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.Status.ComputeNodeName)
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	vm, err := cli.RebootVirtualMachineAgent(context.Background(), &RebootVirtualMachineAgentRequest{
		Name: req.Name,
		Hard: req.Hard,
	})
	if err != nil {
		log.Printf("Fail to reboot on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to reboot volume on node") // TODO #89
	}
	res.Status.State = GetAPIStateFromAgentState(vm.State)

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *VirtualMachineAPI) ShutdownVirtualMachine(ctx context.Context, req *pprovisioning.ShutdownVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	prev := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, prev); err == nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Name)
	} else if status.Code(err) != codes.NotFound {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	res := &pprovisioning.VirtualMachine{
		Metadata: prev.Metadata,
		Spec:     prev.Spec,
		Status:   prev.Status,
	}

	conn, err := a.nodeConnections.GetConnection(prev.Status.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.Status.ComputeNodeName)
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	vm, err := cli.ShutdownVirtualMachineAgent(context.Background(), &ShutdownVirtualMachineAgentRequest{
		Name: req.Name,
		Hard: req.Hard,
	})
	if err != nil {
		log.Printf("Fail to shutdown on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to shutdown volume on node") // TODO #89
	}
	res.Status.State = GetAPIStateFromAgentState(vm.State)

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *VirtualMachineAPI) SaveVirtualMachine(ctx context.Context, req *pprovisioning.SaveVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}
