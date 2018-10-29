package provisioning

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/koding/websocketproxy"
	"github.com/n0stack/n0stack/n0proto/pool/v0"
	"github.com/n0stack/n0stack/n0proto/provisioning/v0"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"

	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const AnnotationVNCWebSocketPort = "n0core/provisioning/virtual_machine_vnc_websocket_port"
const AnnotationVirtualMachineReserve = "n0core/provisioning/virtual_machine_name"

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

func (a *VirtualMachineAPI) CreateVirtualMachine(ctx context.Context, req *pprovisioning.CreateVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	if req.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set 'Name'")
	}

	prev := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if prev.Name != "" {
		return nil, grpc.Errorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", req.Name)
	}

	if req.LimitCpuMilliCore%1000 != 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Make limit_cpu_milli_core '%d' a multiple of 1000", req.LimitCpuMilliCore)
	}

	res := &pprovisioning.VirtualMachine{
		Name:                req.Name,
		Annotations:         req.Annotations,
		RequestCpuMilliCore: req.RequestCpuMilliCore,
		LimitCpuMilliCore:   req.LimitCpuMilliCore,
		RequestMemoryBytes:  req.RequestMemoryBytes,
		LimitMemoryBytes:    req.LimitMemoryBytes,
		BlockStorageNames:   req.BlockStorageNames,
		Nics:                req.Nics,
	}
	var err error
	var blockdev []*BlockDev
	if res.Annotations == nil {
		res.Annotations = make(map[string]string)
	}

	res.ComputeNodeName, res.ComputeName, err = a.reserveCompute(
		req.Name,
		req.Annotations,
		req.RequestCpuMilliCore,
		req.LimitCpuMilliCore,
		req.RequestMemoryBytes,
		req.LimitMemoryBytes,
	)
	if err != nil {
		log.Printf("Failed to reserve compute: err=%v.", err.Error())
		return nil, err
	}

	// errorについて考える
	conn, err := a.nodeConnections.GetConnection(res.ComputeNodeName)
	cli := NewVirtualMachineAgentServiceClient(conn)
	if err != nil {
		log.Printf("Failed to dial to node: err=%v.", err.Error())
		goto ReleaseCompute
	}
	if conn == nil {
		// TODO: goto ReleaseCompute
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.ComputeNodeName)
	}
	defer conn.Close()

	if blockdev, err = a.reserveBlockStorage(req.BlockStorageNames); err != nil {
		log.Printf("Failed to reserve block storage: err=%v.", err.Error())
		goto ReleaseCompute
	}

	res.Nics, res.NetworkInterfaceNames, err = a.reserveNics(req.Name, req.Nics)
	if err != nil {
		log.Printf("Failed to reserve nics: err=%v.", err.Error())
		goto ReleaseBlockStorage
	}

	if vm, err := cli.CreateVirtualMachineAgent(context.Background(), &CreateVirtualMachineAgentRequest{
		Name:        req.Name,
		Vcpus:       req.LimitCpuMilliCore / 1000,
		MemoryBytes: req.LimitMemoryBytes,
		Netdev:      StructNetDev(req.Nics, res.NetworkInterfaceNames),
		Blockdev:    blockdev,
	}); err != nil {
		log.Printf("Failed to create virtual machine on node '%s': err='%s'", res.ComputeNodeName, err.Error()) // TODO: #89
		goto ReleaseNetworkInterface
	} else {
		res.Annotations[AnnotationVNCWebSocketPort] = strconv.Itoa(int(vm.WebsocketPort))
		res.State = GetAPIStateFromAgentState(vm.State)
		res.Uuid = vm.Uuid
	}

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteVirtualMachine
	}

	return res, nil

DeleteVirtualMachine:
	_, err = cli.DeleteVirtualMachineAgent(context.Background(), &DeleteVirtualMachineAgentRequest{
		Name:   req.Name,
		Netdev: StructNetDev(res.Nics, res.NetworkInterfaceNames),
	})
	if err != nil {
		log.Printf("Fail to delete virtual machine on node: err=%s.", err.Error())
	}

ReleaseNetworkInterface:
	if err := a.releaseNics(res.Nics, res.NetworkInterfaceNames); err != nil {
		log.Printf("Fail to release network interfaces on API: err=%s.", err.Error())
	}

ReleaseBlockStorage:
	if err := a.relaseBlockStorages(res.BlockStorageNames); err != nil {
		log.Printf("Fail to release block storage on API: err=%s.", err.Error())
	}

ReleaseCompute:
	if err := a.releaseCompute(res.ComputeNodeName, res.ComputeName); err != nil {
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
	if res.Name == "" {
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
	} else if prev.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	conn, err := a.nodeConnections.GetConnection(prev.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.ComputeNodeName)
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	_, err = cli.DeleteVirtualMachineAgent(context.Background(), &DeleteVirtualMachineAgentRequest{
		Name:   req.Name,
		Netdev: StructNetDev(prev.Nics, prev.NetworkInterfaceNames),
	})
	if err != nil {
		log.Printf("Fail to delete virtual machine on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete virtual machine on node") // TODO #89
	}

	if err := a.releaseCompute(prev.ComputeNodeName, prev.ComputeName); err != nil {
		return nil, err
	}

	if err := a.relaseBlockStorages(prev.BlockStorageNames); err != nil {
		return nil, err
	}

	if err := a.releaseNics(prev.Nics, prev.NetworkInterfaceNames); err != nil {
		return nil, err
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	return &empty.Empty{}, nil
}

func (a VirtualMachineAPI) reserveCompute(name string, annotations map[string]string, reqCpu, limitCpu uint32, reqMem, limitMem uint64) (string, string, error) {
	var n *ppool.Node
	var err error
	if node, ok := annotations[AnnotationRequestNodeName]; !ok {
		n, err = a.nodeAPI.ScheduleCompute(context.Background(), &ppool.ScheduleComputeRequest{
			ComputeName: name,
			Annotations: map[string]string{
				AnnotationVirtualMachineReserve: name,
			},
			RequestCpuMilliCore: reqCpu,
			LimitCpuMilliCore:   limitCpu,
			RequestMemoryBytes:  reqMem,
			LimitMemoryBytes:    limitMem,
		})
	} else {
		n, err = a.nodeAPI.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
			NodeName:    node,
			ComputeName: name,
			Annotations: map[string]string{
				AnnotationVirtualMachineReserve: name,
			},
			RequestCpuMilliCore: reqCpu,
			LimitCpuMilliCore:   limitCpu,
			RequestMemoryBytes:  reqMem,
			LimitMemoryBytes:    limitMem,
		})
	}
	if err != nil {
		return "", "", grpc.Errorf(codes.Internal, "") // TODO: #89
	}

	return n.Name, name, nil
}

func (a VirtualMachineAPI) releaseCompute(node, compute string) error {
	_, err := a.nodeAPI.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		NodeName:    node,
		ComputeName: compute,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", compute, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", compute)
		}
	}

	return nil
}

func (a VirtualMachineAPI) reserveNics(name string, nics []*pprovisioning.VirtualMachineNIC) ([]*pprovisioning.VirtualMachineNIC, []string, error) {
	// res.Status.NetworkInterfaceNames = make([]string, 0, len(req.Spec.Nics))
	networkInterfaceNames := make([]string, 0, len(nics))

	for i, nic := range nics {
		niname := name + strconv.Itoa(i)
		network, err := a.networkAPI.ReserveNetworkInterface(context.Background(), &ppool.ReserveNetworkInterfaceRequest{
			NetworkName:          nic.NetworkName,
			NetworkInterfaceName: niname,
			Annotations: map[string]string{
				AnnotationVirtualMachineReserve: name,
			},
			HardwareAddress: nics[i].HardwareAddress,
			Ipv4Address:     nics[i].Ipv4Address,
			Ipv6Address:     nics[i].Ipv6Address,
		})
		if err != nil {
			log.Printf("Failed to relserve network interface '%s' from API: %s", name+strconv.Itoa(i), err.Error())
			return nil, nil, err // TODO: #89
		}

		nics[i].HardwareAddress = network.ReservedNetworkInterfaces[niname].HardwareAddress
		nics[i].Ipv4Address = network.ReservedNetworkInterfaces[niname].Ipv4Address
		nics[i].Ipv6Address = network.ReservedNetworkInterfaces[niname].Ipv6Address
		networkInterfaceNames = append(networkInterfaceNames, niname)
	}

	return nics, networkInterfaceNames, nil
}

func (a VirtualMachineAPI) releaseNics(nics []*pprovisioning.VirtualMachineNIC, networkInterfaces []string) error {
	for i, nic := range nics {
		_, err := a.networkAPI.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
			NetworkName:          nic.NetworkName,
			NetworkInterfaceName: networkInterfaces[i],
		})
		if err != nil {
			log.Printf("[ERROR] Failed to release network interface '%s': %s", networkInterfaces[i], err.Error())

			// Notfound でもとりあえず問題ないため、処理を続ける
			if grpc.Code(err) != codes.NotFound {
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
			if grpc.Code(err) != codes.NotFound {
				return nil, grpc.Errorf(codes.Internal, "Failed to set block storage '%s' as in use from API", n)
			}

			return nil, grpc.Errorf(codes.InvalidArgument, "BlockStorage '%s' is not found", n)
		}

		bd = append(bd, &BlockDev{
			Name:      names[i],
			Url:       v.Annotations[AnnotationBlockStorageURL],
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

			if grpc.Code(err) != codes.NotFound {
				return grpc.Errorf(codes.Internal, "Failed to get block storage '%s' as in use from API", n)
			}
		}
	}

	return nil
}

func (a *VirtualMachineAPI) BootVirtualMachine(ctx context.Context, req *pprovisioning.BootVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	res := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	conn, err := a.nodeConnections.GetConnection(res.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", res.ComputeNodeName)
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	vm, err := cli.BootVirtualMachineAgent(context.Background(), &BootVirtualMachineAgentRequest{Name: req.Name})
	if err != nil {
		log.Printf("Fail to boot on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to boot block storage on node") // TODO #89
	}
	res.State = GetAPIStateFromAgentState(vm.State)

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *VirtualMachineAPI) RebootVirtualMachine(ctx context.Context, req *pprovisioning.RebootVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	res := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	conn, err := a.nodeConnections.GetConnection(res.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", res.ComputeNodeName)
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
	res.State = GetAPIStateFromAgentState(vm.State)

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *VirtualMachineAPI) ShutdownVirtualMachine(ctx context.Context, req *pprovisioning.ShutdownVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	res := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	conn, err := a.nodeConnections.GetConnection(res.ComputeNodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", res.ComputeNodeName)
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
	res.State = GetAPIStateFromAgentState(vm.State)

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *VirtualMachineAPI) SaveVirtualMachine(ctx context.Context, req *pprovisioning.SaveVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

// TODO: めんどくさいので n0core コマンドで定義した URL に一時的に依存している、治す必要あり
func (a *VirtualMachineAPI) OpenConsole(ctx context.Context, req *pprovisioning.OpenConsoleRequest) (*pprovisioning.OpenConsoleResponse, error) {
	vm := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, vm); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if vm.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	u := &url.URL{
		Scheme:   "http",
		Path:     "/static/virtual_machines/vnc.html",
		RawQuery: fmt.Sprintf("path=api/v0/virtual_machines/%s/vncwebsocket", vm.Name),
	}

	return &pprovisioning.OpenConsoleResponse{
		ConsoleUrl: u.String(),
	}, nil
}

func (a *VirtualMachineAPI) ProxyWebsocket() func(echo.Context) error {
	return func(c echo.Context) error {
		vmName := c.Param("name")

		vm := &pprovisioning.VirtualMachine{}
		if err := a.dataStore.Get(vmName, vm); err != nil {
			log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
			return fmt.Errorf("db error")
			// return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
		} else if vm.Name == "" {
			return err
			// return nil, grpc.Errorf(codes.NotFound, "")
		}

		node, err := a.nodeAPI.GetNode(context.Background(), &ppool.GetNodeRequest{Name: vm.ComputeNodeName})
		if err != nil {
			return err
		}

		nodeIP := node.Address
		websocketPort, err := strconv.Atoi(vm.Annotations[AnnotationVNCWebSocketPort])
		if err != nil {
			return err
		}

		u := &url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%s:%d", nodeIP, websocketPort),
			Path:   "/",
		}

		ws := &websocketproxy.WebsocketProxy{
			Backend: func(r *http.Request) *url.URL {
				return u
			},
		}
		// TODO: セキュリティを無視して、とりあえず動かす https://github.com/koding/websocketproxy/issues/9
		delete(c.Request().Header, "Origin")
		log.Printf("[DEBUG] websocket proxy requesting to backend '%s'", ws.Backend(c.Request()))
		ws.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
