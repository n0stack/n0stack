package piam

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type MockedUserAPI struct {
	API UserServiceServer
}

func (m MockedUserAPI) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error) {
	return m.API.GetUser(ctx, in)
}
func (m MockedUserAPI) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error) {
	return m.API.CreateUser(ctx, in)
}
func (m MockedUserAPI) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error) {
	return m.API.UpdateUser(ctx, in)
}
func (m MockedUserAPI) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return m.API.DeleteUser(ctx, in)
}
