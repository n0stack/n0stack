// +build medium
// +build !without_external

package etcd

import (
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
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

func TestApplyAndDelete(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	e := base.AddPrefix("test")

	k := "test"
	v := &datastore.Test{Name: "hoge"}

	if !e.Lock(k) {
		t.Errorf("Failed to lock: key=%s", k)
	}

	if err := e.Apply(k, v); err != nil {
		t.Errorf("Failed to apply: key='%s', value='%v', err='%s'", k, v, err.Error())
	}

	if err := e.Delete(k); err != nil {
		t.Errorf("Failed to delete: key='%s', err='%s'", k, err.Error())
	}

	if err := base.Close(); err != nil {
		t.Errorf("Failed to close: err='%s'", err.Error())
	}

	if !e.Unlock(k) {
		t.Errorf("Failed to unlock: key=%s", k)
	}
}

func TestGet(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	defer base.Close()
	e := base.AddPrefix("test")

	k := "test"
	v := &datastore.Test{Name: "hoge"}

	e.Lock(k)
	defer e.Unlock(k)

	if err := e.Apply(k, v); err != nil {
		t.Fatalf("Failed to apply: key='%s', value='%v', err='%s'", k, v, err.Error())
	}
	defer e.Delete(k) // TODO: 返り値が複数であるため正しく動作するか要確認

	res := &datastore.Test{}
	if err := e.Get(k, res); err != nil {
		t.Errorf("Failed to get: key='%s', value='%v', err='%s'", k, v, err.Error())
	}
	if res.Name != v.Name {
		t.Errorf("Get 'Name' is wrong: key='%s', have='%s', want='%s'", k, res.Name, v.Name)
	}
}

func TestList(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	defer base.Close()
	e := base.AddPrefix("test")

	k := "test"
	v := &datastore.Test{Name: "hoge"}

	e.Lock(k)
	defer e.Unlock(k)

	if err := e.Apply(k, v); err != nil {
		t.Fatalf("Failed to apply: key='%s', value='%v', err='%s'", k, v, err.Error())
	}
	defer e.Delete(k) // TODO: 返り値が複数であるため正しく動作するか要確認

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
	if err := e.List(f); err != nil {
		t.Errorf("Failed to list: key='%s', value='%v', err='%s'", k, v, err.Error())
	}
	if len(res) != 1 {
		t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	}
	if res[0].Name != v.Name {
		t.Errorf("Get 'Name' is wrong: key='%s', have='%s', want='%s'", k, res[0].Name, v.Name)
	}
}

func TestEmpty(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	defer base.Close()
	e := base.AddPrefix("test")

	k := "test"
	resGet := &datastore.Test{}
	if err := e.Get(k, resGet); err != nil {
		t.Errorf("Failed to get: key='%s', err='%s'", k, err.Error())
	}
	if resGet.Name != "" {
		t.Errorf("Response is not nil on Get: key='%s', have='%s'", k, resGet.Name)
	}

	resList := []*datastore.Test{}
	f := func(s int) []proto.Message {
		resList = make([]*datastore.Test, s)
		for i := range resList {
			resList[i] = &datastore.Test{}
		}

		m := make([]proto.Message, s)
		for i, v := range resList {
			m[i] = v
		}

		return m
	}
	if err := e.List(f); err != nil {
		t.Errorf("Failed to list: key='%s', err='%s'", k, err.Error())
	}
	if len(resList) != 0 {
		t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(resList), 0)
	}
}

func TestUpdateSystemBeforeLock(t *testing.T) {
	base, err := NewEtcdDatastore(getEndpoint())
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	e := base.AddPrefix("test")
	defer base.Close()

	k := "test"
	v := &datastore.Test{Name: "hoge"}

	if err := e.Apply(k, v); err == nil {
		t.Errorf("Failed to apply: key='%s', value='%v'", k, v)
	}

	if err := e.Delete(k); err == nil {
		t.Errorf("Failed to delete: key='%s'", k)
	}
}
