package blockstorage

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"n0st.ac/n0stack/n0core/pkg/datastore/memory"
	ppool "n0st.ac/n0stack/n0proto.go/pool/v0"
	pprovisioning "n0st.ac/n0stack/n0proto.go/provisioning/v0"
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

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Errorf("Failed to create block storage: err='%s'", err.Error())
	}

	expected := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
			AnnotationBlockStorageURL:             "file:///tmp/test-block-storage",
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	createRes.XXX_sizecache = 0
	if diff := cmp.Diff(expected, createRes); diff != "" {
		t.Errorf("CreateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := bsa.ListBlockStorages(ctx, &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 1 {
		t.Errorf("ListBlockStorages return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: expected.Name})
	if err != nil {
		t.Errorf("GetBlockStorage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(expected, getRes); diff != "" {
		t.Errorf("GetBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := bsa.DeleteBlockStorage(ctx, &pprovisioning.DeleteBlockStorageRequest{Name: expected.Name}); err != nil {
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

	createRes, err := bsa.FetchBlockStorage(ctx, &pprovisioning.FetchBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		SourceUrl:    "http://test.local",
	})
	if err != nil {
		t.Errorf("Failed to create block storage: err='%s'", err.Error())
	}

	expected := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
			AnnotationBlockStorageURL:             "file:///tmp/test-block-storage",
			AnnotationBlockStorageFetchFrom:       "http://test.local",
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	createRes.XXX_sizecache = 0
	if diff := cmp.Diff(expected, createRes); diff != "" {
		t.Errorf("CreateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := bsa.ListBlockStorages(ctx, &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 1 {
		t.Errorf("ListBlockStorages return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: expected.Name})
	if err != nil {
		t.Errorf("GetBlockStorage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(expected, getRes); diff != "" {
		t.Errorf("GetBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := bsa.DeleteBlockStorage(ctx, &pprovisioning.DeleteBlockStorageRequest{Name: expected.Name}); err != nil {
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

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	inuse, err := bsa.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: createRes.Name})
	if err != nil {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage got error: err='%s'", err.Error())
	}
	if inuse.State != pprovisioning.BlockStorage_IN_USE {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage response wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_IN_USE, inuse.State)
	}
	getRes, _ := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: createRes.Name})
	if getRes.State != pprovisioning.BlockStorage_IN_USE {
		t.Errorf("[Valid: AVAILABLE -> IN_USE] SetInuseBlockStorage store wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_IN_USE, getRes.State)
	}

	_, err = bsa.SetProtectedBlockStorage(ctx, &pprovisioning.SetProtectedBlockStorageRequest{Name: createRes.Name})
	if err != nil && grpc.Code(err) != codes.FailedPrecondition {
		t.Errorf("[Invalid: IN_USE -> PROTECTED] SetProtectedBlockStorage got error, not FailedPrecondition: err='%s'", err.Error())
	}
	_, err = bsa.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: createRes.Name})
	if err == nil {
		t.Errorf("[Inalid: IN_USE -> IN_USE] SetInuseBlockStorage got no error")
	}
	available, err := bsa.SetAvailableBlockStorage(ctx, &pprovisioning.SetAvailableBlockStorageRequest{Name: createRes.Name})
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

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	protected, err := bsa.SetProtectedBlockStorage(ctx, &pprovisioning.SetProtectedBlockStorageRequest{Name: createRes.Name})
	if err != nil {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetProtectedBlockStorage got error: err='%s'", err.Error())
	}
	if protected.State != pprovisioning.BlockStorage_PROTECTED {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetProtectedBlockStorage response wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_PROTECTED, protected.State)
	}
	getRes, _ := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: createRes.Name})
	if getRes.State != pprovisioning.BlockStorage_PROTECTED {
		t.Errorf("[Valid: AVAILABLE -> PROTECTED] SetInuseBlockStorage store wrong state: want=%+v, have=%+v", pprovisioning.BlockStorage_PROTECTED, getRes.State)
	}

	_, err = bsa.SetInuseBlockStorage(ctx, &pprovisioning.SetInuseBlockStorageRequest{Name: createRes.Name})
	if err != nil && grpc.Code(err) != codes.FailedPrecondition {
		t.Errorf("[InValid: PROTECTED -> IN_USE] SetInuseBlockStorage got error: err='%s'", err.Error())
	}
	_, err = bsa.SetProtectedBlockStorage(ctx, &pprovisioning.SetProtectedBlockStorageRequest{Name: createRes.Name})
	if err == nil {
		t.Errorf("[Valid: PROTECTED -> PROTECTED] SetProtectedBlockStorage got no error")
	}
	available, err := bsa.SetAvailableBlockStorage(ctx, &pprovisioning.SetAvailableBlockStorageRequest{Name: createRes.Name})
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

	src, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "source",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		Labels: map[string]string{
			"test-duplicate":     "should be deleted",
			"test-not-duplicate": "should be deleted",
		},

		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	copyRes, err := bsa.CopyBlockStorage(ctx, &pprovisioning.CopyBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		Labels: map[string]string{
			"test-duplicate": "correct",
		},

		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,

		SourceBlockStorage: src.Name,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	expected := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
			AnnotationBlockStorageURL:             "file:///tmp/test-block-storage",
			AnnotationBlockStorageCopyFrom:        "source",
		},
		Labels: map[string]string{
			"test-duplicate": "correct",
		},

		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	copyRes.XXX_sizecache = 0
	if diff := cmp.Diff(expected, copyRes); diff != "" {
		t.Errorf("CreateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := bsa.ListBlockStorages(ctx, &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		t.Errorf("ListBlockStorages got error: err='%s'", err.Error())
	}
	if len(listRes.BlockStorages) != 2 {
		t.Errorf("ListBlockStorages return wrong length: res='%s', want=1", listRes)
	}

	getRes, err := bsa.GetBlockStorage(ctx, &pprovisioning.GetBlockStorageRequest{Name: expected.Name})
	if err != nil {
		t.Errorf("GetBlockStorage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(expected, getRes); diff != "" {
		t.Errorf("GetBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := bsa.DeleteBlockStorage(ctx, &pprovisioning.DeleteBlockStorageRequest{Name: expected.Name}); err != nil {
		t.Errorf("DeleteBlockStorage got error: err='%s'", err.Error())
	}
}

func TestCopyBlockStorageByLocal(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "src",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Errorf("Failed to create block storage: err='%s'", err.Error())
	}

	bs, err := bsa.CopyBlockStorage(ctx, &pprovisioning.CopyBlockStorageRequest{
		Name: "dst",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		RequestBytes:       1 * bytefmt.GIGABYTE,
		LimitBytes:         1 * bytefmt.GIGABYTE,
		SourceBlockStorage: createRes.Name,
	})
	if err != nil {
		t.Errorf("Failed to copy block storage: err='%s' %v", err.Error(), bs)
	}
}

func TestCopyBlockStorageByRemote(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}
	dnode, err := bsa.NodeAPI.ApplyNode(ctx, &ppool.ApplyNodeRequest{
		Name:          "dst",
		Address:       "127.0.20.181",
		CpuMilliCores: 16000,
		MemoryBytes:   64 * bytefmt.GIGABYTE,
		StorageBytes:  100 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "src",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		RequestBytes: 1 * bytefmt.GIGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Errorf("Failed to create block storage: err='%s'", err.Error())
	}

	bs, err := bsa.CopyBlockStorage(ctx, &pprovisioning.CopyBlockStorageRequest{
		Name: "dst",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: dnode.Name,
		},
		RequestBytes:       1 * bytefmt.GIGABYTE,
		LimitBytes:         1 * bytefmt.GIGABYTE,
		SourceBlockStorage: createRes.Name,
	})
	if err != nil {
		t.Errorf("Failed to copy block storage: err='%s' %v", err.Error(), bs)
	}
}

func TestUpdateBlockStorageByLocal(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		RequestBytes: bytefmt.GIGABYTE,
		LimitBytes:   bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	updateRes, err := bsa.UpdateBlockStorage(ctx, &pprovisioning.UpdateBlockStorageRequest{
		Name:         createRes.Name,
		RequestBytes: 2 * bytefmt.GIGABYTE,
		LimitBytes:   2 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to update block storage: err='%s'", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
			AnnotationBlockStorageURL:             "file:///tmp/test-block-storage",
		},
		RequestBytes: 2 * bytefmt.GIGABYTE,
		LimitBytes:   2 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     mnode.Name,
		StorageName:  "test-block-storage",
	}

	updateRes.XXX_sizecache = 0
	if diff := cmp.Diff(bs, updateRes); diff != "" {
		t.Errorf("UpdateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
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

func TestUpdateBlockStorageByRemote(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	bsa := NewMockBlcokStorageAPI(m)

	mnode, err := bsa.NodeAPI.SetupMockNode(ctx)
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}
	dnode, err := bsa.NodeAPI.ApplyNode(ctx, &ppool.ApplyNodeRequest{
		Name:          "dst",
		Address:       "127.0.20.181",
		CpuMilliCores: 16000,
		MemoryBytes:   64 * bytefmt.GIGABYTE,
		StorageBytes:  100 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to set up mocked node: err=%s", err.Error())
	}

	createRes, err := bsa.CreateBlockStorage(ctx, &pprovisioning.CreateBlockStorageRequest{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: mnode.Name,
		},
		RequestBytes: bytefmt.GIGABYTE,
		LimitBytes:   bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to create block storage: err='%s'", err.Error())
	}

	updateRes, err := bsa.UpdateBlockStorage(ctx, &pprovisioning.UpdateBlockStorageRequest{
		Name: createRes.Name,
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: dnode.Name,
		},
		RequestBytes: 2 * bytefmt.GIGABYTE,
		LimitBytes:   2 * bytefmt.GIGABYTE,
	})
	if err != nil {
		t.Fatalf("Failed to update block storage: err='%s'", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-block-storage",
		Annotations: map[string]string{
			AnnotationBlockStorageRequestNodeName: dnode.Name,
			AnnotationBlockStorageURL:             "file:///tmp/test-block-storage",
		},
		RequestBytes: 2 * bytefmt.GIGABYTE,
		LimitBytes:   2 * bytefmt.GIGABYTE,
		State:        pprovisioning.BlockStorage_AVAILABLE,
		NodeName:     dnode.Name,
		StorageName:  "test-block-storage",
	}

	updateRes.XXX_sizecache = 0
	if diff := cmp.Diff(bs, updateRes); diff != "" {
		t.Errorf("UpdateBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
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
