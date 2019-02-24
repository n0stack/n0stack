package etcd

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"

	"github.com/golang/protobuf/proto"
)

type EtcdDatastore struct {
	client clientv3.KV
	conn   *clientv3.Client

	mutex lock.MutexTable
}

const (
	etcdDialTimeout    = 5 * time.Second
	etcdRequestTimeout = 10 * time.Second
)

func NewEtcdDatastore(endpoints []string) (*EtcdDatastore, error) {
	e := &EtcdDatastore{
		mutex: lock.NewMemoryMutexTable(10000),
	}

	var err error
	e.conn, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: etcdDialTimeout,
	})
	if err != nil {
		return nil, err
	}

	e.client = e.conn.KV

	return e, nil
}

func (d *EtcdDatastore) AddPrefix(prefix string) datastore.Datastore {
	return &EtcdDatastore{
		client: namespace.NewKV(d.client, prefix+"/"),
		conn:   d.conn,
		mutex:  lock.NewMemoryMutexTable(10000),
	}
}

func (d EtcdDatastore) List(f func(int) []proto.Message) error {
	c, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	defer cancel()

	resp, err := d.client.Get(c, "/", clientv3.WithFromKey())
	if err != nil {
		return err
	}

	pb := f(len(resp.Kvs))

	for i, ev := range resp.Kvs {
		err = proto.Unmarshal(ev.Value, pb[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (d EtcdDatastore) Get(key string, pb proto.Message) error {
	c, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	defer cancel()

	resp, err := d.client.Get(c, key)
	if err != nil {
		return err
	}
	if resp.Count == 0 {
		pb = nil
		return nil
	}

	err = proto.Unmarshal(resp.Kvs[0].Value, pb)
	if err != nil {
		return err
	}

	return nil
}

func (d EtcdDatastore) Apply(key string, pb proto.Message) error {
	if !d.mutex.IsLocked(key) {
		return errors.New("key is not locked")
	}

	s, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	defer cancel()

	_, err = d.client.Put(c, key, string(s))
	if err != nil {
		return err
	}

	return nil
}

func (d EtcdDatastore) Delete(key string) error {
	if !d.mutex.IsLocked(key) {
		return errors.New("key is not locked")
	}

	c, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	defer cancel()

	_, err := d.client.Delete(c, key)
	if err != nil {
		return err
	}

	return nil
}

func (d EtcdDatastore) Close() error {
	return d.conn.Close()
}

func (m *EtcdDatastore) Lock(key string) bool {
	return m.mutex.Lock(key)
}
func (m *EtcdDatastore) Unlock(key string) bool {
	return m.mutex.Unlock(key)
}
func (m *EtcdDatastore) IsLocked(key string) bool {
	return m.mutex.IsLocked(key)
}
