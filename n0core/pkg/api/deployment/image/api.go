package image

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0proto.go/deployment/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type ImageAPI struct {
	dataStore       datastore.Datastore
	blockstorageAPI pprovisioning.BlockStorageServiceClient
}

func CreateImageAPI(ds datastore.Datastore, bsa pprovisioning.BlockStorageServiceClient) (*ImageAPI, error) {
	a := &ImageAPI{
		dataStore:       ds,
		blockstorageAPI: bsa,
	}
	a.dataStore.AddPrefix("image")

	return a, nil
}

func (a ImageAPI) ListImages(ctx context.Context, req *pdeployment.ListImagesRequest) (*pdeployment.ListImagesResponse, error) {
	res := &pdeployment.ListImagesResponse{}
	f := func(s int) []proto.Message {
		res.Images = make([]*pdeployment.Image, s)
		for i := range res.Images {
			res.Images[i] = &pdeployment.Image{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Images {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		log.Printf("[WARNING] Failed to list data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.Images) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a ImageAPI) GetImage(ctx context.Context, req *pdeployment.GetImageRequest) (*pdeployment.Image, error) {
	res := &pdeployment.Image{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a ImageAPI) ApplyImage(ctx context.Context, req *pdeployment.ApplyImageRequest) (*pdeployment.Image, error) {

	res := &pdeployment.Image{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	// var err error
	res.Version, _ = datastore.CheckVersion(res.Version, req.Version)
	// if err != nil {
	// 	return nil, grpc.Errorf(codes.InvalidArgument, "Failed to check version: %s", err.Error())
	// }

	res.Name = req.Name
	res.Annotations = req.Annotations

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a ImageAPI) DeleteImage(ctx context.Context, req *pdeployment.DeleteImageRequest) (*empty.Empty, error) {
	prev := &pdeployment.Image{}
	if err := a.dataStore.Get(req.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if prev.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Image '%s' is not found", req.Name)
	}

	for _, bs := range prev.RegisteredBlockStorages {
		_, err := a.blockstorageAPI.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: bs.BlockStorageName})
		if err != nil {
			log.Printf("[WARNING] Failed to set blockstorage as available: err='%s'", err.Error())
			return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
		}
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		log.Printf("[WARNING] Failed to delete data from db: err='%s'", err.Error())
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "Failed to delete '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return &empty.Empty{}, nil
}

func (a ImageAPI) RegisterBlockStorage(ctx context.Context, req *pdeployment.RegisterBlockStorageRequest) (*pdeployment.Image, error) {
	res := &pdeployment.Image{}
	if err := a.dataStore.Get(req.ImageName, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.ImageName)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Image '%s' is not found", req.ImageName)
	}
	if res.Tags == nil {
		res.Tags = make(map[string]string)
	}

	bs, err := a.blockstorageAPI.SetProtectedBlockStorage(context.Background(), &pprovisioning.SetProtectedBlockStorageRequest{Name: req.BlockStorageName})
	if err != nil {
		log.Printf("[WARNING] Failed to set blockstorage as protected: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.ImageName)
	}

	res.RegisteredBlockStorages = append(res.RegisteredBlockStorages, &pdeployment.Image_RegisteredBlockStorage{
		BlockStorageName: bs.Name,
		RegisteredAt:     ptypes.TimestampNow(),
	})
	for _, t := range req.Tags {
		res.Tags[t] = bs.Name
	}

	if err := a.dataStore.Apply(req.ImageName, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.ImageName)
	}

	return res, nil
}

func (a ImageAPI) UnregisterBlockStorage(ctx context.Context, req *pdeployment.UnregisterBlockStorageRequest) (*pdeployment.Image, error) {
	res := &pdeployment.Image{}
	if err := a.dataStore.Get(req.ImageName, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.ImageName)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Image '%s' is not found", req.ImageName)
	}
	if res.Tags == nil {
		res.Tags = make(map[string]string)
	}

	_, err := a.blockstorageAPI.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: req.BlockStorageName})
	if err != nil {
		log.Printf("[WARNING] Failed to set blockstorage as protected: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.ImageName)
	}

	for i, bs := range res.RegisteredBlockStorages {
		if bs.BlockStorageName == req.BlockStorageName {
			res.RegisteredBlockStorages = append(res.RegisteredBlockStorages[:i], res.RegisteredBlockStorages[i+1:]...)

			break
		}
	}

	for k, v := range res.Tags {
		if v == req.BlockStorageName {
			delete(res.Tags, k)
		}
	}

	if err := a.dataStore.Apply(req.ImageName, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.ImageName)
	}

	return res, nil
}

func (a ImageAPI) GenerateBlockStorage(ctx context.Context, req *pdeployment.GenerateBlockStorageRequest) (*pprovisioning.BlockStorage, error) {
	prev := &pdeployment.Image{}
	if err := a.dataStore.Get(req.ImageName, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.ImageName)
	}
	if prev.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Image '%s' is not found", req.ImageName)
	}
	if prev.Tags == nil {
		return nil, grpc.Errorf(codes.NotFound, "Tag '%s' is not found", req.Tag)
	}
	if _, ok := prev.Tags[req.Tag]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Tag '%s' is not found", req.Tag)
	}

	bs, err := a.blockstorageAPI.CopyBlockStorage(context.Background(), &pprovisioning.CopyBlockStorageRequest{
		SourceBlockStorage: prev.Tags[req.Tag],
		Name:               req.BlockStorageName,
		Annotations:        req.Annotations,
		RequestBytes:       req.RequestBytes,
		LimitBytes:         req.LimitBytes,
	})
	if err != nil {
		log.Printf("[WARNING] Failed to copy blockstorage: err='%s'", err.Error())
		// already exists を判定したほうがいいような気がするが、とりあえず err をそのまま返すことで対処する
		return nil, err
	}

	return bs, nil
}

func (a ImageAPI) TagImage(ctx context.Context, req *pdeployment.TagImageRequest) (*pdeployment.Image, error) {
	res := &pdeployment.Image{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Image '%s' is not found", req.Name)
	}
	if res.Tags == nil {
		res.Tags = make(map[string]string)
	}

	exists := false
	for _, bs := range res.RegisteredBlockStorages {
		if bs.BlockStorageName == req.BlockStorageName {
			exists = true
			break
		}
	}
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "BlockStorage '%s' is not in RegisteredBlockStorages", req.BlockStorageName)
	}

	for _, t := range req.Tags {
		res.Tags[t] = req.BlockStorageName
	}

	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a ImageAPI) UntagImage(ctx context.Context, req *pdeployment.UntagImageRequest) (*pdeployment.Image, error) {
	res := &pdeployment.Image{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Image '%s' is not found", req.Name)
	}
	if res.Tags == nil {
		return nil, grpc.Errorf(codes.NotFound, "Tag '%s' is not found", req.Tag)
	}

	if _, ok := res.Tags[req.Tag]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Tag '%s' is not found", req.Tag)
	}

	delete(res.Tags, req.Tag)
	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}
