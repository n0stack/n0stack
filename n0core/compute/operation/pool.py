# coding: UTF-8
from kvmconnect.base import BaseOpen
from operation.xmllib.pool import PoolGen

import os


class Create(BaseOpen):
    def __init__(self):
        super().__init__()

    def __call__(self, pool_name, pool_path):
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
        super().__init__()

    def __call__(self, name):
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
