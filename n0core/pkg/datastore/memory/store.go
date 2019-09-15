package memory

import (
	"context"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
)

type data struct {
	Version int64
	Message []byte
}

type MemoryDatastore struct {
	// 本当は `proto.Message` を入れたいが、何故か中身がなかったのでとりあえずシリアライズする
	Data map[string]data

	prefix string
}

func NewMemoryDatastore() *MemoryDatastore {
	return &MemoryDatastore{
		Data: map[string]data{},
	}
}

func (m *MemoryDatastore) AddPrefix(prefix string) datastore.Datastore {
	return &MemoryDatastore{
		Data:   m.Data,
		prefix: m.prefix + prefix + "/",
	}
}

func (m MemoryDatastore) List(ctx context.Context, f func(length int) []proto.Message) error {
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

		err := proto.Unmarshal(v.Message, pb[i])
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

func (m MemoryDatastore) Get(ctx context.Context, key string, pb proto.Message) (int64, error) {
	v, ok := m.Data[m.prefix+key]
	if !ok {
		pb = nil
		return 0, datastore.NewNotFound(key)
	}

	err := proto.Unmarshal(v.Message, pb)
	if err != nil {
		return 0, err
	}

	return v.Version, nil
}

func (m *MemoryDatastore) Apply(ctx context.Context, key string, pb proto.Message, currentVersion int64) (int64, error) {
	var nextVersion int64 = 1
	if v, ok := m.Data[m.getKey(key)]; ok {
		if currentVersion < v.Version {
			return 0, datastore.ConflictedError{}
		}

		nextVersion = v.Version + 1
	}

	s, err := proto.Marshal(pb)
	if err != nil {
		return 0, err
	}

	m.Data[m.getKey(key)] = data{
		Message: s,
		Version: nextVersion,
	}

	return nextVersion, nil
}

func (m *MemoryDatastore) Delete(ctx context.Context, key string, currentVersion int64) error {
	v, ok := m.Data[m.getKey(key)]
	if ok {
		if currentVersion < v.Version {
			return datastore.ConflictedError{}
		}

		delete(m.Data, m.getKey(key))
		return nil
	}

	return datastore.NewNotFound(key)
}

func (m MemoryDatastore) getKey(key string) string {
	return m.prefix + key
}
