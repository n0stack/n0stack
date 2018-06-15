package iproute2

import (
	fmt "fmt"

	"github.com/vishvananda/netlink"
)

// ip link list $bridge_name
func (a Iproute2Agent) getBridge(name string) (*netlink.Bridge, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		// return nil, err // エラーの場合はLinkがない場合のみだと考えておく
		return nil, nil
	}

	b, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("Failed to get bridge because link is not bridge, LinkName:%s", l.Attrs().Name)
	}

	return b, nil
}

// ip link add type bridge name $name
// ip link set dev $bridge_name master $uplink promisc on up
func (a Iproute2Agent) createBridge(name string, uplink netlink.Link) error {
	l := netlink.NewLinkAttrs()
	l.Name = name
	b := &netlink.Bridge{LinkAttrs: l}

	if err := netlink.LinkAdd(b); err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(uplink, b); err != nil {
		return err
	}

	if err := netlink.SetPromiscOn(uplink); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(b); err != nil {
		return err
	}

	return nil
}

// ip link list $tap_name
func (a Iproute2Agent) getTap(name string) (netlink.Link, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		// return nil, err // エラーの場合はLinkがない場合のみだと考えておく
		return nil, nil
	}

	// 何故かできない
	// t, ok := l.(*netlink.Tuntap)
	// if !ok {
	// 	return nil, fmt.Errorf("Failed to get tuntap because link is not tuntap, LinkName:%s", l.Attrs().Name)
	// }

	// if t.Mode != netlink.TUNTAP_MODE_TAP {
	// 	return nil, fmt.Errorf("Failed to get tuntap because link is not tap, LinkName:%s, Mode:%v", l.Attrs().Name, t.Mode)
	// }

	if l.Type() != "tun" { // なんでtapではなくtunになっているは謎
		return nil, fmt.Errorf("Failed to get tuntap because link is not tuntap, LinkName:%s", l.Type())
	}

	return l, nil
}

// ip tuntap add name $tap_name mode tap
// ip link set dev $tap_name master $bridge_name
func (a Iproute2Agent) createTap(name string, master *netlink.Bridge) error {
	l := netlink.NewLinkAttrs()
	l.Name = name
	t := &netlink.Tuntap{
		LinkAttrs: l,
		Mode:      netlink.TUNTAP_MODE_TAP,
	}

	if err := netlink.LinkAdd(t); err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(t, master); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(t); err != nil {
		return err
	}

	return nil
}
