package tap

import (
	pnotification "github.com/n0stack/go.proto/notification/v0"
	tap "github.com/n0stack/go.proto/tap/v0"
	"github.com/n0stack/n0core/notification"
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
func (f *flat) getInterface(interfaceName string) *pnotification.Notification {
	var err error
	f.interfaceLink, err = netlink.LinkByName(interfaceName)
	if err != nil {
		// Notificationは仕様をまだ決めていないので、後で実装する
		return notification.MakeNotification("getInterface.LinkByName", false, err.Error()) // error文をnotificationを返すようにする
	}

	return notification.MakeNotification("getInterface", true, "") // error文をnotificationを返すようにする
}

// ip link show dev $bridge
func (f *flat) getBridge() *pnotification.Notification {
	m := f.interfaceLink.Attrs().MasterIndex
	if m == 0 {
		return notification.MakeNotification("getBridge", true, "not exists")
	}

	b, err := netlink.LinkByIndex(f.interfaceLink.Attrs().MasterIndex)
	if err != nil {
		return notification.MakeNotification("getBridge.LinkByIndex", false, err.Error()) // error文をnotificationを返すようにする
	}

	switch b.Type() {
	case "bridge":
		f.Spec.Type = tap.Spec_FLAT
	}

	var ok bool
	f.bridgeLink, ok = b.(*netlink.Bridge)
	if !ok {
		return notification.MakeNotification("getBridge.toBridge", false, err.Error()) // error文をnotificationを返すようにする
	}

	f.Status.Bridge = b.Attrs().Name

	return notification.MakeNotification("getBridge", true, "") // error文をnotificationを返すようにする
}

// ip link add name $bridge type bridge
// ip link set dev $interface master $bridge
func (f *flat) createBridge() *pnotification.Notification {
	l := netlink.NewLinkAttrs()
	l.Name = f.getBridgeName()
	f.bridgeLink = &netlink.Bridge{LinkAttrs: l}

	if err := netlink.LinkAdd(f.bridgeLink); err != nil {
		return notification.MakeNotification("createBridge.LinkAdd", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetMaster(f.interfaceLink, f.bridgeLink); err != nil {
		return notification.MakeNotification("createBridge.LinkSetMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	return notification.MakeNotification("createBridge", true, "") // error文をnotificationを返すようにする
}

// ip link set dev $interface promisc on
// ip link set dev $bridge up
//
// 管理上インターフェイスをダウンしている可能性があるので、安全のために自動的にインターフェイスをアップしない
func (f *flat) applyBridge() *pnotification.Notification {
	if err := netlink.SetPromiscOn(f.interfaceLink); err != nil {
		return notification.MakeNotification("applyBridge", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetUp(f.bridgeLink); err != nil {
		return notification.MakeNotification("applyBridge", false, err.Error()) // error文をnotificationを返すようにする
	}

	return notification.MakeNotification("applyBridge", true, "") // error文をnotificationを返すようにする
}

// ip link set dev $interface promisc off
// ip link set dev $interface nomaster
// ip link delete $bridge type bridge
func (f *flat) deleteBridge() *pnotification.Notification {
	if err := netlink.SetPromiscOff(f.interfaceLink); err != nil {
		return notification.MakeNotification("deleteBridge.SetPromiscOff", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetNoMaster(f.interfaceLink); err != nil {
		return notification.MakeNotification("deleteBridge.LinkSetNoMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkDel(f.bridgeLink); err != nil {
		return notification.MakeNotification("deleteBridge.LinkDel", false, err.Error()) // error文をnotificationを返すようにする
	}

	return notification.MakeNotification("deleteBridge", true, "") // error文をnotificationを返すようにする
}

// ip link show dev $tap
func (f *flat) getTap() *pnotification.Notification {
	var err error
	f.tapLink, err = netlink.LinkByName(f.getTapName())
	if err != nil {
		// Notificationは仕様をまだ決めていないので、後で実装する
		return notification.MakeNotification("getTap.LinkByName", true, err.Error()) // error文をnotificationを返すようにする
	}

	f.Status.Tap = f.tapLink.Attrs().Name

	return notification.MakeNotification("getTap", true, "") // error文をnotificationを返すようにする
}

func (f *flat) createTap() *pnotification.Notification {
	l := netlink.NewLinkAttrs()
	l.Name = f.getTapName()
	t := &netlink.Tuntap{
		LinkAttrs: l,
		Mode:      netlink.TUNTAP_MODE_TAP,
	}

	if err := netlink.LinkAdd(t); err != nil {
		return notification.MakeNotification("createTap.LinkAdd", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetMaster(t, f.bridgeLink); err != nil {
		return notification.MakeNotification("createTap.LinkSetMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	return notification.MakeNotification("createTap", true, "")
}

// ip link set dev $tap up
//
// 管理上インターフェイスをダウンしている可能性があるので、安全のために自動的にインターフェイスをアップしない
func (f *flat) applyTap() *pnotification.Notification {
	if err := netlink.LinkSetUp(f.tapLink); err != nil {
		return notification.MakeNotification("applyTap.LinkSetUp", false, err.Error()) // error文をnotificationを返すようにする
	}

	return notification.MakeNotification("applyTap", true, "") // error文をnotificationを返すようにする
}

// ip link delete $bridge type bridge
func (f *flat) deleteTap() *pnotification.Notification {
	if err := netlink.LinkDel(f.tapLink); err != nil {
		return notification.MakeNotification("deleteTap.LinkDel", false, err.Error()) // error文をnotificationを返すようにする
	}

	return notification.MakeNotification("deleteTap", true, "") // error文をnotificationを返すようにする
}
