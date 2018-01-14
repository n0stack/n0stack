package flat

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/n0stack/n0core/model"
	"github.com/vishvananda/netlink"
)

const FlatType = "flat"

type Flat struct {
	InterfaceName string
}

func (f Flat) ManagingType() string {
	return filepath.Join(model.NICType, FlatType)
}

func (f Flat) Apply(m model.AbstractModel) (string, bool) {
	n, ok := m.(*model.Network)
	if !ok {
		return "Failed to cast model to network.", false
	}

	switch n.State {
	case "UP":
		b, err := netlink.LinkByName(f.GetBridgeName(n.ID))
		if err != nil {
			err = nil

			i, err := netlink.LinkByName(f.InterfaceName)
			if err != nil {
				return "Failed to get interface by name", false
			}

			ok := f.createBridge(i)
			if !ok {
				return "Failed to create bridge", false
			}

			b, err = netlink.LinkByName(f.GetBridgeName(n.ID))
			if err != nil {
				return "Failed to get bridge by name", false
			}

			n.Bridge = b.Attrs().Name
		}

	case "DOWN":
		b, err := netlink.LinkByName(f.GetBridgeName(n.ID))
		if err != nil {
			err = nil

			i, err := netlink.LinkByName(f.InterfaceName)
			if err != nil {
				return "Failed to get interface by name", false
			}

			ok := f.createBridge(i)
			if !ok {
				return "Failed to create bridge", false
			}

			b, err = netlink.LinkByName(f.GetBridgeName(n.ID))
			if err != nil {
				return "Failed to get bridge by name", false
			}

			n.Bridge = b.Attrs().Name
		}
		netlink.LinkSetDown(b)

	case "DELETED":
		b, err := netlink.LinkByName(f.GetBridgeName(n.ID))
		if err != nil {
			return "Failed to get bridge by name", false
		}

		netlink.LinkDel(b)
	}

	return "", true
}

func (f Flat) GetBridgeName(id uuid.UUID) string {
	i := strings.Split(id.String(), "-")
	return fmt.Sprintf("nbr%s", i[0])
}

func (f Flat) createBridge(i netlink.Link) bool {
	la := netlink.NewLinkAttrs()
	la.Name = ""
	b := &netlink.Bridge{LinkAttrs: la}

	err := netlink.LinkAdd(b)
	if err != nil {
		return false
	}
	netlink.LinkSetMaster(i, b)

	return true
}
