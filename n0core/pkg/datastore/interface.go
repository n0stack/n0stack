package datastore

import (
	"github.com/golang/protobuf/proto"
)

type Datastore interface {
	AddPrefix(prefix string) Datastore

	List(f func(length int) []proto.Message) error

	Get(key string, pb proto.Message) (int64, error)

	Apply(key string, pb proto.Message, currentVersion int64) (int64, error)
	Delete(key string, currentVersion int64) error
}
