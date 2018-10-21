package flavor

import (
	"context"
	"os"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto/deployment/v0"
	"github.com/n0stack/n0stack/n0proto/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func GetTestVirtualMachineAPI() (pprovisioning.VirtualMachineServiceClient, *grpc.ClientConn, error) {
	endpoint := ""
	if value, ok := os.LookupEnv("VIRTUAL_MACHINE_API_ENDPOINT"); ok {
		endpoint = value
	} else {
		endpoint = "localhost:20180"
	}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	return pprovisioning.NewVirtualMachineServiceClient(conn), conn, nil
}

func GetTestImageAPI() (pdeployment.ImageServiceClient, *grpc.ClientConn, error) {
	endpoint := ""
	if value, ok := os.LookupEnv("IMAGE_API_ENDPOINT"); ok {
		endpoint = value
	} else {
		endpoint = "localhost:20180"
	}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	return pdeployment.NewImageServiceClient(conn), conn, nil
}

func TestEmptyFlavor(t *testing.T) {
	m := memory.NewMemoryDatastore()
	vma, vmconn, err := GetTestVirtualMachineAPI()
	if err != nil {
		t.Fatalf("Failed to connect virtual machine api: err='%s'", err.Error())
	}
	defer vmconn.Close()
	ia, iconn, err := GetTestImageAPI()
	if err != nil {
		t.Fatalf("Failed to connect image api: err='%s'", err.Error())
	}
	defer iconn.Close()

	fa, err := CreateFlavorAPI(m, vma, ia)
	if err != nil {
		t.Fatalf("Failed to create Flavor API: err='%s'", err.Error())
	}

	listRes, err := fa.ListFlavors(context.Background(), &pdeployment.ListFlavorsRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("ListFlavors got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListFlavors do not return nil: res='%s'", listRes)
	}

	getRes, err := fa.GetFlavor(context.Background(), &pdeployment.GetFlavorRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("GetFlavor got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetFlavor do not return nil: res='%s'", getRes)
	}
}

func TestApplyImage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	vma, vmconn, err := GetTestVirtualMachineAPI()
	if err != nil {
		t.Fatalf("Failed to connect virtual machine api: err='%s'", err.Error())
	}
	defer vmconn.Close()
	ia, iconn, err := GetTestImageAPI()
	if err != nil {
		t.Fatalf("Failed to connect image api: err='%s'", err.Error())
	}
	defer iconn.Close()

	fa, err := CreateFlavorAPI(m, vma, ia)
	if err != nil {
		t.Fatalf("Failed to create Flavor API: err='%s'", err.Error())
	}

	f := &pdeployment.Flavor{
		Name: "test",
		// Annotations:       ,
		Version:           1,
		LimitCpuMilliCore: 1000,
		LimitMemoryBytes:  1 * bytefmt.GIGABYTE,
		NetworkName:       "test-network",
	}

	applyRes, err := fa.ApplyFlavor(context.Background(), &pdeployment.ApplyFlavorRequest{
		Name:              f.Name,
		LimitCpuMilliCore: f.LimitCpuMilliCore,
		LimitMemoryBytes:  f.LimitMemoryBytes,
		NetworkName:       f.NetworkName,
	})
	if err != nil {
		t.Fatalf("ApplyFlavor got error: err='%s'", err.Error())
	}
	// diffが取れないので
	applyRes.XXX_sizecache = 0
	if diff := cmp.Diff(f, applyRes); diff != "" {
		t.Fatalf("ApplyFlavor response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := fa.ListFlavors(context.Background(), &pdeployment.ListFlavorsRequest{})
	if err != nil {
		t.Errorf("ListFlavors got error: err='%s'", err.Error())
	}
	if len(listRes.Flavors) != 1 {
		t.Errorf("ListFlavors response is wrong: have='%d', want='%d'", len(listRes.Flavors), 1)
	}

	getRes, err := fa.GetFlavor(context.Background(), &pdeployment.GetFlavorRequest{Name: f.Name})
	if err != nil {
		t.Errorf("GetFlavor got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(f, getRes); diff != "" {
		t.Errorf("GetFlavor response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := fa.DeleteFlavor(context.Background(), &pdeployment.DeleteFlavorRequest{Name: f.Name}); err != nil {
		t.Errorf("DeleteFlavor got error: err='%s'", err.Error())
	}
}

// func TestGenerateVirtualMachine(t *testing.T) {
// 	m := memory.NewMemoryDatastore()
// 	vma, vmconn, err := GetTestVirtualMachineAPI()
// 	if err != nil {
// 		t.Fatalf("Failed to connect virtual machine api: err='%s'", err.Error())
// 	}
// 	defer vmconn.Close()
// 	ia, iconn, err := GetTestImageAPI()
// 	if err != nil {
// 		t.Fatalf("Failed to connect image api: err='%s'", err.Error())
// 	}
// 	defer iconn.Close()

// 	fa, err := CreateFlavorAPI(m, vma, ia)
// 	if err != nil {
// 		t.Fatalf("Failed to create Flavor API: err='%s'", err.Error())
// 	}

// 	f := &pdeployment.Flavor{
// 		Name: "test",
// 		// Annotations:       ,
// 		Version:           1,
// 		LimitCpuMilliCore: 1000,
// 		LimitMemoryBytes:  1 * bytefmt.GIGABYTE,
// 		NetworkName:       "test-network",
// 	}

// 	_, err = fa.ApplyFlavor(context.Background(), &pdeployment.ApplyFlavorRequest{
// 		Name:              f.Name,
// 		LimitCpuMilliCore: f.LimitCpuMilliCore,
// 		LimitMemoryBytes:  f.LimitMemoryBytes,
// 		NetworkName:       f.NetworkName,
// 	})
// 	if err != nil {
// 		t.Fatalf("ApplyFlavor got error: err='%s'", err.Error())
// 	}

// 	res, err := fa.GenerateVirtualMachine(context.Background(), &pdeployment.GenerateVirtualMachineRequest{
// 		FlavorName: f.Name,
// 		VirtualMachineName: "generated_vm",
// 		RequestCpuMilliCore: 100,
// 		RequestMemoryBytes: 100,
// 		RequestStorageBytes: 1000000,
// 		ImageName: "test-image",
// 		ImageTag: "test-tag",
// 	})
// }
