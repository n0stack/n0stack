package memory

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"
)

type MemoryDatastore struct {
	// 本当は `proto.Message` を入れたいが、何故か中身がなかったのでとりあえずシリアライズする
	Data map[string][]byte

	mutex  lock.MutexTable
	prefix string
}

func NewMemoryDatastore() *MemoryDatastore {
	return &MemoryDatastore{
		Data:  map[string][]byte{},
		mutex: lock.NewMemoryMutexTable(10000),
	}
}

func (m *MemoryDatastore) AddPrefix(prefix string) datastore.Datastore {
	return &MemoryDatastore{
		Data:   m.Data,
		prefix: m.prefix + prefix + "/",
		mutex:  lock.NewMemoryMutexTable(10000),
	}
}

func (m MemoryDatastore) List(f func(length int) []proto.Message) error {
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

func (m MemoryDatastore) Get(key string, pb proto.Message) error {
	v, ok := m.Data[m.prefix+key]
	if !ok {
		pb = nil
		return nil
	}

	err := proto.Unmarshal(v, pb)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemoryDatastore) Apply(key string, pb proto.Message) error {
	if !m.mutex.IsLocked(m.getKey(key)) {
		return errors.New("key is not locked")
	}

	s, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	m.Data[m.getKey(key)] = s

	return nil
}

func (m *MemoryDatastore) Delete(key string) error {
	if !m.mutex.IsLocked(m.getKey(key)) {
		return errors.New("key is not locked")
	}

	var ok bool
	_, ok = m.Data[m.getKey(key)]
	if ok {
		delete(m.Data, m.getKey(key))
		return nil
	}

	return nil
}

func (m MemoryDatastore) getKey(key string) string {
	return m.prefix + key
}

func (m *MemoryDatastore) Lock(key string) bool {
	return m.mutex.Lock(m.getKey(key))
}
func (m *MemoryDatastore) Unlock(key string) bool {
	return m.mutex.Unlock(m.getKey(key))
}
func (m *MemoryDatastore) IsLocked(key string) bool {
	return m.mutex.IsLocked(m.getKey(key))
}
