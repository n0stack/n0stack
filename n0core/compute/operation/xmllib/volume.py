# coding: UTF-8
from xml.etree.ElementTree import Element
import xml.etree.ElementTree as ET


class VolumeGen:
    """
    Create volume
    """
    def __init__(self):
        pass

    def __call__(self, volume_name, size):
        volume = Element('volume')
        name = Element('name')
        name.text = volume_name+".img"

        unit = size[-1]
        if unit not in ['B', 'K', 'M', 'G']:
            unit = 'B'
        else:
            size = size[:-1]

        capacity = Element('capacity', attrib={'unit': unit})
        capacity.text = size

        volume.append(name)
        volume.append(capacity)

        self.xml = ET.tostring(volume).decode('utf-8').replace('\n', '')
