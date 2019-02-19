package datastore

import (
	"github.com/golang/protobuf/proto"
)

type Datastore interface {
	AddPrefix(prefix string) Datastore

	List(f func(length int) []proto.Message) error

	// if result is empty, set pb as nil.
	Get(key string, pb proto.Message) error

	// update system requires locking in advance
	Apply(key string, pb proto.Message) error
	Delete(key string) error

	Lock(key string) bool
	Unlock(key string) bool
	IsLocked(key string) bool
}
