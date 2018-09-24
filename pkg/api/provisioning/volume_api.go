package provisioning

import (
	"context"
	"log"
	"reflect"

	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/provisioning/v0"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0core/pkg/datastore"
)

type VolumeAPI struct {
	dataStore datastore.Datastore

	// dependency APIs
	nodeAPI ppool.NodeServiceClient

	nodeConnections *node.NodeConnections
}

const (
	// Create のときに自動生成、消されると困る
	AnnotationVolumePath = "n0core/provisioning/volume_url"
)

func CreateVolumeAPI(ds datastore.Datastore, na ppool.NodeServiceClient) (*VolumeAPI, error) {
	nc := &node.NodeConnections{
		NodeAPI: na,
	}

	a := &VolumeAPI{
		dataStore:       ds,
		nodeAPI:         na,
		nodeConnections: nc,
	}

	return a, nil
}

func (a *VolumeAPI) CreateEmptyVolume(ctx context.Context, req *pprovisioning.CreateEmptyVolumeRequest) (*pprovisioning.Volume, error) {
	prev := &pprovisioning.Volume{}
	if err := a.dataStore.Get(req.Metadata.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	} else if !reflect.ValueOf(prev.Metadata).IsNil() {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Metadata.Name)
	}

	res := &pprovisioning.Volume{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &pprovisioning.VolumeStatus{},
	}

	var err error
	if res.Status.NodeName, res.Status.StorageName, err = a.reserveStorage(
		req.Metadata.Name,
		req.Metadata.Annotations,
		req.Spec.RequestBytes,
		req.Spec.LimitBytes,
	); err != nil {
		return nil, err
	}
	var v *VolumeAgent

	conn, err := a.nodeConnections.GetConnection(res.Status.NodeName) // errorについて考える
	cli := NewVolumeAgentServiceClient(conn)
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		goto ReleaseStorage
	}
	defer conn.Close()

	v, err = cli.CreateEmptyVolumeAgent(context.Background(), &CreateEmptyVolumeAgentRequest{
		Name:  req.Metadata.Name,
		Bytes: req.Spec.LimitBytes,
	})
	if err != nil && status.Code(err) != codes.AlreadyExists {
		log.Printf("Fail to create volume on node '%s': err='%s'", "", err.Error()) // TODO: #89
		goto ReleaseStorage
	}

	res.Metadata.Annotations[AnnotationVolumePath] = v.Path
	res.Status.State = pprovisioning.VolumeStatus_AVAILABLE

	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteVolume
	}

	return res, nil

DeleteVolume:
	_, err = cli.DeleteVolumeAgent(context.Background(), &DeleteVolumeAgentRequest{Path: prev.Metadata.Annotations[AnnotationVolumePath]})
	if err != nil {
		log.Printf("Fail to delete volume on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete volume on node") // TODO #89
	}

ReleaseStorage:
	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		Name:        prev.Status.NodeName,
		StorageName: prev.Status.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", prev.Status.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if status.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", prev.Status.StorageName)
		}
	}

	return nil, grpc.Errorf(codes.Internal, "")
}

func (a *VolumeAPI) CreateVolumeWithDownloading(ctx context.Context, req *pprovisioning.CreateVolumeWithDownloadingRequest) (*pprovisioning.Volume, error) {
	prev := &pprovisioning.Volume{}
	if err := a.dataStore.Get(req.Metadata.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	} else if !reflect.ValueOf(prev.Metadata).IsNil() {
		return nil, grpc.Errorf(codes.AlreadyExists, "Volume '%s' is already exists", req.Metadata.Name)
	}

	res := &pprovisioning.Volume{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &pprovisioning.VolumeStatus{},
	}
	var v *VolumeAgent

	var err error
	if res.Status.NodeName, res.Status.StorageName, err = a.reserveStorage(
		req.Metadata.Name,
		req.Metadata.Annotations,
		req.Spec.RequestBytes,
		req.Spec.LimitBytes,
	); err != nil {
		return nil, err
	}

	conn, err := a.nodeConnections.GetConnection(res.Status.NodeName) // errorについて考える
	cli := NewVolumeAgentServiceClient(conn)
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		goto ReleaseStorage
	}
	defer conn.Close()

	v, err = cli.CreateVolumeAgentWithDownloading(context.Background(), &CreateVolumeAgentWithDownloadingRequest{
		Name:      req.Metadata.Name,
		Bytes:     req.Spec.LimitBytes,
		SourceUrl: req.SourceUrl,
	})
	if err != nil && status.Code(err) != codes.AlreadyExists {
		log.Printf("Fail to create volume on node '%s': err='%s'", "", err.Error()) // TODO: #89
		goto ReleaseStorage
	}

	res.Metadata.Annotations[AnnotationVolumePath] = v.Path
	res.Status.State = pprovisioning.VolumeStatus_AVAILABLE

	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteVolume
	}

	return res, nil

DeleteVolume:
	_, err = cli.DeleteVolumeAgent(context.Background(), &DeleteVolumeAgentRequest{Path: prev.Metadata.Annotations[AnnotationVolumePath]})
	if err != nil {
		log.Printf("Fail to delete volume on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete volume on node") // TODO #89
	}

ReleaseStorage:
	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		Name:        prev.Status.NodeName,
		StorageName: prev.Status.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", prev.Status.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if status.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", prev.Status.StorageName)
		}
	}

	return nil, grpc.Errorf(codes.Internal, "")
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
		log.Printf("[WARNING] Failed to list data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.Volumes) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *VolumeAPI) GetVolume(ctx context.Context, req *pprovisioning.GetVolumeRequest) (*pprovisioning.Volume, error) {
	res := &pprovisioning.Volume{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if reflect.ValueOf(res.Metadata).IsNil() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *VolumeAPI) UpdateVolume(ctx context.Context, req *pprovisioning.UpdateVolumeRequest) (*pprovisioning.Volume, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")

	// res := &pprovisioning.Volume{
	// 	Metadata: req.Metadata,
	// 	Spec:     req.Spec,
	// 	Status:   &pprovisioning.VolumeStatus{},
	// }

	// prev := &pprovisioning.Volume{}
	// if err := a.dataStore.Get(req.Metadata.Name, prev); err != nil {
	// 	log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
	// 	return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	// }
	// var err error
	// res.Metadata.Version, err = datastore.CheckVersion(prev, req)
	// if err != nil {
	// 	return nil, grpc.Errorf(codes.InvalidArgument, "Failed to check version: %s", err.Error())
	// }

	// if prev.Spec.RequestBytes != req.Spec.RequestBytes {
	// 	return nil, grpc.Errorf(codes.Unimplemented, "Not supported change size: from='%s', to='%s'", prev.Spec.RequestBytes, req.Spec.RequestBytes)
	// }
	// if prev.Spec.LimitBytes != req.Spec.LimitBytes {
	// 	return nil, grpc.Errorf(codes.Unimplemented, "Not supported change size: from='%s', to='%s'", prev.Spec.LimitBytes, req.Spec.LimitBytes)
	// }

	// 	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
	// 		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
	// 		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	// 	}

	// return res, nil
}

func (a *VolumeAPI) DeleteVolume(ctx context.Context, req *pprovisioning.DeleteVolumeRequest) (*empty.Empty, error) {
	prev := &pprovisioning.Volume{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if !reflect.ValueOf(prev.Metadata).IsNil() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	if prev.Status.State != pprovisioning.VolumeStatus_AVAILABLE {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Volume '%s' is not available", req.Name)
	}

	conn, err := a.nodeConnections.GetConnection(prev.Status.NodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.Status.NodeName)
	}
	defer conn.Close()
	cli := NewVolumeAgentServiceClient(conn)

	_, err = cli.DeleteVolumeAgent(context.Background(), &DeleteVolumeAgentRequest{Path: prev.Metadata.Annotations[AnnotationVolumePath]})
	if err != nil {
		log.Printf("Fail to delete volume on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete volume on node") // TODO #89
	}

	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		Name:        prev.Status.NodeName,
		StorageName: prev.Status.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", prev.Status.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if status.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", prev.Status.StorageName)
		}
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	return &empty.Empty{}, nil
}

func (a *VolumeAPI) SetInuseVolume(ctx context.Context, req *pprovisioning.SetInuseVolumeRequest) (*pprovisioning.Volume, error) {
	res := &pprovisioning.Volume{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	res.Status.State = pprovisioning.VolumeStatus_IN_USE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *VolumeAPI) SetAvailableVolume(ctx context.Context, req *pprovisioning.SetAvailableVolumeRequest) (*pprovisioning.Volume, error) {
	res := &pprovisioning.Volume{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	res.Status.State = pprovisioning.VolumeStatus_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return nil, nil
}
