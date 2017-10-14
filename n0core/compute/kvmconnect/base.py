import libvirt
import sys


class BaseReadOnly:
    def __init__(self):
        # type: (...) -> None
        self.QEMU_URL = "qemu:///system"
        self.connect_to_qemu()

    def connect_to_qemu(self):
        # type: (...) -> None
        self.connection = libvirt.openReadOnly(self.QEMU_URL)
        if self.connection is None:
            print("Failed to connect to the hypervisor.")
            sys.exit(1)


class BaseOpen:
    def __init__(self):
        # type: (...) -> None
        self.QEMU_URL = "qemu:///system"
        self.connect_to_qemu()

    def connect_to_qemu(self):
        # type: (...) -> None
        self.connection = libvirt.open(self.QEMU_URL)
        if self.connection is None:
            print("Failed to connect to the hypervisor.")
            sys.exit(1)

    def volumeLookupByName(self, volname):
        # type: (str) -> bool
        allvols = [vol for pool in self.connection.listAllStoragePools()
                   for vol in pool.listAllVolumes()]
        vols = list(filter(lambda x: volname+'.img' == x.name(), allvols))

        if len(vols) == 0:
            return False

        assert len(vols) == 1, 'multiple volume found'
        return True
