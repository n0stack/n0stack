package etcd

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"

	"github.com/golang/protobuf/proto"
)

type EtcdDatastore struct {
	client *clientv3.Client
}

const (
	etcdDialTimeout    = 5 * time.Second
	etcdRequestTimeout = 10 * time.Second
)

func NewEtcdDatastore(endpoints []string) (*EtcdDatastore, error) {
	e := &EtcdDatastore{}

	var err error
	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: etcdDialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (d *EtcdDatastore) AddPrefix(prefix string) {
	d.client.KV = namespace.NewKV(d.client.KV, prefix+"/")
	d.client.Watcher = namespace.NewWatcher(d.client.Watcher, prefix+"/")
	d.client.Lease = namespace.NewLease(d.client.Lease, prefix+"/")
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

func (d EtcdDatastore) Get(name string, pb proto.Message) error {
	c, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	defer cancel()

	resp, err := d.client.Get(c, name)
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

func (d EtcdDatastore) Apply(name string, pb proto.Message) error {
	s, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	c, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	defer cancel()

	_, err = d.client.Put(c, name, string(s))
	if err != nil {
		return err
	}

	return nil
}

func (d EtcdDatastore) Delete(name string) error {
	c, cancel := context.WithTimeout(context.Background(), etcdRequestTimeout)
	defer cancel()

	_, err := d.client.Delete(c, name, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

func (d EtcdDatastore) Close() error {
	return d.client.Close()
}
