package store

import "github.com/golang/protobuf/proto"

type Store interface {
	AddPrefix(prefix string) Store

	List(f func(length int) []proto.Message) error
	Get(key string, pb proto.Message) error
	Apply(key string, pb proto.Message) error
	Delete(key string) error
}
