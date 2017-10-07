from n0core.compute.kvmconnect.base import BaseOpen
from n0core.compute.operation.xmllib.volume import VolumeGen


class Create(BaseOpen):
    def __init__(self):
        # type: (...) -> None
        super().__init__()

    def __call__(self, pool_name, volume_name, size):
        # type: (str, str, str) -> bool
        volgen = VolumeGen()
        volgen(volume_name, size)

        try:
            pool = self.connection.storagePoolLookupByName(pool_name)
            status = pool.createXML(volgen.xml)
        except:
            return False

        if not status:
            return False
        else:
            return True


class Delete(BaseOpen):
    def __init__(self):
        #type: (...) -> None
        super().__init__()

    def __call__(self, pool_name, volume_name):
        # type: (str, str) -> bool
        try:
            pool = self.connection.storagePoolLookupByName(pool_name)
            storage = pool.storageVolLookupByName(volume_name)
            storage.delete()
        except:
            return False

        return True
