package flat

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/n0stack/n0core/model"
	"github.com/vishvananda/netlink"
)

const (
	NetworkType = "flat"
)

type Flat struct {
	InterfaceName string
	interfaceLink netlink.Link
	bridgeLink    netlink.Link
}

// ip link list $interface_name
// ip link list $bridge_name
// ip link add $bridge_name type brdige
// sysctl net.bridge.bridge-nf-call-ip6tables=0
// sysctl net.bridge.bridge-nf-call-iptables=0
// sysctl net.bridge.bridge-nf-call-arptables=0

func (f Flat) ManagingType() string {
	return filepath.Join(model.NetworkType, NetworkType)
}

func (f *Flat) Operations(state, task string) ([]string, error) {
	var err error
	if !model.NetworkStateMachine[state][task] {
		err = fmt.Errorf("Not allowed to operate task '%v' from state '%v'", task, state)
	}

	switch task {
	case "Up":
		return []string{
			"ParseModel",
			"CheckState",
			"CheckInterface",
			"CreateBridge",
			"UpBridge",
		}, err
	case "Down":
		return []string{
			"ParseModel",
			"CheckState",
			"CheckInterface",
			"CreateBridge",
			"DownBridge",
		}, err
	case "Delete":
		return []string{
			"ParseModel",
			"CheckState",
			"DeleteBridge",
		}, err
	}

	return nil, fmt.Errorf("Unsupported task '%v'", task)
}

func (f Flat) getBridgeName(id uuid.UUID) string {
	i := strings.Split(id.String(), "-")
	return fmt.Sprintf("nbr%s", i[0])
}

func (f *Flat) ParseModel(m model.AbstractModel) (bool, string) {
	n, ok := m.(*model.Network)
	if !ok {
		return false, fmt.Sprintf("Failed to parse AbstractModel to *Network model")
	}

	if n.Bridge == "" {
		n.Bridge = f.getBridgeName(n.ID)
	}

	return true, fmt.Sprintf("Succeeded to parse AbstractModel to Network model")
}

// func (f *Flat) CheckState(m model.AbstractModel) (bool, string) {
// 	n := m.(*model.Network)

// 	var err error
// 	if f.bridgeLink, err = netlink.LinkByName(n.Bridge); err != nil {
// 		return false, fmt.Sprintf("Succeeded to create a bridge %v because %v to interface %v is already created", n.Bridge, n.Bridge, f.InterfaceName)
// 	}

// 	f.bridgeLink.Attrs().Flags & unix.IFF_UP

// 	return true, fmt.Sprintf("Succeeded to create a bridge %v because %v to interface %v is already created", n.Bridge, n.Bridge, f.InterfaceName)
// }

// CheckInterface get interface attributes
//
// Commands:
//   ip link list $interface_name
func (f *Flat) CheckInterface(m model.AbstractModel) (bool, string) {
	var err error
	f.interfaceLink, err = netlink.LinkByName(f.InterfaceName)
	if err != nil {
		return false, fmt.Sprintf("Failed to check interface by name '%v', error message '%v'", f.InterfaceName, err.Error())
	}

	return true, fmt.Sprintf("Succeeded to check interface by name '%v'", f.InterfaceName)
}

// Commands:
//   ip link list $bridge_name
//     ip link add $bridge_name type bridge
//     ip link list $bridge_name
func (f *Flat) CreateBridge(m model.AbstractModel) (bool, string) {
	n := m.(*model.Network)

	var e error
	if f.bridgeLink, e = netlink.LinkByName(n.Bridge); e == nil {
		return true, fmt.Sprintf("Succeeded to create a bridge %v because %v to interface %v is already created", n.Bridge, n.Bridge, f.InterfaceName)
	}

	la := netlink.NewLinkAttrs()
	la.Name = n.Bridge
	bl := &netlink.Bridge{LinkAttrs: la}

	if err := netlink.LinkAdd(bl); err != nil {
		return false, fmt.Sprintf("Failed to add a bridge, error message '%v'", err.Error())
	}

	var err error
	if f.bridgeLink, err = netlink.LinkByName(n.Bridge); err != nil {
		return false, fmt.Sprintf("Failed to check a created bridge '%v' to interface '%v': error message '%v'", n.Bridge, f.InterfaceName, err.Error())
	}

	return true, fmt.Sprintf("Succeeded to create a bridge '%v' to interface '%v'", n.Bridge, f.InterfaceName)
}

// SetMasterOfBridge means `ip link set dev $bridge_name master $interface_name`.
func (f *Flat) SetMasterOfBridge(m model.AbstractModel) (bool, string) {
	if err := netlink.LinkSetMaster(f.bridgeLink, bl); err != nil {
		return false, fmt.Sprintf("Failed to set master of the bridge, error message '%v'", err.Error())
	}

	return true, fmt.Sprintf()
}

// Commands:
//   ip link set up dev $bridge_name
func (f *Flat) UpBridge(m model.AbstractModel) (bool, string) {
	n := m.(*model.Network)

	if err := netlink.LinkSetUp(f.bridgeLink); err != nil {
		return false, fmt.Sprintf("Failed to up bridge, error message '%v'", err.Error())
	}

	return true, fmt.Sprintf("Succeeded to `ip set up dev %v`", n.Bridge)
}

// Commands:
//   ip link set down dev $bridge_name
func (f *Flat) DownBridge(m model.AbstractModel) (bool, string) {
	n := m.(*model.Network)

	if err := netlink.LinkSetUp(f.bridgeLink); err != nil {
		return false, fmt.Sprintf("Failed to up bridge, error message '%v'", err.Error())
	}

	return true, fmt.Sprintf("Succeeded to `ip set down dev %v`", n.Bridge)
}

// Commands:
//   ip link list $bridge_name
//     ip link delete $bridge_name
func (f *Flat) DeleteBridge(m model.AbstractModel) (bool, string) {
	n := m.(*model.Network)

	var err error
	if f.bridgeLink, err = netlink.LinkByName(n.Bridge); err != nil {
		return true, fmt.Sprintf("Succeeded to delete because the bridge '%v' is already deleted: error message '%v'", n.Bridge, err.Error())
	}

	if err := netlink.LinkDel(f.bridgeLink); err != nil {
		return false, fmt.Sprintf("Failed to delete bridge, '%v': error message '%v'", n.Bridge, err.Error())
	}

	return true, fmt.Sprintf("Succeeded to delete bridge, '%v'", n.Bridge)
}
