package memory

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
)

type MemoryDatastore struct {
	// 本当は `proto.Message` を入れたいが、何故か中身がなかったのでとりあえずシリアライズする
	Data map[string][]byte

	prefix string
}

func NewMemoryDatastore() *MemoryDatastore {
	return &MemoryDatastore{Data: map[string][]byte{}}
}

func (m *MemoryDatastore) AddPrefix(prefix string) datastore.Datastore {
	return &MemoryDatastore{
		Data:   m.Data,
		prefix: m.prefix + prefix + "/",
	}
}

func (m MemoryDatastore) List(f func(length int) []proto.Message) error {
	l := 0
	for k, _ := range m.Data {
		if !strings.HasPrefix(k, m.prefix) {
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

func (m MemoryDatastore) Get(name string, pb proto.Message) error {
	v, ok := m.Data[m.prefix+name]
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

func (m MemoryDatastore) Apply(name string, pb proto.Message) error {
	s, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	m.Data[m.prefix+name] = s

	return nil
}

func (m MemoryDatastore) Delete(name string) error {
	var ok bool
	_, ok = m.Data[m.prefix+name]
	if ok {
		delete(m.Data, m.prefix+name)
		return nil
	}

	return nil
}
