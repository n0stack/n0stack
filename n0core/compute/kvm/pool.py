import os

from n0core.compute.kvmconnect.base import BaseOpen
from n0core.compute.operation.xmllib.pool import PoolGen


class Create(BaseOpen):
    def __init__(self):
        # type: () -> None
        super().__init__()

    def __call__(self, pool_name, pool_path):
        # type: (str, str) -> bool
        path = os.path.expandvars(pool_path)
        if not os.path.exists(path):
            try:
                os.mkdir(path)
            except PermissionError:
                return False

        poolgen = PoolGen()
        poolgen(pool_name, path)

        status = self.connection.storagePoolDefineXML(poolgen.xml, 0)
        pool = self.connection.storagePoolLookupByName(pool_name)
        pool.setAutostart(True)
        pool.create()

        if not status:
            return False
        else:
            return True


class Delete(BaseOpen):
    def __init__(self):
        # type: () -> None
        super().__init__()

    def __call__(self, name):
        # type: (str) -> bool
        try:
            pool = self.connection.storagePoolLookupByName(name)
            try:
                pool.destroy()
            except:
                pass
            pool.undefine()
        except:
            return False

        return True
