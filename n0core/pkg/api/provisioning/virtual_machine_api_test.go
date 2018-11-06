// +build medium

package provisioning

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestEmptyVirtualMachine(t *testing.T) {
	m := memory.NewMemoryDatastore()
	noa, noconn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer noconn.Close()
	nea, neconn, err := getTestNetworkAPI()
	if err != nil {
		t.Fatalf("Failed to connect network api: err='%s'", err.Error())
	}
	defer neconn.Close()
	bsa, bsconn, err := getTestBlockStorageAPI()
	if err != nil {
		t.Fatalf("Failed to connect block storage api: err='%s'", err.Error())
	}
	defer bsconn.Close()

	vma, err := CreateVirtualMachineAPI(m, noa, nea, bsa)
	if err != nil {
		t.Fatalf("Failed to create virtual machine API: err='%s'", err.Error())
	}

	listRes, err := vma.ListVirtualMachines(context.Background(), &pprovisioning.ListVirtualMachinesRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("ListVirtualMachines got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListVirtualMachines do not return nil: res='%s'", listRes)
	}

	getRes, err := vma.GetVirtualMachine(context.Background(), &pprovisioning.GetVirtualMachineRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("GetVirtualMachine got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetVirtualMachine do not return nil: res='%s'", listRes)
	}
}

func TestCreateVirtualMachine(t *testing.T) {
	m := memory.NewMemoryDatastore()
	noa, noconn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer noconn.Close()
	nea, neconn, err := getTestNetworkAPI()
	if err != nil {
		t.Fatalf("Failed to connect network api: err='%s'", err.Error())
	}
	defer neconn.Close()
	bsa, bsconn, err := getTestBlockStorageAPI()
	if err != nil {
		t.Fatalf("Failed to connect block storage api: err='%s'", err.Error())
	}
	defer bsconn.Close()

	vma, err := CreateVirtualMachineAPI(m, noa, nea, bsa)
	if err != nil {
		t.Fatalf("Failed to create virtual machine API: err='%s'", err.Error())
	}

	ne, err := nea.ApplyNetwork(context.Background(), &ppool.ApplyNetworkRequest{
		Name:     "test-network",
		Ipv4Cidr: "192.168.0.0/30",
		Domain:   "test.local",
	})
	if err != nil {
		t.Fatalf("Failed to apply network: err='%s'", err.Error())
	}
	defer nea.DeleteNetwork(context.Background(), &ppool.DeleteNetworkRequest{Name: ne.Name})

	bs, err := bsa.CreateBlockStorage(context.Background(), &pprovisioning.CreateBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: "mock-node",
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}
	defer bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name})

	vm := &pprovisioning.VirtualMachine{
		Name: "test-virtual-machine",
		Annotations: map[string]string{
			AnnotationRequestNodeName: "mock-node",
		},

		LimitCpuMilliCore:   1000,
		RequestCpuMilliCore: 100,
		LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
		RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		BlockStorageNames:   []string{"test-block-storage"},
		Nics: []*pprovisioning.VirtualMachineNIC{
			{
				NetworkName: "test-network",
				// TODO: 決め打ちなので恒常的に正しいとは限らない
				Ipv4Address:     "192.168.0.1",
				HardwareAddress: "52:54:78:fe:71:fd",
			},
		},
		Uuid:                  "1d5fd196-b6c9-4f58-86f2-3ef227018e47",

		State:                 pprovisioning.VirtualMachine_RUNNING,
		ComputeNodeName:       "mock-node",
		ComputeName:           "test-virtual-machine",
		NetworkInterfaceNames: []string{"test-virtual-machine0"},
	}

	createRes, err := vma.CreateVirtualMachine(context.Background(), &pprovisioning.CreateVirtualMachineRequest{
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
		t.Errorf("Failed to create virtual machine: err='%s'", err.Error())
	}
	createRes.XXX_sizecache = 0
	createRes.Nics[0].XXX_sizecache = 0
	if diff := cmp.Diff(vm, createRes); diff != "" {
		t.Errorf("CreateVirtualMachine response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := vma.ListVirtualMachines(context.Background(), &pprovisioning.ListVirtualMachinesRequest{})
	if err != nil {
		t.Errorf("ListVirtualMachines got error: err='%s'", err.Error())
	}
	if len(listRes.VirtualMachines) != 1 {
		t.Errorf("ListVirtualMachines return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := vma.GetVirtualMachine(context.Background(), &pprovisioning.GetVirtualMachineRequest{Name: vm.Name})
	if err != nil {
		t.Errorf("GetVirtualMachine got error: err='%s'", err.Error())
	}
	getRes.Nics[0].XXX_sizecache = 0
	if diff := cmp.Diff(vm, getRes); diff != "" {
		t.Errorf("GetVirtualMachine response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := vma.DeleteVirtualMachine(context.Background(), &pprovisioning.DeleteVirtualMachineRequest{Name: vm.Name}); err != nil {
		t.Errorf("DeleteVirtualMachine got error: err='%s'", err.Error())
	}
}
