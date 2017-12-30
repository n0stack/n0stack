from n0core.model.resource.network import Network
from n0core.model.resource.nic import NIC
from n0core.model.resource.volume import Volume
from n0core.model.resource.vm import VM as VM_MODEL
from n0core.target.network.flat import Flat
from n0core.target.compute.libvirt_kvm import VM


INTERFACE_NAME = "enp0s25"
NETWORK_ID = "00f967bf-0cc6-443e-b2a9-277b158f42f0"


f = Flat(INTERFACE_NAME)
up_network = Network(NETWORK_ID,
                     "flat",
                     Network.STATES.UP,
                     "test_network")
up_network.apply_subnet("172.16.100.0/24",
                        "172.16.100.3-172.16.100.127",
                        ["172.16.100.254"],
                        "172.16.100.254")
test_network, _, _ = f.apply(up_network)

nic = NIC("755cedc3-1fe2-44a8-bc45-93d8edf2ad82",
          "dhcp",
          NIC.STATES.ATTACHED,
          "test_nic",
          "ae:b4:b1:e3:91:8a",
          ["192.168.3.109"])
nic.add_dependency(up_network, "n0stack/n0core/resource/nic/network")
test_nic, _, _ = f.apply(nic)


v = VM()
vm = VM_MODEL("fc470967-bae8-4afa-af11-3eb7d7d51f42",
              "kvm",
              VM_MODEL.STATES.RUNNING,
              "vm_ubuntu-10G-04",
              "x86_64",
              1,
              1*1024*1024,
              "hogehoge")
vm.add_dependency(test_nic, "n0stack/n0core/resource/vm/attachments")

volume = Volume("a5793fbb-7031-4bd1-b1aa-ca915c56cf46",
                "file",
                Volume.STATES.ALLOCATED,
                "volume_ubuntu-10G-04",
                10*1024*1024*1024,
                "file:///var/lib/n0stack/ubuntu-10G-04.qcow2")
vm.add_dependency(volume, "n0stack/n0core/resource/vm/attachments")
test_vm, _, _ = v.apply(vm)
