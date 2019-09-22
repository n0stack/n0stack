package auth

import (
	"context"

	"google.golang.org/grpc"

	jwtutil "n0st.ac/n0stack/n0core/pkg/util/jwt"

	"github.com/pkg/errors"
	pauth "n0st.ac/n0stack/auth/v1alpha"
)

type AuthenticationClient struct {
	privateKey *jwtutil.PrivateKey
	audience   string

	authnAPI pauth.AuthenticationServiceClient
}

func NewAuthenticationClient(conn *grpc.ClientConn, privateKey *jwtutil.PrivateKey) *AuthenticationClient {
	authnAPI := pauth.NewAuthenticationServiceClient(conn)
	return &AuthenticationClient{
		privateKey: privateKey,
		authnAPI:   authnAPI,
		audience:   conn.Target(),
	}
}

func (c AuthenticationClient) GetAuthenticationToken(ctx context.Context, user string) (string, error) {
	stream, err := c.authnAPI.PublicKeyAuthenricate(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "")
	}

	err = stream.Send(&pauth.PublicKeyAuthenricateRequest{
		Message: &pauth.PublicKeyAuthenricateRequest_Start_{
			Start: &pauth.PublicKeyAuthenricateRequest_Start{
				UserName:  user,
				PublicKey: string(c.privateKey.PublicKey().MarshalAuthorizedKey()),
			},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "")
	}

	res, err := stream.Recv()
	if err != nil {
		return "", errors.Wrapf(err, "")
	}
	challenge := res.GetChallenge()

	challengeToken, err := c.privateKey.GenerateChallengeToken(user, c.audience, challenge.Challenge)
	if err != nil {
		return "", errors.Wrap(err, "")
	}

	err = stream.Send(&pauth.PublicKeyAuthenricateRequest{
		Message: &pauth.PublicKeyAuthenricateRequest_Response_{
			Response: &pauth.PublicKeyAuthenricateRequest_Response{
				ChallengeToken: challengeToken,
			},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "")
	}

	res, err = stream.Recv()
	if err != nil {
		return "", errors.Wrapf(err, "")
	}
	result := res.GetResult()

	return result.AuthenticationToken, nil
}
