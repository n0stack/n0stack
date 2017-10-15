from .xmllib import build_volume
from .base import BaseOpen


class Volume(BaseOpen):
    def __init__(self):
        super().__init__()

    def create(self, name, size):
        # type: (str, str) -> bool
        xml = build_volume(name, size)

        if self.pool.createXML(xml) is None:
            # TODO: error log
            return False

        return True

    def delete(self, name, wipe=True):
        # type: (str) -> bool
        storage = self.pool.storageVolLookupByName(name)

        if storage is None:
            # TODO: error log
            return False

        if wipe:
            storage.wipe(0)
        storage.delete(0)

        return True
