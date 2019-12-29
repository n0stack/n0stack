package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
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
func NewAuthenticationServiceProvider(ctx context.Context, authnAPI pauth.AuthenticationServiceClient, host string) (*AuthenticationServiceProvier, error) {
	sp := &AuthenticationServiceProvier{
		authnAPI: authnAPI,
		host:     host,
	}

	err := sp.renewPublicKey(ctx)
	if err != nil {
		log.Printf("[CRITICAL] failed to renew a public key, err=%+v", err)
		return nil, errors.Wrap(err, "renewPublicKey() is failed")
	}
	log.Printf("[INFO] got a public key to verify authentication tokens")

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
	for {
		select {
		case <-time.After(20 * time.Minute):
			log.Printf("[INFO] renewing a public key to verify authentication tokens")
			err := sp.renewPublicKey(ctx)
			if err != nil {
				log.Printf("[CRITICAL] failed to renew a public key, err=%+v", err)
			}
			log.Printf("[INFO] renewed a public key to verify authentication tokens")

		case <-ctx.Done():
			return
		}
	}
}

// TODO: grpc.Error を返してもいいのか…？
func (sp *AuthenticationServiceProvier) GetConnectingAccountName(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("n0stack-authentication metadata is required")
	}

	token := md.Get("n0stack-authentication")
	if len(token) < 1 {
		return "", fmt.Errorf("n0stack-authentication metadata is required")
	}

	serviceClient, err := sp.publicKey.VerifyAuthenticationToken(token[0], sp.host)
	if err != nil {
		return "", fmt.Errorf("authentication token is invalid: %s", err.Error())
	}

	return serviceClient, nil
}
