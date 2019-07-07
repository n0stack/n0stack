package iuser

import (
	"context"

	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	auser "github.com/n0stack/n0stack/n0core/pkg/api/iam/user"
	stdapi "github.com/n0stack/n0stack/n0core/pkg/api/standard_api"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	piam "github.com/n0stack/n0stack/n0proto.go/iam/v0"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
)

type UserAPI struct {
	dataStore datastore.Datastore
}

func (a *UserAPI) ListUsers(ctx context.Context, req *piam.ListUsersRequest) (*piam.ListUsersResponse, error) {
	return auser.ListUsers(ctx, req, a.dataStore)
}

func (a *UserAPI) GetUser(ctx context.Context, req *piam.GetUserRequest) (*piam.User, error) {
	return auser.GetUser(ctx, req, a.dataStore)
}

func (a *UserAPI) CreateUser(ctx context.Context, req *piam.CreateUserRequest) (*piam.User, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	if err := auser.PendNewUser(tx, a.dataStore, req.Name); err != nil {
		return nil, err
	}

	user := &piam.User{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		State: piam.User_AVAILABLE,
	}

	if err := auser.ApplyUser(a.dataStore, user); err != nil {
		return nil, err
	}

	tx.Commit()
	return user, nil
}

func (a *UserAPI) DeleteUser(ctx context.Context, req *piam.DeleteUserRequest) (*empty.Empty, error) {
	if !a.dataStore.Lock(req.Name) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	user, err := auser.GetAndPendExistingUser(tx, a.dataStore, req.Name)
	if err != nil {
		return nil, err
	}

	if err := auser.DeleteUser(a.dataStore, user.Name); err != nil {
		return nil, err
	}

	tx.Commit()
	return &empty.Empty{}, nil
}

func (a *UserAPI) AddSshPublicKey(ctx context.Context, req *piam.AddSshPublicKeyRequest) (*piam.User, error) {
	if !a.dataStore.Lock(req.UserName) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.UserName)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	user, err := auser.GetAndPendExistingUser(tx, a.dataStore, req.UserName)
	if err != nil {
		return nil, err
	}

	if user.SshPublicKeys == nil {
		user.SshPublicKeys = make(map[string]string)
	}
	user.SshPublicKeys[req.SshPublicKeyName] = req.SshPublicKey

	if err := auser.ApplyUser(a.dataStore, user); err != nil {
		return nil, err
	}

	tx.Commit()
	return user, nil
}

func (a *UserAPI) DeleteSshPublicKey(ctx context.Context, req *piam.DeleteSshPublicKeyRequest) (*piam.User, error) {
	if !a.dataStore.Lock(req.UserName) {
		return nil, stdapi.LockError()
	}
	defer a.dataStore.Unlock(req.UserName)

	tx := transaction.Begin()
	defer tx.RollbackWithLog()

	user, err := auser.GetAndPendExistingUser(tx, a.dataStore, req.UserName)
	if err != nil {
		return nil, err
	}

	if user.SshPublicKeys == nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "publicKey '%s' does not exist", req.SshPublicKeyName)
	}
	if _, ok := user.SshPublicKeys[req.SshPublicKeyName]; !ok {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "publicKey '%s' does not exist", req.SshPublicKeyName)
	}

	delete(user.SshPublicKeys, req.SshPublicKeyName)

	if err := auser.ApplyUser(a.dataStore, user); err != nil {
		return nil, err
	}

	tx.Commit()
	return user, nil
}
