package grpcutil

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// WrapGrpcErrorf returns grpc.Errorf
// in the case of 'Internal', logging message because the server has failed
func WrapGrpcErrorf(c codes.Code, format string, a ...interface{}) error {
	err := grpc.Errorf(c, format, a...)

	if c == codes.Internal {
		log.Printf("[WARNING] "+format, a...)
	}

	return err
}
