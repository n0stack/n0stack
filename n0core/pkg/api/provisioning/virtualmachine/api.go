package virtualmachine

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/koding/websocketproxy"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"n0st.ac/n0stack/n0core/pkg/api/pool/network"
	"n0st.ac/n0stack/n0core/pkg/api/pool/node"
	"n0st.ac/n0stack/n0core/pkg/api/provisioning/blockstorage"
	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/datastore"
	"n0st.ac/n0stack/n0core/pkg/datastore/lock"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	netutil "n0st.ac/n0stack/n0core/pkg/util/net"
	"n0st.ac/n0stack/n0proto.go/pkg/transaction"
	ppool "n0st.ac/n0stack/n0proto.go/pool/v0"
	pprovisioning "n0st.ac/n0stack/n0proto.go/provisioning/v0"
)

var N0coreVirtualMachineNamespace uuid.UUID

func init() {
	N0coreVirtualMachineNamespace, _ = uuid.FromString("a015d18d-b2c3-4181-8028-6f707ef31c95")
}

type VirtualMachineAPI struct {
	dataStore datastore.Datastore

	// dependency APIs
	nodeAPI         ppool.NodeServiceClient
	networkAPI      ppool.NetworkServiceClient
	blockstorageAPI pprovisioning.BlockStorageServiceClient

	getAgent func(ctx context.Context, nodeName string) (VirtualMachineAgentServiceClient, func() error, error)
}

func CreateVirtualMachineAPI(ds datastore.Datastore, noa ppool.NodeServiceClient, nea ppool.NetworkServiceClient, bsa pprovisioning.BlockStorageServiceClient) *VirtualMachineAPI {
	a := &VirtualMachineAPI{
		dataStore:       ds.AddPrefix("virtual_machine"),
		nodeAPI:         noa,
		networkAPI:      nea,
		blockstorageAPI: bsa,
	}

	a.getAgent = func(ctx context.Context, nodeName string) (VirtualMachineAgentServiceClient, func() error, error) {
		conn, err := node.GetConnection(ctx, a.nodeAPI, nodeName)
		cli := NewVirtualMachineAgentServiceClient(conn)
		if err != nil {
			return nil, nil, grpcutil.Errorf(codes.Internal, "Failed to dial to node: err=%s", err.Error())
		}
		if conn == nil {
			return nil, nil, grpcutil.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", nodeName)
		}

		return cli, conn.Close, nil
	}

	return a
}

func (a *VirtualMachineAPI) addDefaultGateway(ctx context.Context, nw *ppool.Network) (string, error) {
	ipn := netutil.ParseCIDR(nw.Ipv4Cidr)
	ip := netutil.GetEndIP(ipn.Network())

	_, err := a.networkAPI.ReserveNetworkInterface(ctx, &ppool.ReserveNetworkInterfaceRequest{
		NetworkName:          nw.Name,
		NetworkInterfaceName: "default-gateway",
		Ipv4Address:          ip.String(),
		Annotations: map[string]string{
			AnnotationNetworkInterfaceIsGateway:                   "true",
			network.AnnotationNetworkInterfaceDisableDeletionLock: "true",
		},
	})

	return ip.String(), err
}

// PENDINGステートにすることで楽観的なロックを行う
func (a *VirtualMachineAPI) lockOptimistically(vm *pprovisioning.VirtualMachine) (func() error, error) {
	// PENDINGステートにすることで楽観的なロックを行う
	if vm.State == pprovisioning.VirtualMachine_PENDING {
		return nil, grpcutil.Errorf(codes.ResourceExhausted, "State is PENDING, so cannnot do any actions") // これで State がいいのか自信ない
	}

	vm.State = pprovisioning.VirtualMachine_PENDING
	if err := a.dataStore.Apply(vm.Name, vm); err != nil {
		return nil, grpcutil.Errorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}

	f := func() error {
		if vm.State == pprovisioning.VirtualMachine_PENDING {
			vm.State = pprovisioning.VirtualMachine_VIRTUAL_MACHINE_UNSPECIFIED
		}

		return a.dataStore.Apply(vm.Name, vm)
	}

	return f, nil
}

func (a *VirtualMachineAPI) CreateVirtualMachine(ctx context.Context, req *pprovisioning.CreateVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	// validation
	var id uuid.UUID
	{
		switch {
		case req.Name == "":
			return nil, grpcutil.Errorf(codes.InvalidArgument, "Set name")

		case req.LimitCpuMilliCore%1000 != 0:
			return nil, grpcutil.Errorf(codes.InvalidArgument, "Make limit_cpu_milli_core '%d' a multiple of 1000", req.LimitCpuMilliCore)

		case req.RequestCpuMilliCore == 0 || req.RequestMemoryBytes == 0:
			return nil, grpcutil.Errorf(codes.InvalidArgument, "Set request_*")
		}

		var err error
		if req.Uuid == "" {
			id = uuid.NewV5(N0coreVirtualMachineNamespace, req.Name)
		} else if id, err = uuid.FromString(req.Uuid); err != nil {
			return nil, grpcutil.Errorf(codes.InvalidArgument, "Set valid uuid")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	vm := &pprovisioning.VirtualMachine{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		Uuid:                id.String(),
		RequestCpuMilliCore: req.RequestCpuMilliCore,
		LimitCpuMilliCore:   req.LimitCpuMilliCore,
		RequestMemoryBytes:  req.RequestMemoryBytes,
		LimitMemoryBytes:    req.LimitMemoryBytes,

		BlockStorageNames: req.BlockStorageNames,
		Nics:              req.Nics,

		LoginUsername:     req.LoginUsername,
		SshAuthorizedKeys: req.SshAuthorizedKeys,
	}
	if vm.Annotations == nil {
		vm.Annotations = make(map[string]string)
	}

	{
		prev := &pprovisioning.VirtualMachine{}
		if err := a.dataStore.Get(vm.Name, prev); err == nil {
			return nil, grpcutil.Errorf(codes.AlreadyExists, "VirtualMachine '%s' is already exists", vm.Name)
		} else if !datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
		}

		vm.State = pprovisioning.VirtualMachine_PENDING
		if err := a.dataStore.Apply(vm.Name, vm); err != nil {
			return nil, grpcutil.Errorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
		}
		tx.PushRollback("free optimistic lock", func() error {
			return a.dataStore.Delete(vm.Name)
		})

		vm.State = pprovisioning.VirtualMachine_VIRTUAL_MACHINE_UNSPECIFIED
	}

	{
		vm.ComputeName = vm.Name

		var err error
		var n *ppool.Node
		if node, ok := vm.Annotations[AnnotationVirtualMachineRequestNodeName]; !ok {
			n, err = a.nodeAPI.ScheduleCompute(ctx, &ppool.ScheduleComputeRequest{
				ComputeName: vm.ComputeName,
				Annotations: map[string]string{
					AnnotationComputeReservedBy: vm.Name,
				},
				RequestCpuMilliCore: vm.RequestCpuMilliCore,
				LimitCpuMilliCore:   vm.LimitCpuMilliCore,
				RequestMemoryBytes:  vm.RequestMemoryBytes,
				LimitMemoryBytes:    vm.LimitMemoryBytes,
			})
			if err != nil {
				return nil, grpcutil.Errorf(grpc.Code(err), "Failed to ScheduleCompute: desc=%s", grpc.ErrorDesc(err))
			}
		} else {
			n, err = a.nodeAPI.ReserveCompute(ctx, &ppool.ReserveComputeRequest{
				NodeName:    node,
				ComputeName: vm.ComputeName,
				Annotations: map[string]string{
					AnnotationComputeReservedBy: vm.Name,
				},
				RequestCpuMilliCore: vm.RequestCpuMilliCore,
				LimitCpuMilliCore:   vm.LimitCpuMilliCore,
				RequestMemoryBytes:  vm.RequestMemoryBytes,
				LimitMemoryBytes:    vm.LimitMemoryBytes,
			})
			if err != nil {
				return nil, grpcutil.Errorf(grpc.Code(err), "Failed to ReserveCompute: desc=%s", grpc.ErrorDesc(err))
			}
		}

		vm.ComputeNodeName = n.Name

		tx.PushRollback(fmt.Sprintf("ReleaseCompute '%s'", vm.Name), func() error {
			_, err := a.nodeAPI.ReleaseCompute(ctx, &ppool.ReleaseComputeRequest{
				NodeName:    vm.ComputeNodeName,
				ComputeName: vm.ComputeName,
			})
			return err
		})
	}

	{
		tx.PushRollback("SetAvailableBlockStorage", func() error {
			for _, n := range vm.BlockStorageNames {
				_, err := a.blockstorageAPI.SetAvailableBlockStorage(ctx, &pprovisioning.SetAvailableBlockStorageRequest{Name: n})
				if err != nil {
					return err // errをスタックする必要がある
				}
			}

			return nil
		})
		for _, n := range vm.BlockStorageNames {
			if _, err := a.blockstorageAPI.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: n}); err != nil {
				return nil, grpcutil.Errorf(grpc.Code(err), "Failed to SetInuseBlockStorage: desc=%s", grpc.ErrorDesc(err))
			}
		}
	}

	{
		vm.NetworkInterfaceNames = make([]string, len(vm.Nics))

		tx.PushRollback("ReleaseNetworkInterface", func() error {
			for i := range vm.Nics {
				_, err := a.networkAPI.ReleaseNetworkInterface(ctx, &ppool.ReleaseNetworkInterfaceRequest{
					NetworkName:          vm.Nics[i].NetworkName,
					NetworkInterfaceName: vm.NetworkInterfaceNames[i],
				})

				if err != nil {
					if grpc.Code(err) == codes.NotFound { // When NotFound, failed until processing all
						break
					}

					return err
				}
			}

			return nil
		})
		for i, nic := range vm.Nics {
			vm.NetworkInterfaceNames[i] = vm.Name + strconv.Itoa(i)

			annotations := nic.Annotations
			if annotations == nil {
				annotations = make(map[string]string)
			}
			annotations[AnnotationComputeReservedBy] = vm.Name
			network, err := a.networkAPI.ReserveNetworkInterface(ctx, &ppool.ReserveNetworkInterfaceRequest{
				NetworkName:          nic.NetworkName,
				NetworkInterfaceName: vm.NetworkInterfaceNames[i],
				Annotations:          annotations,
				HardwareAddress:      nic.HardwareAddress,
				Ipv4Address:          nic.Ipv4Address,
				Ipv6Address:          nic.Ipv6Address,
			})
			if err != nil {
				return nil, grpcutil.Errorf(grpc.Code(err), "Failed to ReserveNetworkInterface: desc=%s", grpc.ErrorDesc(err))
			}

			vm.Nics[i].HardwareAddress = network.ReservedNetworkInterfaces[vm.NetworkInterfaceNames[i]].HardwareAddress
			vm.Nics[i].Ipv4Address = network.ReservedNetworkInterfaces[vm.NetworkInterfaceNames[i]].Ipv4Address
			vm.Nics[i].Ipv6Address = network.ReservedNetworkInterfaces[vm.NetworkInterfaceNames[i]].Ipv6Address

			if network.Ipv4Cidr != "" {
				havingGateway := false
				for _, ni := range network.ReservedNetworkInterfaces {
					if _, ok := ni.Annotations[AnnotationNetworkInterfaceIsGateway]; ok {
						havingGateway = true
					}
				}
				if !havingGateway {
					if _, err = a.addDefaultGateway(ctx, network); err != nil {
						return nil, grpcutil.Errorf(codes.Internal, errors.Wrapf(err, "Failed to add default gateway").Error())
					}
				}
			}
		}
	}

	if err := a.dataStore.Apply(vm.Name, vm); err != nil {
		return nil, grpcutil.Errorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}

	res, err := a.bootVirtualMachine(ctx, &pprovisioning.BootVirtualMachineRequest{Name: vm.Name})
	if err != nil {
		return nil, grpcutil.Errorf(grpc.Code(err), errors.Wrapf(err, "Failed to BootVirtualMachineRequest").Error())
	}

	tx.Commit()
	return res, nil
}

func GetAPIStateFromAgentState(s VirtualMachineState) pprovisioning.VirtualMachine_VirtualMachineState {
	switch s {
	case VirtualMachineState_SHUTDOWN:
		return pprovisioning.VirtualMachine_SHUTDOWN

	case VirtualMachineState_RUNNING:
		return pprovisioning.VirtualMachine_RUNNING

	case VirtualMachineState_PAUSED:
		return pprovisioning.VirtualMachine_PAUSED
	}

	return pprovisioning.VirtualMachine_VIRTUAL_MACHINE_UNSPECIFIED
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
		return nil, grpcutil.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.VirtualMachines) == 0 {
		return nil, grpcutil.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *VirtualMachineAPI) GetVirtualMachine(ctx context.Context, req *pprovisioning.GetVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	vm := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, vm); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	return vm, nil
}

func (a *VirtualMachineAPI) UpdateVirtualMachine(ctx context.Context, req *pprovisioning.UpdateVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	return nil, grpcutil.Errorf(codes.Unimplemented, "")
}

func (a *VirtualMachineAPI) DeleteVirtualMachine(ctx context.Context, req *pprovisioning.DeleteVirtualMachineRequest) (*empty.Empty, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	vm := &pprovisioning.VirtualMachine{}
	{
		if err := a.dataStore.Get(req.Name, vm); err != nil {
			if datastore.IsNotFound(err) {
				return nil, grpcutil.Errorf(codes.NotFound, err.Error())
			}

			return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
		}

		if vm.State == pprovisioning.VirtualMachine_PENDING {
			return nil, grpcutil.Errorf(codes.FailedPrecondition, "VirtualMachine '%s' is pending", req.Name)
		}

		current := vm.State
		vm.State = pprovisioning.VirtualMachine_PENDING
		if err := a.dataStore.Apply(vm.Name, vm); err != nil {
			return nil, grpcutil.Errorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
		}
		vm.State = current
		tx.PushRollback("free optimistic lock", func() error {
			vm.State = current
			return a.dataStore.Apply(vm.Name, vm)
		})
	}

	{
		cli, done, err := a.getAgent(ctx, vm.ComputeNodeName)
		if err != nil {
			return nil, err
		}
		defer done()

		netdevs := make([]*NetDev, len(vm.Nics))
		for i, n := range vm.Nics {
			netdevs[i] = &NetDev{
				Name:            vm.NetworkInterfaceNames[i],
				NetworkName:     n.NetworkName,
				HardwareAddress: n.HardwareAddress,
			}
		}

		_, err = cli.DeleteVirtualMachine(context.Background(), &DeleteVirtualMachineRequest{
			Name:    vm.Name,
			Netdevs: netdevs,
		})
		if err != nil {
			return nil, grpcutil.Errorf(grpc.Code(err), "Failed to DeleteVirtualMachineAgent: desc=%s", grpc.ErrorDesc(err))
		}
	}

	_, err := a.nodeAPI.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		NodeName:    vm.ComputeNodeName,
		ComputeName: vm.ComputeName,
	})
	if err != nil {
		return nil, grpcutil.Errorf(grpc.Code(err), "Failed to ReleaseCompute: desc=%s", grpc.ErrorDesc(err))
	}
	tx.PushRollback(fmt.Sprintf("ReserveCompute '%s'", vm.Name), func() error {
		_, err = a.nodeAPI.ReserveCompute(ctx, &ppool.ReserveComputeRequest{
			NodeName:    vm.ComputeNodeName,
			ComputeName: vm.ComputeName,
			Annotations: map[string]string{
				AnnotationComputeReservedBy: vm.Name,
			},
			RequestCpuMilliCore: vm.RequestCpuMilliCore,
			LimitCpuMilliCore:   vm.LimitCpuMilliCore,
			RequestMemoryBytes:  vm.RequestMemoryBytes,
			LimitMemoryBytes:    vm.LimitMemoryBytes,
		})

		return err
	})

	tx.PushRollback("ReserveNetworkInterface", func() error {
		for i := range vm.Nics {
			_, err := a.networkAPI.ReserveNetworkInterface(ctx, &ppool.ReserveNetworkInterfaceRequest{
				NetworkName:          vm.Nics[i].NetworkName,
				HardwareAddress:      vm.Nics[i].HardwareAddress,
				Ipv4Address:          vm.Nics[i].Ipv4Address,
				Ipv6Address:          vm.Nics[i].Ipv6Address,
				NetworkInterfaceName: vm.NetworkInterfaceNames[i],
			})

			if err != nil {
				if grpc.Code(err) == codes.AlreadyExists { // When AlreadyExists, failed until processing all
					break
				}

				return err
			}
		}

		return nil
	})

	for i, nic := range vm.Nics {
		_, err := a.networkAPI.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
			NetworkName:          nic.NetworkName,
			NetworkInterfaceName: vm.NetworkInterfaceNames[i],
		})
		if err != nil {
			return nil, grpcutil.Errorf(grpc.Code(err), "Failed to ReleaseNetworkInterface: desc=%s", grpc.ErrorDesc(err))
		}
	}

	tx.PushRollback("SetInuseBlockStorage", func() error {
		for _, n := range vm.BlockStorageNames {
			_, err := a.blockstorageAPI.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: n})
			if err != nil && grpc.Code(err) != codes.FailedPrecondition {
				return err // errをスタックする必要がある
			}
		}

		return nil
	})
	for _, n := range vm.BlockStorageNames {
		_, err := a.blockstorageAPI.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: n})
		if err != nil {
			return nil, grpcutil.Errorf(grpc.Code(err), "Failed to SetAvailableBlockStorage: desc=%s", grpc.ErrorDesc(err))
		}
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, grpcutil.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	tx.Commit()
	return &empty.Empty{}, nil
}

func (a *VirtualMachineAPI) BootVirtualMachine(ctx context.Context, req *pprovisioning.BootVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	// validation
	{
		switch {
		case req.Name == "":
			return nil, grpcutil.Errorf(codes.InvalidArgument, "Set name")
		}
	}

	if !lock.WaitUntilLock(a.dataStore, req.Name, 5*time.Second, 10*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	return a.bootVirtualMachine(ctx, req)
}

// TODO: こうする
// func (a *VirtualMachineAPI) bootVirtualMachine(ctx context.Context, vm *pprovisioning.VirtualMachine) error {
func (a *VirtualMachineAPI) bootVirtualMachine(ctx context.Context, req *pprovisioning.BootVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	vm := &pprovisioning.VirtualMachine{}
	{
		if err := a.dataStore.Get(req.Name, vm); err != nil {
			if datastore.IsNotFound(err) {
				return nil, grpcutil.Errorf(codes.NotFound, err.Error())
			}

			return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
		}

		if vm.State == pprovisioning.VirtualMachine_PENDING {
			return nil, grpcutil.Errorf(codes.FailedPrecondition, "VirtualMachine '%s' is pending", vm.Name)
		}

		current := vm.State
		vm.State = pprovisioning.VirtualMachine_PENDING
		if err := a.dataStore.Apply(vm.Name, vm); err != nil {
			return nil, grpcutil.Errorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
		}
		vm.State = current
		tx.PushRollback("free optimistic lock", func() error {
			return a.dataStore.Apply(vm.Name, vm)
		})
	}

	blockdevs := make([]*BlockDev, len(vm.BlockStorageNames))
	{
		for i, n := range vm.BlockStorageNames {
			bs, err := a.blockstorageAPI.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: n})
			if err != nil {
				return nil, grpcutil.Errorf(grpc.Code(err), "Failed to GetBlockStorage: desc=%s", grpc.ErrorDesc(err))
			}

			blockdevs[i] = &BlockDev{
				Name:      n,
				Url:       bs.Annotations[blockstorage.AnnotationBlockStorageURL],
				BootIndex: uint32(i),
			}
		}
	}

	netdevs := make([]*NetDev, len(vm.Nics))
	{
		gatewayIndex := 0
		for i, nic := range vm.Nics {
			if nic.Annotations != nil {
				if _, ok := nic.Annotations[AnnotationVirtualMachineNICIsGateway]; ok {
					gatewayIndex = i
				}
			}
		}

		for i, nic := range vm.Nics {
			network, err := a.networkAPI.GetNetwork(ctx, &ppool.GetNetworkRequest{Name: nic.NetworkName})
			if err != nil {
				return nil, grpcutil.Errorf(grpc.Code(err), "Failed to GetNetworkInterface: desc=%s", grpc.ErrorDesc(err))
			}

			netdevs[i] = &NetDev{
				Name:            vm.NetworkInterfaceNames[i],
				NetworkName:     vm.Nics[i].NetworkName,
				HardwareAddress: vm.Nics[i].HardwareAddress,
			}

			ip := netutil.ParseCIDR(network.Ipv4Cidr)
			if ip != nil && vm.Nics[i].Ipv4Address != "" {
				gateway := ""
				for _, ni := range network.ReservedNetworkInterfaces {
					if _, ok := ni.Annotations[AnnotationNetworkInterfaceIsGateway]; ok {
						gateway = ni.Ipv4Address
					}
				}

				netdevs[i].Ipv4AddressCidr = fmt.Sprintf("%s/%d", vm.Nics[i].Ipv4Address, ip.SubnetMaskBits())

				if gatewayIndex == i {
					netdevs[i].Ipv4Gateway = gateway
					netdevs[i].Nameservers = []string{"8.8.8.8"} // TODO: 取るようにする
					// TODO: domain searchはnetworkのdomainから取る
				}
			}
		}
	}

	{
		cli, done, err := a.getAgent(ctx, vm.ComputeNodeName)
		if err != nil {
			return nil, err
		}
		defer done()

		res, err := cli.BootVirtualMachine(ctx, &BootVirtualMachineRequest{
			Name:              vm.Name,
			Uuid:              vm.Uuid,
			Vcpus:             vm.LimitCpuMilliCore / 1000,
			MemoryBytes:       vm.LimitMemoryBytes,
			Netdevs:           netdevs,
			Blockdevs:         blockdevs,
			LoginUsername:     vm.LoginUsername,
			SshAuthorizedKeys: vm.SshAuthorizedKeys,
		})
		if err != nil {
			return nil, grpcutil.Errorf(grpc.Code(err), "Failed to CreateVirtualMachineAgent: desc=%s", grpc.ErrorDesc(err))
		}
		tx.PushRollback("", func() error {
			cli, done, err := a.getAgent(ctx, vm.ComputeNodeName)
			if err != nil {
				return err
			}
			defer done()

			if _, err := cli.DeleteVirtualMachine(ctx, &DeleteVirtualMachineRequest{Name: vm.Name}); err != nil {
				return err
			}

			return nil
		})

		vm.Annotations[AnnotationVirtualMachineVncWebSocketPort] = strconv.Itoa(int(res.WebsocketPort))
		vm.State = GetAPIStateFromAgentState(res.State)
	}

	if err := a.dataStore.Apply(vm.Name, vm); err != nil {
		return nil, grpcutil.Errorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}

	tx.Commit()
	return vm, nil
}

func (a *VirtualMachineAPI) RebootVirtualMachine(ctx context.Context, req *pprovisioning.RebootVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	return nil, grpcutil.Errorf(codes.Unimplemented, "")
}

func (a *VirtualMachineAPI) ShutdownVirtualMachine(ctx context.Context, req *pprovisioning.ShutdownVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	return nil, grpcutil.Errorf(codes.Unimplemented, "")
}

func (a *VirtualMachineAPI) SaveVirtualMachine(ctx context.Context, req *pprovisioning.SaveVirtualMachineRequest) (*pprovisioning.VirtualMachine, error) {
	return nil, grpcutil.Errorf(codes.Unimplemented, "")
}

// TODO: めんどくさいので n0core コマンドで定義した URL に一時的に依存している、治す必要あり
func (a *VirtualMachineAPI) OpenConsole(ctx context.Context, req *pprovisioning.OpenConsoleRequest) (*pprovisioning.OpenConsoleResponse, error) {
	vm := &pprovisioning.VirtualMachine{}
	if err := a.dataStore.Get(req.Name, vm); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.Errorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	u := &url.URL{
		Scheme:   "http",
		Path:     "/n0core/static/virtual_machines/vnc.html",
		RawQuery: fmt.Sprintf("path=n0core/api/v0/virtual_machines/%s/vncwebsocket", vm.Name),
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
			if datastore.IsNotFound(err) {
				return err
			}

			log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
			return fmt.Errorf("db error")
		}

		node, err := a.nodeAPI.GetNode(context.Background(), &ppool.GetNodeRequest{Name: vm.ComputeNodeName})
		if err != nil {
			return err
		}

		nodeIP := node.Address
		websocketPort, err := strconv.Atoi(vm.Annotations[AnnotationVirtualMachineVncWebSocketPort])
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
