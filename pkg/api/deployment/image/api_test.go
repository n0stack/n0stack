package image

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestApplyNetwork(t *testing.T) {
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
