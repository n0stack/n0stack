package provisioning

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu_img"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
)

type BlockStorageAgentAPI struct {
	baseDirectory string
}

func CreateBlockStorageAgentAPI(basedir string) (*BlockStorageAgentAPI, error) {
	b, err := filepath.Abs(basedir)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get absolute path")
	}

	if _, err := os.Stat(b); os.IsNotExist(err) {
		if err := os.MkdirAll(b, 0644); err != nil { // TODO: check permission
			return nil, errors.Wrapf(err, "Failed to mkdir '%s'", b)
		}
	}

	return &BlockStorageAgentAPI{
		baseDirectory: b,
	}, nil
}

func (a *BlockStorageAgentAPI) structPath(name string) string {
	return filepath.Join(a.baseDirectory, name)
}

func (a *BlockStorageAgentAPI) CreateEmptyBlockStorage(ctx context.Context, req *CreateEmptyBlockStorageRequest) (*CreateEmptyBlockStorageResponse, error) {
	path := a.structPath(req.Name)
	i, err := img.OpenQemuImg(path)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot open '%s': err='%s'", path, err.Error())
	}
	if i.IsExists() {
		return nil, grpc.Errorf(codes.AlreadyExists, "")
	}

	if err := i.Create(req.Bytes); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to create image: err='%s'", err.Error())
	}

	return &CreateEmptyBlockStorageResponse{
		Name:  req.Name,
		Path:  path,
		Bytes: req.Bytes,
	}, nil
}

// タイムアウトが心配
func (a *BlockStorageAgentAPI) FetchBlockStorage(ctx context.Context, req *FetchBlockStorageRequest) (*FetchBlockStorageResponse, error) {
	path := a.structPath(req.Name)
	i, err := img.OpenQemuImg(path)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot open '%s': err='%s'", path, err.Error())
	}
	if i.IsExists() {
		return nil, grpc.Errorf(codes.AlreadyExists, "")
	}

	u, err := url.Parse(req.SourceUrl)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Parsing source_url '%s' is invalid url: err='%s'", req.SourceUrl, err.Error())
	}

	switch u.Scheme {
	case "http", "https":
		if err := i.Download(u); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to download image: err='%s'", err.Error())
		}

	case "file":
		src, err := img.OpenQemuImg(u.Path)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to open source image: err='%s'", err.Error())
		}

		if err := i.Copy(src); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to download image: err='%s'", err.Error())
		}
	}

	return &FetchBlockStorageResponse{
		Name:  req.Name,
		Path:  path,
		Bytes: req.Bytes,
	}, nil
}

func (a *BlockStorageAgentAPI) DeleteBlockStorage(ctx context.Context, req *DeleteBlockStorageRequest) (*empty.Empty, error) {
	i, err := img.OpenQemuImg(req.Path)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Cannot open '%s': err='%s'", req.Path, err.Error())
	}
	if !i.IsExists() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	if err := i.Delete(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to delete image: err='%s'", err.Error())
	}

	return &empty.Empty{}, nil
}
