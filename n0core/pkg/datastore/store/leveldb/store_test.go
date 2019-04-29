package leveldb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"
)

const dbDir = "test.db"

func TestLeveldbStore(t *testing.T) {
	ds, err := NewLeveldbStore(dbDir)
	if err != nil {
		t.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "key"
	v := []byte("value")

	if _, err := ds.Get(k); err == nil {
		t.Errorf("Get() does not return error, want NotFound")
	} else if !store.IsNotFound(err) {
		t.Errorf("Get() return wrong error, want NotFound: %s", err.Error())
	}

	if err := ds.Apply(k, v); err != nil {
		t.Fatalf("failed to apply data: %s", err.Error())
	}

	if b, err := ds.Get(k); err != nil {
		t.Errorf("failed to get stored data: %s", err.Error())
	} else if string(v) != string(b) {
		t.Errorf("Get result is wrong: want=%s, have=%s", string(v), string(b))
	}

	if b, err := ds.List(); err != nil {
		t.Errorf("failed to list: %s", err.Error())
	} else if len(b) != 1 {
		t.Errorf("list length is wrong: want=%d, have=%d", 1, len(b))
	}

	if err := ds.Delete(k); err != nil {
		t.Errorf("failed to delete data: %s", err.Error())
	}
	if _, err := ds.Get(k); err == nil {
		t.Errorf("Get() does not return error, want NotFound")
	} else if !store.IsNotFound(err) {
		t.Errorf("Get() return wrong error, want NotFound: %s", err.Error())
	}

	if err := ds.Close(); err != nil {
		t.Errorf("failed to close db: %s", err.Error())
	}
}

func TestCheckDataIsSame(t *testing.T) {
	ds, err := NewLeveldbStore(dbDir)
	if err != nil {
		t.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	prefix := "prefix"
	withPrefix := ds.AddPrefix(prefix)

	k := "key"
	v := []byte("value")

	if err := withPrefix.Apply(k, v); err != nil {
		t.Fatalf("Failed to apply: err='%s'", err.Error())
	}

	if b, err := ds.Get(filepath.Join("prefix", k)); err != nil {
		t.Errorf("Failed to get: err=%s", err.Error())
	} else if string(b) != string(v) {
		t.Errorf("Response is invalid: want=%s, have=%s", string(v), string(b))
	}

	k2 := "key"
	v2 := []byte("value")

	if err := ds.Apply(k2, v2); err != nil {
		t.Fatalf("Failed to apply secondary: err='%s'", err.Error())
	}

	if res, err := withPrefix.List(); err != nil {
		t.Errorf("Failed to list: err='%s'", err.Error())
	} else if len(res) != 1 {
		t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	}
}

func BenchmarkLeveldbStoreApply(b *testing.B) {
	m, err := NewLeveldbStore("test.db")
	if err != nil {
		b.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "key"
	v := []byte("value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := k + string(i)
		m.Apply(key, v)
	}
}

func BenchmarkLeveldbStoreDeleteAfterApply(b *testing.B) {
	m, err := NewLeveldbStore("test.db")
	if err != nil {
		b.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "key"
	v := []byte("value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := k + string(i)
		m.Apply(key, v)
		m.Delete(key)
	}
}

func BenchmarkLeveldbStoreGet(b *testing.B) {
	m, err := NewLeveldbStore("test.db")
	if err != nil {
		b.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

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

func BenchmarkLeveldbStoreList(b *testing.B) {
	m, err := NewLeveldbStore("test.db")
	if err != nil {
		b.Fatalf("failed to generate leveldb datastore: %s", err.Error())
	}
	defer os.RemoveAll(dbDir)

	k := "key"
	v := []byte("value")

	for i := 0; i < 100; i++ {
		key := k + string(i)
		m.Apply(key, v)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.List()
	}
}
