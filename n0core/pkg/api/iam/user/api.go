package userapi

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	piam "github.com/n0stack/n0stack/n0proto.go/iam/v0"
)

type UserAPI struct {
	dataStore datastore.Datastore
}

func (a *UserAPI) ListUsers(ctx context.Context, req *piam.ListUsersRequest) (*piam.ListUsersResponse, error) {
	return ListUsers(ctx, req, a.dataStore)
}

func (a *UserAPI) GetUser(ctx context.Context, req *piam.GetUserRequest) (*piam.User, error) {
	return GetUser(ctx, req, a.dataStore)
}
func (a *UserAPI) CreateUser(context.Context, *piam.CreateUserRequest) (*piam.User, error) {
	return nil, nil
}
func (a *UserAPI) DeleteUser(context.Context, *piam.DeleteUserRequest) (*empty.Empty, error) {
	return nil, nil
}
func (a *UserAPI) AddSshPublicKey(context.Context, *piam.AddSshPublicKeyRequest) (*piam.User, error) {
	return nil, nil
}
func (a *UserAPI) DeleteSshPublicKey(context.Context, *piam.DeleteSshPublicKeyRequest) (*piam.User, error) {
	return nil, nil
}
