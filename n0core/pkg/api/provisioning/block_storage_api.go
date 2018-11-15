package provisioning

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
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
	AnnotationBlockStorageURL = "n0core/provisioning/block_storage_url"

	AnnotationBlockStorageReserve = "n0core/provisioning/block_storage_name"

	DownloadBlockStorageHTTPPrefix = "/api/block_storage/files"
)

func CreateBlockStorageAPI(ds datastore.Datastore, na ppool.NodeServiceClient) *BlockStorageAPI {
	nc := &node.NodeConnections{
		NodeAPI: na,
	}

	a := &BlockStorageAPI{
		dataStore:       ds,
		nodeAPI:         na,
		nodeConnections: nc,
	}
	a.dataStore.AddPrefix("block_storage")

	return a
}

func (a BlockStorageAPI) getBlockStorageAgent(n *ppool.Node) (BlockStorageAgentServiceClient, *grpc.ClientConn, error) {
	if n.State == ppool.Node_NotReady {
		return nil, nil, fmt.Errorf("node is not ready")
	}

	// port を何かから取れるようにする
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.Address, 20181), grpc.WithInsecure())
	if err != nil {
		return nil, nil, errors.Wrap(err, "Fail to dial to node")
	}

	cli := NewBlockStorageAgentServiceClient(conn)
	if err != nil {
		return nil, nil, err
	}

	return cli, conn, nil
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
	var v *CreateEmptyBlockStorageResponse

	conn, err := a.nodeConnections.GetConnection(res.NodeName) // errorについて考える
	cli := NewBlockStorageAgentServiceClient(conn)
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		goto ReleaseStorage
	}
	defer conn.Close()

	v, err = cli.CreateEmptyBlockStorage(context.Background(), &CreateEmptyBlockStorageRequest{
		Name:  req.Name,
		Bytes: req.LimitBytes,
	})
	if err != nil {
		log.Printf("Failed to create block_storage on node '%s': err='%s'", res.NodeName, err.Error()) // TODO: #89
		goto ReleaseStorage
	}

	res.Annotations[AnnotationBlockStorageURL] = (&url.URL{
		Scheme: "file",
		Path:   v.Path,
	}).String()
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteBlockStorage
	}

	return res, nil

DeleteBlockStorage:
	_, err = cli.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: v.Path})
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
	if _, err := url.Parse(req.SourceUrl); req.SourceUrl == "" || err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid source url")
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
	var v *FetchBlockStorageResponse

	conn, err := a.nodeConnections.GetConnection(res.NodeName) // errorについて考える
	cli := NewBlockStorageAgentServiceClient(conn)
	if err != nil {
		log.Printf("Fail to dial to node: err=%v.", err.Error())
		goto ReleaseStorage
	}
	defer conn.Close()

	v, err = cli.FetchBlockStorage(context.Background(), &FetchBlockStorageRequest{
		Name:      req.Name,
		Bytes:     req.LimitBytes,
		SourceUrl: req.SourceUrl,
	})
	if err != nil {
		log.Printf("Fail to create block_storage on node '%s': err='%s'", "", err.Error()) // TODO: #89
		goto ReleaseStorage
	}

	res.Annotations[AnnotationBlockStorageURL] = (&url.URL{
		Scheme: "file",
		Path:   v.Path,
	}).String()
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		goto DeleteBlockStorage
	}

	return res, nil

DeleteBlockStorage:
	_, err = cli.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: v.Path})
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
	switch {
	case req.Name == "":
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set 'name'")

	case req.RequestBytes == 0:
		return nil, grpc.Errorf(codes.InvalidArgument, "Set 'request_bytes'")
	}

	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if res.Name != "" {
		return nil, grpc.Errorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", req.Name)
	}

	src := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.SourceBlockStorage, src); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", src.Name)
	} else if src.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "BlockStorage '%s' is not exists", src.Name)
	}
	srcNode, err := a.nodeAPI.GetNode(ctx, &ppool.GetNodeRequest{Name: src.NodeName})
	if err != nil {
		return nil, WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to get node '%s' which is hosting source block storage", src.NodeName).Error())
	}
	srcUrl := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", srcNode.Address, 8081), // TODO: agentのポートをハードコードしている
		Path:   filepath.Join(DownloadBlockStorageHTTPPrefix, src.Name),
	}

	res.Name = req.Name
	res.Annotations = req.Annotations
	res.RequestBytes = req.RequestBytes
	res.LimitBytes = req.LimitBytes

	tx := transaction.Begin()

	res.StorageName = res.Name
	var dstNode *ppool.Node
	var ok bool
	if res.NodeName, ok = res.Annotations[AnnotationRequestNodeName]; !ok {
		dstNode, err = a.nodeAPI.ScheduleStorage(ctx, &ppool.ScheduleStorageRequest{
			StorageName: res.StorageName,
			Annotations: map[string]string{
				AnnotationBlockStorageReserve: res.Name,
			},
			RequestBytes: res.RequestBytes,
			LimitBytes:   res.LimitBytes,
		})
	} else {
		dstNode, err = a.nodeAPI.ReserveStorage(ctx, &ppool.ReserveStorageRequest{
			NodeName:    res.NodeName,
			StorageName: res.StorageName,
			Annotations: map[string]string{
				AnnotationBlockStorageReserve: res.Name,
			},
			RequestBytes: res.RequestBytes,
			LimitBytes:   res.LimitBytes,
		})
	}
	if err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, errors.Wrap(err, "Failed to reserve storage").Error())
	}
	tx.PushRollback("Release storage", func() error {
		_, err = a.nodeAPI.ReleaseStorage(ctx, &ppool.ReleaseStorageRequest{
			NodeName:    res.NodeName,
			StorageName: res.StorageName,
		})
		return err
	})

	cli, conn, err := a.getBlockStorageAgent(dstNode)
	if err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, errors.Wrap(err, "Failed to connect to agent").Error())
	}
	defer conn.Close()

	v, err := cli.FetchBlockStorage(ctx, &FetchBlockStorageRequest{
		Name:      req.Name,
		Bytes:     req.LimitBytes,
		SourceUrl: srcUrl.String(),
	})
	if err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to fetch block storage on node '%s'", res.NodeName).Error())
	}
	tx.PushRollback("Delete created block storage", func() error {
		_, err = cli.DeleteBlockStorage(ctx, &DeleteBlockStorageRequest{Path: v.Path})
		return err
	})

	res.Annotations[AnnotationBlockStorageURL] = (&url.URL{
		Scheme: "file",
		Path:   v.Path,
	}).String()
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		tx.Rollback()
		return nil, WrapGrpcErrorf(codes.Internal, errors.Wrap(err, "[WARNING] Failed to apply data for db").Error())
	}

	return res, nil
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

	u, _ := url.Parse(prev.Annotations[AnnotationBlockStorageURL]) // TODO: エラー処理
	_, err = cli.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: u.Path})
	if err != nil {
		log.Printf("Fail to delete block_storage on node: err=%s, req=%s", err.Error(), &DeleteBlockStorageRequest{Path: u.Path})
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
