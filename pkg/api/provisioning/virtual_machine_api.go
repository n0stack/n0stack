package provisioning

import (
	"context"
	"log"
	"net/url"
	"reflect"
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
	nodeAPI         ppool.NodeServiceClient
	networkAPI      ppool.NetworkServiceClient
	blockstorageAPI pprovisioning.BlockStorageServiceClient

	nodeConnections *node.NodeConnections
}

func CreateVirtualMachineAPI(ds datastore.Datastore, noa ppool.NodeServiceClient, nea ppool.NetworkServiceClient, bsa pprovisioning.BlockStorageServiceClient) (*VirtualMachineAPI, error) {
	nc := &node.NodeConnections{
		NodeAPI: noa,
	}

	a := &VirtualMachineAPI{
		dataStore:       ds,
		nodeAPI:         noa,
		networkAPI:      nea,
		blockstorageAPI: bsa,
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
		return nil, grpc.Errorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", req.Metadata.Name)
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

	if blockdev, err = a.reserveBlockStorage(req.Spec.BlockStorageNames); err != nil {
		log.Printf("Failed to reserve block storage: err=%v.", err.Error())
		goto ReleaseCompute
	}

	res.Spec.Nics, res.Status.NetworkInterfaceNames, err = a.reserveNics(req.Metadata.Name, req.Spec.Nics)
	if err != nil {
		log.Printf("Failed to reserve nics: err=%v.", err.Error())
		goto ReleaseBlockStorage
	}

	if vm, err := cli.CreateVirtualMachineAgent(context.Background(), &CreateVirtualMachineAgentRequest{
		Name:        req.Metadata.Name,
		Vcpus:       req.Spec.LimitCpuMilliCore / 1000,
		MemoryBytes: req.Spec.LimitMemoryBytes,
		Netdev:      StructNetDev(req.Spec.Nics, res.Status.NetworkInterfaceNames),
		Blockdev:    blockdev,
	}); err != nil {
		log.Printf("Failed to create block storage on node '%s': err='%s'", res.Status.ComputeNodeName, err.Error()) // TODO: #89
		goto ReleaseNetworkInterface
	} else {
		log.Printf("[DEBUG] after CreateVirtualMachineAgent: vm='%+v'", vm)
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

ReleaseBlockStorage:
	if err := a.relaseBlockStorages(res.Spec.BlockStorageNames); err != nil {
		log.Printf("Fail to release block storage on API: err=%s.", err.Error())
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
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if reflect.ValueOf(prev.Metadata).IsNil() {
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

	if err := a.relaseBlockStorages(prev.Spec.BlockStorageNames); err != nil {
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
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if reflect.ValueOf(prev.Metadata).IsNil() {
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
		return nil, grpc.Errorf(codes.Internal, "Fail to boot block storage on node") // TODO #89
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
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if reflect.ValueOf(prev.Metadata).IsNil() {
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
		return nil, grpc.Errorf(codes.Internal, "Fail to reboot block storage on node") // TODO #89
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
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if reflect.ValueOf(prev.Metadata).IsNil() {
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
		return nil, grpc.Errorf(codes.Internal, "Fail to shutdown block storage on node") // TODO #89
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

func (a VirtualMachineAPI) reserveCompute(name string, annotations map[string]string, reqCpu, limitCpu uint32, reqMem, limitMem uint64) (string, string, error) {
	var rcr *ppool.ReserveComputeResponse
	var err error
	if node, ok := annotations[AnnotationRequestNodeName]; !ok {
		rcr, err = a.nodeAPI.ScheduleCompute(context.Background(), &ppool.ScheduleComputeRequest{
			ComputeName: name,
			Compute: &pbudget.Compute{
				RequestCpuMilliCore: reqCpu,
				LimitCpuMilliCore:   limitCpu,
				RequestMemoryBytes:  reqMem,
				LimitMemoryBytes:    limitMem,
			},
		})
	} else {
		rcr, err = a.nodeAPI.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
			Name:        node,
			ComputeName: name,
			Compute: &pbudget.Compute{
				RequestCpuMilliCore: reqCpu,
				LimitCpuMilliCore:   limitCpu,
				RequestMemoryBytes:  reqMem,
				LimitMemoryBytes:    limitMem,
			},
		})
	}
	if err != nil {
		return "", "", grpc.Errorf(codes.Internal, "") // TODO: #89
	}

	return rcr.Name, rcr.ComputeName, nil
}

func (a VirtualMachineAPI) releaseCompute(node, compute string) error {
	_, err := a.nodeAPI.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		Name:        node,
		ComputeName: compute,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", compute, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if status.Code(err) != codes.NotFound {
			return grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", compute)
		}
	}

	return nil
}

func (a VirtualMachineAPI) reserveNics(name string, nics []*pprovisioning.VirtualMachineSpec_NIC) ([]*pprovisioning.VirtualMachineSpec_NIC, []string, error) {
	// res.Status.NetworkInterfaceNames = make([]string, 0, len(req.Spec.Nics))
	networkInterfaceNames := make([]string, 0, len(nics))

	for i, n := range nics {
		ni, err := a.networkAPI.ReserveNetworkInterface(context.Background(), &ppool.ReserveNetworkInterfaceRequest{
			Name:                 n.NetworkName,
			NetworkInterfaceName: name + strconv.Itoa(i),
			NetworkInterface: &pbudget.NetworkInterface{
				HardwareAddress: nics[i].HardwareAddress,
				Ipv4Address:     nics[i].Ipv4Address,
				Ipv6Address:     nics[i].Ipv6Address,
			},
		})
		if err != nil {
			log.Printf("Failed to relserve network interface '%s' from API: %s", name+strconv.Itoa(i), err.Error())
			return nil, nil, err // TODO: #89
		}

		nics[i].HardwareAddress = ni.NetworkInterface.HardwareAddress
		nics[i].Ipv4Address = ni.NetworkInterface.Ipv4Address
		nics[i].Ipv6Address = ni.NetworkInterface.Ipv6Address
		networkInterfaceNames = append(networkInterfaceNames, ni.NetworkInterfaceName)
	}

	return nics, networkInterfaceNames, nil
}

func (a VirtualMachineAPI) releaseNics(nics []*pprovisioning.VirtualMachineSpec_NIC, networkInterfaces []string) error {
	for i, n := range nics {
		_, err := a.networkAPI.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
			Name:                 n.NetworkName,
			NetworkInterfaceName: networkInterfaces[i],
		})
		if err != nil {
			log.Printf("[ERROR] Failed to release network interface '%s': %s", networkInterfaces[i], err.Error())

			// Notfound でもとりあえず問題ないため、処理を続ける
			if status.Code(err) != codes.NotFound {
				return grpc.Errorf(codes.Internal, "Failed to release network interface '%s': please check network interface on your own", networkInterfaces[i])
			}
		}
	}

	return nil
}

func (a VirtualMachineAPI) reserveBlockStorage(names []string) ([]*BlockDev, error) {
	bd := make([]*BlockDev, 0, len(names))
	for i, n := range names {
		v, err := a.blockstorageAPI.SetInuseBlockStorage(context.Background(), &pprovisioning.SetInuseBlockStorageRequest{Name: n})
		if err != nil {
			log.Printf("Failed to get block storage '%s' from API: %s", n, err.Error())
			if status.Code(err) != codes.NotFound {
				return nil, grpc.Errorf(codes.Internal, "Failed to set block storage '%s' as in use from API", n)
			}

			return nil, grpc.Errorf(codes.InvalidArgument, "BlockStorage '%s' is not found", n)
		}

		u := url.URL{
			Scheme: "file",
			Path:   v.Metadata.Annotations[AnnotationBlockStoragePath],
		}
		bd = append(bd, &BlockDev{
			Name:      names[i],
			Url:       u.String(),
			BootIndex: uint32(i),
		})
	}

	return bd, nil
}

func (a VirtualMachineAPI) relaseBlockStorages(names []string) error {
	for _, n := range names {
		_, err := a.blockstorageAPI.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: n})
		if err != nil {
			log.Printf("Failed to get block storage '%s' from API: %s", n, err.Error())

			if status.Code(err) != codes.NotFound {
				return grpc.Errorf(codes.Internal, "Failed to get block storage '%s' as in use from API", n)
			}
		}
	}

	return nil
}
