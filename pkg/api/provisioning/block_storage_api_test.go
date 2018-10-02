package provisioning

import (
	"context"
	"os"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getNodeAPI() (ppool.NodeServiceClient, *grpc.ClientConn, error) {
	endpoint := ""
	if value, ok := os.LookupEnv("NODE_API_ENDPOINT"); ok {
		endpoint = value
	} else {
		endpoint = "localhost:20181"
	}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	noc := ppool.NewNodeServiceClient(conn)

	return noc, conn, nil
}

func TestEmptyBlockStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, conn, err := getNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer conn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create Network API: err='%s'", err.Error())
	}

	listRes, err := bsa.ListBlockStorages(context.Background(), &pprovisioning.ListBlockStoragesRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("ListBlockStorages got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListBlockStorages do not return nil: res='%s'", listRes)
	}

	getRes, err := bsa.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("GetBlockStorage got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetBlockStorage do not return nil: res='%s'", listRes)
	}
}

func TestApplyBlockStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, nconn, err := getNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create Network API: err='%s'", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: "mock-node",
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     "mock-node",
		StorageName:  "test-block-storage",
	}

	createRes, err := bsa.CreateBlockStorage(context.Background(), &pprovisioning.CreateBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Errorf("Failed to create block storage: err='%s'", err.Error())
	}

	createRes.XXX_sizecache = 0
	if diff := cmp.Diff(bs, createRes); diff != "" {
		t.Errorf("CreateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := bsa.ListBlockStorages(context.Background(), &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 1 {
		t.Errorf("ListBlockStorages return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := bsa.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("GetBlockStorage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(bs, getRes); diff != "" {
		t.Errorf("GetBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name}); err != nil {
		t.Errorf("DeleteBlockStorage got error: err='%s'", err.Error())
	}
}

// func TestApplyBlockStorageAboutErrors(t *testing.T) {}

func TestBlockStorageAboutInUseState(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, nconn, err := getNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create Network API: err='%s'", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: "mock-node",
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     "mock-node",
		StorageName:  "test-block-storage",
	}

	_, err = bsa.CreateBlockStorage(context.Background(), &pprovisioning.CreateBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}
	defer bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name})

	inuse, err := bsa.SetInuseBlockStorage(context.Background(), &pprovisioning.SetInuseBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage got error: err='%s'", err.Error())
	}
	if inuse.State != pprovisioning.BlockStorage_IN_USE {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage response wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_IN_USE, inuse.State)
	}
	getRes, _ := bsa.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if getRes.State != pprovisioning.BlockStorage_IN_USE {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage store wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_IN_USE, getRes.State)
	}

	_, err = bsa.SetProtectedBlockStorage(context.Background(), &pprovisioning.SetProtectedBlockStorageRequest{Name: bs.Name})
	if err != nil && status.Code(err) != codes.FailedPrecondition {
		t.Errorf("[Invalid: IN_USE -> PROTECTED] SetProtectedBlockStorage got error, not FailedPrecondition: err='%s'", err.Error())
	}
	_, err = bsa.SetInuseBlockStorage(context.Background(), &pprovisioning.SetInuseBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: IN_USE -> IN_USE] SetInuseBlockStorage got error: err='%s'", err.Error())
	}
	available, err := bsa.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: IN_USE -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
	if available.State != pprovisioning.BlockStorage_AVAILABLE {
		t.Errorf("[Valid: IN_USE -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
}

func TestBlockStorageAboutProtectedState(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, nconn, err := getNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create Network API: err='%s'", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: "mock-node",
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     "mock-node",
		StorageName:  "test-block-storage",
	}

	_, err = bsa.CreateBlockStorage(context.Background(), &pprovisioning.CreateBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}
	defer bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name})

	protected, err := bsa.SetProtectedBlockStorage(context.Background(), &pprovisioning.SetProtectedBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetProtectedBlockStorage got error: err='%s'", err.Error())
	}
	if protected.State != pprovisioning.BlockStorage_PROTECTED {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetProtectedBlockStorage response wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_PROTECTED, protected.State)
	}
	getRes, _ := bsa.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if getRes.State != pprovisioning.BlockStorage_PROTECTED {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetInuseBlockStorage store wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_PROTECTED, getRes.State)
	}

	_, err = bsa.SetInuseBlockStorage(context.Background(), &pprovisioning.SetInuseBlockStorageRequest{Name: bs.Name})
	if err != nil && status.Code(err) != codes.FailedPrecondition {
		t.Errorf("[InValid: PROTECTED -> IN_USE] SetInuseBlockStorage got error: err='%s'", err.Error())
	}
	_, err = bsa.SetProtectedBlockStorage(context.Background(), &pprovisioning.SetProtectedBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: PROTECTED -> PROTECTED] SetProtectedBlockStorage got error, not FailedPrecondition: err='%s'", err.Error())
	}
	available, err := bsa.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: PROTECTED -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
	if available.State != pprovisioning.BlockStorage_AVAILABLE {
		t.Errorf("[Valid: PROTECTED -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
}
