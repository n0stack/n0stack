package embed

import (
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/store/leveldb"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type EmbedDatastore struct {
	db    *leveldb.LeveldbStore
	mutex lock.MutexTable
}

func NewEmbedDatastore(dbDirectory string) (*EmbedDatastore, error) {
	l, err := leveldb.NewLeveldbStore(dbDirectory)
	if err != nil {
		return nil, err
	}

	return &EmbedDatastore{
		db: l,

		// TODO: prefix で同期されない
		mutex: lock.NewMemoryMutexTable(10000),
	}, nil
}

func (m *EmbedDatastore) AddPrefix(prefix string) datastore.Datastore {
	return &EmbedDatastore{
		db:    m.db.AddPrefix(prefix),
		mutex: lock.NewMemoryMutexTable(10000),
	}
}

func (ds EmbedDatastore) List(f func(length int) []proto.Message) error {
	b, err := ds.db.List()
	if err != nil {
		return errors.Wrap(err, "failed to list by snapshot")
	}

	pb := f(len(b))
	for i, v := range b {
		if err := proto.Unmarshal(v, pb[i]); err != nil {
			return err
		}
	}

	return nil
}

func (ds EmbedDatastore) Get(key string, pb proto.Message) error {
	v, err := ds.db.Get(key)
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(v, pb); err != nil {
		return err
	}

	return nil
}

func (ds *EmbedDatastore) Apply(key string, pb proto.Message) error {
	if !ds.mutex.IsLocked(key) {
		return errors.New("key is not locked")
	}

	s, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	if err := ds.db.Apply(key, s); err != nil {
		return errors.Wrap(err, "failed to apply for log")
	}

	return nil
}

func (ds *EmbedDatastore) Delete(key string) error {
	if !ds.mutex.IsLocked(key) {
		return errors.New("key is not locked")
	}

	if err := ds.db.Delete(key); err != nil {
		return err
	}

	return nil
}

func (ds *EmbedDatastore) Lock(key string) bool {
	return ds.mutex.Lock(key)
}
func (ds *EmbedDatastore) Unlock(key string) bool {
	return ds.mutex.Unlock(key)
}
func (ds *EmbedDatastore) IsLocked(key string) bool {
	return ds.mutex.IsLocked(key)
}
