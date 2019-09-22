package auth

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	pauth "n0st.ac/n0stack/auth/v1alpha"
	jwtutil "n0st.ac/n0stack/n0core/pkg/util/jwt"
)

type AuthenticationServiceProvier struct {
	authnAPI pauth.AuthenticationServiceClient
	host     string

	publicKey *jwtutil.PublicKey
}

// host: URL is "http://example.com:8080/hoge", then host is "example.com:8080"
func NewAuthenticationServiceProvider(ctx context.Context, conn *grpc.ClientConn, host string) (*AuthenticationServiceProvier, error) {
	authnAPI := pauth.NewAuthenticationServiceClient(conn)
	sp := &AuthenticationServiceProvier{
		authnAPI: authnAPI,
		host:     host,
	}

	err := sp.renewPublicKey(ctx)
	if err != nil {
		log.Printf("[CRITICAL] failed to renew a public key, err=%+v", err)
		return nil, errors.Wrap(err, "renewPublicKey() is failed")
	}

	go sp.loop(ctx)

	return sp, nil
}

func (sp *AuthenticationServiceProvier) renewPublicKey(ctx context.Context) error {
	res, err := sp.authnAPI.GetAuthenticationTokenPublicKey(ctx, &pauth.GetAuthenticationTokenPublicKeyRequest{})
	if err != nil {
		return err
	}

	key, err := jwtutil.ParsePublicKey([]byte(res.PublicKey))
	if err != nil {
		return err
	}

	sp.publicKey = key
	return nil
}

func (sp *AuthenticationServiceProvier) loop(ctx context.Context) {
	reqCtx := context.Background()
	reqCtx, cancel := context.WithCancel(reqCtx)

	for {
		select {
		case <-time.After(20 * time.Minute):
			err := sp.renewPublicKey(reqCtx)
			if err != nil {
				log.Printf("[CRITICAL] failed to renew a public key, err=%+v", err)
			}

		case <-ctx.Done():
			cancel()
		}
	}
}

func (sp *AuthenticationServiceProvier) GetConnectingAccountName(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "authentication metadata is required")
	}

	token := md.Get("authentication")
	if len(token) < 1 {
		return "", grpc.Errorf(codes.Unauthenticated, "authentication metadata is required")
	}

	serviceClient, err := sp.publicKey.VerifyAuthenticationToken(token[0], sp.host)
	if err != nil {
		return "", grpc.Errorf(codes.Unauthenticated, "authentication token is invalid: %s", err.Error())
	}

	return serviceClient, nil
}
