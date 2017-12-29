from n0core.model.resource.volume import Volume
from n0core.model.resource.vm import VM
from n0core.model.resource.nic import NIC
from n0core.target.compute.libvirt_kvm import VM as VM_MODEL


nic = NIC("c579329b-6021-2b11-78cd-ca915c56c1f1",
          "virtio",
          NIC.STATES.ATTACHED,
          "virbr0",
          "52:54:00:55:c6:78",
          ["192.168.122.101", "2400:2410:cc40:fc00:7b97:7eaf:8386:859d"])

volume = Volume("a5793fbb-7031-4bd1-b1aa-ca915c56cf45",
                "file",
                Volume.STATES.ALLOCATED,
                "volume_ubuntu-5G-01",
                5*1024*1024*1024,
                "/var/lib/n0stack/ubuntu-5G-01.qcow2")

vm = VM("fc470966-bae8-4afa-af11-3eb7d7d51f42",
        "kvm",
        VM.STATES.POWER_OFF,
        "vm_ubuntu-5G-01",
        "x86_64",
        1,
        1*1024*1024,
        "hogehoge")


vm.add_dependency(nic, "n0stack/n0core/resource/vm/attachments")
vm.add_dependency(volume, "n0stack/n0core/resource/vm/attachments")

c = VM_MODEL()
print(c.apply(vm))
