package sqlite

import (
	"path/filepath"
	"time"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"

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

func (ds *SqliteStore) AddPrefix(prefix string) *SqliteStore {
	return &SqliteStore{
		db:     ds.db,
		prefix: filepath.Join(ds.prefix, prefix),
	}
}

func (ds *SqliteStore) List() ([][]byte, error) {
	return nil, nil
}

func (ds *SqliteStore) Get(key string) ([]byte, error) {
	l := &Log{}
	if err := ds.db.Last(&l).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, store.NewNotFound(key)
		}

		return nil, err
	}

	if l.Type == LOGTYPE_DELETE {
		return nil, store.NewNotFound(key)
	}

	return l.Value, nil
}

func (ds *SqliteStore) Apply(key string, value []byte) error {
	return ds.db.Create(&Log{
		Type:   LOGTYPE_APPLY,
		Prefix: ds.prefix,
		Key:    key,
		Value:  value,
	}).Error
}

func (ds *SqliteStore) Delete(key string) error {
	return ds.db.Create(&Log{
		Type:   LOGTYPE_DELETE,
		Prefix: ds.prefix,
		Key:    key,
	}).Error
}
