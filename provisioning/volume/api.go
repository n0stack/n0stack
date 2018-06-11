package volume

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/n0stack/n0core/provisioning/node/qcow2"

	"github.com/n0stack/n0core/provisioning/node"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/datastore"
	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"
)

// Annotations:
//     n0core/url: Path of qcow2 that is stored on node.
//     n0core/node_name: Scheduled node. (空の場合は実装していない)
type VolumeAPI struct {
	dataStore datastore.Datastore
	nodeAPI   *node.NodeAPI

	baseDirectory string
}

func CreateVolumeAPI(ds datastore.Datastore, na *node.NodeAPI, baseDirectory string) (*VolumeAPI, error) {
	a := &VolumeAPI{
		dataStore: ds,
		nodeAPI:   na,

		baseDirectory: baseDirectory,
	}

	return a, nil
}

func (a *VolumeAPI) structureURL(name string) *url.URL {
	return &url.URL{
		Scheme: "file",
		Path:   filepath.Join(a.baseDirectory, name+".qcow2"),
	}
}

func (a *VolumeAPI) ListVolumes(ctx context.Context, req *pprovisioning.ListVolumesRequest) (*pprovisioning.ListVolumesResponse, error) {
	res := &pprovisioning.ListVolumesResponse{}
	f := func(s int) []proto.Message {
		res.Volumes = make([]*pprovisioning.Volume, s)
		for i := range res.Volumes {
			res.Volumes[i] = &pprovisioning.Volume{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Volumes {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}
	if len(res.Volumes) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *VolumeAPI) GetVolume(ctx context.Context, req *pprovisioning.GetVolumeRequest) (*pprovisioning.Volume, error) {
	res := &pprovisioning.Volume{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}

	if res.Metadata == nil {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *VolumeAPI) ApplyVolume(ctx context.Context, req *pprovisioning.ApplyVolumeRequest) (*pprovisioning.Volume, error) {
	res := &pprovisioning.Volume{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &pprovisioning.VolumeStatus{},
	}

	prev := &pprovisioning.Volume{}
	err := a.dataStore.Get(req.Metadata.Name, prev)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get db, got:%v.", err.Error())
	}
	if prev.Metadata == nil && req.Metadata.Version != 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set the same version as GetVolume result, have:%d, want:0.", req.Metadata.Version)
	}
	if prev.Metadata != nil && req.Metadata.Version != prev.Metadata.Version {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set the same version as GetVolume result, have:%d, want:%d.", req.Metadata.Version, prev.Metadata.Version)
	}

	nn, ok := res.Metadata.Annotations["n0core/node_name"]
	if !ok {
		if res.Status == nil {
			res.Status = &pprovisioning.VolumeStatus{}
		}

		res.Status.State = pprovisioning.VolumeStatus_PENDING

		return res, nil
	}

	// 切り出したい、こっから
	n, err := a.nodeAPI.GetNode(context.Background(), &pprovisioning.GetNodeRequest{Name: nn})
	if err != nil {
		return nil, err
	}

	// nodeがreadyか

	// portはendpointから取る
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Spec.Address, 20181), grpc.WithInsecure())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Fail to dial to node, err:%v.", err.Error())
	}
	defer conn.Close()
	cli := qcow2.NewQcow2ServiceClient(conn)
	// ここまで

	// urlのパースとInvalidArguments
	u, ok := res.Metadata.Annotations["n0core/url"]
	if !ok {
		u = a.structureURL(res.Metadata.Name).String()
	}

	_, err = cli.ApplyQcow2(context.Background(), &qcow2.ApplyQcow2Request{Qcow2: &qcow2.Qcow2{
		Bytes: res.Spec.Bytes,
		Url:   u,
	}})
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return nil, grpc.Errorf(codes.Internal, "Fail to apply qcow2 on node, err:%v.", err.Error())
	}

	res.Metadata.Annotations["n0core/url"] = u
	res.Status.State = pprovisioning.VolumeStatus_AVAILABLE

	res.Metadata.Version++
	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to apply for db, got:%v.", err.Error())
	}
	log.Printf("[INFO] On Applly, applied Volume:%v", res)

	return res, nil
}

func (a *VolumeAPI) DeleteVolume(ctx context.Context, req *pprovisioning.DeleteVolumeRequest) (*empty.Empty, error) {
	v := &pprovisioning.Volume{}

	if err := a.dataStore.Get(req.Name, v); err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "Failed to get from db.\tgot:%v", err.Error())
	}

	if v.Metadata == nil {
		return &empty.Empty{}, grpc.Errorf(codes.NotFound, "")
	}

	nn, ok := v.Metadata.Annotations["n0core/node_name"]
	if !ok {
		return &empty.Empty{}, grpc.Errorf(codes.InvalidArgument, "Set n0core/node_name in annotations.")
	}

	// 切り出したい、こっから
	n, err := a.nodeAPI.GetNode(context.Background(), &pprovisioning.GetNodeRequest{Name: nn})
	if err != nil {
		return &empty.Empty{}, err
	}

	// portはendpointから取る
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Spec.Address, 20181), grpc.WithInsecure())
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "Fail to dial to node, err:%v.", err.Error())
	}
	defer conn.Close()
	cli := qcow2.NewQcow2ServiceClient(conn)
	// ここまで

	u, ok := v.Metadata.Annotations["n0core/url"]
	if !ok {
		return &empty.Empty{}, grpc.Errorf(codes.InvalidArgument, "Set n0core/url in annotations.")
	}
	pu, err := url.Parse(u)
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.InvalidArgument, "Invalid n0core/url in annotations.")
	}

	_, err = cli.DeleteQcow2(context.Background(), &qcow2.DeleteQcow2Request{Qcow2: &qcow2.Qcow2{
		Bytes: v.Spec.Bytes,
		Url:   pu.String(),
	}})
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "Fail to apply qcow2 on node, err:%v.", err.Error())
	}

	d, err := a.dataStore.Delete(req.Name)
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}
	if d < 1 {
		return &empty.Empty{}, grpc.Errorf(codes.NotFound, "")
	}

	return &empty.Empty{}, nil
}
