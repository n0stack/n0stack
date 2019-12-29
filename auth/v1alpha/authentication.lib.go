package pauth

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
)

type AuthenticationClientByServer struct {
	API AuthenticationServiceServer
}

func (a AuthenticationClientByServer) GetAuthenticationTokenPublicKey(ctx context.Context, in *GetAuthenticationTokenPublicKeyRequest, opts ...grpc.CallOption) (*GetAuthenticationTokenPublicKeyResponse, error) {
	return a.API.GetAuthenticationTokenPublicKey(ctx, in)
}
func (a AuthenticationClientByServer) PublicKeyAuthenticate(ctx context.Context, opts ...grpc.CallOption) (AuthenticationService_PublicKeyAuthenticateClient, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}
