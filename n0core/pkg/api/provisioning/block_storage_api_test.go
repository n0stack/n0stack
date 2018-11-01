// +build medium

package provisioning

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestEmptyBlockStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, conn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer conn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create block storage API: err='%s'", err.Error())
	}

	listRes, err := bsa.ListBlockStorages(context.Background(), &pprovisioning.ListBlockStoragesRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("ListBlockStorages got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListBlockStorages do not return nil: res='%s'", listRes)
	}

	getRes, err := bsa.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("GetBlockStorage got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetBlockStorage do not return nil: res='%s'", listRes)
	}
}

func TestCreateBlockStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, nconn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create block storage API: err='%s'", err.Error())
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

func TestFetchBlockStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, nconn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create block storage API: err='%s'", err.Error())
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

	createRes, err := bsa.FetchBlockStorage(context.Background(), &pprovisioning.FetchBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
		SourceUrl:    "http://test.local",
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

// url is blank
// func TestFetchBlockStorageAboutErrors(t *testing.T) {}

func TestBlockStorageAboutInUseState(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, nconn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create block storage API: err='%s'", err.Error())
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
	if err != nil && grpc.Code(err) != codes.FailedPrecondition {
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
	na, nconn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create block storage API: err='%s'", err.Error())
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
	if err != nil && grpc.Code(err) != codes.FailedPrecondition {
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

func TestBlockCopyStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, nconn, err := getTestNodeAPI()
	if err != nil {
		t.Fatalf("Failed to connect node api: err='%s'", err.Error())
	}
	defer nconn.Close()

	bsa, err := CreateBlockStorageAPI(m, na)
	if err != nil {
		t.Fatalf("Failed to create block storage API: err='%s'", err.Error())
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

	src, err := bsa.CreateBlockStorage(context.Background(), &pprovisioning.CreateBlockStorageRequest{
		Name:         "source",
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}
	defer bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: src.Name})

	copyRes, err := bsa.CopyBlockStorage(context.Background(), &pprovisioning.CopyBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,

		SourceBlockStorage: src.Name,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}
	defer bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name})

	copyRes.XXX_sizecache = 0
	if diff := cmp.Diff(bs, copyRes); diff != "" {
		t.Errorf("CreateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := bsa.ListBlockStorages(context.Background(), &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 2 {
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
