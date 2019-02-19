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

	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/labstack/echo"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"
	"github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
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

func (a *BlockStorageAPI) ReserveStorage(ctx context.Context, tx *transaction.Transaction, bs *pprovisioning.BlockStorage) error {
	bs.StorageName = bs.Name

	var n *ppool.Node
	var err error
	if node, ok := bs.Annotations[AnnotationBlockStorageRequestNodeName]; !ok {
		n, err = a.nodeAPI.ScheduleStorage(ctx, &ppool.ScheduleStorageRequest{
			StorageName: bs.StorageName,
			Annotations: map[string]string{
				AnnotationStorageReservedBy: bs.Name,
			},
			RequestBytes: bs.RequestBytes,
			LimitBytes:   bs.LimitBytes,
		})
	} else {
		n, err = a.nodeAPI.ReserveStorage(ctx, &ppool.ReserveStorageRequest{
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
		_, err := a.nodeAPI.ReleaseStorage(ctx, &ppool.ReleaseStorageRequest{
			NodeName:    bs.NodeName,
			StorageName: bs.StorageName,
		})

		return err
	})

	return nil
}

func (a *BlockStorageAPI) CheckAndLock(tx *transaction.Transaction, bs *pprovisioning.BlockStorage) error {
	prev := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(bs.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", bs.Name)
	} else if prev.Name != "" {
		return grpcutil.WrapGrpcErrorf(codes.AlreadyExists, "BlockStorage '%s' is already exists", bs.Name)
	}

	bs.State = pprovisioning.BlockStorage_PENDING
	if err := a.dataStore.Apply(bs.Name, bs); err != nil {
		return grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}
	tx.PushRollback("free optimistic lock", func() error {
		return a.dataStore.Delete(bs.Name)
	})

	return nil
}

func (a *BlockStorageAPI) GetAndLock(tx *transaction.Transaction, name string) (*pprovisioning.BlockStorage, error) {
	bs := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(name, bs); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", name)
	} else if bs.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	if bs.State == pprovisioning.BlockStorage_PENDING {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "BlockStorage '%s' is pending", name)
	}

	current := bs.State
	bs.State = pprovisioning.BlockStorage_PENDING
	if err := a.dataStore.Apply(bs.Name, bs); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}
	bs.State = current
	tx.PushRollback("free optimistic lock", func() error {
		return a.dataStore.Apply(bs.Name, bs)
	})

	return bs, nil
}

func (a *BlockStorageAPI) CreateBlockStorage(ctx context.Context, req *pprovisioning.CreateBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	{
		if req.RequestBytes == 0 {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set 'request_bytes'")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs := &pprovisioning.BlockStorage{
		Name:         req.Name,
		Annotations:  req.Annotations,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}

	if err := a.CheckAndLock(tx, bs); err != nil {
		return nil, err
	}

	if err := a.ReserveStorage(ctx, tx, bs); err != nil {
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

	if err := a.dataStore.Apply(req.Name, bs); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
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
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs := &pprovisioning.BlockStorage{
		Name:         req.Name,
		Annotations:  req.Annotations,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}
	if bs.Annotations == nil {
		bs.Annotations = make(map[string]string)
	}
	bs.Annotations[AnnotationBlockStorageFetchFrom] = req.SourceUrl

	if err := a.CheckAndLock(tx, bs); err != nil {
		return nil, err
	}

	if err := a.ReserveStorage(ctx, tx, bs); err != nil {
		return nil, err
	}

	{
		cli, done, err := a.getAgent(ctx, bs.NodeName)
		if err != nil {
			return nil, err
		}
		defer done()

		v, err := cli.FetchBlockStorage(context.Background(), &FetchBlockStorageRequest{
			Name:      bs.Name,
			Bytes:     bs.LimitBytes,
			SourceUrl: bs.Annotations[AnnotationBlockStorageFetchFrom],
		})
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to FetchBlockStorage on node '%s': err='%s'", bs.NodeName, err.Error())
		}
		tx.PushRollback("DeleteBlockStorage", func() error {
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
	}

	if err := a.dataStore.Apply(req.Name, bs); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}

	tx.Commit()
	return bs, nil
}

func (a *BlockStorageAPI) CopyBlockStorage(ctx context.Context, req *pprovisioning.CopyBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	{
		if req.RequestBytes == 0 {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set 'request_bytes'")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs := &pprovisioning.BlockStorage{
		Name:         req.Name,
		Annotations:  req.Annotations,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	}
	if bs.Annotations == nil {
		bs.Annotations = make(map[string]string)
	}
	bs.Annotations[AnnotationBlockStorageCopyFrom] = req.SourceBlockStorage

	if err := a.CheckAndLock(tx, bs); err != nil {
		return nil, err
	}

	var srcUrl *url.URL
	{
		src := &pprovisioning.BlockStorage{}
		if err := a.dataStore.Get(req.SourceBlockStorage, src); err != nil {
			log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", src.Name)
		} else if src.Name == "" {
			return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "BlockStorage '%s' is not exists", src.Name)
		}

		srcNode, err := a.nodeAPI.GetNode(ctx, &ppool.GetNodeRequest{Name: src.NodeName})
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to get node '%s' which is hosting source block storage", src.NodeName).Error())
		}

		srcUrl = &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", srcNode.Address, 8081), // TODO: agentのポートをハードコードしている
			Path:   filepath.Join(DownloadBlockStorageHTTPPrefix, src.Name),
		}
	}

	if err := a.ReserveStorage(ctx, tx, bs); err != nil {
		return nil, err
	}

	{
		cli, done, err := a.getAgent(ctx, bs.NodeName)
		if err != nil {
			return nil, err
		}
		defer done()

		v, err := cli.FetchBlockStorage(context.Background(), &FetchBlockStorageRequest{
			Name:      req.Name,
			Bytes:     req.LimitBytes,
			SourceUrl: srcUrl.String(),
		})
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to FetchBlockStorage on node '%s': err='%s'", bs.NodeName, err.Error())
		}
		tx.PushRollback("DeleteBlockStorage", func() error {
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
	}

	if err := a.dataStore.Apply(req.Name, bs); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}

	tx.Commit()
	return bs, nil
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
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.BlockStorages) == 0 {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *BlockStorageAPI) GetBlockStorage(ctx context.Context, req *pprovisioning.GetBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	return res, nil
}

func (a *BlockStorageAPI) UpdateBlockStorage(ctx context.Context, req *pprovisioning.UpdateBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	return nil, grpcutil.WrapGrpcErrorf(codes.Unimplemented, "")
}

func (a *BlockStorageAPI) DeleteBlockStorage(ctx context.Context, req *pprovisioning.DeleteBlockStorageRequest) (*empty.Empty, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	bs, err := a.GetAndLock(tx, req.Name)
	if err != nil {
		return nil, err
	}

	if bs.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "BlockStorage '%s' is not available: now=%d", req.Name, bs.State)
	}

	cli, done, err := a.getAgent(ctx, bs.NodeName)
	if err != nil {
		return nil, err
	}
	defer done()

	u, _ := url.Parse(bs.Annotations[AnnotationBlockStorageURL]) // TODO: エラー処理
	_, err = cli.DeleteBlockStorage(context.Background(), &DeleteBlockStorageRequest{Path: u.Path})
	if err != nil {
		log.Printf("Fail to delete block_storage on node: err=%s, req=%s", err.Error(), &DeleteBlockStorageRequest{Path: u.Path})
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Fail to delete block_storage on node") // TODO #89
	}

	_, err = a.nodeAPI.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		NodeName:    bs.NodeName,
		StorageName: bs.StorageName,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", bs.StorageName, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if grpc.Code(err) != codes.NotFound {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to release compute '%s': please retry", bs.StorageName)
		}
	}

	if err := a.dataStore.Delete(bs.Name); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}

	tx.Commit()
	return &empty.Empty{}, nil
}

func (a *BlockStorageAPI) SetAvailableBlockStorage(ctx context.Context, req *pprovisioning.SetAvailableBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !lock.WaitUntilLock(a.dataStore, req.Name, 1*time.Second, 50*time.Millisecond) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if res.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	if res.State == pprovisioning.BlockStorage_UNKNOWN {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Can change state to AVAILABLE when UNKNOWN")
	}
	res.State = pprovisioning.BlockStorage_AVAILABLE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *BlockStorageAPI) SetInuseBlockStorage(ctx context.Context, req *pprovisioning.SetInuseBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !lock.WaitUntilLock(a.dataStore, req.Name, 1*time.Second, 50*time.Millisecond) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if res.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	if res.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Can change state to IN_USE when only AVAILABLE")
	}
	res.State = pprovisioning.BlockStorage_IN_USE

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *BlockStorageAPI) SetProtectedBlockStorage(ctx context.Context, req *pprovisioning.SetProtectedBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	if !lock.WaitUntilLock(a.dataStore, req.Name, 1*time.Second, 50*time.Millisecond) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	res := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if res.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	if res.State != pprovisioning.BlockStorage_AVAILABLE {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Can change state to PROTECTED when only AVAILABLE")
	}
	res.State = pprovisioning.BlockStorage_PROTECTED

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a *BlockStorageAPI) DownloadBlockStorage(ctx context.Context, req *pprovisioning.DownloadBlockStorageRequest) (*pprovisioning.DownloadBlockStorageResponse, error) {
	bs := &pprovisioning.BlockStorage{}
	if err := a.dataStore.Get(req.Name, bs); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	} else if bs.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	u := &url.URL{
		Scheme: "http",
		Path:   fmt.Sprintf("/api/v0/block_storage/download/%s", req.Name),
	}
	res := &pprovisioning.DownloadBlockStorageResponse{
		DownloadUrl: u.String(),
	}

	return res, nil
}

func NewReverseProxyStrippingAllPath(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

// TODO: agentPort は Node から取れるようにしたい
func (a *BlockStorageAPI) ProxyDownloadBlockStorage(agentPort int) func(echo.Context) error {
	return func(c echo.Context) error {
		name := c.Param("name")

		bs := &pprovisioning.BlockStorage{}
		if err := a.dataStore.Get(name, bs); err != nil {
			log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
			return fmt.Errorf("db error")
		} else if bs.Name == "" {
			return err
		}

		ctx := context.Background()
		node, err := a.nodeAPI.GetNode(ctx, &ppool.GetNodeRequest{Name: bs.NodeName})
		if err != nil {
			return err
		}

		u := &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", node.Address, agentPort),
			Path:   filepath.Join(DownloadBlockStorageHTTPPrefix, name),
		}

		log.Printf("[DEBUG] ProxyDownloadBlockStorage: url=%s", u.String())
		proxy := NewReverseProxyStrippingAllPath(u)
		proxy.ServeHTTP(c.Response(), c.Request())

		return nil
	}
}
