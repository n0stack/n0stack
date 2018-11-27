package provisioning

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestEmptyBlockStorage(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	listRes, err := bsa.ListBlockStorages(ctx, &pprovisioning.ListBlockStoragesRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("ListBlockStorages got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListBlockStorages do not return nil: res='%s'", listRes)
	}

	getRes, err := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("GetBlockStorage got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetBlockStorage do not return nil: res='%s'", listRes)
	}
}

func TestCreateBlockStorage(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}
	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
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

	listRes, err := bsa.ListBlockStorages(ctx, &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 1 {
		t.Errorf("ListBlockStorages return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("GetBlockStorage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(bs, getRes); diff != "" {
		t.Errorf("GetBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := bsa.DeleteBlockStorage(ctx, &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name}); err != nil {
		t.Errorf("DeleteBlockStorage got error: err='%s'", err.Error())
	}
}

func TestFetchBlockStorage(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	createRes, err := bsa.FetchBlockStorage(ctx, &pprovisioning.FetchBlockStorageRequest{
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

	listRes, err := bsa.ListBlockStorages(ctx, &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 1 {
		t.Errorf("ListBlockStorages return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("GetBlockStorage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(bs, getRes); diff != "" {
		t.Errorf("GetBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := bsa.DeleteBlockStorage(ctx, &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name}); err != nil {
		t.Errorf("DeleteBlockStorage got error: err='%s'", err.Error())
	}
}

// url is blank
// func TestFetchBlockStorageAboutErrors(t *testing.T) {}

func TestBlockStorageAboutInUseState(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	_, err = bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	inuse, err := bsa.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage got error: err='%s'", err.Error())
	}
	if inuse.State != pprovisioning.BlockStorage_IN_USE {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage response wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_IN_USE, inuse.State)
	}
	getRes, _ := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if getRes.State != pprovisioning.BlockStorage_IN_USE {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage store wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_IN_USE, getRes.State)
	}

	_, err = bsa.SetProtectedBlockStorage(ctx, &pprovisioning.SetProtectedBlockStorageRequest{Name: bs.Name})
	if err != nil && grpc.Code(err) != codes.FailedPrecondition {
		t.Errorf("[Invalid: IN_USE -> PROTECTED] SetProtectedBlockStorage got error, not FailedPrecondition: err='%s'", err.Error())
	}
	_, err = bsa.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: bs.Name})
	if err == nil {
		t.Errorf("[Inalid: IN_USE -> IN_USE] SetInuseBlockStorage got no error")
	}
	available, err := bsa.SetAvailableBlockStorage(ctx, &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: IN_USE -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
	if available.State != pprovisioning.BlockStorage_AVAILABLE {
		t.Errorf("[Valid: IN_USE -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
}

func TestBlockStorageAboutProtectedState(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	_, err = bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	protected, err := bsa.SetProtectedBlockStorage(ctx, &pprovisioning.SetProtectedBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetProtectedBlockStorage got error: err='%s'", err.Error())
	}
	if protected.State != pprovisioning.BlockStorage_PROTECTED {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetProtectedBlockStorage response wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_PROTECTED, protected.State)
	}
	getRes, _ := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if getRes.State != pprovisioning.BlockStorage_PROTECTED {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetInuseBlockStorage store wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_PROTECTED, getRes.State)
	}

	_, err = bsa.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: bs.Name})
	if err != nil && grpc.Code(err) != codes.FailedPrecondition {
		t.Errorf("[InValid: PROTECTED -> IN_USE] SetInuseBlockStorage got error: err='%s'", err.Error())
	}
	_, err = bsa.SetProtectedBlockStorage(ctx, &pprovisioning.SetProtectedBlockStorageRequest{Name: bs.Name})
	if err == nil {
		t.Errorf("[Valid: PROTECTED -> PROTECTED] SetProtectedBlockStorage got no error")
	}
	available, err := bsa.SetAvailableBlockStorage(ctx, &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("[Valid: PROTECTED -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
	if available.State != pprovisioning.BlockStorage_AVAILABLE {
		t.Errorf("[Valid: PROTECTED -> AVAILABLE] SetAvailableBlockStorage got error: err='%s'", err.Error())
	}
}

func TestCopyBlockStorage(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	src, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name:         "source",
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	copyRes, err := bsa.CopyBlockStorage(ctx, &pprovisioning.CopyBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,

		SourceBlockStorage: src.Name,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	copyRes.XXX_sizecache = 0
	if diff := cmp.Diff(bs, copyRes); diff != "" {
		t.Errorf("CreateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := bsa.ListBlockStorages(ctx, &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 2 {
		t.Errorf("ListBlockStorages return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("GetBlockStorage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(bs, getRes); diff != "" {
		t.Errorf("GetBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := bsa.DeleteBlockStorage(ctx, &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name}); err != nil {
		t.Errorf("DeleteBlockStorage got error: err='%s'", err.Error())
	}
}
