package datastore

import (
	"context"

	"github.com/golang/protobuf/proto"
)

type Datastore interface {
	AddPrefix(prefix string) Datastore

	List(ctx context.Context, f func(length int) []proto.Message) error

	Get(ctx context.Context, key string, pb proto.Message) (int64, error)

	Apply(ctx context.Context, key string, pb proto.Message, currentVersion int64) (int64, error)
	Delete(ctx context.Context, key string, currentVersion int64) error
}
