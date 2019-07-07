package blockstorage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/labstack/echo"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	stdapi "github.com/n0stack/n0stack/n0core/pkg/api/standard_api"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
	pprovisioning "github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const DownloadBlockStorageHTTPPrefix = "/api/block_storage/files"

type BlockStorageAPI struct {
	dataStore datastore.Datastore

	// dependency APIs
	nodeAPI ppool.NodeServiceClient

	getAgent func(ctx context.Context, nodeName string) (BlockStorageAgentServiceClient, func() error, error)
}

func CreateBlockStorageAPI(ds datastore.Datastore, na ppool.NodeServiceClient) *BlockStorageAPI {
	a := &BlockStorageAPI{
		dataStore: ds.AddPrefix("block_storage"),
		nodeAPI:   na,
	}

	a.getAgent = func(ctx context.Context, nodeName string) (BlockStorageAgentServiceClient, func() error, error) {
		conn, err := node.GetConnection(ctx, a.nodeAPI, nodeName)
		cli := NewBlockStorageAgentServiceClient(conn)
		if err != nil {
			return nil, nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to dial to node: err=%s", err.Error())
		}
		if conn == nil {
			return nil, nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Node '%s' is not ready, so cannot delete: please wait a moment", nodeName)
		}

		return cli, conn.Close, nil
	}

	return a
}

// Reserve a block storage on a node using NodeServiceClient.
func ReserveStorage(ctx context.Context, tx *transaction.Transaction, na ppool.NodeServiceClient, bs *pprovisioning.BlockStorage) error {
	bs.StorageName = bs.Name

	var n *ppool.Node
	var err error
	if node, ok := bs.Annotations[AnnotationBlockStorageRequestNodeName]; !ok {
		n, err = na.ScheduleStorage(ctx, &ppool.ScheduleStorageRequest{
			StorageName: bs.StorageName,
			Annotations: map[string]string{
				AnnotationStorageReservedBy: bs.Name,
			},
			RequestBytes: bs.RequestBytes,
			LimitBytes:   bs.LimitBytes,
		})
	} else {
		n, err = na.ReserveStorage(ctx, &ppool.ReserveStorageRequest{
			NodeName:    node,
			StorageName: bs.StorageName,
			Annotations: map[string]string{
				AnnotationStorageReservedBy: bs.Name,
			},
			RequestBytes: bs.RequestBytes,
			LimitBytes:   bs.LimitBytes,
		})
	}
	if err != nil {
		return grpcutil.WrapGrpcErrorf(grpc.Code(err), "Failed to ReserveStorage: desc=%s", grpc.ErrorDesc(err))
	}

	bs.NodeName = n.Name

	tx.PushRollback("release stroage", func() error {
		_, err := na.ReleaseStorage(ctx, &ppool.ReleaseStorageRequest{
			NodeName:    bs.NodeName,
			StorageName: bs.StorageName,
		})

		return err
	})

	return nil
}

// Release a block storage on a node using NodeServiceClient.
func ReleaseStorage(ctx context.Context, tx *transaction.Transaction, na ppool.NodeServiceClient, bs *pprovisioning.BlockStorage) error {
	_, err := na.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		NodeName:    bs.NodeName,
		StorageName: bs.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", bs.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to release compute '%s': please retry", bs.StorageName)
		}
	}
	tx.PushRollback("reserve storage", func() error {
		_, err = na.ReserveStorage(context.Background(), &ppool.ReserveStorageRequest{
			NodeName:    bs.NodeName,
			StorageName: bs.StorageName,
			Annotations: map[string]string{
				AnnotationStorageReservedBy: bs.Name,
			},
			RequestBytes: bs.RequestBytes,
			LimitBytes:   bs.LimitBytes,
		})

		return err
	})

	return nil
}

func (a *BlockStorageAPI) CreateBlockStorage(ctx context.Context, req *pprovisioning.CreateBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	{
		if req.RequestBytes == 0 {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set 'request_bytes'")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	if err := PendNewBlockStorage(tx, a.dataStore, req.Name); err != nil {
		return nil, err
	}

	bs := &pprovisioning.BlockStorage{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}

	if err := ReserveStorage(ctx, tx, a.nodeAPI, bs); err != nil {
		return nil, err
	}

	{
		cli, done, err := a.getAgent(ctx, bs.NodeName)
		if err != nil {
			return nil, err
		}
		defer done()

		v, err := cli.CreateEmptyBlockStorage(ctx, &CreateEmptyBlockStorageRequest{
			Name:  bs.Name,
			Bytes: bs.LimitBytes,
		})
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to create block_storage on node '%s': err='%s'", bs.NodeName, err.Error())
		}
		tx.PushRollback("DeleteBlockStorage", func() error {
			cli, done, err := a.getAgent(ctx, bs.NodeName)
			if err != nil {
				return err
			}
			defer done()

			_, err = cli.DeleteBlockStorage(ctx, &DeleteBlockStorageRequest{Path: v.Path})
			if err != nil {
				return err
			}

			return nil
		})

		bs.Annotations[AnnotationBlockStorageURL] = (&url.URL{
			Scheme: "file",
			Path:   v.Path,
		}).String()
		bs.State = pprovisioning.BlockStorage_AVAILABLE
	}

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) FetchBlockStorage(ctx context.Context, req *pprovisioning.FetchBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	{
		u, err := url.Parse(req.SourceUrl)
		if req.SourceUrl == "" || err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "invalid source url")
		}
		switch u.Scheme {
		case "http", "https", "ftp", "ftps", "file":
		default:
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "invalid source url: %s is not supported", u.Scheme)
		}

		if req.RequestBytes == 0 {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set 'request_bytes'")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	if err := PendNewBlockStorage(tx, a.dataStore, req.Name); err != nil {
		return nil, err
	}

	bs := &pprovisioning.BlockStorage{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}
	if bs.Annotations == nil {
		bs.Annotations = make(map[string]string)
	}
	bs.Annotations[AnnotationBlockStorageFetchFrom] = req.SourceUrl

	if err := ReserveStorage(ctx, tx, a.nodeAPI, bs); err != nil {
		return nil, err
	}

	if err := a.fetchBlockStorage(ctx, tx, bs, &FetchBlockStorageRequest{
		Name:      bs.Name,
		Bytes:     bs.LimitBytes,
		SourceUrl: bs.Annotations[AnnotationBlockStorageFetchFrom],
	}); err != nil {
		return nil, err
	}

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) fetchBlockStorage(ctx context.Context, tx *transaction.Transaction, bs *pprovisioning.BlockStorage, req *FetchBlockStorageRequest) error {
	cli, done, err := a.getAgent(ctx, bs.NodeName)
	if err != nil {
		return err
	}
	defer done()

	v, err := cli.FetchBlockStorage(context.Background(), req)
	if err != nil {
		return grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to FetchBlockStorage on node '%s': err='%s'", bs.NodeName, err.Error())
	}
	tx.PushRollback("DeleteBlockStorage", func() error {
		cli, done, err := a.getAgent(ctx, bs.NodeName)
		if err != nil {
			return err
		}
		defer done()

		_, err = cli.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: v.Path})
		if err != nil {
			return err
		}

		return nil
	})

	bs.Annotations[AnnotationBlockStorageURL] = (&url.URL{
		Scheme: "file",
		Path:   v.Path,
	}).String()
	bs.State = pprovisioning.BlockStorage_AVAILABLE

	return nil
}

func (a *BlockStorageAPI) CopyBlockStorage(ctx context.Context, req *pprovisioning.CopyBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	{
		if req.RequestBytes == 0 {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set 'request_bytes'")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	if err := PendNewBlockStorage(tx, a.dataStore, req.Name); err != nil {
		return nil, err
	}

	bs := &pprovisioning.BlockStorage{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}
	if bs.Annotations == nil {
		bs.Annotations = make(map[string]string)
	}
	bs.Annotations[AnnotationBlockStorageCopyFrom] = req.SourceBlockStorage

	if err := ReserveStorage(ctx, tx, a.nodeAPI, bs); err != nil {
		return nil, err
	}

	var srcUrl *url.URL
	{
		src := &pprovisioning.BlockStorage{}
		if err := a.dataStore.Get(req.SourceBlockStorage, src); err != nil {
			if datastore.IsNotFound(err) {
				return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, err.Error())
			}

			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, datastore.DefaultErrorMessage(err))
		}

		srcNode, err := a.nodeAPI.GetNode(ctx, &ppool.GetNodeRequest{Name: src.NodeName})
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to get node '%s' which is hosting source block storage", src.NodeName).Error())
		}

		srcUrl, err = url.Parse(src.Annotations[AnnotationBlockStorageURL])
		if err != nil || src.NodeName != bs.NodeName {
			srcUrl = &url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", srcNode.Address, 8081), // TODO: agentのポートをハードコードしている
				Path:   filepath.Join(DownloadBlockStorageHTTPPrefix, src.Name),
			}
		}
	}

	if err := a.fetchBlockStorage(ctx, tx, bs, &FetchBlockStorageRequest{
		Name:      req.Name,
		Bytes:     req.LimitBytes,
		SourceUrl: srcUrl.String(),
	}); err != nil {
		return nil, err
	}

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) ListBlockStorages(ctx context.Context, req *pprovisioning.ListBlockStoragesRequest) (*pprovisioning.ListBlockStoragesResponse, error) {
	return ListBlockStorages(ctx, req, a.dataStore)
}

func (a *BlockStorageAPI) GetBlockStorage(ctx context.Context, req *pprovisioning.GetBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	return GetBlockStorage(ctx, req, a.dataStore)
}

// Update a blocok storage: update its request bytes, limit bytes, and request node, a node to host the block storage.
func (a *BlockStorageAPI) UpdateBlockStorage(ctx context.Context, req *pprovisioning.UpdateBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := GetAndPendExistingBlockStorage(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "BlockStorage '%s' is not available: now=%d", req.Name, bs.State)
	}

	newBs := &pprovisioning.BlockStorage{
		Name:         bs.Name,
		Annotations:  make(map[string]string),
		RequestBytes: bs.RequestBytes,
		LimitBytes:   bs.LimitBytes,
	}

	if req.RequestBytes > 0 {
		newBs.RequestBytes = req.RequestBytes
	} else {
		newBs.RequestBytes = bs.RequestBytes
	}

	if req.LimitBytes > 0 {
		newBs.LimitBytes = req.LimitBytes
	} else {
		newBs.LimitBytes = bs.LimitBytes
	}

	if an, ok := bs.Annotations[AnnotationBlockStorageFetchFrom]; ok {
		newBs.Annotations[AnnotationBlockStorageFetchFrom] = an
	}

	if an, ok := bs.Annotations[AnnotationBlockStorageCopyFrom]; ok {
		newBs.Annotations[AnnotationBlockStorageCopyFrom] = an
	}

	if node, ok := req.Annotations[AnnotationBlockStorageRequestNodeName]; ok {
		newBs.Annotations[AnnotationBlockStorageRequestNodeName] = node
	} else if node, ok = bs.Annotations[AnnotationBlockStorageRequestNodeName]; ok {
		newBs.Annotations[AnnotationBlockStorageRequestNodeName] = node
	}

	if node, ok := newBs.Annotations[AnnotationBlockStorageRequestNodeName]; ok && node != bs.NodeName {
		if err := ReserveStorage(ctx, tx, a.nodeAPI, newBs); err != nil {
			return nil, err
		}

		srcNode, err := a.nodeAPI.GetNode(ctx, &ppool.GetNodeRequest{Name: bs.NodeName})
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to get node '%s' which is hosting source block storage", bs.NodeName).Error())
		}

		srcUrl := &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", srcNode.Address, 8081),
			Path:   filepath.Join(DownloadBlockStorageHTTPPrefix, bs.Name),
		}

		if err := a.fetchBlockStorage(ctx, tx, newBs, &FetchBlockStorageRequest{
			Name:      newBs.Name,
			Bytes:     newBs.LimitBytes,
			SourceUrl: srcUrl.String(),
		}); err != nil {
			return nil, err
		}

		if err := a.deleteBlockStorage(ctx, tx, bs); err != nil {
			return nil, err
		}
	} else {
		newBs.NodeName = bs.NodeName
		newBs.StorageName = bs.StorageName
		newBs.Annotations[AnnotationBlockStorageURL] = bs.Annotations[AnnotationBlockStorageURL]

		if newBs.LimitBytes != bs.LimitBytes {
			/*
				TODO: Implement ResizeStorage on NodeAPI:
				Calling ReserveStorage after ReleaseStorage may result in failure:
				Another thread can reserve a storage of the same name before ReserveStorage is called by this thread.
			*/
			if err := ReleaseStorage(ctx, tx, a.nodeAPI, bs); err != nil {
				return nil, err
			}
			if err := ReserveStorage(ctx, tx, a.nodeAPI, newBs); err != nil {
				return nil, err
			}

			cli, done, err := a.getAgent(ctx, newBs.NodeName)
			if err != nil {
				return nil, err
			}
			defer done()

			u, _ := url.Parse(newBs.Annotations[AnnotationBlockStorageURL]) // TODO: エラー処理
			if _, err := cli.ResizeBlockStorage(ctx, &ResizeBlockStorageRequest{
				Bytes: newBs.LimitBytes,
				Path:  u.Path,
			}); err != nil {
				return nil, err
			}

		}

		newBs.State = pprovisioning.BlockStorage_AVAILABLE
	}

	if err := ApplyBlockStorage(a.dataStore, newBs); err != nil {
		return nil, err
	}

	tx.Commit()
	return newBs, nil
}

func (a *BlockStorageAPI) DeleteBlockStorage(ctx context.Context, req *pprovisioning.DeleteBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := GetAndPendExistingBlockStorage(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Cannot change state to DELETED when not AVAILABLE")
	}
	bs.State = pprovisioning.BlockStorage_DELETED

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) UndeleteBlockStorage(ctx context.Context, req *pprovisioning.UndeleteBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := GetAndPendExistingBlockStorage(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State != pprovisioning.BlockStorage_DELETED {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Cannot change state to AVAILABLE when not DELETED")
	}
	bs.State = pprovisioning.BlockStorage_AVAILABLE

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) PurgeBlockStorage(ctx context.Context, req *pprovisioning.PurgeBlockStorageRequest) (*empty.Empty, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := GetAndPendExistingBlockStorage(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State != pprovisioning.BlockStorage_DELETED {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "BlockStorage '%s' is not deleted: now=%d", req.Name, bs.State)
	}

	if err := a.deleteBlockStorage(ctx, tx, bs); err != nil {
		return nil, err
	}

	if err := DeleteBlockStorage(a.dataStore, bs.Name); err != nil {
		return nil, err
	}

	tx.Commit()
	return &empty.Empty{}, nil
}

func (a *BlockStorageAPI) deleteBlockStorage(ctx context.Context, tx *transaction.Transaction, bs *pprovisioning.BlockStorage) error {
	cli, done, err := a.getAgent(ctx, bs.NodeName)
	if err != nil {
		return err
	}
	defer done()

	u, _ := url.Parse(bs.Annotations[AnnotationBlockStorageURL]) // TODO: エラー処理
	_, err = cli.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: u.Path})
	if err != nil { // 多分ロールバックが必要
		log.Printf("Fail to delete block_storage on node: err=%s, req=%s", err.Error(), &DeleteBlockStorageRequest{Path: u.Path})
		return grpcutil.WrapGrpcErrorf(codes.Internal, "Fail to delete block_storage on node") // TODO #89
	}

	if err := ReleaseStorage(ctx, tx, a.nodeAPI, bs); err != nil {
		return err
	}

	return nil
}

func (a *BlockStorageAPI) SetAvailableBlockStorage(ctx context.Context, req *pprovisioning.SetAvailableBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !lock.WaitUntilLock(a.dataStore, req.Name, 5*time.Second, 10*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := GetAndPendExistingBlockStorage(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State == pprovisioning.BlockStorage_PENDING {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Cannot change state to AVAILABLE when not UNKNOWN")
	}
	bs.State = pprovisioning.BlockStorage_AVAILABLE

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) SetInuseBlockStorage(ctx context.Context, req *pprovisioning.SetInuseBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !lock.WaitUntilLock(a.dataStore, req.Name, 5*time.Second, 10*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := GetAndPendExistingBlockStorage(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Cannot change state to IN_USE when not AVAILABLE")
	}
	bs.State = pprovisioning.BlockStorage_IN_USE

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) SetProtectedBlockStorage(ctx context.Context, req *pprovisioning.SetProtectedBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !lock.WaitUntilLock(a.dataStore, req.Name, 5*time.Second, 10*time.Millisecond) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := GetAndPendExistingBlockStorage(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Cannot change state to PROTECTED when not AVAILABLE")
	}
	bs.State = pprovisioning.BlockStorage_PROTECTED

	if err := ApplyBlockStorage(a.dataStore, bs); err != nil {
		return nil, err
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) DownloadBlockStorage(ctx context.Context, req *pprovisioning.DownloadBlockStorageRequest) (*pprovisioning.DownloadBlockStorageResponse, error) {
	bs := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, bs); err != nil {
		if datastore.IsNotFound(err) {
			return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, err.Error())
		}

		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	u := &url.URL{
		Scheme: "http",
		Path:   fmt.Sprintf("/n0core/api/v0/block_storage/download/%s", req.Name),
	}
	res := &pprovisioning.DownloadBlockStorageResponse{
		DownloadUrl: u.String(),
	}

	return res, nil
}

// TODO: agentPort は Node から取れるようにしたい
func (a *BlockStorageAPI) ProxyDownloadBlockStorage(agentPort int, basePath string) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")

		bs := &pprovisioning.BlockStorage{}
		if err := a.dataStore.Get(name, bs); err != nil {
			if datastore.IsNotFound(err) {
				return c.String(http.StatusNotFound, err.Error())
			}

			log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
			return c.String(http.StatusInternalServerError, datastore.DefaultErrorMessage(err))
		}

		ctx := context.Background()
		node, err := a.nodeAPI.GetNode(ctx, &ppool.GetNodeRequest{Name: bs.NodeName})
		if err != nil {
			return c.String(http.StatusInternalServerError, errors.Wrap(err, "Failed to get node").Error())
		}

		u := &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", node.Address, agentPort),
			Path:   DownloadBlockStorageHTTPPrefix,
		}

		log.Printf("[DEBUG] ProxyDownloadBlockStorage: url=%s", u.String())
		proxy := http.StripPrefix(basePath, httputil.NewSingleHostReverseProxy(u))
		proxy.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
