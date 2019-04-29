package embed

import (
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"
)

const dbDir = "test.db"

func TestEmbedDatastore(t *testing.T) {
	m, err := NewEmbedDatastore(dbDir)
	if err != nil {
		t.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "test"
	v := &datastore.Test{Name: "value"}

	if !m.Lock(k) {
		t.Errorf("Failed to lock")
	}
	defer m.Unlock(k)

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

func TestEmbedDatastoreNotFound(t *testing.T) {
	m, err := NewEmbedDatastore(dbDir)
	if err != nil {
		t.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "test"

	e := &datastore.Test{}
	if err := m.Get(k, e); err == nil || !store.IsNotFound(err) {
		t.Errorf("error is wrong, required NotFoundError")
	}
}

func TestCheckDataIsSame(t *testing.T) {
	m, err := NewEmbedDatastore(dbDir)
	if err != nil {
		t.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	prefix := "prefix"
	withPrefix := m.AddPrefix(prefix)

	k := "test"
	v := &datastore.Test{Name: "value"}

	withPrefix.Lock(k)
	defer withPrefix.Unlock(k)

	if err := withPrefix.Apply(k, v); err != nil {
		t.Fatalf("Failed to apply: err='%s'", err.Error())
	}
	e := &datastore.Test{}
	if err := m.Get("prefix/"+k, e); err != nil {
		t.Errorf("Failed to get: err=%s", err.Error())
	}
	if e == nil || e.Name != v.Name {
		t.Errorf("Response is invalid")
	}

	k2 := "test"
	v2 := &datastore.Test{Name: "value"}

	m.Lock(k2)
	defer m.Unlock(k2)

	if err := m.Apply(k2, v2); err != nil {
		t.Fatalf("Failed to apply secondaly: err='%s'", err.Error())
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

func TestUpdateSystemBeforeLock(t *testing.T) {
	m, err := NewEmbedDatastore(dbDir)
	if err != nil {
		t.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "test"
	v := &datastore.Test{Name: "value"}

	if err := m.Apply(k, v); err == nil {
		t.Errorf("applied before lock")
	}

	if err := m.Delete(k); err == nil {
		t.Errorf("deleted before lock")
	}
}

func BenchmarkEmbedDatastore(b *testing.B) {
	m, err := NewEmbedDatastore("test.db")
	if err != nil {
		b.Fatalf("failed to generate datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "key"
	v := &datastore.Test{Name: "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := k + string(i)
		m.Lock(key)
		defer m.Unlock(key)

		m.Apply(key, v)
	}
}
