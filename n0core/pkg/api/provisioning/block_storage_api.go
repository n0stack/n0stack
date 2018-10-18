package provisioning

import (
	"context"
	"log"
	"net/url"

	"github.com/n0stack/n0stack/n0proto/pool/v0"
	"github.com/n0stack/n0stack/n0proto/provisioning/v0"
	"github.com/pkg/errors"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
)

type BlockStorageAPI struct {
	dataStore datastore.Datastore

	// dependency APIs
	nodeAPI ppool.NodeServiceClient

	nodeConnections *node.NodeConnections
}

const (
	// Create のときに自動生成、消されると困る
	AnnotationBlockStoragePath = "n0core/provisioning/block_storage_url"

	AnnotationBlockStorageReserve = "n0core/provisioning/block_storage_name"
)

func CreateBlockStorageAPI(ds datastore.Datastore, na ppool.NodeServiceClient) (*BlockStorageAPI, error) {
	nc := &node.NodeConnections{
		NodeAPI: na,
	}

	a := &BlockStorageAPI{
		dataStore:       ds,
		nodeAPI:         na,
		nodeConnections: nc,
	}

	return a, nil
}

func (a *BlockStorageAPI) CreateBlockStorage(ctx context.Context, req *pprovisioning.CreateBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if req.RequestBytes == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set 'request_bytes'")
	}

	prev := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if prev.Name != "" {
		return nil, grpc.Errorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", req.Name)
	}

	res := &pprovisioning.BlockStorage{
		Name:         req.Name,
		Annotations:  req.Annotations,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}

	var err error
	if res.NodeName, res.StorageName, err = a.reserveStorage(
		req.Name,
		req.Annotations,
		req.RequestBytes,
		req.LimitBytes,
	); err != nil {
		return nil, errors.Wrap(err, "Failed to reserve storage")
	}
	var v *BlockStorageAgent

	conn, err := a.nodeConnections.GetConnection(res.NodeName) // errorについて考える
	cli := NewBlockStorageAgentServiceClient(conn)
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		goto ReleaseStorage
	}
	defer conn.Close()

	v, err = cli.CreateEmptyBlockStorageAgent(context.Background(), &CreateEmptyBlockStorageAgentRequest{
		Name:  req.Name,
		Bytes: req.LimitBytes,
	})
	if err != nil && grpc.Code(err) != codes.AlreadyExists {
		log.Printf("Failed to create block_storage on node '%s': err='%s'", res.NodeName, err.Error()) // TODO: #89
		goto ReleaseStorage
	}

	res.Annotations[AnnotationBlockStoragePath] = v.Path
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteBlockStorage
	}

	return res, nil

DeleteBlockStorage:
	_, err = cli.DeleteBlockStorageAgent(context.Background(), &DeleteBlockStorageAgentRequest{Path: res.Annotations[AnnotationBlockStoragePath]})
	if err != nil {
		log.Printf("Fail to delete block_storage on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete block_storage on node") // TODO #89
	}

ReleaseStorage:
	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		NodeName:    res.NodeName,
		StorageName: res.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", res.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", res.StorageName)
		}
	}

	return nil, grpc.Errorf(codes.Internal, "")
}

func (a *BlockStorageAPI) FetchBlockStorage(ctx context.Context, req *pprovisioning.FetchBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if req.RequestBytes == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set 'request_bytes'")
	}
	// parse url

	prev := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if prev.Name != "" {
		return nil, grpc.Errorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", req.Name)
	}

	res := &pprovisioning.BlockStorage{
		Name:         req.Name,
		Annotations:  req.Annotations,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}

	var err error
	if res.NodeName, res.StorageName, err = a.reserveStorage(
		req.Name,
		req.Annotations,
		req.RequestBytes,
		req.LimitBytes,
	); err != nil {
		return nil, errors.Wrap(err, "Failed to reserve storage")
	}
	var v *BlockStorageAgent

	conn, err := a.nodeConnections.GetConnection(res.NodeName) // errorについて考える
	cli := NewBlockStorageAgentServiceClient(conn)
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		goto ReleaseStorage
	}
	defer conn.Close()

	v, err = cli.CreateBlockStorageAgentWithDownloading(context.Background(), &CreateBlockStorageAgentWithDownloadingRequest{
		Name:      req.Name,
		Bytes:     req.LimitBytes,
		SourceUrl: req.SourceUrl,
	})
	if err != nil && grpc.Code(err) != codes.AlreadyExists {
		log.Printf("Fail to create block_storage on node '%s': err='%s'", "", err.Error()) // TODO: #89
		goto ReleaseStorage
	}

	res.Annotations[AnnotationBlockStoragePath] = v.Path
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteBlockStorage
	}

	return res, nil

DeleteBlockStorage:
	_, err = cli.DeleteBlockStorageAgent(context.Background(), &DeleteBlockStorageAgentRequest{Path: res.Annotations[AnnotationBlockStoragePath]})
	if err != nil {
		log.Printf("Fail to delete block_storage on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete block_storage on node") // TODO #89
	}

ReleaseStorage:
	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		NodeName:    res.NodeName,
		StorageName: res.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", res.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", res.StorageName)
		}
	}

	return nil, grpc.Errorf(codes.Internal, "")
}

func (a *BlockStorageAPI) CopyBlockStorage(ctx context.Context, req *pprovisioning.CopyBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if req.RequestBytes == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set 'request_bytes'")
	}

	prev := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if prev.Name != "" {
		return nil, grpc.Errorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", req.Name)
	}

	src := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.SourceBlockStorage, src); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", src.Name)
	} else if src.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "BlockStorage '%s' is not exists", src.Name)
	}

	res := &pprovisioning.BlockStorage{
		Name:         req.Name,
		Annotations:  req.Annotations,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}

	if res.Annotations == nil {
		res.Annotations = make(map[string]string)
	}
	if v, ok := res.Annotations[AnnotationRequestNodeName]; ok {
		if src.Annotations[AnnotationRequestNodeName] != v {
			return nil, grpc.Errorf(codes.InvalidArgument, "Set annotation, about request, node same as src")
		}
	} else {
		res.Annotations[AnnotationRequestNodeName] = src.Annotations[AnnotationRequestNodeName]
	}

	var err error
	if res.NodeName, res.StorageName, err = a.reserveStorage(
		req.Name,
		res.Annotations,
		req.RequestBytes,
		req.LimitBytes,
	); err != nil {
		return nil, errors.Wrap(err, "Failed to reserve storage")
	}
	var v *BlockStorageAgent
	srcUrl := url.URL{
		Scheme: "file",
		Path:   res.Annotations[AnnotationBlockStoragePath],
	}

	conn, err := a.nodeConnections.GetConnection(res.NodeName) // errorについて考える
	cli := NewBlockStorageAgentServiceClient(conn)
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		goto ReleaseStorage
	}
	defer conn.Close()

	v, err = cli.CreateBlockStorageAgentWithDownloading(context.Background(), &CreateBlockStorageAgentWithDownloadingRequest{
		Name:      req.Name,
		Bytes:     req.LimitBytes,
		SourceUrl: srcUrl.String(),
	})
	if err != nil && grpc.Code(err) != codes.AlreadyExists {
		log.Printf("Fail to create block_storage on node '%s': err='%s'", "", err.Error()) // TODO: #89
		goto ReleaseStorage
	}

	res.Annotations[AnnotationBlockStoragePath] = v.Path
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteBlockStorage
	}

	return res, nil

DeleteBlockStorage:
	_, err = cli.DeleteBlockStorageAgent(context.Background(), &DeleteBlockStorageAgentRequest{Path: res.Annotations[AnnotationBlockStoragePath]})
	if err != nil {
		log.Printf("Fail to delete block_storage on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete block_storage on node") // TODO #89
	}

ReleaseStorage:
	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		NodeName:    res.NodeName,
		StorageName: res.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", res.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", res.StorageName)
		}
	}

	return nil, grpc.Errorf(codes.Internal, "")
}

func (a *BlockStorageAPI) ListBlockStorages(ctx context.Context, req *pprovisioning.ListBlockStoragesRequest) (*pprovisioning.ListBlockStoragesResponse, error) {
	res := &pprovisioning.ListBlockStoragesResponse{}
	f := func(s int) []proto.Message {
		res.BlockStorages = make([]*pprovisioning.BlockStorage, s)
		for i := range res.BlockStorages {
			res.BlockStorages[i] = &pprovisioning.BlockStorage{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.BlockStorages {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		log.Printf("[WARNING] Failed to list data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.BlockStorages) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *BlockStorageAPI) GetBlockStorage(ctx context.Context, req *pprovisioning.GetBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *BlockStorageAPI) UpdateBlockStorage(ctx context.Context, req *pprovisioning.UpdateBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")

	// res := &pprovisioning.BlockStorage{
	// 	Metadata: req.Metadata,
	// 	Spec:     req.Spec,
	// 	Status:   &pprovisioning.BlockStorageStatus{},
	// }

	// prev := &pprovisioning.BlockStorage{}
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

func (a *BlockStorageAPI) DeleteBlockStorage(ctx context.Context, req *pprovisioning.DeleteBlockStorageRequest) (*empty.Empty, error) {
	prev := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if prev.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	if prev.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpc.Errorf(codes.FailedPrecondition, "BlockStorage '%s' is not available", req.Name)
	}

	conn, err := a.nodeConnections.GetConnection(prev.NodeName)
	if err != nil {
		log.Printf("[WARNING] Fail to dial to node: err=%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "") // TODO: #89
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", prev.NodeName)
	}
	defer conn.Close()
	cli := NewBlockStorageAgentServiceClient(conn)

	_, err = cli.DeleteBlockStorageAgent(context.Background(), &DeleteBlockStorageAgentRequest{Path: prev.Annotations[AnnotationBlockStoragePath]})
	if err != nil {
		log.Printf("Fail to delete block_storage on node, err:%v.", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Fail to delete block_storage on node") // TODO #89
	}

	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		NodeName:    prev.NodeName,
		StorageName: prev.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", prev.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return nil, grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", prev.StorageName)
		}
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	return &empty.Empty{}, nil
}

func (a *BlockStorageAPI) SetAvailableBlockStorage(ctx context.Context, req *pprovisioning.SetAvailableBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	if res.State == pprovisioning.BlockStorage_UNKNOWN {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Can change state to AVAILABLE when UNKNOWN")
	}
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *BlockStorageAPI) SetInuseBlockStorage(ctx context.Context, req *pprovisioning.SetInuseBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	if res.State != pprovisioning.BlockStorage_AVAILABLE && res.State != pprovisioning.BlockStorage_IN_USE {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Can change state to IN_USE when only AVAILABLE")
	}
	res.State = pprovisioning.BlockStorage_IN_USE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *BlockStorageAPI) SetProtectedBlockStorage(ctx context.Context, req *pprovisioning.SetProtectedBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	if res.State != pprovisioning.BlockStorage_AVAILABLE && res.State != pprovisioning.BlockStorage_PROTECTED {
		return nil, grpc.Errorf(codes.FailedPrecondition, "Can change state to PROTECTED when only AVAILABLE")
	}
	res.State = pprovisioning.BlockStorage_PROTECTED

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

// func (a *BlockStorageAPI) UploadBlockStorage(req pprovisioning.BlockStorageService_UploadBlockStorageServer) error {
// 	return grpc.Errorf(codes.Unimplemented, "")
// }

// func (a *BlockStorageAPI) DownloadBlockStorage(req *pprovisioning.DownloadBlockStorageRequest, stream pprovisioning.BlockStorageService_DownloadBlockStorageServer) error {
// 	return grpc.Errorf(codes.Unimplemented, "")
// }

func (a BlockStorageAPI) reserveStorage(name string, annotations map[string]string, req, limit uint64) (string, string, error) {
	var n *ppool.Node
	var err error
	if node, ok := annotations[AnnotationRequestNodeName]; !ok {
		n, err = a.nodeAPI.ScheduleStorage(context.Background(), &ppool.ScheduleStorageRequest{
			StorageName: name,
			Annotations: map[string]string{
				AnnotationBlockStorageReserve: name,
			},
			RequestBytes: req,
			LimitBytes:   limit,
		})
	} else {
		n, err = a.nodeAPI.ReserveStorage(context.Background(), &ppool.ReserveStorageRequest{
			NodeName: node,
			Annotations: map[string]string{
				AnnotationBlockStorageReserve: name,
			},
			StorageName:  name,
			RequestBytes: req,
			LimitBytes:   limit,
		})
	}
	if err != nil {
		return "", "", err // TODO: #89
	}

	return n.Name, name, nil
}
