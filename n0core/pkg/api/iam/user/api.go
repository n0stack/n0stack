package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	piam "n0st.ac/n0stack/iam/v1alpha"
	stdapi "n0st.ac/n0stack/n0core/pkg/api/standard_api"
	"n0st.ac/n0stack/n0core/pkg/datastore"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
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

	if _, err := ApplyUser(ctx, a.datastore, user, 0); err != nil {
		return nil, err
	}

	return user, nil
}

func (a *UserAPI) UpdateUser(ctx context.Context, req *piam.UpdateUserRequest) (*piam.User, error) {
	return nil, grpcutil.Errorf(codes.Unimplemented, "")
	// user, version, err := GetUser(ctx, a.datastore, req.User.Name)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, mask := range req.UpdateMask.Paths {
	// 	path := strings.Split(mask, ".")

	// 	for

	// 	switch path[0] {
	// 	case strings.HasPrefix(mask, "name"):
	// 		user.Name = req.User.Name

	// 	case strings.HasPrefix(mask, "annotations"):
	// 		switch {
	// 			case mask
	// 		}
	// 	// case strings.HasPrefix(mask, "labels"):

	// 	case strings.HasPrefix(mask, "display_name"):
	// 		user.DisplayName = req.User.DisplayName

	// 		// case strings.HasPrefix(mask, "public_keys"):
	// 	}
	// }

	// if _, err := ApplyUser(ctx, a.datastore, user, version); err != nil {
	// 	return nil, err
	// }

	// return user, nil
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
