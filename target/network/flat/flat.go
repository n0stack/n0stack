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

func (f Flat) ManagingType() string {
	return filepath.Join(model.NetworkType, NetworkType)
}

func (f *Flat) Operations(state, task string) ([]func(n model.AbstractModel) (string, bool, string), bool) { // ディスカバリにも利用するはずなので実装については後で考える
	if !model.NetworkStateMachine[state][task] {
		return nil, false
	}

	switch task {
	case "Up":
		return []func(n model.AbstractModel) (string, bool, string){
			f.ParseModel,
			f.CheckInterface,
			f.CreateBridge,
			f.UpBridge,
		}, true
	case "Down":
		return []func(n model.AbstractModel) (string, bool, string){
			f.ParseModel,
			f.CheckInterface,
			f.CreateBridge,
			f.DownBridge,
		}, true
	case "Delete":
		return []func(n model.AbstractModel) (string, bool, string){
			f.ParseModel,
			f.DeleteBridge,
		}, true
	}

	return nil, false
}

func (f *Flat) GetBridgeName(id uuid.UUID) string {
	i := strings.Split(id.String(), "-")
	return fmt.Sprintf("nbr%s", i[0])
}

func (f *Flat) ParseModel(m model.AbstractModel) (string, bool, string) {
	const OPERATION = "ParseModel"

	n, ok := m.(*model.Network)
	if !ok {
		return OPERATION, false, fmt.Sprintf("Failed to parse AbstractModel to *Network model")
	}

	if n.Bridge == "" {
		n.Bridge = f.GetBridgeName(n.ID)
	}

	return OPERATION, true, fmt.Sprintf("Succeeded to parse AbstractModel to Network model")
}

func (f *Flat) CheckInterface(m model.AbstractModel) (string, bool, string) {
	const OPERATION = "CheckInterface"

	var err error
	f.interfaceLink, err = netlink.LinkByName(f.InterfaceName)
	if err != nil {
		return OPERATION, false, fmt.Sprintf("Failed to check interface by name '%v', error message '%v'", f.InterfaceName, err.Error())
	}

	return OPERATION, true, fmt.Sprintf("Succeeded to check interface by name '%v'", f.InterfaceName)
}

func (f *Flat) CreateBridge(m model.AbstractModel) (string, bool, string) {
	const OPERATION = "CreateBridge"
	n := m.(*model.Network)

	var e error
	if f.bridgeLink, e = netlink.LinkByName(n.Bridge); e == nil {
		return OPERATION, true, fmt.Sprintf("Succeeded to create a bridge %v because %v to interface %v is already created", n.Bridge, n.Bridge, f.InterfaceName)
	}

	la := netlink.NewLinkAttrs()
	la.Name = n.Bridge
	bl := &netlink.Bridge{LinkAttrs: la}

	if err := netlink.LinkAdd(bl); err != nil {
		return OPERATION, false, fmt.Sprintf("Failed to add a bridge, error message '%v'", err.Error())
	}

	if err := netlink.LinkSetMaster(f.bridgeLink, bl); err != nil {
		return OPERATION, false, fmt.Sprintf("Failed to set master of the bridge, error message '%v'", err.Error())
	}

	var err error
	if f.bridgeLink, err = netlink.LinkByName(n.Bridge); err != nil {
		return OPERATION, false, fmt.Sprintf("Failed to check a created bridge '%v' to interface '%v': error message '%v'", n.Bridge, f.InterfaceName, err.Error())
	}

	return OPERATION, true, fmt.Sprintf("Succeeded to create a bridge '%v' to interface '%v'", n.Bridge, f.InterfaceName)
}

func (f *Flat) UpBridge(m model.AbstractModel) (string, bool, string) {
	const OPERATION = "UpBridge"
	n := m.(*model.Network)

	if err := netlink.LinkSetUp(f.bridgeLink); err != nil {
		return OPERATION, false, fmt.Sprintf("Failed to up bridge, error message '%v'", err.Error())
	}

	return OPERATION, true, fmt.Sprintf("Succeeded to `ip set up dev %v`", n.Bridge)
}

func (f *Flat) DownBridge(m model.AbstractModel) (string, bool, string) {
	const OPERATION = "DownBridge"
	n := m.(*model.Network)

	if err := netlink.LinkSetUp(f.bridgeLink); err != nil {
		return OPERATION, false, fmt.Sprintf("Failed to up bridge, error message '%v'", err.Error())
	}

	return OPERATION, true, fmt.Sprintf("Succeeded to `ip set down dev %v`", n.Bridge)
}

func (f *Flat) DeleteBridge(m model.AbstractModel) (string, bool, string) {
	const OPERATION = "DeleteBridge"
	n := m.(*model.Network)

	var err error
	if f.bridgeLink, err = netlink.LinkByName(n.Bridge); err != nil {
		return OPERATION, true, fmt.Sprintf("Succeeded to delete because the bridge '%v' is already deleted: error message '%v'", n.Bridge, err.Error())
	}

	if err := netlink.LinkDel(f.bridgeLink); err != nil {
		return OPERATION, false, fmt.Sprintf("Failed to delete bridge, '%v': error message '%v'", n.Bridge, err.Error())
	}

	return OPERATION, true, fmt.Sprintf("Succeeded to delete bridge, '%v'", n.Bridge)
}
