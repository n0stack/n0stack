package virtualmachine

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"n0st.ac/n0stack/n0core/pkg/datastore/memory"
	netutil "n0st.ac/n0stack/n0core/pkg/util/net"
	ppool "n0st.ac/n0stack/n0proto.go/pool/v0"
	pprovisioning "n0st.ac/n0stack/n0proto.go/provisioning/v0"
)

func TestCreateVirtualMachine(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	vma := NewMockVirtualMachineAPI(m)

	mnode, err := vma.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	network, err := vma.NetworkAPI.FactoryNetwork(ctx)
	if err != nil {
		t.Fatalf("Failed to factory network: err='%s'", err.Error())
	}
	ip := netutil.ParseCIDR(network.Ipv4Cidr)

	bs, err := vma.BlockStorageAPI.FactoryBlockStorage(ctx, mnode.Name)
	if err != nil {
		t.Fatalf("Failed to factory bloclstorage: err='%s'", err.Error())
	}

	createRes, err := vma.CreateVirtualMachine(ctx, &pprovisioning.CreateVirtualMachineRequest{
		Name: "test-virtual-machine",
		Annotations: map[string]string{
			AnnotationVirtualMachineRequestNodeName: mnode.Name,
			"test-annotation":                       "testing",
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		LimitCpuMilliCore:   1000,
		RequestCpuMilliCore: 100,
		LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
		RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		BlockStorageNames:   []string{bs.Name},
		Nics: []*pprovisioning.VirtualMachineNIC{
			{
				NetworkName: network.Name,
				// TODO: 決め打ちなので恒常的に正しいとは限らない
				Ipv4Address:     ip.Next().IP().String(),
				HardwareAddress: "52:54:78:fe:71:fd",
			},
		},
		Uuid: "1d5fd196-b6c9-4f58-86f2-3ef227018e47",
	})
	if err != nil {
		t.Fatalf("Failed to create virtual machine: err='%s'", err.Error())
	}

	expected := &pprovisioning.VirtualMachine{
		Name: "test-virtual-machine",
		Annotations: map[string]string{
			AnnotationVirtualMachineRequestNodeName:  mnode.Name,
			AnnotationVirtualMachineVncWebSocketPort: "6900",
			"test-annotation":                        "testing",
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		LimitCpuMilliCore:   1000,
		RequestCpuMilliCore: 100,
		LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
		RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		BlockStorageNames:   []string{bs.Name},
		Nics: []*pprovisioning.VirtualMachineNIC{
			{
				NetworkName: network.Name,
				// TODO: 決め打ちなので恒常的に正しいとは限らない
				Ipv4Address:     ip.Next().IP().String(),
				HardwareAddress: "52:54:78:fe:71:fd",
			},
		},
		Uuid: "1d5fd196-b6c9-4f58-86f2-3ef227018e47",

		State:                 pprovisioning.VirtualMachine_RUNNING,
		ComputeNodeName:       mnode.Name,
		ComputeName:           "test-virtual-machine",
		NetworkInterfaceNames: []string{"test-virtual-machine0"},
	}

	createRes.XXX_sizecache = 0
	createRes.Nics[0].XXX_sizecache = 0
	expected.Nics[0].XXX_sizecache = 0
	if diff := cmp.Diff(expected, createRes); diff != "" {
		t.Errorf("CreateVirtualMachine response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := vma.ListVirtualMachines(ctx, &pprovisioning.ListVirtualMachinesRequest{})
	if err != nil {
		t.Errorf("ListVirtualMachines got error: err='%s'", err.Error())
	}
	if len(listRes.VirtualMachines) != 1 {
		t.Errorf("ListVirtualMachines return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := vma.GetVirtualMachine(ctx, &pprovisioning.GetVirtualMachineRequest{Name: expected.Name})
	if err != nil {
		t.Errorf("GetVirtualMachine got error: err='%s'", err.Error())
	}
	getRes.Nics[0].XXX_sizecache = 0
	if diff := cmp.Diff(expected, getRes); diff != "" {
		t.Errorf("GetVirtualMachine response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := vma.DeleteVirtualMachine(ctx, &pprovisioning.DeleteVirtualMachineRequest{Name: expected.Name}); err != nil {
		t.Errorf("DeleteVirtualMachine got error: err='%s'", err.Error())
	}
}

func TestCreateVirtualMachineFailedOnNetworkInterface(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	vma := NewMockVirtualMachineAPI(m)

	mnode, err := vma.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	network, err := vma.NetworkAPI.FactoryNetwork(ctx)
	if err != nil {
		t.Fatalf("Failed to factory network: err='%s'", err.Error())
	}
	ip := netutil.ParseCIDR(network.Ipv4Cidr)

	bs, err := vma.BlockStorageAPI.FactoryBlockStorage(ctx, mnode.Name)
	if err != nil {
		t.Fatalf("Failed to factory bloclstorage: err='%s'", err.Error())
	}

	if _, err = vma.CreateVirtualMachine(ctx, &pprovisioning.CreateVirtualMachineRequest{
		Name: "test-virtual-machine",
		Annotations: map[string]string{
			AnnotationVirtualMachineRequestNodeName: mnode.Name,
		},

		LimitCpuMilliCore:   1000,
		RequestCpuMilliCore: 100,
		LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
		RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		BlockStorageNames:   []string{bs.Name},
		Nics: []*pprovisioning.VirtualMachineNIC{
			{
				NetworkName:     network.Name,
				Ipv4Address:     ip.Next().IP().String(),
				HardwareAddress: "52:54:78:fe:71:fd",
			},
			{
				Ipv4Address:     ip.Next().IP().String(),
				HardwareAddress: "52:54:78:fe:71:fd",
			},
		},
		Uuid: "1d5fd196-b6c9-4f58-86f2-3ef227018e47",
	}); err == nil {
		t.Fatalf("Create virtual machine do not failed")
	}

	network, _ = vma.NetworkAPI.GetNetwork(ctx, &ppool.GetNetworkRequest{
		Name: network.Name,
	})
	if len(network.ReservedNetworkInterfaces) >= 2 { // there is rollbacked interface and default-gateway
		t.Errorf("Failed to rollback about network interface")
	}
}

func TestCreateVirtualMachineFailedOnBlockStorage(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	vma := NewMockVirtualMachineAPI(m)

	mnode, err := vma.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	network, err := vma.NetworkAPI.FactoryNetwork(ctx)
	if err != nil {
		t.Fatalf("Failed to factory network: err='%s'", err.Error())
	}
	ip := netutil.ParseCIDR(network.Ipv4Cidr)

	bs, err := vma.BlockStorageAPI.FactoryBlockStorage(ctx, mnode.Name)
	if err != nil {
		t.Fatalf("Failed to factory bloclstorage: err='%s'", err.Error())
	}

	if _, err := vma.CreateVirtualMachine(ctx, &pprovisioning.CreateVirtualMachineRequest{
		Name: "test-virtual-machine",
		Annotations: map[string]string{
			AnnotationVirtualMachineRequestNodeName: mnode.Name,
		},

		LimitCpuMilliCore:   1000,
		RequestCpuMilliCore: 100,
		LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
		RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		BlockStorageNames: []string{
			bs.Name,
			"not found",
		},
		Nics: []*pprovisioning.VirtualMachineNIC{
			{
				NetworkName:     network.Name,
				Ipv4Address:     ip.Next().IP().String(),
				HardwareAddress: "52:54:78:fe:71:fd",
			},
		},
		Uuid: "1d5fd196-b6c9-4f58-86f2-3ef227018e47",
	}); err == nil {
		t.Fatalf("Create virtual machine do not failed")
	}

	bs, _ = vma.BlockStorageAPI.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{
		Name: bs.Name,
	})
	if bs.State != pprovisioning.BlockStorage_AVAILABLE {
		t.Errorf("Failed to rollback about block storage")
	}
}

func TestDeleteVirtualMachineFailedOnBlockStorage(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	vma := NewMockVirtualMachineAPI(m)

	mnode, err := vma.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	network, err := vma.NetworkAPI.FactoryNetwork(ctx)
	if err != nil {
		t.Fatalf("Failed to factory network: err='%s'", err.Error())
	}
	ip := netutil.ParseCIDR(network.Ipv4Cidr)

	bs, err := vma.BlockStorageAPI.FactoryBlockStorage(ctx, mnode.Name)
	if err != nil {
		t.Fatalf("Failed to factory bloclstorage: err='%s'", err.Error())
	}

	vm, err := vma.CreateVirtualMachine(ctx, &pprovisioning.CreateVirtualMachineRequest{
		Name: "test-virtual-machine",
		Annotations: map[string]string{
			AnnotationVirtualMachineRequestNodeName: mnode.Name,
		},

		LimitCpuMilliCore:   1000,
		RequestCpuMilliCore: 100,
		LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
		RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		BlockStorageNames: []string{
			bs.Name,
		},
		Nics: []*pprovisioning.VirtualMachineNIC{
			{
				NetworkName:     network.Name,
				Ipv4Address:     ip.Next().IP().String(),
				HardwareAddress: "52:54:78:fe:71:fd",
			},
		},
		Uuid: "1d5fd196-b6c9-4f58-86f2-3ef227018e47",
	})
	if err != nil {
		t.Fatalf("failed to create virtual machine")
	}

	if bs, err = vma.BlockStorageAPI.SetAvailableBlockStorage(ctx, &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.Name}); err != nil {
		t.Fatalf("failed to set available (precondition)")
	}
	if _, err = vma.BlockStorageAPI.DeleteBlockStorage(ctx, &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name}); err != nil {
		t.Fatalf("failed to delete block storage (precondition)")
	}
	if _, err = vma.BlockStorageAPI.PurgeBlockStorage(ctx, &pprovisioning.PurgeBlockStorageRequest{Name: bs.Name}); err != nil {
		t.Fatalf("failed to delete block storage (precondition)")
	}

	if _, err := vma.DeleteVirtualMachine(ctx, &pprovisioning.DeleteVirtualMachineRequest{Name: vm.Name}); err == nil {
		t.Errorf("completed to delete virtual machine")
	}

	if network, err = vma.NetworkAPI.GetNetwork(ctx, &ppool.GetNetworkRequest{Name: network.Name}); err != nil {
		t.Errorf("failed to get network")
	}

	if len(network.ReservedNetworkInterfaces) != 2 {
		t.Errorf("failed to rollback network interface")
	}
}
