from xml.etree.ElementTree import Element
import xml.etree.ElementTree as ET


def build_pool(name, path):
    # type: (str, str) -> str
    el_pool = Element('pool', attrib={'type': 'dir'})
    el_name = Element('name')
    el_name.text = name

    el_target = Element('target')
    el_path = Element('path')
    el_path.text = path

    el_target.append(el_path)

    el_pool.append(el_name)
    el_pool.append(el_target)

    return ET.tostring(el_pool).decode('utf-8')


def build_volume(name, size):
    # type: (str, str) -> str
    el_volume = Element('volume')
    el_name = Element('name')
    el_name.text = name + ".img"

    el_capacity = Element('capacity', attrib={'unit': 'M'})
    el_capacity.text = str(size)

    el_volume.append(el_name)
    el_volume.append(el_capacity)

    return ET.tostring(el_volume).decode('utf-8')
