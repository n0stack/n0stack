package memory

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0core/pkg/datastore"
)

func TestMemoryDatastore(t *testing.T) {
	m := NewMemoryDatastore()

	k := "test"
	v := &datastore.Test{Name: "value"}

	if err := m.Apply(k, v); err != nil {
		t.Errorf("Failed to apply: err='%s'", err.Error())
	}

	e := &datastore.Test{}
	if err := m.Get(k, e); err != nil {
		t.Errorf("Failed to get: err='%s'", err.Error())
	} else if e == nil {
		t.Errorf("Failed to get: result is nil")
	}

	res := []*datastore.Test{}
	f := func(s int) []proto.Message {
		res = make([]*datastore.Test, s)
		for i := range res {
			res[i] = &datastore.Test{}
		}

		m := make([]proto.Message, s)
		for i, v := range res {
			m[i] = v
		}

		return m
	}
	if err := m.List(f); err != nil {
		t.Errorf("Failed to list: key='%s', value='%v', err='%s'", k, v, err.Error())
	}
	if len(res) != 1 {
		t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	}
	if res[0].Name != v.Name {
		t.Errorf("Get 'Name' is wrong: key='%s', have='%s', want='%s'", k, res[0].Name, v.Name)
	}

	if err := m.Delete(k); err != nil {
		t.Errorf("Failed to delete: err='%s'", err.Error())
	}
}
