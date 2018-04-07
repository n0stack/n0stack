package flat

import (
	"github.com/n0stack/n0core/lib"
	n0stack "github.com/n0stack/proto"
	tap "github.com/n0stack/proto/device/tap/v0"
	"github.com/n0stack/proto/resource/networkid/v0"
	"github.com/satori/go.uuid"
	"github.com/vishvananda/netlink"
)

type flat struct {
	tap.Spec
	tap.Status

	id            uuid.UUID
	interfaceLink netlink.Link
	bridgeLink    *netlink.Bridge
	tapLink       netlink.Link
}

// ip link show dev $interface
func (f *flat) getInterface(interfaceName string) *n0stack.Notification {
	var err error
	f.interfaceLink, err = netlink.LinkByName(interfaceName)
	if err != nil {
		// Notificationは仕様をまだ決めていないので、後で実装する
		return lib.MakeNotification("getInterface.LinkByName", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("getInterface", true, "") // error文をnotificationを返すようにする
}

// ip link show dev $bridge
func (f *flat) getBridge() *n0stack.Notification {
	m := f.interfaceLink.Attrs().MasterIndex
	if m == 0 {
		return lib.MakeNotification("getBridge", true, "not exists")
	}

	b, err := netlink.LinkByIndex(f.interfaceLink.Attrs().MasterIndex)
	if err != nil {
		return lib.MakeNotification("getBridge.LinkByIndex", false, err.Error()) // error文をnotificationを返すようにする
	}

	switch b.Type() {
	case "bridge":
		f.Spec.NetworkID = &networkid.Spec{
			Type: networkid.Spec_FLAT,
		}
	}

	var ok bool
	f.bridgeLink, ok = b.(*netlink.Bridge)
	if !ok {
		return lib.MakeNotification("getBridge.toBridge", false, err.Error()) // error文をnotificationを返すようにする
	}

	f.Status.Bridge = b.Attrs().Name

	return lib.MakeNotification("getBridge", true, "") // error文をnotificationを返すようにする
}

// ip link add name $bridge type bridge
// ip link set dev $interface master $bridge
func (f *flat) createBridge() *n0stack.Notification {
	l := netlink.NewLinkAttrs()
	l.Name = f.getBridgeName()
	f.bridgeLink = &netlink.Bridge{LinkAttrs: l}

	if err := netlink.LinkAdd(f.bridgeLink); err != nil {
		return lib.MakeNotification("createBridge.LinkAdd", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetMaster(f.interfaceLink, f.bridgeLink); err != nil {
		return lib.MakeNotification("createBridge.LinkSetMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("createBridge", true, "") // error文をnotificationを返すようにする
}

// ip link set dev $interface promisc on
// ip link set dev $bridge up
//
// 管理上インターフェイスをダウンしている可能性があるので、安全のために自動的にインターフェイスをアップしない
func (f *flat) applyBridge() *n0stack.Notification {
	if err := netlink.SetPromiscOn(f.interfaceLink); err != nil {
		return lib.MakeNotification("applyBridge", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetUp(f.bridgeLink); err != nil {
		return lib.MakeNotification("applyBridge", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("applyBridge", true, "") // error文をnotificationを返すようにする
}

// ip link set dev $interface promisc off
// ip link set dev $interface nomaster
// ip link delete $bridge type bridge
func (f *flat) deleteBridge() *n0stack.Notification {
	if err := netlink.SetPromiscOff(f.interfaceLink); err != nil {
		return lib.MakeNotification("deleteBridge.SetPromiscOff", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetNoMaster(f.interfaceLink); err != nil {
		return lib.MakeNotification("deleteBridge.LinkSetNoMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkDel(f.bridgeLink); err != nil {
		return lib.MakeNotification("deleteBridge.LinkDel", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("deleteBridge", true, "") // error文をnotificationを返すようにする
}

// ip link show dev $tap
func (f *flat) getTap() *n0stack.Notification {
	var err error
	f.tapLink, err = netlink.LinkByName(f.getTapName())
	if err != nil {
		// Notificationは仕様をまだ決めていないので、後で実装する
		return lib.MakeNotification("getTap.LinkByName", true, err.Error()) // error文をnotificationを返すようにする
	}

	f.Status.Tap = f.tapLink.Attrs().Name

	return lib.MakeNotification("getTap", true, "") // error文をnotificationを返すようにする
}

func (f *flat) createTap() *n0stack.Notification {
	l := netlink.NewLinkAttrs()
	l.Name = f.getTapName()
	t := &netlink.Tuntap{
		LinkAttrs: l,
		Mode:      netlink.TUNTAP_MODE_TAP,
	}

	if err := netlink.LinkAdd(t); err != nil {
		return lib.MakeNotification("createTap.LinkAdd", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetMaster(t, f.bridgeLink); err != nil {
		return lib.MakeNotification("createTap.LinkSetMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("createTap", true, "")
}

// ip link set dev $tap up
//
// 管理上インターフェイスをダウンしている可能性があるので、安全のために自動的にインターフェイスをアップしない
func (f *flat) applyTap() *n0stack.Notification {
	if err := netlink.LinkSetUp(f.tapLink); err != nil {
		return lib.MakeNotification("applyTap.LinkSetUp", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("applyTap", true, "") // error文をnotificationを返すようにする
}

// ip link delete $bridge type bridge
func (f *flat) deleteTap() *n0stack.Notification {
	if err := netlink.LinkDel(f.tapLink); err != nil {
		return lib.MakeNotification("deleteTap.LinkDel", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("deleteTap", true, "") // error文をnotificationを返すようにする
}
