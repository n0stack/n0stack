package memory

import (
	"path/filepath"
	"testing"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"
)

func TestMemoryStore(t *testing.T) {
	m := NewMemoryStore()

	k := "key"
	v := []byte("value")

	if err := m.Apply(k, v); err != nil {
		t.Errorf("Failed to apply: err='%s'", err.Error())
	}

	if b, err := m.Get(k); err != nil {
		t.Errorf("Failed to get: err='%s'", err.Error())
	} else if string(b) != string(v) {
		t.Errorf("get result is wrong: want=%s, have=%s", string(b), string(v))
	}

	if res, err := m.List(); err != nil {
		t.Errorf("Failed to list: key='%s', value='%v', err='%s'", k, v, err.Error())
	} else if len(res) != 1 {
		t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	} else if string(res[0]) != string(v) {
		t.Errorf("Get 'Name' is wrong: key='%s', have='%s', want='%s'", k, string(res[0]), string(v))
	}

	if err := m.Delete(k); err != nil {
		t.Errorf("Failed to delete: err='%s'", err.Error())
	}
}

func TestMemoryStoreNotFound(t *testing.T) {
	m := NewMemoryStore()

	k := "key"

	if _, err := m.Get(k); err == nil || !store.IsNotFound(err) {
		t.Errorf("error is wrong, required NotFoundError")
	}

	if err := m.Delete(k); err == nil || !store.IsNotFound(err) {
		t.Errorf("error is wrong, required NotFoundError")
	}
}

func TestCheckDataIsSame(t *testing.T) {
	m := NewMemoryStore()

	prefix := "prefix"
	withPrefix := m.AddPrefix(prefix)

	k := "key"
	v := []byte("value")

	if err := withPrefix.Apply(k, v); err != nil {
		t.Fatalf("Failed to apply: err='%s'", err.Error())
	}

	if _, err := m.Get(filepath.Join("prefix", k)); err == nil {
		t.Errorf("not got error, want NotFound")
	} else if !store.IsNotFound(err) {
		t.Errorf("got error which is not NotFound: err=%s", err.Error())
	}

	k2 := "key"
	v2 := []byte("value")

	if err := m.Apply(k2, v2); err != nil {
		t.Fatalf("Failed to apply secondary: err='%s'", err.Error())
	}

	if res, err := withPrefix.List(); err != nil {
		t.Errorf("Failed to list: err='%s'", err.Error())
	} else if len(res) != 1 {
		t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	}
}

func BenchmarkMemoryStoreApply(b *testing.B) {
	m := NewMemoryStore()

	k := "key"
	v := []byte("value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := k + string(i)
		m.Apply(key, v)
	}
}

func BenchmarkMemoryStoreDeleteAfterApply(b *testing.B) {
	m := NewMemoryStore()

	k := "key"
	v := []byte("value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := k + string(i)
		m.Apply(key, v)
		m.Delete(key)
	}
}

func BenchmarkMemoryStoreGet(b *testing.B) {
	m := NewMemoryStore()

	k := "key"
	v := []byte("value")

	for i := 0; i < b.N; i++ {
		key := k + string(i)
		m.Apply(key, v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := k + string(i)
		m.Get(key)
	}
}

func BenchmarkMemoryStoreList(b *testing.B) {
	m := NewMemoryStore()

	k := "key"
	v := []byte("value")

	for i := 0; i < 1000; i++ {
		key := k + string(i)
		m.Apply(key, v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.List()
	}
}
