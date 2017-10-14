import libvirt
import sys


class BaseReadOnly:
    QEMU_URL = "qemu:///system"

    def __init__(self):
        try:
            conn = libvirt.openReadOnly(self.QEMU_URL)
        except libvirt.libvirtError as e:
            print(e)
            sys.exit(1)
        self.conn = conn


class BaseOpen:
    QEMU_URL = "qemu:///system"

    def __init__(self):
        try:
            conn = libvirt.open(self.QEMU_URL)
        except libvirt.libvirtError as e:
            print(e)
            sys.exit(1)
        self.conn = conn
