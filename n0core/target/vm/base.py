from typing import Any  # NOQA
import os
import sys

import libvirt

from n0library.logger import Logger
from .xml_generator import build_pool


POOL_NAME = 'n0stack'
POOL_PATH = '/var/lib/{}'.format(POOL_NAME)
QEMU_URL = 'qemu:///system'
logger = Logger()


class QemuReadOnly:
    def __init__(self):
        # type: () -> None
        try:
            conn = libvirt.openReadOnly(QEMU_URL)
        except libvirt.libvirtError as e:
            logger.error('unable to connect to libvirt')
            sys.exit(1)

        try:
            pool = conn.storagePoolLookupByName(POOL_NAME)
        except libvirt.libvirtError:
            _init_pool(conn)
            pool = conn.storagePoolLookupByName(POOL_NAME)

        self.conn = conn
        self.pool = pool


class QemuOpen:
    def __init__(self):
        # type: () -> None
        try:
            conn = libvirt.open(QEMU_URL)
        except libvirt.libvirtError as e:
            logger.error('unable to connect to libvirt')
            sys.exit(1)

        try:
            pool = conn.storagePoolLookupByName(POOL_NAME)
        except libvirt.libvirtError:
            _init_pool(conn)
            pool = conn.storagePoolLookupByName(POOL_NAME)
        self.conn = conn
        self.pool = pool


def _init_pool(conn):
    # type: (Any) -> None
    path = POOL_PATH

    if not os.path.exists(path):
        logger.critical('not such pool path {}'.format(path))
        sys.exit(1)

    xml = build_pool(POOL_NAME, path)
    pool = conn.storagePoolDefineXML(xml, 0)
    if pool is None:
        logger.critical('failed to define pool')
        sys.exit(1)
    pool.setAutostart(True)
    pool.create()
