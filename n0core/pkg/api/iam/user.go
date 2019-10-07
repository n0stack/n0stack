package iam

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	piam "n0st.ac/n0stack/iam/v1alpha"
	stdapi "n0st.ac/n0stack/n0core/pkg/api/stdapi"
	"n0st.ac/n0stack/n0core/pkg/datastore"
	"n0st.ac/n0stack/n0core/pkg/driver/n0stack/auth"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	structutil "n0st.ac/n0stack/n0core/pkg/util/struct"
)

func GetUser(ctx context.Context, ds datastore.Datastore, name string) (*piam.User, int64, error) {
	resourse := &piam.User{}
	version, err := ds.Get(ctx, name, resourse)
	if err != nil {
		if datastore.IsNotFound(err) {
			return nil, 0, grpcutil.Errorf(codes.NotFound, err.Error())
		}

		return nil, 0, grpcutil.Errorf(codes.Internal, "failed to get User %!s(MISSING) from db: err='%!s(MISSING)'", name, err.Error())
	}

	return resourse, version, nil
}

type UserAPI struct {
	datastore datastore.Datastore

	auth *auth.AuthenticationServiceProvier
}

func CreateUserAPI(datastore datastore.Datastore, auth *auth.AuthenticationServiceProvier) *UserAPI {
	return &UserAPI{
		datastore: datastore.AddPrefix("iam/user"),
		auth:      auth,
	}
}

func (a *UserAPI) GetUser(ctx context.Context, req *piam.GetUserRequest) (*piam.User, error) {
	u, _, err := GetUser(ctx, a.datastore, req.Name)
	return u, err
}

func (a *UserAPI) CreateUser(ctx context.Context, req *piam.CreateUserRequest) (*piam.User, error) {
	if req.User == nil {
		return nil, grpcutil.Errorf(codes.InvalidArgument, "set user")
	}

	if err := stdapi.ValidateName(req.User.Name); err != nil {
		return nil, err
	}

	if len(req.User.PublicKeys) < 1 {
		return nil, grpcutil.Errorf(codes.InvalidArgument, "set public key")
	}
	for k, v := range req.User.PublicKeys {
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(v))
		if err != nil {
			return nil, grpcutil.Errorf(codes.InvalidArgument, "public key %s is invalid", k)
		}
	}

	if _, _, err := GetUser(ctx, a.datastore, req.User.Name); err != nil {
		if grpc.Code(err) != codes.NotFound {
			return nil, err
		}
	}

	user := &piam.User{
		Name:        req.User.Name,
		Annotations: req.User.Annotations,
		Labels:      req.User.Labels,
		DisplayName: req.User.DisplayName,
		PublicKeys:  req.User.PublicKeys,
	}

	if _, err := a.datastore.Apply(ctx, user.Name, user, 0); err != nil {
		return nil, stdapi.DatastoreApplyError(err, "User", user.Name)
	}

	return user, nil
}

func (a *UserAPI) UpdateUser(ctx context.Context, req *piam.UpdateUserRequest) (*piam.User, error) {
	if err := stdapi.CheckAuthenticatedUserName(ctx, a.auth, req.User.Name); err != nil {
		return nil, err
	}

	user, version, err := GetUser(ctx, a.datastore, req.User.Name)
	if err != nil {
		return nil, err
	}

	if err := structutil.UpdateWithMaskUsingJson(user, req.User, req.UpdateMask.Paths); err != nil {
		return nil, stdapi.UpdateMaskError(err)
	}

	if _, err := a.datastore.Apply(ctx, user.Name, user, version); err != nil {
		return nil, stdapi.DatastoreApplyError(err, "User", user.Name)
	}

	return user, nil
}

func (a *UserAPI) DeleteUser(ctx context.Context, req *piam.DeleteUserRequest) (*empty.Empty, error) {
	if err := stdapi.CheckAuthenticatedUserName(ctx, a.auth, req.Name); err != nil {
		return nil, err
	}

	user, version, err := GetUser(ctx, a.datastore, req.Name)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return &empty.Empty{}, nil
		}

		return nil, err
	}

	if err := a.datastore.Delete(ctx, user.Name, version); err != nil {
		return nil, stdapi.DatastoreDeleteError(err, "User", user.Name)
	}

	return &empty.Empty{}, nil
}
