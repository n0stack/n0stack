package iproute2

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

type Bridge struct {
	name string
	link *netlink.Bridge // TODO: exportする必要があるか？
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

// ip link set dev $name up
func (b *Bridge) Up() error {
	if err := netlink.LinkSetUp(b.link); err != nil {
		return fmt.Errorf("Failed 'ip link set dev %s up': err='%s'", b.name, err.Error())
	}

	return nil
}

// ip addr replace $addr dev $name
func (b *Bridge) SetAddress(addr string) error {
	a, err := netlink.ParseAddr(addr)
	if err != nil {
		return fmt.Errorf("Failed to parse ip address: addr='%s', err='%s'", addr, err.Error())
	}

	if err := netlink.AddrReplace(b.link, a); err != nil {
		return fmt.Errorf("Failed to add address: addr='%s', err='%s'", a.String(), err.Error())
	}

	return nil
}

// ip link del name $name
func (b *Bridge) Delete() error {
	if err := netlink.LinkDel(b.link); err != nil {
		return fmt.Errorf("Failed 'ip link del %s type bridge': err='%s'", b.name, err.Error())
	}

	b.link = nil
	return nil
}
