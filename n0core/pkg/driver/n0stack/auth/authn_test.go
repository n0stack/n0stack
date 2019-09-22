package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"testing"

	"google.golang.org/grpc"
	pauth "n0st.ac/n0stack/auth/v1alpha"
	piam "n0st.ac/n0stack/iam/v1alpha"
	authn "n0st.ac/n0stack/n0core/pkg/api/auth/authentication"
	"n0st.ac/n0stack/n0core/pkg/api/iam/user"
	"n0st.ac/n0stack/n0core/pkg/datastore/memory"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	jwtutil "n0st.ac/n0stack/n0core/pkg/util/jwt"
)

func TestGetAuthnToken(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKey, _ := jwtutil.NewPrivateKey(key)

	datastore := memory.NewMemoryDatastore()
	userapi := &piam.MockedUserAPI{
		API: user.CreateUserAPI(datastore),
	}
	_, err := userapi.CreateUser(context.Background(), &piam.CreateUserRequest{
		User: &piam.User{
			Name: "test-user",
			PublicKeys: map[string]string{
				"test_key": string(privateKey.PublicKey().MarshalAuthorizedKey()),
			},
		},
	})
	if err != nil {
		t.Errorf("CreateUser returns err=%+v", err)
	}

	conn, close := grpcutil.PrepareMockedGRPC(t, func(grpcServer *grpc.Server) {
		secret := make([]byte, 256)
		io.ReadFull(rand.Reader, secret)
		authnapi := authn.CreateAuthenticationAPI(userapi, secret)
		pauth.RegisterAuthenticationServiceServer(grpcServer, authnapi)
	})
	defer close()

	a := NewAuthenticationClient(conn, privateKey)
	_, err = a.GetAuthenticationToken(context.Background(), "test-user")
	if err != nil {
		t.Errorf("GetAuthenticationToken() returns err=%+v", err)
	}
}
