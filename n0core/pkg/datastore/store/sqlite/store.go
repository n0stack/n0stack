package sqlite

import (
	"path/filepath"
	"time"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type LogType int

const (
	LOGTYPE_APPLY LogType = iota
	LOGTYPE_DELETE
)

type Log struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"index"`

	Type   LogType
	Prefix string `gorm:"index:idx_prefix_key"`
	Key    string `gorm:"index:idx_prefix_key"`
	Value  []byte
}

type SqliteStore struct {
	db *gorm.DB

	prefix string
}

func NewSqliteStore(file string) (*SqliteStore, error) {
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	db.AutoMigrate(&Log{})

	return &SqliteStore{
		db: db,
	}, nil
}

func (ds *SqliteStore) Close() error {
	return ds.db.Close()
}

func (ds *SqliteStore) AddPrefix(prefix string) store.Store {
	return &SqliteStore{
		db:     ds.db,
		prefix: filepath.Join(ds.prefix, prefix),
	}
}

func (ds *SqliteStore) List(f func(length int) []proto.Message) error {
	return nil
}

func (ds *SqliteStore) Get(key string, pb proto.Message) error {
	l := &Log{}
	if err := ds.db.Last(&l).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return store.NewNotFound(key)
		}

		return err
	}

	if l.Type == LOGTYPE_DELETE {
		return store.NewNotFound(key)
	}

	err := proto.Unmarshal(l.Value, pb)
	if err != nil {
		return err
	}

	return nil
}

func (ds *SqliteStore) Apply(key string, pb proto.Message) error {
	s, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	return ds.db.Create(&Log{
		Type:   LOGTYPE_APPLY,
		Prefix: ds.prefix,
		Key:    key,
		Value:  s,
	}).Error
}

func (ds *SqliteStore) Delete(key string) error {
	return ds.db.Create(&Log{
		Type:   LOGTYPE_DELETE,
		Prefix: ds.prefix,
		Key:    key,
	}).Error
}
