// +build ignore
// +build medium

package etcd

import (
	"os"
	"testing"
)

func TestApplyAndDelete(t *testing.T) {
	e, err := NewEtcdDatastore("test", []string{os.Getenv("ETCD_ENDPOINT")})
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}

	k := "test"
	// v :=
	// if err := e.Apply(k, v); err != nil {
	// 	t.Errorf("Failed to apply: key='%s', value='%v', err='%s'", k, v, err.Error())
	// }

	n, err := e.Delete(k)
	if err != nil {
		t.Errorf("Failed to delete: key='%s', err='%s'", k, err.Error())
	}
	if n != 1 {
		t.Errorf("Number of deleted keys is mismatch: have='%d', want='%d'", n, 1)
	}

	if err := e.Close(); err != nil {
		t.Errorf("Failed to close: err='%s'", err.Error())
	}
}

func TestGet(t *testing.T) {
	e, err := NewEtcdDatastore("test", []string{os.Getenv("ETCD_ENDPOINT")})
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	defer e.Close()

	k := "test"
	// v :=
	// if err := e.Apply(k, v); err != nil {
	// 	t.Errorf("Failed to apply: key='%s', value='%v', err='%s'", k, v, err.Error())
	// }
	defer e.Delete(k) // TODO: 返り値が複数であるため正しく動作するか要確認

	// if err := e.Get(res); err != nil {
	// 	t.Errorf("Failed to get: key='%s', value='%v', err='%s'", k, v, err.Error())
	// }
	// if assert.DeepEqual(v, res)
}

func TestList(t *testing.T) {
	e, err := NewEtcdDatastore("test", []string{os.Getenv("ETCD_ENDPOINT")})
	if err != nil {
		t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	}
	defer e.Close()

	k := "test"
	// v :=
	// if err := e.Apply(k, v); err != nil {
	// 	t.Errorf("Failed to apply: key='%s', value='%v', err='%s'", k, v, err.Error())
	// }
	defer e.Delete(k) // TODO: 返り値が福通であるため正しく動作するか要確認

	// if err := e.List(); err != nil {
	// 	t.Errorf("Failed to list: key='%s', value='%v', err='%s'", k, v, err.Error())
	// }
	// if len(res) != 1 {
	// 	t.Errorf("Number of listed keys is mismatch: have='%d', want='%d'", len(res), 1)
	// }
	// if assert.DeepEqual(v, res[0])
}

func TestWatch(t *testing.T) {
	// e, err := NewEtcdDatastore("test", []string{os.Getenv("ETCD_ENDPOINT")})
	// if err != nil {
	// 	t.Fatalf("Failed to connect etcd: err='%s'", err.Error())
	// }
	// defer e.Close()

	// watch

	// k := "test"
	// v :=
	// if err := e.Apply(k, v); err != nil {
	// 	t.Errorf("Failed to apply: key='%s', value='%v', err='%s'", k, v, err.Error())
	// }
	// defer e.Delete(k) // TODO: 返り値が福通であるため正しく動作するか要確認
}
