package grpcutil

import (
	"log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const InternalErrorMessage = "please try again later or contact your administrator"

// Errorf returns grpc.Errorf
// in the case of 'Internal', logging message because the server has failed
func Errorf(c codes.Code, format string, a ...interface{}) error {
	var err error

	if c == codes.Internal {
		// TODO: 遡った行数を出す
		log.Printf("[CRITICAL] "+format, a...)

		err = grpc.Errorf(c, InternalErrorMessage)
	} else {
		err = grpc.Errorf(c, format, a...)
	}

	return err
}

func Wrapf(c codes.Code, e error, format string, a ...interface{}) error {
	wrapped := errors.Wrapf(e, format, a)

	return Errorf(c, wrapped.Error())
}
