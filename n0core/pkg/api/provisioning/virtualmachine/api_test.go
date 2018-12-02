package virtualmachine

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0core/pkg/util/net"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
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

	vm := &pprovisioning.VirtualMachine{
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

	createRes, err := vma.CreateVirtualMachine(ctx, &pprovisioning.CreateVirtualMachineRequest{
		Name:                vm.Name,
		Annotations:         vm.Annotations,
		LimitCpuMilliCore:   vm.LimitCpuMilliCore,
		RequestCpuMilliCore: vm.RequestCpuMilliCore,
		LimitMemoryBytes:    vm.LimitMemoryBytes,
		RequestMemoryBytes:  vm.RequestMemoryBytes,
		BlockStorageNames:   vm.BlockStorageNames,
		Nics:                vm.Nics,
		Uuid:                vm.Uuid,
	})
	if err != nil {
		t.Fatalf("Failed to create virtual machine: err='%s'", err.Error())
	}
	createRes.XXX_sizecache = 0
	createRes.Nics[0].XXX_sizecache = 0
	if diff := cmp.Diff(vm, createRes); diff != "" {
		t.Errorf("CreateVirtualMachine response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := vma.ListVirtualMachines(ctx, &pprovisioning.ListVirtualMachinesRequest{})
	if err != nil {
		t.Errorf("ListVirtualMachines got error: err='%s'", err.Error())
	}
	if len(listRes.VirtualMachines) != 1 {
		t.Errorf("ListVirtualMachines return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := vma.GetVirtualMachine(ctx, &pprovisioning.GetVirtualMachineRequest{Name: vm.Name})
	if err != nil {
		t.Errorf("GetVirtualMachine got error: err='%s'", err.Error())
	}
	getRes.Nics[0].XXX_sizecache = 0
	if diff := cmp.Diff(vm, getRes); diff != "" {
		t.Errorf("GetVirtualMachine response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := vma.DeleteVirtualMachine(ctx, &pprovisioning.DeleteVirtualMachineRequest{Name: vm.Name}); err != nil {
		t.Errorf("DeleteVirtualMachine got error: err='%s'", err.Error())
	}
}
