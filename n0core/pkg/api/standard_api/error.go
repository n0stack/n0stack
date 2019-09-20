package stdapi

import (
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
	"google.golang.org/grpc/codes"
)

func ValidationError(field, validationFormat string) error {
	return grpcutil.Errorf(codes.InvalidArgument, "the %s filed validation is failed: the format is %s", field, validationFormat)
}
