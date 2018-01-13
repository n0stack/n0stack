package model

import (
	"net"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

const NetworkType = "resource/network"

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
		Subnets []Subnet
	}

	Subnet struct {
		Cidr net.IPNet
		DHCP DHCP
	}

	DHCP struct {
		RangeStart  net.IPAddr
		RangeEnd    net.IPAddr
		Gateway     net.IPAddr
		Nameservers []net.IPAddr
	}
)

func (n Network) ToModel() *Model {
	return &n.Model
}

func NewNetwork(id uuid.UUID, specificType, state, name string, meta map[string]string, dependencies Dependencies, bridge string, subnets []Subnet) *Network {
	return &Network{
		Model: Model{
			ID:           id,
			Type:         filepath.Join(NetworkType, specificType),
			State:        state,
			Name:         name,
			Meta:         meta,
			Dependencies: Dependencies{},
		},
		Bridge:  bridge,
		Subnets: subnets,
	}
}
