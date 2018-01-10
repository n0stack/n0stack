package resource

import "github.com/n0stack/n0core/model"

// VM manage memory and CPU resource.
//
// Example:
//
// 	.. code-block:: yaml
//
// 	id: 13bae4ae-67f3-456a-ab05-a217d7cf0861
// 	type: resource/vm/kvm
// 	name: hogehoge
// 	state: running
// 	name:
// 	arch: amd64
// 	vcpus: 2
// 	memory: 4 * 1024 * 1024 * 1024
// 	vnc_password: hogehoge
// 	dependencies:
// 	- model:
// 		id: 0a0615bf-8d26-4e9f-bfbc-bbd0890fcd4f
// 		type: resource/nic
// 		name: port
// 		state: attached
// 		hw_addr: ffffffffffff
// 		ip_addrs:
// 		  - 192.168.0.1
// 		  - fe08::1
// 		dependencies:
// 		  - model:
// 			  id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
// 			  type: resource/network/vlan
// 			  name: hogehoge
// 			  meta:
// 				n0stack/n0core/resource/network/vlan_id: 100
// 			  state: up
// 			  bridge: nvlan0f97b5a3
// 			  subnets:
// 			  - cidr: 192.168.0.0/24
// 				dhcp:
// 				  range: 192.168.0.1-192.168.0.127
// 				  nameservers:
// 					- 192.168.0.254
// 				  gateway: 192.168.0.254
// 		  label: n0stack/n0core/port/network
// 	  label: n0stack/n0core/resource/vm/attachments
// 	- model:
// 		type: resource/volume/local
// 		id: 486274b2-49e4-4bcd-a60d-4f627ce8c041
// 		state: allocated
// 		name: hogehoge
// 		size: 10gb
// 		url: file:///data/hoge
// 	  label: n0stack/n0core/resource/vm/attachments
// 	  property:
// 		n0stack/n0core/resource/vm/boot_priority: 1
//
// States:
// 	POWEROFF: Shutdowned VM.
// 	RUNNING: Running VM.
// 	SAVED: Suspended VM.
// 	DELETED: Deleted VM.
//
// Meta:
//
// Labels:
// 	n0stack/n0core/resource/vm/attachments: Attachemt resource
//
// Properties:
// 	n0stack/n0core/resource/vm/boot_priority: Boot priority of volume.
//
// Args:
// 	id: UUID
// 	type:
// 	state:
// 	name:
// 	arch: CPU architecture.
// 	vcpus: Number of CPU cores.
// 	memory: Size of memory bytes.
// 	vnc_password: VNC Password.
// 	meta:
// 	dependencies: List of dependency to
type VM struct {
	model.Model

	Arch        string // enumにしたい
	VCPUs       uint
	Memory      uint64
	VNCPassword string
}

func (v VM) GetModel() *model.Model {
	return &v.Model
}
