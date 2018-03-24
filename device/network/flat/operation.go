package flat

import (
	"github.com/n0stack/n0core/lib"
	n0stack "github.com/n0stack/proto"
	network "github.com/n0stack/proto/device/network/v0"
	"github.com/n0stack/proto/resource/networkid/v0"
	"github.com/satori/go.uuid"
	"github.com/vishvananda/netlink"
)

type flat struct {
	network.Spec
	network.Status

	id         uuid.UUID
	intLink    netlink.Link
	bridgeLink *netlink.Bridge
}

// func (t tap) getTapName() string {
// 	i := strings.Split(t.id.String(), "-")
// 	return fpmt.Sprintf("ntap%s", i[0])
// }

// ip link show dev $interface
// ip link show dev $bridge
func (f *flat) getBridge(interfaceName string) *n0stack.Notification {
	var err error
	f.intLink, err = netlink.LinkByName(interfaceName)
	if err != nil {
		// Notificationは仕様をまだ決めていないので、後で実装する
		return lib.MakeNotification("getBridge.LinkByName", false, err.Error()) // error文をnotificationを返すようにする
	}

	m := f.intLink.Attrs().MasterIndex
	if m == 0 {
		return lib.MakeNotification("getBridge", true, "not exists")
	}

	b, err := netlink.LinkByIndex(f.intLink.Attrs().MasterIndex)
	if err != nil {
		return lib.MakeNotification("getBridge.LinkByIndex", false, err.Error()) // error文をnotificationを返すようにする
	}

	switch b.Type() {
	case "bridge":
		f.Spec.NetworkId = &networkid.Spec{
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
	l.Name = f.getBridgeName(f.id)
	f.bridgeLink = &netlink.Bridge{LinkAttrs: l}

	if err := netlink.LinkAdd(f.bridgeLink); err != nil {
		return lib.MakeNotification("createBridge.LinkAdd", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetMaster(f.intLink, f.bridgeLink); err != nil {
		return lib.MakeNotification("createBridge.LinkSetMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("createBridge", true, "") // error文をnotificationを返すようにする
}

// ip link set dev $interface promisc on
// ip link set dev $bridge up
//
// 管理上インターフェイスをダウンしている可能性があるので、安全のために自動的にインターフェイスをアップしない
func (f *flat) applyBridge() *n0stack.Notification {
	if err := netlink.SetPromiscOn(f.intLink); err != nil {
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
	if err := netlink.SetPromiscOff(f.intLink); err != nil {
		return lib.MakeNotification("deleteBridge.SetPromiscOff", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkSetNoMaster(f.intLink); err != nil {
		return lib.MakeNotification("deleteBridge.LinkSetNoMaster", false, err.Error()) // error文をnotificationを返すようにする
	}

	if err := netlink.LinkDel(f.bridgeLink); err != nil {
		return lib.MakeNotification("deleteBridge.LinkDel", false, err.Error()) // error文をnotificationを返すようにする
	}

	return lib.MakeNotification("deleteBridge", true, "") // error文をnotificationを返すようにする
}
