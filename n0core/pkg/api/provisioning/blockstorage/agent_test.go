// +build medium
// +build !without_external

package blockstorage

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
		Path: path,
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
		Path: path,
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
		Path: fpath,
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

// // errors
// // からの場合を
// func TestAgentFetchBlockStorageAboutErrors(t *testing.T) {
// 	bsaa, err := CreateBlockStorageAgentAPI(".")
// 	if err != nil {
// 		t.Fatalf("CreateBlockStorageAgentAPI is failed: err=%s", err.Error())
// 	}

// 	cases := []struct {
// 		name string
// 		req  *FetchBlockStorageRequest
// 		res  *FetchBlockStorageResponse
// 		code codes.Code
// 	}{
// 		{
// 			"no url",
// 			&FetchBlockStorageRequest{},
// 			nil,
// 			codes.Internal,
// 		},
// 	}

// 	for _, c := range cases {
// 		res, err := bsaa.FetchBlockStorage(context.Background(), c.req)
// 		if diff := cmp.Diff(c.res, res); diff != "" {
// 			t.Errorf("")
// 		}

// 		if c.code == 0 && err != nil {
// 			t.Errorf("")
// 		}

// 		if grpc.Code(err) != c.code {
// 			t.Errorf("")
// 		}
// 	}
// }

func TestAgentResizeBlockStorage(t *testing.T) {
	bsaa, err := CreateBlockStorageAgentAPI(".")
	if err != nil {
		t.Fatalf("CreateBlockStorageAgentAPI is failed: err=%s", err.Error())
	}

	name := "test-empty"
	size := uint64(20 * bytefmt.MEGABYTE)

	createRes, err := bsaa.CreateEmptyBlockStorage(context.Background(), &CreateEmptyBlockStorageRequest{
		Name:  name,
		Bytes: uint64(10 * bytefmt.MEGABYTE),
	})
	if err != nil {
		t.Fatalf("CreateEmptyBlockStorage got error: err=%s", err.Error())
	}

	if _, err := bsaa.ResizeBlockStorage(context.Background(), &ResizeBlockStorageRequest{
		Bytes: size,
		Path:  createRes.Path,
	}); err != nil {
		t.Errorf("ResizeBlockStorage got error: err=%s", err.Error())
	}

	_, err = bsaa.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: createRes.Path})
	if err != nil {
		t.Errorf("DeleteBlockStorage got error: err=%s", err.Error())
	}
}
