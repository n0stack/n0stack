package provisioning

import (
	"context"
	"log"
	"strconv"

	"github.com/n0stack/proto.go/budget/v0"
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

func NewVirtualMachineAPI(ds datastore.Datastore, noa ppool.NodeServiceClient, nea ppool.NetworkServiceClient, va pprovisioning.VolumeServiceClient) (*VirtualMachineAPI, error) {
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
	if err := a.dataStore.Get(req.Metadata.Name, prev); err == nil {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Metadata.Name)
	} else if status.Code(err) != codes.NotFound {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	}

	if req.Spec.LimitCpuMilliCore%1000 != 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Make limit_cpu_milli_core '%s' a multiple of 1000", req.Spec.LimitCpuMilliCore)
	}

	res := &pprovisioning.VirtualMachine{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &pprovisioning.VirtualMachineStatus{},
	}

	// reserve compute
	var rcr *ppool.ReserveComputeResponse
	var err error
	if node, ok := req.Metadata.Annotations[AnnotationRequestNodeName]; !ok {
		rcr, err = a.nodeAPI.ScheduleCompute(context.Background(), &ppool.ScheduleComputeRequest{
			ComputeName: req.Metadata.Name,
			Compute: &pbudget.Compute{
				RequestCpuMilliCore: req.Spec.RequestCpuMilliCore,
				LimitCpuMilliCore:   req.Spec.LimitCpuMilliCore,
				RequestMemoryBytes:  req.Spec.RequestMemoryBytes,
				LimitMemoryBytes:    req.Spec.LimitMemoryBytes,
			},
		})
	} else {
		rcr, err = a.nodeAPI.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
			Name:        node,
			ComputeName: req.Metadata.Name,
			Compute: &pbudget.Compute{
				RequestCpuMilliCore: req.Spec.RequestCpuMilliCore,
				LimitCpuMilliCore:   req.Spec.LimitCpuMilliCore,
				RequestMemoryBytes:  req.Spec.RequestMemoryBytes,
				LimitMemoryBytes:    req.Spec.LimitMemoryBytes,
			},
		})
	}
	if err != nil {
		return nil, err // TODO: #89
	}
	res.Status.ComputeNodeName = rcr.Name
	res.Status.ComputeName = rcr.ComputeName

	// reserve network
	res.Status.NetworkInterfaceNames = make([]string, 0, len(req.Spec.Nics))
	for i, n := range req.Spec.Nics {
		ni, err := a.networkAPI.ReserveNetworkInterface(context.Background(), &ppool.ReserveNetworkInterfaceRequest{
			Name:                 n.NetworkName,
			NetworkInterfaceName: req.Metadata.Name + strconv.Itoa(i),
			NetworkInterface: &pbudget.NetworkInterface{
				HardwareAddress: req.Spec.Nics[i].HardwareAddress,
				Ipv4Address:     req.Spec.Nics[i].Ipv4Address,
				Ipv6Address:     req.Spec.Nics[i].Ipv6Address,
			},
		})
		if err != nil {
			return nil, err // TODO: #89
		}

		res.Spec.Nics[i].HardwareAddress = ni.NetworkInterface.HardwareAddress
		res.Spec.Nics[i].Ipv4Address = ni.NetworkInterface.Ipv4Address
		res.Spec.Nics[i].Ipv6Address = ni.NetworkInterface.Ipv6Address
		res.Status.NetworkInterfaceNames = append(res.Status.NetworkInterfaceNames, ni.NetworkInterfaceName)
	}

	// check blockdev
	var blockdev []*BlockDev
	if blockdev, err = StructBlockDev(req.Spec.VolumeNames, a.volumeAPI); err != nil {
		return nil, err
	}

	conn, err := a.nodeConnections.GetConnection(res.Status.ComputeNodeName) // errorについて考える
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "")
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	vm, err := cli.CreateVirtualMachineAgent(context.Background(), &CreateVirtualMachineAgentRequest{
		Name:        req.Metadata.Name,
		Vcpus:       req.Spec.LimitCpuMilliCore / 1000,
		MemoryBytes: req.Spec.LimitMemoryBytes,
		Netdev:      StructNetDev(req.Spec.Nics, res.Status.NetworkInterfaceNames),
		Blockdev:    blockdev,
	})
	if err != nil && status.Code(err) != codes.AlreadyExists {
		log.Printf("Fail to create volume on node '%s': err='%s'", "", err.Error()) // TODO: #89
		return nil, grpc.Errorf(codes.Internal, "")
	}

	res.Metadata.Annotations[AnnotationVNCWebSocketPort] = strconv.Itoa(int(vm.WebsocketPort))
	res.Status.State = GetAPIStateFromAgentState(vm.State)
	res.Status.Uuid = vm.Uuid
	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	}

	return res, nil
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
	if res == nil {
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

	_, err = a.nodeAPI.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		Name:        prev.Status.ComputeNodeName,
		ComputeName: prev.Status.ComputeName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", prev.Status.ComputeName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if status.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", prev.Status.ComputeName)
		}
	}

	for i, n := range prev.Spec.Nics {
		_, err := a.networkAPI.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
			Name:                 n.NetworkName,
			NetworkInterfaceName: prev.Status.NetworkInterfaceNames[i],
		})
		if err != nil {
			log.Printf("[ERROR] Failed to release network interface '%s': %s", prev.Status.NetworkInterfaceNames[i], err.Error())

			// Notfound でもとりあえず問題ないため、処理を続ける
			if status.Code(err) != codes.NotFound {
				return nil, grpc.Errorf(codes.Internal, "Failed to release network interface '%s': please check network interface on your own", prev.Status.NetworkInterfaceNames[i])
			}
		}
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	return &empty.Empty{}, nil
}

func GetAPIStateFromAgentState(s VirtualMachineAgentState) pprovisioning.VirtualMachineStatus_VirtualMachineState {
	switch s {
	case VirtualMachineAgentState_SHUTDOWN:
		return pprovisioning.VirtualMachineStatus_SHUTDOWN

	case VirtualMachineAgentState_RUNNING:
		return pprovisioning.VirtualMachineStatus_RUNNING

	case VirtualMachineAgentState_PAUSED:
		return pprovisioning.VirtualMachineStatus_PAUSED
	}

	return pprovisioning.VirtualMachineStatus_SHUTDOWN
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
