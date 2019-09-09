package memory

import (
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
)

func TestMemoryDatastore(t *testing.T) {
	m := NewMemoryDatastore()

	k := "test"
	v := &datastore.Test{Name: "value"}

	if version, err := m.Apply(k, v, 0); err != nil {
		t.Fatalf("Apply('%v', '%v', '%v') err='%s'", k, v, 0, err.Error())
	} else if version != 1 {
		t.Errorf("Apply('%v', '%v', '%v') return wrong version '%v'", k, v, 0, version)
	}

	e := &datastore.Test{}
	if v, err := m.Get(k, e); err != nil {
		t.Errorf("Get() err='%s'", err.Error())
	} else if e == nil {
		t.Errorf("Get() result is nil")
	} else if v != 1 {
		t.Errorf("Get() got wrong version just after created")
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
		t.Errorf("List() key='%s', value='%v', err='%s'", k, v, err.Error())
	}
	if len(res) != 1 {
		t.Errorf("List() number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	}
	if res[0].Name != v.Name {
		t.Errorf("List() got 'Name' is wrong: key='%s', have='%s', want='%s'", k, res[0].Name, v.Name)
	}

	if err := m.Delete(k, 1); err != nil {
		t.Errorf("Delete() err='%s'", err.Error())
	}
}

func TestMemoryDatastoreNotFound(t *testing.T) {
	m := NewMemoryDatastore()
	k := "test"

	e := &datastore.Test{}
	if _, err := m.Get(k, e); err == nil || !datastore.IsNotFound(err) {
		t.Errorf("Get() error is wrong, required NotFoundError")
	}

	if err := m.Delete(k, 0); err == nil || !datastore.IsNotFound(err) {
		t.Errorf("Delete() error is wrong, required NotFoundError")
	}
}

func TestConfliction(t *testing.T) {
	m := NewMemoryDatastore()
	k := "test"
	v := &datastore.Test{Name: "value"}

	if _, err := m.Apply(k, v, 0); err != nil {
		t.Fatalf("Apply('%v', '%v', '%v') err='%s'", k, v, 0, err.Error())
	}
	if _, err := m.Apply(k, v, 1); err != nil {
		t.Fatalf("Apply('%v', '%v', '%v') err='%s'", k, v, 1, err.Error())
	}
	if _, err := m.Apply(k, v, 1); err == nil {
		t.Errorf("Apply('%v', '%v', '%v') no error on applying confliction", k, v, 1)
	} else if _, ok := err.(datastore.ConflictedError); !ok {
		t.Errorf("Apply('%v', '%v', '%v') wrong error on applying confliction: err=%+v", k, v, 1, err)
	}

	if err := m.Delete(k, 1); err == nil {
		t.Errorf("Delete('%v', '%v') no error on applying confliction", k, 1)
	}
	if err := m.Delete(k, 2); err != nil {
		t.Errorf("Apply('%v', '%v') err='%s'", k, 2, err.Error())
	}
}

func TestPrefixCollision(t *testing.T) {
	m := NewMemoryDatastore()

	prefix := "prefix"
	withPrefix := m.AddPrefix(prefix)

	k := "test"
	v := &datastore.Test{Name: "value"}

	if _, err := withPrefix.Apply(k, v, 0); err != nil {
		t.Fatalf("Failed to apply: err='%s'", err.Error())
	}
	e := &datastore.Test{}
	if _, err := m.Get(filepath.Join(prefix, k), e); err != nil {
		t.Errorf("Failed to get: err=%s", err.Error())
	}
	if e == nil || e.Name != v.Name {
		t.Errorf("Response is invalid")
	}

	k2 := "test"
	v2 := &datastore.Test{Name: "value"}

	if _, err := m.Apply(k2, v2, 0); err != nil {
		t.Fatalf("Failed to apply secondary: err='%s'", err.Error())
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
	if err := withPrefix.List(f); err != nil {
		t.Errorf("Failed to list: err='%s'", err.Error())
	}
	if len(res) != 1 {
		t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	}
}
