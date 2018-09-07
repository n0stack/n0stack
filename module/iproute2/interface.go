package iproute2

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

type Interface struct {
	name string
	link netlink.Link // TODO
}

func GetInterface(name string) (*Interface, error) {
	i := &Interface{
		name: name,
	}

	var err error
	i.link, err = netlink.LinkByName(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to find interface: name='%s', err='%s'", i.name, err.Error())
	}

	return i, nil
}

func (i Interface) Up() error {
	if err := netlink.LinkSetUp(i.link); err != nil {
		return fmt.Errorf("Failed 'ip link set dev %s up': err='%s'", i.name, err.Error())
	}

	return nil
}

func (i *Interface) SetMaster(b *Bridge) error {
	if err := netlink.SetPromiscOn(i.link); err != nil {
		return fmt.Errorf("Failed 'ip link set dev %s promisc on': err='%s'", i.name, err.Error())
	}

	if err := netlink.LinkSetMaster(i.link, b.link); err != nil {
		return fmt.Errorf("Failed ip link set dev %s master %s': err='%s'", i.name, b.name, err.Error())
	}

	return nil
}
