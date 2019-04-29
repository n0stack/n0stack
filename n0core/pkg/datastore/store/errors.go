package store

import (
	"fmt"
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
