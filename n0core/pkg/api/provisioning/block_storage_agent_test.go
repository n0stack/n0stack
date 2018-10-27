// +build medium
// +build !without_external

package provisioning

import (
	"context"
	"net/url"
	"path/filepath"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
)

func TestAgentCreateEmptyBlockStorage(t *testing.T) {
	bsaa, err := CreateBlockStorageAgentAPI(".")
	if err != nil {
		t.Fatalf("CreateBlockStorageAgentAPI is failed: err=%s", err.Error())
	}

	name := "test-empty"
	size := uint64(10 * bytefmt.MEGABYTE)
	pwd, _ := filepath.Abs(".")
	path := filepath.Join(pwd, name)

	createRes, err := bsaa.CreateEmptyBlockStorage(context.Background(), &CreateEmptyBlockStorageRequest{
		Name:  name,
		Bytes: size,
	})
	if err != nil {
		t.Fatalf("CreateEmptyBlockStorage got error: err=%s", err.Error())
	}
	if diff := cmp.Diff(&CreateEmptyBlockStorageResponse{
		Name:  name,
		Bytes: size,
		Path:  path,
	}, createRes); diff != "" {
		t.Errorf("CreateEmptyBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	_, err = bsaa.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: createRes.Path})
	if err != nil {
		t.Errorf("DeleteBlockStorage got error: err=%s", err.Error())
	}
}

func TestAgentFetchBlockStorage(t *testing.T) {
	bsaa, err := CreateBlockStorageAgentAPI(".")
	if err != nil {
		t.Fatalf("CreateBlockStorageAgentAPI is failed: err=%s", err.Error())
	}

	name := "test-empty"
	size := uint64(10 * bytefmt.MEGABYTE)
	pwd, _ := filepath.Abs(".")
	path := filepath.Join(pwd, name)

	fetchRes, err := bsaa.FetchBlockStorage(context.Background(), &FetchBlockStorageRequest{
		Name:      name,
		Bytes:     size,
		SourceUrl: "http://archive.ubuntu.com/ubuntu/dists/bionic-updates/main/installer-amd64/current/images/netboot/mini.iso",
	})
	if err != nil {
		t.Fatalf("[http] FetchBlockStorage got error: err=%s", err.Error())
	}
	if diff := cmp.Diff(&FetchBlockStorageResponse{
		Name:  name,
		Bytes: size,
		Path:  path,
	}, fetchRes); diff != "" {
		t.Errorf("[http] FetchBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	fname := name + "-copied"
	fpath := filepath.Join(pwd, fname)
	ffetchRes, err := bsaa.FetchBlockStorage(context.Background(), &FetchBlockStorageRequest{
		Name:  fname,
		Bytes: size,
		SourceUrl: (&url.URL{
			Scheme: "file",
			Path:   fetchRes.Path,
		}).String(),
	})
	if err != nil {
		t.Fatalf("[file] FetchBlockStorage got error: err=%s", err.Error())
	}
	if diff := cmp.Diff(&FetchBlockStorageResponse{
		Name:  fname,
		Bytes: size,
		Path:  fpath,
	}, ffetchRes); diff != "" {
		t.Errorf("[file] FetchBlockStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	_, err = bsaa.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: fetchRes.Path})
	if err != nil {
		t.Errorf("DeleteBlockStorage got error: err=%s", err.Error())
	}
	_, err = bsaa.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: ffetchRes.Path})
	if err != nil {
		t.Errorf("DeleteBlockStorage got error: err=%s", err.Error())
	}
}
