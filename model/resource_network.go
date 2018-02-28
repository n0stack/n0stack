package model

import (
	"fmt"
	"net"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

const NetworkType = "resource/network"

var NetworkStateMachine = map[string]map[string]bool{
	"UP": map[string]bool{
		"Up":     true,
		"Down":   true,
		"Delete": true,
	},
	"DOWN": map[string]bool{
		"Up":     true,
		"Down":   true,
		"Delete": true,
	},
	"DELETED": map[string]bool{
		"Up":     false,
		"Down":   false,
		"Delete": true,
	},
}

type (
	// Network manage network range resource.
	//
	// Example:
	// 	.. code-block:: yaml
	// 	id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
	// 	type: resource/network/vlan
	// 	name: hogehoge
	// 	state: up
	// 	bridge: br-flat
	// 	subnets:
	// 	  - cidr: 192.168.0.0/24
	// 		dhcp:
	// 		  range: 192.168.0.1-192.168.0.127
	// 		  nameservers:
	// 			- 192.168.0.254
	// 		  gateway: 192.168.0.254
	// 	meta:
	// 	  n0stack/n0core/resource/network/vlan/id: 100
	//
	// States:
	// 	UP: Up network.
	// 	DOWN: Down network.
	// 	DELETED: Delete network.
	//
	// Meta:
	// 	n0stack/n0core/resource/network/vlan/id: VLAN ID on vlan network type.
	// 	n0stack/n0core/resource/network/vxlan/id: VXLAN ID on vxlan network type.
	//
	// Labels:
	//
	// Property:
	//
	// Args:
	// 	id: UUID
	// 	type:
	// 	state:
	// 	name: Name of resource.
	// 	bridge: Bridge which provide service network in Linux.
	// 	subnets: Subnets which manage IP range.
	// 	meta:
	// 	dependencies: List of dependency to
	Network struct {
		Model `yaml:",inline"`

		Bridge  string
		Subnets []*Subnet
	}

	Subnet struct {
		Cidr *net.IPNet `yaml:"-"`
		// DHCP *DHCP
	}

	// DHCP struct {
	// 	RangeStart  net.IP `yaml:"rangeStart"`
	// 	RangeEnd    net.IP `yaml:"rangeEnd"`
	// 	Gateway     net.IP
	// 	Nameservers []net.IP
	// }
)

func (n Network) ToModel() *Model {
	return &n.Model
}

func NewNetwork(id, specificType, state, name string, meta map[string]string, dependencies Dependencies, bridge string, subnets []*Subnet) (*Network, error) {
	i, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	return &Network{
		Model: Model{
			ID:           i,
			Type:         filepath.Join(NetworkType, specificType),
			State:        state,
			Name:         name,
			Meta:         meta,
			Dependencies: Dependencies{},
		},
		Bridge:  bridge,
		Subnets: subnets,
	}, nil
}

func NewSubnet(c string, d *DHCP) (*Subnet, error) {
	_, i, err := net.ParseCIDR(c)
	if err != nil {
		return nil, err
	}

	return &Subnet{
		Cidr: i,
		DHCP: d,
	}, nil
}

func NewDHCP(start, end, gateway string, nameservers []string) (*DHCP, error) {
	s := net.ParseIP(start)
	if s == nil {
		return nil, fmt.Errorf("Failed to parse IP address of range start: got %v", start)
	}

	e := net.ParseIP(end)
	if e == nil {
		return nil, fmt.Errorf("Failed to parse IP address of range end: got %v", end)
	}

	g := net.ParseIP(gateway)
	if g == nil {
		return nil, fmt.Errorf("Failed to parse IP address of gateway: got %v", gateway)
	}

	n := make([]net.IP, len(nameservers))
	for i, ns := range nameservers {
		n[i] = net.ParseIP(ns)
		if g == nil {
			return nil, fmt.Errorf("Failed to parse IP address of nameservers: got %v", ns)
		}
	}

	return &DHCP{
		RangeStart:  s,
		RangeEnd:    e,
		Gateway:     g,
		Nameservers: n,
	}, nil
}

func (d *Subnet) MarshalYAML() (interface{}, error) {
	return &struct {
		Cidr string
		DHCP *DHCP
	}{
		Cidr: d.Cidr.String(),
		DHCP: d.DHCP,
	}, nil
}

func (d *Subnet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Alias Subnet
	a := &struct {
		Cidr string
		DHCP **DHCP
	}{
		DHCP: &d.DHCP,
	}

	err := unmarshal(a)
	if err != nil {
		return err
	}

	_, d.Cidr, err = net.ParseCIDR(a.Cidr)
	if err != nil {
		return err
	}

	return nil
}
