package leveldb

import (
	"path/filepath"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	lerrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LeveldbStore struct {
	db *leveldb.DB

	prefix string
}

func NewLeveldbStore(directory string) (*LeveldbStore, error) {
	db, err := leveldb.OpenFile(directory, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	return &LeveldbStore{
		db: db,
	}, nil
}

func (ds *LeveldbStore) Close() error {
	return ds.db.Close()
}

func (ds *LeveldbStore) AddPrefix(prefix string) *LeveldbStore {
	return &LeveldbStore{
		db:     ds.db,
		prefix: filepath.Join(ds.prefix, prefix),
	}
}

func (ds *LeveldbStore) List() ([][]byte, error) {
	res := make([][]byte, 0)

	iter := ds.db.NewIterator(util.BytesPrefix([]byte(ds.prefix)), nil)
	for iter.Next() {
		res = append(res, iter.Value())
	}
	iter.Release()

	if err := iter.Error(); err != nil {
		return nil, err
	}

	return res, nil
}

func (ds *LeveldbStore) Get(key string) ([]byte, error) {
	v, err := ds.db.Get(ds.getKey(key), nil)
	if err != nil {
		if err == lerrors.ErrNotFound {
			return nil, store.NewNotFound(key)
		}

		return nil, err
	}

	return v, nil
}

func (ds *LeveldbStore) Apply(key string, value []byte) error {
	return ds.db.Put(ds.getKey(key), value, nil)
}

func (ds *LeveldbStore) Delete(key string) error {
	return ds.db.Delete(ds.getKey(key), nil)
}

func (ds LeveldbStore) getKey(key string) []byte {
	return []byte(filepath.Join(ds.prefix, key))
}
