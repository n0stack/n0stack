package image

import (
	"context"
	"os"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0core/pkg/api/provisioning"
	"github.com/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/proto.go/deployment/v0"
	"github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getTestBlockStorageAPI() (pprovisioning.BlockStorageServiceClient, *grpc.ClientConn, error) {
	endpoint := ""
	if value, ok := os.LookupEnv("BLOCK_STORAGE_API_ENDPOINT"); ok {
		endpoint = value
	} else {
		endpoint = "localhost:20183"
	}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	return pprovisioning.NewBlockStorageServiceClient(conn), conn, nil
}

func TestEmptyImage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	bsa, bsconn, err := getTestBlockStorageAPI()
	if err != nil {
		t.Fatalf("Failed to connect block storage api: err='%s'", err.Error())
	}
	defer bsconn.Close()

	ia, err := CreateImageAPI(m, bsa)
	if err != nil {
		t.Fatalf("Failed to create Image API: err='%s'", err.Error())
	}

	listRes, err := ia.ListImages(context.Background(), &pdeployment.ListImagesRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("ListImages got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListImages do not return nil: res='%s'", listRes)
	}

	getRes, err := ia.GetImage(context.Background(), &pdeployment.GetImageRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("GetImage got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetImage do not return nil: res='%s'", getRes)
	}
}

func TestApplyImage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	bsa, bsconn, err := getTestBlockStorageAPI()
	if err != nil {
		t.Fatalf("Failed to connect block storage api: err='%s'", err.Error())
	}
	defer bsconn.Close()

	ia, err := CreateImageAPI(m, bsa)
	if err != nil {
		t.Fatalf("Failed to create Image API: err='%s'", err.Error())
	}

	i := &pdeployment.Image{
		Name:    "test-network",
		Version: 1,
	}

	applyRes, err := ia.ApplyImage(context.Background(), &pdeployment.ApplyImageRequest{
		Name: i.Name,
	})
	if err != nil {
		t.Fatalf("ApplyImage got error: err='%s'", err.Error())
	}
	// diffが取れないので
	applyRes.XXX_sizecache = 0
	if diff := cmp.Diff(i, applyRes); diff != "" {
		t.Fatalf("ApplyImage response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := ia.ListImages(context.Background(), &pdeployment.ListImagesRequest{})
	if err != nil {
		t.Errorf("ListImages got error: err='%s'", err.Error())
	}
	if len(listRes.Images) != 1 {
		t.Errorf("ListImages response is wrong: have='%d', want='%d'", len(listRes.Images), 1)
	}

	getRes, err := ia.GetImage(context.Background(), &pdeployment.GetImageRequest{Name: i.Name})
	if err != nil {
		t.Errorf("GetImage got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(i, getRes); diff != "" {
		t.Errorf("GetImage response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := ia.DeleteImage(context.Background(), &pdeployment.DeleteImageRequest{Name: i.Name}); err != nil {
		t.Errorf("DeleteImage got error: err='%s'", err.Error())
	}
}

func TestImageAboutRegister(t *testing.T) {
	m := memory.NewMemoryDatastore()
	bsa, bsconn, err := getTestBlockStorageAPI()
	if err != nil {
		t.Fatalf("Failed to connect block storage api: err='%s'", err.Error())
	}
	defer bsconn.Close()

	ia, err := CreateImageAPI(m, bsa)
	if err != nil {
		t.Fatalf("Failed to create Image API: err='%s'", err.Error())
	}

	i := &pdeployment.Image{
		Name:    "test-network",
		Version: 1,
		RegisteredBlockStorages: []*pdeployment.Image_RegisteredBlockStorage{
			&pdeployment.Image_RegisteredBlockStorage{
				BlockStorageName: "test-image",
			},
		},
		Tags: map[string]string{
			"test-tag": "test-image",
		},
	}
	_, err = ia.ApplyImage(context.Background(), &pdeployment.ApplyImageRequest{
		Name: i.Name,
	})
	if err != nil {
		t.Fatalf("ApplyImage got error: err='%s'", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-image",
		Annotations: map[string]string{
			provisioning.AnnotationRequestNodeName: "mock-node",
		},
		RequestBytes: 10 * bytefmt.MEGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	}
	_, err = bsa.CreateBlockStorage(context.Background(), &pprovisioning.CreateBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create test-image on BlockStorageAPI got error: err='%s'", err.Error())
	}
	defer bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name})
	defer bsa.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.Name})

	regRes, err := ia.RegisterBlockStorage(context.Background(), &pdeployment.RegisterBlockStorageRequest{
		ImageName:        i.Name,
		BlockStorageName: bs.Name,
		Tags: []string{
			"test-tag",
		},
	})
	if err != nil {
		t.Errorf("RegisterBlockStorage got error: err='%s'", err.Error())
	}
	if len(regRes.RegisteredBlockStorages) != 1 {
		t.Errorf("RegisterBlockStorage response of len(RegisteredBlockStorages) is wrong: have=%d, want=%d", len(regRes.RegisteredBlockStorages), 1)
	}
	regRes.XXX_sizecache = 0
	regRes.RegisteredBlockStorages[0].XXX_sizecache = 0
	regRes.RegisteredBlockStorages[0].RegisteredAt = nil
	if diff := cmp.Diff(i, regRes); diff != "" {
		t.Errorf("RegisterBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	rbs, err := bsa.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("Failed to get test-image on BlockStorageAPI got error: err='%s'", err.Error())
	}
	if rbs.State != pprovisioning.BlockStorage_PROTECTED {
		t.Errorf("BlockStorage 'test-image' state is wrong: have=%+v, want=%+v", rbs.State, pprovisioning.BlockStorage_PROTECTED)
	}

	unregRes, err := ia.UnregisterBlockStorage(context.Background(), &pdeployment.UnregisterBlockStorageRequest{
		ImageName:        i.Name,
		BlockStorageName: bs.Name,
	})
	if err != nil {
		t.Errorf("RegisterBlockStorage got error: err='%s'", err.Error())
	}
	if len(unregRes.RegisteredBlockStorages) != 0 {
		t.Errorf("RegisterBlockStorage response of len(RegisteredBlockStorages) is wrong: have=%d, want=%d", len(unregRes.RegisteredBlockStorages), 0)
	}
	if len(unregRes.Tags) != 0 {
		t.Errorf("RegisterBlockStorage response of len(Tags) is wrong: have=%d, want=%d", len(unregRes.Tags), 0)
	}

	rbs, err = bsa.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: bs.Name})
	if err != nil {
		t.Errorf("Failed to get test-image on BlockStorageAPI got error: err='%s'", err.Error())
	}
	if rbs.State != pprovisioning.BlockStorage_AVAILABLE {
		t.Errorf("BlockStorage 'test-image' state is wrong: have=%+v, want=%+v", rbs.State, pprovisioning.BlockStorage_AVAILABLE)
	}
}

func TestImageAboutTag(t *testing.T) {
	m := memory.NewMemoryDatastore()
	bsa, bsconn, err := getTestBlockStorageAPI()
	if err != nil {
		t.Fatalf("Failed to connect block storage api: err='%s'", err.Error())
	}
	defer bsconn.Close()

	ia, err := CreateImageAPI(m, bsa)
	if err != nil {
		t.Fatalf("Failed to create Image API: err='%s'", err.Error())
	}

	i := &pdeployment.Image{
		Name:    "test-network",
		Version: 1,
		RegisteredBlockStorages: []*pdeployment.Image_RegisteredBlockStorage{
			&pdeployment.Image_RegisteredBlockStorage{
				BlockStorageName: "test-image",
			},
		},
		Tags: map[string]string{
			"test-tag": "test-image",
		},
	}
	_, err = ia.ApplyImage(context.Background(), &pdeployment.ApplyImageRequest{
		Name: i.Name,
	})
	if err != nil {
		t.Fatalf("ApplyImage got error: err='%s'", err.Error())
	}

	bs := &pprovisioning.BlockStorage{
		Name: "test-image",
		Annotations: map[string]string{
			provisioning.AnnotationRequestNodeName: "mock-node",
		},
		RequestBytes: 10 * bytefmt.MEGABYTE,
		LimitBytes:   1 * bytefmt.GIGABYTE,
	}
	_, err = bsa.CreateBlockStorage(context.Background(), &pprovisioning.CreateBlockStorageRequest{
		Name:         bs.Name,
		Annotations:  bs.Annotations,
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	})
	if err != nil {
		t.Fatalf("Failed to create test-image on BlockStorageAPI got error: err='%s'", err.Error())
	}
	defer bsa.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: bs.Name})
	defer bsa.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.Name})

	_, err = ia.RegisterBlockStorage(context.Background(), &pdeployment.RegisterBlockStorageRequest{
		ImageName:        i.Name,
		BlockStorageName: bs.Name,
	})
	if err != nil {
		t.Errorf("RegisterBlockStorage got error: err='%s'", err.Error())
	}

	tagRes, err := ia.TagImage(context.Background(), &pdeployment.TagImageRequest{
		Name:             i.Name,
		BlockStorageName: bs.Name,
		Tags: []string{
			"test-tag",
		},
	})
	if err != nil {
		t.Errorf("TagBlockStorage got error: err='%s'", err.Error())
	}
	if len(tagRes.Tags) != 1 {
		t.Errorf("TagBlockStorage response of len(Tags) is wrong: have=%d, want=%d", len(tagRes.Tags), 0)
	}
	tagRes.XXX_sizecache = 0
	i.RegisteredBlockStorages = nil
	tagRes.RegisteredBlockStorages = nil
	if diff := cmp.Diff(i, tagRes); diff != "" {
		t.Errorf("TagBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	untagRes, err := ia.UntagImage(context.Background(), &pdeployment.UntagImageRequest{
		Name: i.Name,
		Tag:  "test-tag",
	})
	if err != nil {
		t.Errorf("UntagImage got error: err='%s'", err.Error())
	}
	if len(untagRes.Tags) != 0 {
		t.Errorf("UntagImage response of len(Tags) is wrong: have=%d, want=%d", len(untagRes.Tags), 0)
	}
}
