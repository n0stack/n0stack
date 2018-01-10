package resource

import (
	"net"

	"github.com/n0stack/n0core/model"
)

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
		model.Model

		Bridge  string
		Subnets []subnet
	}

	subnet struct {
		Cidr net.IPNet
		DHCP dhcp
	}

	dhcp struct {
		RangeStart  net.IPAddr
		RangeEnd    net.IPAddr
		Gateway     net.IPAddr
		Nameservers []net.IPAddr
	}
)

func (n Network) GetModel() *model.Model {
	return &n.Model
}
