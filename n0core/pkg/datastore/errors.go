package datastore

import (
	fmt "fmt"

	"github.com/pkg/errors"
)

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

func IsNotFound(err error) bool {
	if e, ok := err.(NotFoundError); ok {
		return e.IsNotFound()
	}

	return false
}

// func LockErrorMessage() string {
// 	return "this is locked, wait a moment"
// }

func DefaultErrorMessage(err error) string {
	return errors.Wrap(err, "Failed to operate on datastore").Error()
}

type ConflictedError struct{}

func (c ConflictedError) Error() string {
	return "conflicted"
}
