package stdapi

import (
	"google.golang.org/grpc/codes"
	grpcutil "n0st.ac/n0stack/n0core/pkg/util/grpc"
)

func ValidationError(field, validationFormat string) error {
	return grpcutil.Errorf(codes.InvalidArgument, "the %s filed validation is failed: the format is %s", field, validationFormat)
}

func UpdateMaskError(err error) error {
	return grpcutil.Errorf(codes.InvalidArgument, "failed about update_mask, err=%s", err.Error())
}
