# coding: UTF-8
import time
import libvirt
import enum
import xml.etree.ElementTree as ET

from kvmconnect.base import BaseOpen

from operation.xmllib.vm import VmGen
from operation.xmllib.volume import VolumeGen

from operation.volume import Create as VolCreate


class Status(BaseOpen):
    """
    manage vm status
    """

    status = enum.Enum('status', 'poweroff running')

    def __init__(self):
        super().__init__()

    def info(self):
        """Return status of vm
        """
        pass

    def start(self, name):
        domain = self.connection.lookupByName(name)
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
        domain = self.connection.lookupByName(name)
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
        domain = self.connection.lookupByName(name)
        domain.destroy()

        # fail if over 60 seconds
        s = time.time()
        while True:
            if domain.info()[0] != 1:
                break
            if time.time() - s > 60:
                return False

        return True


class Create(BaseOpen):
    """
    Create VM

    parameters:
        name: VM (domain) name
        cpu:
            arch: cpu architecture
            nvcpu: number of vcpus
        memory: memory size of VM
        disk:
            pool: pool name where disk is stored
            size: volume size
        cdrom: iso image path
        mac_addr: mac address
        vnc_password: vnc password
    """
    def __init__(self):
        super().__init__()

    def __call__(self, name, cpu, memory, disk, cdrom, mac_addr, vnc_password):
        vmgen = VmGen()

        # create volume (disk)
        volcreate = VolCreate()
        if not volcreate(disk['pool'], name, disk['size']):
            return False

        # default values of nic
        nic = {'type': 'bridge', 'source': 'virbr0', 'mac_addr': mac_addr, 'model': 'virtio'}

        pool = self.connection.storagePoolLookupByName(disk['pool'])
        vol = pool.storageVolLookupByName(name+'.img')

        vmgen(name, cpu, memory, vol.path(), cdrom, nic, vnc_password)

        dom = self.connection.createXML(vmgen.xml, 0)

        if not dom:
            return False
        else:
            return True


class Delete(BaseOpen):
    """
    Delete VM
    """
    def __init__(self):
        super().__init__()

    def __call__(self, name):
        try:
            vdom = self.connection.lookupByName(name)
            if vdom.isActive():  # vm is up
                vdom.shutdown()

            # delete matched volume
            vol = self.volumeLookupByName(name)
            vol.wipe(0)
            vol.delete(0)

            if vdom.isActive():
                vdom.destroy()
            else:
                vdom.undefine()

        except libvirt.libvirtError as e:
            print(e)
            return False

        return True


class Clone(BaseOpen):
    """
    Clone VM

    parameters:
        src: original vm name
        dst: new vm name
    """
    def __init__(self):
        super().__init__()

    def __call__(self, src, dst, vncpass):
        srcdom = self.connection.lookupByName(src)
        # if srcdom.isActive(): # if vm is up
        #     # TODO: save state or something
        #     return False

        # clone volume from src to dst
        srcvol = self.volumeLookupByName(src)
        volgen = VolumeGen()
        dst_cap = srcvol.info()[1]  # get capacity
        volgen(dst, str(dst_cap)+'B')
        pool = srcvol.storagePoolLookupByVolume()
        status = pool.createXMLFrom(volgen.xml, srcvol)  # do clone

        if not status:
            return False

        dstvol = self.volumeLookupByName(dst)

        # clone VM
        # copy XML from src
        root = ET.fromstring(srcdom.XMLDesc())
        # replace name
        el_name = root.find('./name')
        el_name.text = dst
        # remove uuid
        el_uuid = root.find('./uuid')
        root.remove(el_uuid)
        # replace disk
        el_disk = root.find("./devices/disk[@device='disk']/source")
        el_disk.set('file', dstvol.path())
        # remove mac addr
        el_interface = root.find("./devices/interface[@type='bridge']")
        el_mac = el_interface.find('./mac')
        el_interface.remove(el_mac)
        # remove serial and console
        el_devices = root.find('./devices')
        el_devices.remove(el_devices.find("./serial[@type='pty']"))
        el_devices.remove(el_devices.find("./console[@type='pty']"))
        # reset vnc port and vnc password
        el_graphics = root.find("./devices/graphics")
        el_graphics.set('port', '-1')
        el_graphics.set('passwd', vncpass)
        # remove seclabel
        root.remove(root.find('./seclabel'))

        dst_xml = ET.tostring(root).decode()
        dstdom = self.connection.createXML(dst_xml, 0)

        if not dstdom:
            return False
        else:
            return True
