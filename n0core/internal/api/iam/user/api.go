package iuser

import (
	"context"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	auser "github.com/n0stack/n0stack/n0core/pkg/api/iam/user"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	piam "github.com/n0stack/n0stack/n0proto.go/iam/v1alpha"
)

type UserAPI struct {
	dataStore datastore.Datastore
}

func CreateUserAPI(datastore datastore.Datastore) *UserAPI {
	return &UserAPI{
		dataStore: datastore.AddPrefix("iam/user"),
	}
}

func (a *UserAPI) GetUser(ctx context.Context, req *piam.GetUserRequest) (*piam.User, error) {
	u, _, err := auser.GetUser(ctx, a.dataStore, req.Name)
	return u, err
}

func (a *UserAPI) CreateUser(ctx context.Context, req *piam.CreateUserRequest) (*piam.User, error) {
	if _, _, err := auser.GetUser(ctx, a.dataStore, req.Name); err != nil {
		if grpc.Code(err) != codes.NotFound {
			return nil, err
		}
	}

	user := &piam.User{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		DisplayName: req.DisplayName,
	}

	if _, err := auser.ApplyUser(ctx, a.dataStore, user, 0); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UserAPI) UpdateUser(ctx context.Context, req *piam.UpdateUserRequest) (*piam.User, error) {
	_, version, err := auser.GetUser(ctx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	user := &piam.User{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		DisplayName: req.DisplayName,
	}

	if _, err := auser.ApplyUser(ctx, a.dataStore, user, version); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UserAPI) DeleteUser(ctx context.Context, req *piam.DeleteUserRequest) (*empty.Empty, error) {
	user, version, err := auser.GetUser(ctx, a.dataStore, req.Name)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return &empty.Empty{}, nil
		}

		return nil, err
	}

	if err := auser.DeleteUser(ctx, a.dataStore, user.Name, version); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (a *UserAPI) AddPublicKey(ctx context.Context, req *piam.AddPublicKeyRequest) (*piam.User, error) {
	user, version, err := auser.GetUser(ctx, a.dataStore, req.UserName)
	if err != nil {
		return nil, err
	}

	if user.PublicKeys == nil {
		user.PublicKeys = make(map[string]string)
	}
	user.PublicKeys[req.PublicKeyName] = req.PublicKey

	if _, err := auser.ApplyUser(ctx, a.dataStore, user, version); err != nil {
		return nil, err
	}

	return user, nil
}
func (a *UserAPI) DeletePublicKey(ctx context.Context, req *piam.DeletePublicKeyRequest) (*piam.User, error) {
	user, version, err := auser.GetUser(ctx, a.dataStore, req.UserName)
	if err != nil {
		return nil, err
	}

	if user.PublicKeys == nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "publicKey '%s' does not exist", req.PublicKeyName)
	}
	if _, ok := user.PublicKeys[req.PublicKeyName]; !ok {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "publicKey '%s' does not exist", req.PublicKeyName)
	}

	delete(user.PublicKeys, req.PublicKeyName)

	if _, err := auser.ApplyUser(ctx, a.dataStore, user, version); err != nil {
		return nil, err
	}

	return user, nil
}
