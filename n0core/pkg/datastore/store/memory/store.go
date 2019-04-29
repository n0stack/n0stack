package memory

import (
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"
)

type MemoryStore struct {
	// 本当は `proto.Message` を入れたいが、何故か中身がなかったのでとりあえずシリアライズする
	Data map[string][]byte

	prefix string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Data: map[string][]byte{},
	}
}

func (m *MemoryStore) AddPrefix(prefix string) store.Store {
	return &MemoryStore{
		Data:   m.Data,
		prefix: filepath.Join(m.prefix+prefix, "/"),
	}
}

func (m MemoryStore) List(f func(length int) []proto.Message) error {
	l := 0
	for k, _ := range m.Data {
		if strings.HasPrefix(k, m.prefix) {
			l++
		}
	}

	pb := f(l)
	i := 0
	for k, v := range m.Data {
		if !strings.HasPrefix(k, m.prefix) {
			continue
		}

		err := proto.Unmarshal(v, pb[i])
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

func (m MemoryStore) Get(key string, pb proto.Message) error {
	v, ok := m.Data[m.prefix+key]
	if !ok {
		pb = nil
		return store.NewNotFound(key)
	}

	err := proto.Unmarshal(v, pb)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemoryStore) Apply(key string, pb proto.Message) error {
	s, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	m.Data[m.getKey(key)] = s

	return nil
}

func (m *MemoryStore) Delete(key string) error {
	var ok bool
	_, ok = m.Data[m.getKey(key)]
	if ok {
		delete(m.Data, m.getKey(key))
		return nil
	}

	return store.NewNotFound(key)
}

func (m MemoryStore) getKey(key string) string {
	return filepath.Join(m.prefix, key)
}
