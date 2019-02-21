package iproute2

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

type Vlan struct {
	name string
	id   int
	i    *Interface

	link netlink.Link
}

func NewVlan(i *Interface, id int) (*Vlan, error) {
	v := &Vlan{
		name: fmt.Sprintf("%s.%d", i.Name(), id),
		id:   id,
		i:    i,
	}

	l, err := netlink.LinkByName(v.name)
	if err != nil {
		if err := v.createVlan(); err != nil {
			return nil, err
		}

		return v, nil
	}

	v.link = l

	return v, nil
}

// ip link add $name type vlan id $id
func (v *Vlan) createVlan() error {
	l := netlink.NewLinkAttrs()
	l.Name = v.name
	l.ParentIndex = v.i.link.Attrs().Index
	v.link = &netlink.Vlan{
		LinkAttrs: l,
		VlanId:    v.id,
	}

	if err := netlink.LinkAdd(v.link); err != nil {
		return fmt.Errorf("Failed to command 'ip link add link %s name %s type vlan id %d': err='%s'", v.i.Name(), v.name, v.id, err.Error())
	}

	if err := v.Up(); err != nil {
		return err
	}

	return nil
}

func (v *Vlan) Name() string {
	return v.name
}

// ip link set dev $name up
func (v *Vlan) Up() error {
	if err := netlink.LinkSetUp(v.link); err != nil {
		return fmt.Errorf("Failed to command 'ip link set dev %s up': err='%s'", v.name, err.Error())
	}

	return nil
}

func (v *Vlan) SetMaster(b *Bridge) error {
	if err := netlink.LinkSetMaster(v.link, b.link); err != nil {
		return fmt.Errorf("Failed to command 'ip link set dev %s master %s': err='%s'", v.name, b.name, err.Error())
	}

	return nil
}

func (v *Vlan) Delete() error {
	if err := netlink.LinkDel(v.link); err != nil {
		return fmt.Errorf("Failed 'ip link del %s type bridge': err='%s'", v.name, err.Error())
	}

	v.link = nil
	return nil
}
