package stdapi

import (
	"context"

	"google.golang.org/grpc/codes"

	"n0st.ac/n0stack/n0core/pkg/driver/n0stack/auth"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
)

func CheckAuthenticatedUserName(auth *auth.AuthenticationServiceProvier, ctx context.Context, username string) error {
	luser, err := auth.GetConnectingAccountName(ctx)
	if err != nil {
		return grpcutil.Errorf(codes.Unauthenticated, err.Error())
	}
	if luser != username {
		return grpcutil.Errorf(codes.PermissionDenied, "you must login as %s to call this RPC, but you are %s", username, luser)
	}

	return nil
}
