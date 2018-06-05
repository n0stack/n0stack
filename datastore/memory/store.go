package memory

import (
	"github.com/golang/protobuf/proto"
)

type MemoryDatastore struct {
	// 本当は `proto.Message` を入れたいが、何故か中身がなかったのでとりあえずシリアライズする
	Data map[string][]byte
}

func NewMemoryDatastore() *MemoryDatastore {
	return &MemoryDatastore{Data: map[string][]byte{}}
}

func (m MemoryDatastore) List(f func(length int) []proto.Message) error {
	pb := f(len(m.Data))

	i := 0
	for _, v := range m.Data {
		err := proto.Unmarshal(v, pb[i])
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

func (m MemoryDatastore) Get(name string, pb proto.Message) error {
	err := proto.Unmarshal(m.Data[name], pb)
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

	m.Data[name] = s

	return nil
}

func (m MemoryDatastore) Delete(name string) (int64, error) {
	var ok bool
	_, ok = m.Data[name]

	if ok {
		delete(m.Data, name)
		return 1, nil
	}

	return 0, nil
}
