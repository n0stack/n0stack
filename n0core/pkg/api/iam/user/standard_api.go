package auser

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	piam "github.com/n0stack/n0stack/n0proto.go/iam/v1alpha"
	"google.golang.org/grpc/codes"
)

func ListUsers(ctx context.Context, req *piam.ListUsersRequest, ds datastore.Datastore) (*piam.ListUsersResponse, error) {
	res := &piam.ListUsersResponse{}
	f := func(s int) []proto.Message {
		res.Users = make([]*piam.User, s)
		for i := range res.Users {
			res.Users[i] = &piam.User{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Users {
			m[i] = v
		}

		return m
	}

	if err := ds.List(ctx, f); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.Users) == 0 {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	return res, nil
}

func GetUser(ctx context.Context, ds datastore.Datastore, name string) (*piam.User, int64, error) {
	resourse := &piam.User{}
	v, err := ds.Get(ctx, name, resourse)
	if err != nil {
		if datastore.IsNotFound(err) {
			return nil, 0, grpcutil.WrapGrpcErrorf(codes.NotFound, err.Error())
		}

		return nil, 0, grpcutil.WrapGrpcErrorf(codes.Internal, datastore.DefaultErrorMessage(err))
	}

	return resourse, v, nil
}

func DeleteUser(ctx context.Context, ds datastore.Datastore, name string, version int64) error {
	if err := ds.Delete(ctx, name, version); err != nil {
		return grpcutil.WrapGrpcErrorf(codes.Internal, "failed to delete User %s from db: err='%s'", name, err.Error())
	}

	return nil
}

func ApplyUser(ctx context.Context, ds datastore.Datastore, resource *piam.User, version int64) (int64, error) {
	v, err := ds.Apply(ctx, resource.Name, resource, version)
	if err != nil {
		return 0, grpcutil.WrapGrpcErrorf(codes.Internal, "failed to apply User %s to db: err='%s'", resource.Name, err.Error())
	}

	return v, nil
}
