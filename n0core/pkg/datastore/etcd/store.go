package etcd

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"n0st.ac/n0stack/n0core/pkg/datastore"

	"github.com/golang/protobuf/proto"
)

type EtcdDatastore struct {
	client clientv3.KV
	conn   *clientv3.Client
}

const (
	etcdDialTimeout    = 5 * time.Second
	etcdRequestTimeout = 10 * time.Second
)

func NewEtcdDatastore(endpoints []string) (*EtcdDatastore, error) {
	e := &EtcdDatastore{}

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
	}
}

func (d EtcdDatastore) List(ctx context.Context, f func(int) []proto.Message) error {
	c, cancel := context.WithTimeout(ctx, etcdRequestTimeout)
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

func (d EtcdDatastore) Get(ctx context.Context, key string, pb proto.Message) (int64, error) {
	resp, err := d.client.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if resp.Count == 0 {
		pb = nil
		return 0, datastore.NewNotFound(key)
	}

	err = proto.Unmarshal(resp.Kvs[0].Value, pb)
	if err != nil {
		return 0, err
	}

	return resp.Kvs[0].Version, nil
}

func (d EtcdDatastore) Apply(ctx context.Context, key string, pb proto.Message, version int64) (int64, error) {
	s, err := proto.Marshal(pb)
	if err != nil {
		return 0, err
	}

	txn := d.client.Txn(ctx)
	txnRes, err := txn.
		If(clientv3.Compare(clientv3.Version(key), "=", version)).
		Then(clientv3.OpPut(key, string(s))).
		Commit()
	if err != nil {
		return 0, err
	}

	if !txnRes.Succeeded {
		return 0, datastore.ConflictedError{}
	}

	return version + 1, nil
}

func (d EtcdDatastore) Delete(ctx context.Context, key string, version int64) error {
	txn := d.client.Txn(ctx)
	txnRes, err := txn.
		If(clientv3.Compare(clientv3.Version(key), "=", version)).
		Then(clientv3.OpDelete(key)).
		Commit()

	if !txnRes.Succeeded {
		return datastore.ConflictedError{}
	}

	return err
}

func (d EtcdDatastore) Close() error {
	return d.conn.Close()
}
