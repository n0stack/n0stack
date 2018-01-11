package resource

import (
	"net"

	"github.com/n0stack/n0core/model"
)

// NIC manage IP address resource.
//
// Example:
//
// 	.. code-block:: yaml
//
// 	id: 0a0615bf-8d26-4e9f-bfbc-bbd0890fcd4f
// 	type: resource/nic
// 	name: port
// 	state: attached
// 	hw_addr: ffffffffffff
// 	ip_addrs:
// 	- 192.168.0.1
// 	- fe08::1
// 	dependencies:
// 	- model:
// 		id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
// 		type: resource/network/vlan
// 		state: up
// 		name: hogehoge
// 		meta:
// 		  n0stack/n0stack/resource/network/vlan_id: 100
// 		bridge: nvlan0f97b5a3
// 		subnets:
// 		  - cidr: 192.168.0.0/24
// 			dhcp:
// 			  range: 192.168.0.1-192.168.0.127
// 			  nameservers:
// 				- 192.168.0.254
// 			  gateway: 192.168.0.254
// 		parameters:
// 	  label: n0stack/n0core/resource/nic/network
//
// States:
// 	ATTACHED: Attached NIC.
// 	DELETED: Deleted NIC.
//
// Meta:
//
// Labels:
// 	n0stack/n0core/resource/nic/network: Network to be attached.
//
// Property:
//
// Args:
// 	id: UUID
// 	type:
// 	state:
// 	name: Name of resource.
// 	hw_addr: Hardware address.
// 	ip_addrs: IP addresses to assign.
// 	meta:
// 	dependencies: List of dependency to
type NIC struct {
	model.Model

	HWAddr  net.HardwareAddr
	IPAddrs []net.IP
}

func (n NIC) ToModel() *model.Model {
	return &n.Model
}
