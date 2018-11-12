package provisioning

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/koding/websocketproxy"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"

	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const AnnotationVNCWebSocketPort = "n0core/provisioning/virtual_machine/vnc_websocket_port"
const AnnotationVirtualMachineReserving = "n0core/provisioning/virtual_machine/virtual_machine/name"
const AnnotationNetworkInterfaceGateway = "n0core/provisioning/virtual_machine/gateway"

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
	// validation
	switch {
	case req.Name == "":
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set name")

	case req.LimitCpuMilliCore%1000 != 0:
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Make limit_cpu_milli_core '%d' a multiple of 1000", req.LimitCpuMilliCore)

	case req.RequestCpuMilliCore == 0 || req.RequestMemoryBytes == 0:
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set request_*")
	}

	prev := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to get data for db: err='%s'", err.Error())
	} else if prev.Name != "" {
		return nil, WrapGrpcErrorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", req.Name)
	}

	tx := transaction.Begin()
	res := &pprovisioning.VirtualMachine{
		Name:        req.Name,
		Annotations: req.Annotations,
	}
	if res.Annotations == nil {
		res.Annotations = make(map[string]string)
	}

	var err error
	var n *ppool.Node
	if node, ok := req.Annotations[AnnotationRequestNodeName]; !ok {
		n, err = a.nodeAPI.ScheduleCompute(context.Background(), &ppool.ScheduleComputeRequest{
			ComputeName: req.Name,
			Annotations: map[string]string{
				AnnotationVirtualMachineReserving: req.Name,
			},
			RequestCpuMilliCore: req.RequestCpuMilliCore,
			LimitCpuMilliCore:   req.LimitCpuMilliCore,
			RequestMemoryBytes:  req.RequestMemoryBytes,
			LimitMemoryBytes:    req.LimitMemoryBytes,
		})
		if err != nil {
			return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to ScheduleCompute: desc=%s", grpc.ErrorDesc(err))
		}
	} else {
		n, err = a.nodeAPI.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
			NodeName:    node,
			ComputeName: req.Name,
			Annotations: map[string]string{
				AnnotationVirtualMachineReserving: req.Name,
			},
			RequestCpuMilliCore: req.RequestCpuMilliCore,
			LimitCpuMilliCore:   req.LimitCpuMilliCore,
			RequestMemoryBytes:  req.RequestMemoryBytes,
			LimitMemoryBytes:    req.LimitMemoryBytes,
		})
		if err != nil {
			return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to ReserveCompute: desc=%s", grpc.ErrorDesc(err))
		}
	}
	tx.PushRollback(fmt.Sprintf("ReleaseCompute '%s'", req.Name), func() error {
		_, err := a.nodeAPI.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
			NodeName:    n.Name,
			ComputeName: req.Name,
		})
		return err
	})
	res.ComputeName = req.Name
	res.ComputeNodeName = n.Name
	res.RequestCpuMilliCore = req.RequestCpuMilliCore
	res.LimitCpuMilliCore = req.LimitCpuMilliCore
	res.RequestMemoryBytes = req.RequestMemoryBytes
	res.LimitMemoryBytes = req.LimitMemoryBytes

	// errorについて考える
	conn, err := a.nodeConnections.GetConnection(res.ComputeNodeName)
	cli := NewVirtualMachineAgentServiceClient(conn)
	if err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to dial to node: err=%s", err.Error())
	}
	if conn == nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", res.ComputeNodeName)
	}
	defer conn.Close()

	blockdev := make([]*BlockDev, 0, len(req.BlockStorageNames))
	for i, n := range req.BlockStorageNames {
		v, err := a.blockstorageAPI.SetInuseBlockStorage(context.Background(), &pprovisioning.SetInuseBlockStorageRequest{Name: n})
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to SetInuseBlockStorage: desc=%s", grpc.ErrorDesc(err))
		}
		tx.PushRollback(fmt.Sprintf("SetAvailableBlockStorage '%s'", n), func() error {
			_, err := a.blockstorageAPI.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: n})
			return err
		})

		blockdev = append(blockdev, &BlockDev{
			Name:      n,
			Url:       v.Annotations[AnnotationBlockStorageURL],
			BootIndex: uint32(i),
		})
	}
	res.BlockStorageNames = req.BlockStorageNames

	netdev := make([]*NetDev, len(req.Nics))
	res.NetworkInterfaceNames = make([]string, len(req.Nics))
	res.Nics = make([]*pprovisioning.VirtualMachineNIC, len(req.Nics))
	for i, nic := range req.Nics {
		niname := req.Name + strconv.Itoa(i)
		network, err := a.networkAPI.ReserveNetworkInterface(context.Background(), &ppool.ReserveNetworkInterfaceRequest{
			NetworkName:          nic.NetworkName,
			NetworkInterfaceName: niname,
			Annotations: map[string]string{
				AnnotationVirtualMachineReserving: req.Name,
			},
			HardwareAddress: nic.HardwareAddress,
			Ipv4Address:     nic.Ipv4Address,
			Ipv6Address:     nic.Ipv6Address,
		})
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to ReserveNetworkInterface: desc=%s", grpc.ErrorDesc(err))
		}
		tx.PushRollback("ReleaseNetworkInterface", func() error {
			_, err := a.networkAPI.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
				NetworkName:          nic.NetworkName,
				NetworkInterfaceName: niname,
			})
			return err
		})

		res.NetworkInterfaceNames[i] = niname
		res.Nics[i] = &pprovisioning.VirtualMachineNIC{
			NetworkName:     nic.NetworkName,
			HardwareAddress: network.ReservedNetworkInterfaces[niname].HardwareAddress,
			Ipv4Address:     network.ReservedNetworkInterfaces[niname].Ipv4Address,
			Ipv6Address:     network.ReservedNetworkInterfaces[niname].Ipv6Address,
		}

		gateway := ""
		for _, ni := range network.ReservedNetworkInterfaces {
			if v, ok := ni.Annotations[AnnotationNetworkInterfaceGateway]; ok {
				gateway = v
			}
		}
		_, ipn, _ := net.ParseCIDR(network.Ipv4Cidr)
		m, _ := ipn.Mask.Size()
		netdev[i] = &NetDev{
			Name:            niname,
			NetworkName:     res.Nics[i].NetworkName,
			HardwareAddress: res.Nics[i].HardwareAddress,
			Ipv4AddressCidr: fmt.Sprintf("%s/%d", res.Nics[i].Ipv4Address, m),
			Ipv4Gateway:     gateway,
			Nameservers:     []string{"8.8.8.8"}, // TODO: 取るようにする
			// TODO: domain searchはnetworkのdomainから取る
		}
	}

	vm, err := cli.CreateVirtualMachineAgent(context.Background(), &CreateVirtualMachineAgentRequest{
		Name:              req.Name,
		Uuid:              req.Uuid,
		Vcpus:             req.LimitCpuMilliCore / 1000,
		MemoryBytes:       req.LimitMemoryBytes,
		Netdev:            netdev,
		Blockdev:          blockdev,
		LoginUsername:     req.LoginUsername,
		SshAuthorizedKeys: req.SshAuthorizedKeys,
	})
	if err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to CreateVirtualMachineAgent: desc=%s", grpc.ErrorDesc(err))
	}
	tx.PushRollback("", func() error {
		_, err := cli.DeleteVirtualMachineAgent(context.Background(), &DeleteVirtualMachineAgentRequest{Uuid: req.Uuid})
		return err
	})

	res.Annotations[AnnotationVNCWebSocketPort] = strconv.Itoa(int(vm.WebsocketPort))
	res.State = GetAPIStateFromAgentState(vm.State)
	res.Uuid = vm.Uuid
	res.LoginUsername = req.LoginUsername
	res.SshAuthorizedKeys = req.SshAuthorizedKeys

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
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
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db: err='%s'", req.Name, err.Error())
	} else if prev.Name == "" {
		return nil, WrapGrpcErrorf(codes.NotFound, "")
	}

	conn, err := a.nodeConnections.GetConnection(prev.ComputeNodeName)
	if err != nil {
		return nil, WrapGrpcErrorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, WrapGrpcErrorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.ComputeNodeName)
	}
	defer conn.Close()
	cli := NewVirtualMachineAgentServiceClient(conn)

	netdev := make([]*NetDev, len(prev.Nics))
	for i, n := range prev.Nics {
		netdev[i] = &NetDev{
			Name:            prev.NetworkInterfaceNames[i],
			NetworkName:     n.NetworkName,
			HardwareAddress: n.HardwareAddress,
		}
	}

	_, err = cli.DeleteVirtualMachineAgent(context.Background(), &DeleteVirtualMachineAgentRequest{
		Uuid:   prev.Uuid,
		Netdev: netdev,
	})
	if err != nil {
		return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to DeleteVirtualMachineAgent: desc=%s", grpc.ErrorDesc(err))
	}

	_, err = a.nodeAPI.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		NodeName:    prev.ComputeNodeName,
		ComputeName: prev.ComputeName,
	})
	if err != nil {
		return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to ReleaseCompute: desc=%s", grpc.ErrorDesc(err))
	}

	for i, nic := range prev.Nics {
		_, err := a.networkAPI.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
			NetworkName:          nic.NetworkName,
			NetworkInterfaceName: prev.NetworkInterfaceNames[i],
		})
		if err != nil {
			return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to ReleaseNetworkInterface: desc=%s", grpc.ErrorDesc(err))
		}
	}

	for _, n := range prev.BlockStorageNames {
		_, err := a.blockstorageAPI.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: n})
		if err != nil {
			return nil, WrapGrpcErrorf(grpc.Code(err), "Failed to SetAvailableBlockStorage: desc=%s", grpc.ErrorDesc(err))
		}
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, WrapGrpcErrorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	return &empty.Empty{}, nil
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

	vm, err := cli.BootVirtualMachineAgent(context.Background(), &BootVirtualMachineAgentRequest{Uuid: res.Uuid})
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
		Uuid: res.Uuid,
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
		Uuid: res.Uuid,
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
