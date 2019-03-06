package iproute2

import (
	"fmt"
	"sync"
	"time"

	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

type bridgeCounterTable struct {
	data  map[string]int
	mutex lock.MutexTable
}

var bct *bridgeCounterTable
var once sync.Once

func getBridgeCounterTable() *bridgeCounterTable {
	once.Do(func() {
		bct = &bridgeCounterTable{
			data:  map[string]int{},
			mutex: lock.NewMemoryMutexTable(100000),
		}
	})

	return bct
}

type Bridge struct {
	name string
	link *netlink.Bridge // TODO: exportする必要があるか？
}

func AquireBridge(name string) (*Bridge, error) {
	table := getBridgeCounterTable()
	if !lock.WaitUntilLock(table.mutex, name, 5*time.Second, 10*time.Millisecond) {
		return nil, datastore.LockError()
	}
	defer table.mutex.Unlock(name)

	b, err := NewBridge(name)

	if err != nil {
		return nil, err
	}

	table.data[name] += 1
	return b, nil
}

func (b *Bridge) Release() error {
	table := getBridgeCounterTable()
	if !lock.WaitUntilLock(table.mutex, b.name, 5*time.Second, 10*time.Millisecond) {
		return datastore.LockError()
	}
	defer table.mutex.Unlock(b.name)

	table.data[b.name] -= 1
	return nil
}

func (b *Bridge) DeleteIfNoSlave() (bool, error) {
	table := getBridgeCounterTable()
	if !lock.WaitUntilLock(table.mutex, b.name, 5*time.Second, 10*time.Millisecond) {
		return false, datastore.LockError()
	}
	defer table.mutex.Unlock(b.name)
	table.data[b.name] -= 1

	if table.data[b.name] != 0 {
		return false, nil
	}

	links, err := b.ListSlaves()
	if err != nil {
		return false, err
	}

	// TODO: 以下遅い気がする
	i := 0
	for _, l := range links {
		if _, err := NewTap(l); err == nil {
			i++
		}
	}

	if i != 0 {
		return false, nil
	}

	if err := b.Delete(); err == nil {
		return false, err
	}

	return true, nil
}

func NewBridge(name string) (*Bridge, error) {
	b := &Bridge{
		name: name,
	}

	l, err := netlink.LinkByName(name)
	if err != nil {
		if err := b.createBridge(); err != nil {
			return nil, err
		}

		return b, nil
	}

	var ok bool
	b.link, ok = l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("The interface '%s' is not bridge", b.name)
	}

	return b, nil
}

// ip link add $name type bridge
func (b *Bridge) createBridge() error {
	l := netlink.NewLinkAttrs()
	l.Name = b.name
	b.link = &netlink.Bridge{LinkAttrs: l}

	if err := netlink.LinkAdd(b.link); err != nil {
		return fmt.Errorf("Failed 'ip link add name %s type bridge': err='%s'", b.name, err.Error())
	}

	if err := b.Up(); err != nil {
		return err
	}

	return nil
}

func (b Bridge) Name() string {
	return b.name
}

func (b Bridge) GetIPv4() (string, error) {
	a, err := netlink.AddrList(b.link, nl.FAMILY_V4)
	if err != nil {
		return "", fmt.Errorf("Failed 'ip addr show': err='%s'", err.Error())
	}

	if len(a) < 1 {
		return "", fmt.Errorf("Do not exists IP address")
	}

	return a[0].String(), nil
}

func (b Bridge) GetIPv6() (string, error) {
	a, err := netlink.AddrList(b.link, nl.FAMILY_V6)
	if err != nil {
		return "", fmt.Errorf("Failed 'ip addr show': err='%s'", err.Error())
	}

	if len(a) < 1 {
		return "", fmt.Errorf("Do not exists IP address")
	}

	return a[0].String(), nil
}

// ip link set dev $name up
func (b *Bridge) Up() error {
	if err := netlink.LinkSetUp(b.link); err != nil {
		return fmt.Errorf("Failed 'ip link set dev %s up': err='%s'", b.name, err.Error())
	}

	return nil
}

// ip addr replace $addr dev $name
// Golang の net ライブラリはCIDRを一つの型として扱えるものがないので、stringを受け取る
// Example:
// 		192.168.0.1/24
func (b *Bridge) SetAddress(addr string) error {
	a, err := netlink.ParseAddr(addr)
	if err != nil {
		return fmt.Errorf("Failed to parse ip address: addr='%s', err='%s'", addr, err.Error())
	}

	if err := netlink.AddrReplace(b.link, a); err != nil { // TODO: IPv4, IPv6ひとつずつになるか確認する必要がある
		return fmt.Errorf("Failed to add address: addr='%s', err='%s'", a.String(), err.Error())
	}

	return nil
}

// ip link list
func (b *Bridge) ListSlaves() ([]string, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list links")
	}

	slaves := []string{}
	for _, l := range links {
		if l.Attrs().MasterIndex == b.link.Index {
			slaves = append(slaves, l.Attrs().Name)
		}
	}

	return slaves, nil
}

// ip link del name $name
func (b *Bridge) Delete() error {
	if err := netlink.LinkDel(b.link); err != nil {
		return fmt.Errorf("Failed 'ip link del %s type bridge': err='%s'", b.name, err.Error())
	}

	b.link = nil
	return nil
}
