package sqlite

import (
	"os"
	"testing"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/store"

	"github.com/google/go-cmp/cmp"

	"github.com/n0stack/n0stack/n0core/pkg/datastore"
)

const dbFile = "test.db"

func TestSqliteStore(t *testing.T) {
	ds, err := NewSqliteStore("test.db")
	if err != nil {
		t.Fatalf("failed to generate sqlite datastore: %s", err.Error())
	}
	defer os.Remove(dbFile)

	k := "key"
	v := &datastore.Test{Name: "value"}

	if err := ds.Get(k, &datastore.Test{}); err == nil {
		t.Errorf("Get() does not return error, want NotFound")
	} else if !store.IsNotFound(err) {
		t.Errorf("Get() return wrong error, want NotFound: %s", err.Error())
	}

	if err := ds.Apply(k, v); err != nil {
		t.Fatalf("failed to apply data: %s", err.Error())
	}

	got := &datastore.Test{}
	if err := ds.Get(k, got); err != nil {
		t.Errorf("failed to get stored data: %s", err.Error())
	}
	v.XXX_sizecache = 0
	if diff := cmp.Diff(v, got); diff != "" {
		t.Errorf("Get result is wrong: diff=(-want +got)\n%s", diff)
	}

	if err := ds.Delete(k); err != nil {
		t.Errorf("failed to delete data: %s", err.Error())
	}
	if err := ds.Get(k, &datastore.Test{}); err == nil {
		t.Errorf("Get() does not return error, want NotFound")
	} else if !store.IsNotFound(err) {
		t.Errorf("Get() return wrong error, want NotFound: %s", err.Error())
	}

	if err := ds.Close(); err != nil {
		t.Errorf("failed to close db: %s", err.Error())
	}
}
