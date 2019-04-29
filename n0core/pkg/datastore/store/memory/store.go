package memory

import (
	"path/filepath"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"
)

var data map[string]map[string][]byte

func init() {
	data = make(map[string]map[string][]byte)
	data[""] = make(map[string][]byte)
}

type MemoryStore struct {
	prefix string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) AddPrefix(prefix string) *MemoryStore {
	p := filepath.Join(m.prefix, prefix)
	data[p] = make(map[string][]byte)

	return &MemoryStore{
		prefix: p,
	}
}

func (m MemoryStore) List() ([][]byte, error) {
	l := make([][]byte, len(data[m.prefix]))
	i := 0
	for _, v := range data[m.prefix] {
		l[i] = v
		i++
	}

	return l, nil
}

func (m MemoryStore) Get(key string) ([]byte, error) {
	v, ok := data[m.prefix][key]
	if !ok {
		return nil, store.NewNotFound(key)
	}

	return v, nil
}

func (m *MemoryStore) Apply(key string, value []byte) error {
	data[m.prefix][key] = value

	return nil
}

func (m *MemoryStore) Delete(key string) error {
	if !m.IsExisting(key) {
		return store.NewNotFound(key)
	}

	delete(data[m.prefix], key)
	return nil
}

func (m MemoryStore) IsExisting(key string) bool {
	_, ok := data[m.prefix][key]
	return ok
}
