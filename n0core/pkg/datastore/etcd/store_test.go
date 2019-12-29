// +build medium

package etcd

import (
	"context"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"n0st.ac/n0stack/n0core/pkg/datastore"
)

func getEndpoint() []string {
	endpoint := make([]string, 1)
	if value, ok := os.LookupEnv("ETCD_ENDPOINT"); ok {
		endpoint[0] = value
	} else {
		endpoint[0] = "localhost:2379"
	}

	return endpoint
}

func TestEtcdDatastore(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd; start etcd with docker-compose up: err='%s'", err.Error())
	}
	e := base.AddPrefix("test")

	k := "test"
	v := &datastore.Test{Name: "hoge"}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	version, err := e.Apply(ctx, k, v, 0)
	if err != nil {
		t.Fatalf("Apply(key=%s, value=%v, version=%d) returns err=%+v", k, v, 0, err)
	}

	getRes := &datastore.Test{}
	if got, err := e.Get(ctx, k, getRes); err != nil {
		t.Errorf("Failed to get: key='%s', value='%v', err='%s'", k, v, err.Error())
	} else if got != version {
		t.Errorf("Get(key=%s) returns wrong version: got=%d, wave=%d", k, got, version)
	}
	if getRes.Name != v.Name {
		t.Errorf("Get(key=%s) returns wrong Name: got=%s, want=%s", k, getRes.Name, v.Name)
	}

	listRes := []*datastore.Test{}
	f := func(s int) []proto.Message {
		listRes = make([]*datastore.Test, s)
		for i := range listRes {
			listRes[i] = &datastore.Test{}
		}

		m := make([]proto.Message, s)
		for i, v := range listRes {
			m[i] = v
		}

		return m
	}
	if err := e.List(ctx, f); err != nil {
		t.Errorf("List() returns err=%+v", err)
	}
	if len(listRes) != 1 {
		t.Errorf("List() returns the wrong number of listed keys: got=%d, want=%d", len(listRes), 1)
	}
	if listRes[0].Name != v.Name {
		t.Errorf("List() returns wrong Name: got=%s, want=%s", getRes.Name, v.Name)
	}

	if err := e.Delete(ctx, k, version); err != nil {
		t.Errorf("Delete(key=%s, version=%d) returns err=%+v", k, version, err)
	}

	if err := base.Close(); err != nil {
		t.Errorf("Close() returns err='%s'", err.Error())
	}
}

func TestAboutEmpty(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	defer base.Close()
	e := base.AddPrefix("test")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	k := "test"
	getRes := &datastore.Test{}
	if _, err := e.Get(ctx, k, getRes); err == nil {
		t.Errorf("Get(key=%s) returns no err", k)
	} else if _, ok := err.(datastore.NotFoundError); !ok {
		t.Errorf("Get(key=%s) returns err=%+v", k, err)
	}
	if getRes.Name != "" {
		t.Errorf("Get(key=%s) returns not empty value: have=%s", k, getRes.Name)
	}

	listRes := []*datastore.Test{}
	f := func(s int) []proto.Message {
		listRes = make([]*datastore.Test, s)
		for i := range listRes {
			listRes[i] = &datastore.Test{}
		}

		m := make([]proto.Message, s)
		for i, v := range listRes {
			m[i] = v
		}

		return m
	}
	if err := e.List(ctx, f); err != nil {
		t.Errorf("List() returns err=%s", err)
	}
	if len(listRes) != 0 {
		t.Errorf("List() returns the wrong number of listed keys: got=%d, want=%d", len(listRes), 0)
	}
}

func TestConfliction(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd; start etcd with docker-compose up: err='%s'", err.Error())
	}
	e := base.AddPrefix("test")

	k := "test"
	v := &datastore.Test{Name: "value"}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if _, err := e.Apply(ctx, k, v, 0); err != nil {
		t.Fatalf("Apply('%v', '%v', '%v') err='%s'", k, v, 0, err.Error())
	}
	if _, err := e.Apply(ctx, k, v, 1); err != nil {
		t.Fatalf("Apply('%v', '%v', '%v') err='%s'", k, v, 1, err.Error())
	}
	if _, err := e.Apply(ctx, k, v, 1); err == nil {
		t.Errorf("Apply('%v', '%v', '%v') no error on applying confliction", k, v, 1)
	} else if _, ok := err.(datastore.ConflictedError); !ok {
		t.Errorf("Apply('%v', '%v', '%v') wrong error on applying confliction: err=%+v", k, v, 1, err)
	}

	if err := e.Delete(ctx, k, 1); err == nil {
		t.Errorf("Delete('%v', '%v') no error on applying confliction", k, 1)
	}
	if err := e.Delete(ctx, k, 2); err != nil {
		t.Errorf("Delete('%v', '%v') err='%s'", k, 2, err.Error())
	}
}
