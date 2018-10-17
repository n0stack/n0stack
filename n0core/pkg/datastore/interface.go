package datastore

import (
	"github.com/golang/protobuf/proto"
)

type Datastore interface {
	List(f func(length int) []proto.Message) error

	// if result is empty, set pb as nil.
	Get(name string, pb proto.Message) error

	Apply(name string, pb proto.Message) error

	// Delete returns how many query was deleted and error
	Delete(name string) error // TODO: returnをerrorのみにする
}
