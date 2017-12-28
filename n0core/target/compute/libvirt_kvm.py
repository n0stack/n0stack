import time
import libvirt
from typing import Any  # NOQA

from n0core.target.compute.xml_generator import (define_vm_xml,
                                                 define_volume_xml,
                                                 define_interface_xml,
                                                 build_volume,
                                                 build_network)
from n0core.target.compute.base import QemuOpen
from n0core.target import Target
from n0core.model.resource.vm import VM as VM_MODEL


class VM(QemuOpen, Target):  # NOQA
    """
    manage vm status
    """
    def __init__(self):
        # type: () -> None
        super().__init__()

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        if model.state is VM_MODEL.STATES.RUNNING:
            if not self.start(model.name):
                # TODO: error process
                return model, False, "faild"

            return model, True, "succeeded"

        elif model.state is VM_MODEL.STATES.POWEROFF:
            if not self.force_stop(model.name):
                # TODO: error process
                return model, False, "faild"

            return model, True, "succeeded"

        elif model.state is VM_MODEL.STATES.DELETED:
            if not self.delete(model.name):
                # TODO: error process
                return model, False, "faild"

            return model, True, "succeeded"

    def start(self, name):
        # type: (str) -> bool
        domain = self.conn.lookupByName(name)
        domain.create()

        # fail if over 120 seconds
        s = time.time()
        while True:
            if domain.info()[0] == 1:
                break
            if time.time() - s > 120:
                return False

        return True

    def stop(self, name):
        # type: (str) -> bool
        domain = self.conn.lookupByName(name)
        domain.shutdown()

        # fail if over 120 seconds
        s = time.time()
        while True:
            if domain.info()[0] != 1:
                break
            if time.time() - s > 120:
                return False

        return True

    def force_stop(self, name):
        # type: (str) -> bool
        domain = self.conn.lookupByName(name)
        domain.destroy()

        # fail if over 60 seconds
        s = time.time()
        while True:
            if domain.info()[0] != 1:
                break
            if time.time() - s > 60:
                # TODO: add error log
                return False

        return True

    def create(self,
               name,  # type: str
               cpu,  # type: Any
               memory,  # type: str
               disk_path,  # type: str
               cdrom,  # type: str
               device,  # type: Any
               mac_addr,  # type: str
               vnc_password,  # type: str
               nic_type  # type: Any
               ):
        # type: (...) -> bool

        # default values of nic
        nic = {'type': 'bridge', 'source': device, 'mac_addr': mac_addr, 'model': nic_type}

        vm_xml = define_vm_xml(name,
                               cpu,
                               memory,
                               disk_path,
                               cdrom,
                               nic,
                               vnc_password)

        dom = self.conn.createXML(vm_xml, 0)

        if not dom:
            # TODO: error log
            return False

        return True

    def delete(self, name):
        # type: (str) -> bool
        try:
            vdom = self.conn.lookupByName(name)
            if vdom.isActive():
                vdom.shutdown()

            if vdom.isActive():
                vdom.destroy()
            else:
                vdom.undefine()

        except libvirt.libvirtError as e:
            # TODO: error log
            print(e)
            return False

        return True

    def update(self, name, vcpus, memory):
        # type: (str, int, int) -> bool
        try:
            dom = self.conn.lookupByName(name)
            dom.setMemory(vcpus)
            dom.setMemory(memory)

        except libvirt.libvirtError as e:
            # TODO: error log
            print(e)
            return False

        return True

    def attach_volume(self, name, volume_id):
        # type: (str, str) -> bool
        vm = self.conn.lookupByName(name)

        xml = define_volume_xml("/var/lib/n0stack/" + volume_id)
        vm.attachDevice(xml)

        return True

    def attach_nic(self, name, network_name, mac_address):
        # type: (str, str) -> bool

        # TODO: mac_address
        vm = self.conn.lookupByName(name)

        xml = define_interface_xml(network_name, mac_address)
        vm.attachDevice(xml)

        return True


class Volume(QemuOpen, Target):
    """
    Generate xml of volume
    """
    def __init__(self):
        # type: () -> None
        super().__init__()

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        return model, True, "succeeded"


    def create(self, name, size):
        # type: (str, str) -> bool
        xml = build_volume(name, size)

        if self.pool.createXML(xml) is None:
            # TODO: error log
            return False

        return True

    def delete(self, name, wipe=True):
        # type: (str, bool) -> bool
        storage = self.pool.storageVolLookupByName(name+'.img')

        if storage is None:
            # TODO: error log
            return False

        if wipe:
            storage.wipe(0)
        storage.delete(0)

        return True


class Network(QemuOpen):
    """
    Generate xml of network
    """
    def __init__(self):
        # type: () -> None
        super().__init__()

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        return model, True, "succeeded"


    def create(self,
               network_name,  # type: str
               bridge_name,  # type: str
               address,  # type: str
               netmask,  # type: str
               range_start,  # type: str
               range_end  # type: str
               ):
        # type: (...) -> bool
        xml = build_network(network_name,
                            bridge_name,
                            address,
                            netmask,
                            range_start,
                            range_end)

        network = self.conn.networkCreateXML(xml)

        if network is None:
            # TODO: error log
            print('Failed to create a virtual network')
            return False

        return True
