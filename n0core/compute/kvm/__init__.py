import libvirt
import sys


class QemuReadOnly:
    URI = "qemu:///system"

    def __init__(self):
        # type: () -> None
        try:
            self.conn = libvirt.openReadOnly(QemuReadOnly.URI)

        except libvirt.libvirtError as e:
            print(e)
            sys.exit(1)

class QemuOpen:
    URI = "qemu:///system"

    def __init__(self):
        # type: () -> None
        try:
            self.conn = libvirt.open(QemuOpen.URI)

        except libvirt.libvirtError as e:
            print(e)
            sys.exit(1)
        
