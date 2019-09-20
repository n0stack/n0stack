package user

import (
	"context"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/datastore"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	piam "n0st.ac/n0stack/n0proto.go/iam/v1alpha"
)

type UserAPI struct {
	datastore datastore.Datastore
}

func CreateUserAPI(datastore datastore.Datastore) *UserAPI {
	return &UserAPI{
		datastore: datastore.AddPrefix("iam/user"),
	}
}

func (a *UserAPI) GetUser(ctx context.Context, req *piam.GetUserRequest) (*piam.User, error) {
	u, _, err := GetUser(ctx, a.datastore, req.Name)
	return u, err
}

func (a *UserAPI) CreateUser(ctx context.Context, req *piam.CreateUserRequest) (*piam.User, error) {
	if err := stdapi.ValidateName(req.Name); err != nil {
		return nil, err
	}

	if _, _, err := GetUser(ctx, a.datastore, req.Name); err != nil {
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

	if _, err := ApplyUser(ctx, a.datastore, user, 0); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UserAPI) UpdateUser(ctx context.Context, req *piam.UpdateUserRequest) (*piam.User, error) {
	_, version, err := GetUser(ctx, a.datastore, req.Name)
	if err != nil {
		return nil, err
	}

	user := &piam.User{
		Name:        req.Name,
		Annotations: req.Annotations,
		Labels:      req.Labels,

		DisplayName: req.DisplayName,
	}

	if _, err := ApplyUser(ctx, a.datastore, user, version); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UserAPI) DeleteUser(ctx context.Context, req *piam.DeleteUserRequest) (*empty.Empty, error) {
	user, version, err := GetUser(ctx, a.datastore, req.Name)
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return &empty.Empty{}, nil
		}

		return nil, err
	}

	if err := DeleteUser(ctx, a.datastore, user.Name, version); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (a *UserAPI) AddPublicKey(ctx context.Context, req *piam.AddPublicKeyRequest) (*piam.User, error) {
	user, version, err := GetUser(ctx, a.datastore, req.UserName)
	if err != nil {
		return nil, err
	}

	if user.PublicKeys == nil {
		user.PublicKeys = make(map[string]string)
	}
	user.PublicKeys[req.PublicKeyName] = req.PublicKey

	if _, err := ApplyUser(ctx, a.datastore, user, version); err != nil {
		return nil, err
	}

	return user, nil
}
func (a *UserAPI) DeletePublicKey(ctx context.Context, req *piam.DeletePublicKeyRequest) (*piam.User, error) {
	user, version, err := GetUser(ctx, a.datastore, req.UserName)
	if err != nil {
		return nil, err
	}

	if user.PublicKeys == nil {
		return nil, grpcutil.Errorf(codes.NotFound, "publicKey '%s' does not exist", req.PublicKeyName)
	}
	if _, ok := user.PublicKeys[req.PublicKeyName]; !ok {
		return nil, grpcutil.Errorf(codes.NotFound, "publicKey '%s' does not exist", req.PublicKeyName)
	}

	delete(user.PublicKeys, req.PublicKeyName)

	if _, err := ApplyUser(ctx, a.datastore, user, version); err != nil {
		return nil, err
	}

	return user, nil
}
