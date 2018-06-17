package qcow2

import (
	fmt "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/bytefmt"
	"github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Qcow2Agent is for developer and administrator service.
// Make attention about logs!!
//
// refs: https://github.com/n0stack/n0core/blob/7793ac93917f4dc3f524e6d716174e1dee490173/qcow2/operation.go
type Qcow2Agent struct{}

func (a Qcow2Agent) qcow2IsExist(path *url.URL) bool {
	_, err := os.Stat(path.Path)

	return err == nil
}

// qemu-img create -f qcow2 $image $size
func (a Qcow2Agent) createQcow2(size uint64, path *url.URL) error {
	if err := os.MkdirAll(filepath.Dir(path.Path), os.ModePerm); err != nil {
		return err
	}

	args := []string{"qemu-img", "create", "-f", "qcow2", path.String(), bytefmt.ByteSize(size)}
	cmd := exec.Command(args[0], args[1:]...)
	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error message '%s', args '%s', output '%s'", err.Error(), a, o)
	}

	return nil
}

// rm $image
func (a Qcow2Agent) deleteQcow2(path *url.URL) error {
	if err := os.Remove(path.Path); err != nil {
		return fmt.Errorf("error message '%s'", err.Error())
	}

	return nil
}

func (a Qcow2Agent) ApplyQcow2(ctx context.Context, req *ApplyQcow2Request) (*Qcow2, error) {
	u, err := url.Parse(req.Qcow2.Url)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse url, err:%v.", err.Error())
	}

	if a.qcow2IsExist(u) {
		// need to implement resizing.
		return req.Qcow2, grpc.Errorf(codes.AlreadyExists, "")
	}

	if err := a.createQcow2(req.Qcow2.Bytes, u); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to create qcow2, err:%v.", err.Error())
	}
	log.Printf("[INFO] Applied qcow2, qcow2:%v", req.Qcow2)

	return req.Qcow2, nil
}

func (a Qcow2Agent) DownloadQcow2(ctx context.Context, req *DownloadQcow2Request) (*Qcow2, error) {
	u, err := url.Parse(req.Qcow2.Url)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse url, err:%v.", err.Error())
	}

	if a.qcow2IsExist(u) {
		// need to implement resizing.
		return req.Qcow2, grpc.Errorf(codes.AlreadyExists, "")
	}

	// Download qcow2 img
	response, err := http.Get(req.SourceUrl)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse SourceUrl, err:%v.", err.Error())
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to download img, err:%v.", err.Error())
	}
	file, err := os.OpenFile(req.Qcow2.Url, os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to save img, err:%v.", err.Error())
	}
	file.Write(body)
	log.Printf("[INFO] Download qcow2, qcow2:%v", req.Qcow2)
	return req.Qcow2, nil
}

func (a Qcow2Agent) BuildQcow2WithPacker(context.Context, *BuildQcow2WithPackerRequest) (*Qcow2, error) {
	return nil, nil
}

func (a Qcow2Agent) DeleteQcow2(ctx context.Context, req *DeleteQcow2Request) (*empty.Empty, error) {
	u, err := url.Parse(req.Qcow2.Url)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse url, err:%v.", err.Error())
	}

	if !a.qcow2IsExist(u) {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	a.deleteQcow2(u)
	log.Printf("[INFO] Deleted qcow2, qcow2:%v", req.Qcow2)

	return &empty.Empty{}, nil
}
