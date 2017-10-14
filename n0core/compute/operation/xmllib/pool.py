# coding: UTF-8
from xml.etree.ElementTree import Element
import xml.etree.ElementTree as ET


class PoolGen:
    """
    Create pool and storage
    """
    def __init__(self):
        pass

    def __call__(self, pool_name, pool_path):
        pool = Element('pool', attrib={'type': 'dir'})
        name = Element('name')
        name.text = pool_name

        target = Element('target')
        path = Element('path')
        path.text = pool_path

        target.append(path)

        pool.append(name)
        pool.append(target)

        self.xml = ET.tostring(pool).decode('utf-8').replace('\n', '')
