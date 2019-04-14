package datastore

import (
	fmt "fmt"

	"google.golang.org/grpc/codes"

	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
)

func LockError() error {
	return grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "this is locked, wait a moment")
}

type NotFoundError interface {
	IsNotFound() bool

	Error() string
}

type NotFound struct {
	key string
}

func NewNotFound(key string) *NotFound {
	return &NotFound{
		key: key,
	}
}

func (e NotFound) IsNotFound() bool {
	return true
}

func (e NotFound) Error() string {
	return fmt.Sprintf("Key '%s' is not found", e.key)
}

func IsNotfound(err error) bool {
	if e, ok := err.(NotFoundError); ok {
		return e.IsNotFound()
	}

	return false
}

// func LockErrorMessage() string {
// 	return "this is locked, wait a moment"
// }
