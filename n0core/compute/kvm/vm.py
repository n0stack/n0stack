import time
import libvirt
import enum
import xml.etree.ElementTree as ET
from typing import Any  # NOQA

from xmllib import xml_generate
from base import QemuOpen

class VM(QemuOpen):  # NOQA
    """
    manage vm status
    """
    def __init__(self):
        # type: () -> None
        super().__init__()

    def start(self, name):
        # type: (str) -> bool
        domain = self.conn.lookupByName(name)
        try:
            domain.create()
        except:
            return False

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
                return False

        return True

    def create(self,
               name,  # type: str
               cpu,  # type: str
               memory,  # type: str
               disk_path,  # type: str
               cdrom,  # type: str
               device,
               mac_addr,  # type: str
               vnc_password,  # type: str
               nic_type
               ):
        # type: (...) -> bool
        
        # default values of nic
        nic = {'type': 'bridge', 'source': device, 'mac_addr': mac_addr, 'model': nic_type}

        vm_xml = xml_generate(name,
                              cpu,
                              memory,
                              disk_path,
                              cdrom,
                              nic,
                              vnc_password)

        dom = self.conn.createXML(vm_xml, 0)

        if not dom:
            return False
        else:
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
            print(e)
            return False

        return True
