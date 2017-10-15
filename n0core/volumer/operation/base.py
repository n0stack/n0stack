from .xmllib import build_pool

from n0library.logger import Logger

import libvirt
import sys
import os


POOL_NAME = 'n0stack'
POOL_PATH = '/var/lib'
QEMU_URL = 'qemu:///system'
logger = Logger()


class BaseReadOnly:
    def __init__(self):
        # type: () -> None
        try:
            conn = libvirt.openReadOnly(QEMU_URL)
        except libvirt.libvirtError as e:
            logger.error('unabled to connect to libvirt')
            sys.exit(1)

        pool = conn.storagePoolLookupByName(POOL_NAME)
        if pool is None:
            _init_pool(self)
        self.conn = conn
        self.pool = pool


class BaseOpen:
    def __init__(self):
        # type: () -> None
        try:
            conn = libvirt.open(QEMU_URL)
        except libvirt.libvirtError as e:
            logger.error('unabled to connect to libvirt')
            sys.exit(1)

        pool = conn.storagePoolLookupByName(POOL_NAME)
        if pool is None:
            _init_pool(self)
        self.conn = conn
        self.pool = pool


def _init_pool(cls):
    path = POOL_PATH
    if not os.path.exists(path):
        # TODO: error log
        logger.critical('not such pool path {}'.format(path))
        sys.exit(1)

    xml = build_pool(POOL_NAME, path)
    pool = cls.conn.storagePoolDefineXML(xml, 0)
    if pool is None:
        logger.critical('failed to define pool')
        sys.exit(1)
    pool.setAutostart(True)
    pool.create()
