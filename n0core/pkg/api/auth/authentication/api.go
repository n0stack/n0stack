package authn

import (
	"context"
	"crypto/rand"
	"io"
	"log"
	"reflect"

	"google.golang.org/grpc/metadata"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pauth "n0st.ac/n0stack/auth/v1alpha"
	piam "n0st.ac/n0stack/iam/v1alpha"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	jwtutil "n0st.ac/n0stack/n0core/pkg/util/jwt"
)

type AuthenticationAPI struct {
	userAPI    piam.UserServiceClient
	projectAPI piam.ProjectServiceClient

	secret []byte
}

func CreateAuthenticationAPI(userAPI piam.UserServiceClient, secret []byte) *AuthenticationAPI {
	return &AuthenticationAPI{
		userAPI: userAPI,
		secret:  secret,
	}
}

func CreatePartialAuthenticationAPI(secret []byte) *AuthenticationAPI {
	return &AuthenticationAPI{
		secret: secret,
	}
}

func (a *AuthenticationAPI) GetAuthenticationTokenPublicKey(ctx context.Context, req *pauth.GetAuthenticationTokenPublicKeyRequest) (*pauth.GetAuthenticationTokenPublicKeyResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpcutil.Errorf(codes.InvalidArgument, "set :authority header")
	}
	if len(md.Get(":authority")) < 1 {
		return nil, grpcutil.Errorf(codes.InvalidArgument, "set :authority header")
	}
	hostname := md.Get(":authority")[0]
	if hostname == "" {
		return nil, grpcutil.Errorf(codes.InvalidArgument, "set :authority header")
	}

	kg := jwtutil.NewKeyGenerator(a.secret)
	authnPrivateKey, err := kg.Generate(hostname)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &pauth.GetAuthenticationTokenPublicKeyResponse{
		PublicKey: string(authnPrivateKey.PublicKey().MarshalAuthorizedKey()),
	}, nil
}

func IsContainingPublicKey(keys map[string]string, key *jwtutil.PublicKey) (bool, error) {
	for _, k := range keys {
		parsed, err := jwtutil.ParsePublicKey([]byte(k))
		if err != nil {
			return false, err
		}

		if reflect.DeepEqual(parsed.MarshalAuthorizedKey(), key.MarshalAuthorizedKey()) {
			return true, nil
		}
	}

	return false, nil
}

func (a *AuthenticationAPI) PublicKeyAuthenticate(stream pauth.AuthenticationService_PublicKeyAuthenticateServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "")
	}
	start := req.GetStart()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return grpcutil.Errorf(codes.InvalidArgument, "set :authority header")
	}
	if len(md.Get(":authority")) < 1 {
		return grpcutil.Errorf(codes.InvalidArgument, "set :authority header")
	}
	hostname := md.Get(":authority")[0]
	if hostname == "" {
		return grpcutil.Errorf(codes.InvalidArgument, "set :authority header")
	}

	user, err := a.userAPI.GetUser(stream.Context(), &piam.GetUserRequest{
		Name: start.UserName,
	})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return grpcutil.Errorf(codes.PermissionDenied, "")
		}

		return grpcutil.Errorf(codes.Internal, "GetUser returns err=%+v", err)
	}

	challengePublicKey, err := jwtutil.ParsePublicKey([]byte(start.PublicKey))
	if err != nil {
		return grpcutil.Wrapf(codes.Internal, err, "")
	}

	if ok, err := IsContainingPublicKey(user.PublicKeys, challengePublicKey); err != nil {
		return errors.Wrap(err, "")
	} else if !ok {
		return grpcutil.Errorf(codes.PermissionDenied, "want:%+v\ngot:%s", user.PublicKeys, string(challengePublicKey.MarshalAuthorizedKey()))
	}

	challenge := make([]byte, 256)
	io.ReadFull(rand.Reader, challenge[:])
	err = stream.Send(&pauth.PublicKeyAuthenticateResponse{
		Message: &pauth.PublicKeyAuthenticateResponse_Challenge_{
			Challenge: &pauth.PublicKeyAuthenticateResponse_Challenge{
				Challenge: challenge,
			},
		},
	})
	if err != nil {
		return grpcutil.Wrapf(codes.Internal, err, "")
	}

	req, err = stream.Recv()
	if err != nil {
		return grpcutil.Wrapf(codes.Internal, err, "")
	}
	challengeToken := req.GetResponse()

	if err := challengePublicKey.VerifyChallengeToken(challengeToken.ChallengeToken, user.Name, hostname, challenge); err != nil {
		return grpcutil.Wrapf(codes.Internal, err, "")
	}

	kg := jwtutil.NewKeyGenerator(a.secret)
	authnPrivateKey, err := kg.Generate(hostname)
	if err != nil {
		return grpcutil.Wrapf(codes.Internal, err, "")
	}

	authnToken, err := authnPrivateKey.GenerateAuthenticationToken(user.Name, hostname)
	if err != nil {
		return grpcutil.Wrapf(codes.Internal, err, "")
	}
	stream.Send(&pauth.PublicKeyAuthenticateResponse{
		Message: &pauth.PublicKeyAuthenticateResponse_Result_{
			Result: &pauth.PublicKeyAuthenticateResponse_Result{
				AuthenticationToken: authnToken,
			},
		},
	})

	return nil
}
