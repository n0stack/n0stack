package auth

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	jwtutil "n0st.ac/n0stack/n0core/pkg/util/jwt"

	"github.com/pkg/errors"
	pauth "n0st.ac/n0stack/auth/v1alpha"
)

type AuthenticationClient struct {
	privateKey *jwtutil.PrivateKey
	user       string

	audience string
	authnAPI pauth.AuthenticationServiceClient

	token string
}

func NewAuthenticationClient(ctx context.Context, conn *grpc.ClientConn, user string, privateKey *jwtutil.PrivateKey) (*AuthenticationClient, error) {
	authnAPI := pauth.NewAuthenticationServiceClient(conn)
	c := &AuthenticationClient{
		privateKey: privateKey,
		user:       user,

		authnAPI: authnAPI,
		audience: conn.Target(),
	}

	err := c.renewAuthenticationToken(ctx)
	if err != nil {
		log.Printf("[CRITICAL] failed to renew an authentication token, err=%+v", err)
		return nil, errors.Wrap(err, "renewAuthenticationToken() is failed")
	}

	go c.loop(ctx)

	return c, nil
}

func (c *AuthenticationClient) loop(ctx context.Context) {
	reqToken := context.Background()
	reqToken, cancel := context.WithCancel(reqToken)

	for {
		select {
		case <-time.After(20 * time.Minute):
			err := c.renewAuthenticationToken(reqToken)
			if err != nil {
				log.Printf("[CRITICAL] failed to renew an authentication token, err=%+v", err)
			}

		case <-ctx.Done():
			cancel()
		}
	}
}

func (c *AuthenticationClient) renewAuthenticationToken(ctx context.Context) error {
	stream, err := c.authnAPI.PublicKeyAuthenticate(ctx)
	if err != nil {
		return errors.Wrapf(err, "")
	}

	err = stream.Send(&pauth.PublicKeyAuthenticateRequest{
		Message: &pauth.PublicKeyAuthenticateRequest_Start_{
			Start: &pauth.PublicKeyAuthenticateRequest_Start{
				UserName:  c.user,
				PublicKey: string(c.privateKey.PublicKey().MarshalAuthorizedKey()),
			},
		},
	})
	if err != nil {
		return errors.Wrapf(err, "")
	}

	res, err := stream.Recv()
	if err != nil {
		return errors.Wrapf(err, "")
	}
	challenge := res.GetChallenge()

	challengeToken, err := c.privateKey.GenerateChallengeToken(c.user, c.audience, challenge.Challenge)
	if err != nil {
		return errors.Wrap(err, "")
	}

	err = stream.Send(&pauth.PublicKeyAuthenticateRequest{
		Message: &pauth.PublicKeyAuthenticateRequest_Response_{
			Response: &pauth.PublicKeyAuthenticateRequest_Response{
				ChallengeToken: challengeToken,
			},
		},
	})
	if err != nil {
		return errors.Wrapf(err, "")
	}

	res, err = stream.Recv()
	if err != nil {
		return errors.Wrapf(err, "")
	}
	result := res.GetResult()
	c.token = result.AuthenticationToken

	return nil
}

func (c AuthenticationClient) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"n0stack-authentication": c.token,
	}, nil
}

func (c AuthenticationClient) RequireTransportSecurity() bool {
	return true
}
