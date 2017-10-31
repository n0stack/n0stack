# coding: UTF-8
from xml.etree.ElementTree import Element
import xml.etree.ElementTree as ET
from typing import Dict, Any  # NOQA


def xml_generate(name,  # type: str
                 cpu,  # type: Any
                 memory,  # type: str
                 disk,  # type: str
                 cdrom,  # type: str
                 nic,  # type: Dict[str, Any]
                 vnc_password  # type: str
                 ):
    # type: (...) -> str
    root = Element('domain', attrib={'type': 'kvm'})

    el_name = Element('name')
    el_name.text = name
    root.append(el_name)

    el_memory = Element('memory')
    el_memory.text = memory
    root.append(el_memory)

    el_vcpu = Element('vcpu')
    el_vcpu.text = cpu['nvcpu']
    root.append(el_vcpu)

    # <os>
    # 	<type arch="${arch}">hvm</type>
    # 	<boot dev="cdrom"/>
    # 	<boot dev="hd"/>
    # </os>
    el_os = Element('os')
    el_type = Element('type', attrib={'arch': cpu['arch']})
    el_type.text = "hvm"
    el_boot1 = Element('boot', attrib={'dev': 'cdrom'})
    el_boot2 = Element('boot', attrib={'dev': 'hd'})
    el_os.append(el_type)
    el_os.append(el_boot1)
    el_os.append(el_boot2)
    root.append(el_os)

    # <features>
    # <acpi/>
    # <apic/>
    # </features>
    el_features = Element('features')
    el_acpi = Element('acpi')
    el_apic = Element('apic')
    el_features.append(el_acpi)
    el_features.append(el_apic)
    root.append(el_features)

    # <cpu mode="custom" match="exact">
    #   <model>IvyBridge</model>
    # </cpu>
    el_cpu = Element('cpu', attrib={'mode': 'custom', 'match': 'exact'})
    el_model = Element('model')
    el_cpu.append(el_model)
    root.append(el_cpu)

    # <clock offset="utc">
    #   <timer name="rtc" tickpolicy="catchup"/>
    #   <timer name="pit" tickpolicy="delay"/>
    #   <timer name="hpet" present="no"/>
    # </clock>
    el_clock = Element('clock', attrib={'offset': 'utc'})
    el_timer1 = Element('timer', attrib={'name': 'rtc', 'tickpolicy': 'catchup'})
    el_timer2 = Element('timer', attrib={'name': 'pit', 'tickpolicy': 'delay'})
    el_timer3 = Element('timer', attrib={'name': 'hpet', 'present': 'no'})
    el_clock.append(el_timer1)
    el_clock.append(el_timer2)
    el_clock.append(el_timer3)
    root.append(el_clock)

    # <on_poweroff>destroy</on_poweroff>
    # <on_reboot>restart</on_reboot>
    # <on_crash>restart</on_crash>
    el_on1 = Element('on_poweroff')
    el_on1.text = 'destroy'
    el_on2 = Element('on_reboot')
    el_on2.text = 'restart'
    el_on3 = Element('on_crash')
    el_on3.text = 'restart'
    root.append(el_on1)
    root.append(el_on2)
    root.append(el_on3)

    # <pm>
    #   <suspend-to-mem enabled="no"/>
    #   <suspend-to-disk enabled="no"/>
    # </pm>
    el_pm = Element('pm')
    el_suspend1 = Element('suspend-to-mem', attrib={'enabled': 'no'})
    el_suspend2 = Element('suspend-to-disk', attrib={'enabled': 'no'})
    el_pm.append(el_suspend1)
    el_pm.append(el_suspend2)
    root.append(el_pm)

    # devices
    el_devices = Element('devices')

    # <disk type="file" device="disk">
    #   <driver name="qemu" type="raw"/>
    #   <source file="${disk}"/>
    #   <target dev="vda" bus="virtio"/>
    # </disk>
    el_disk = Element('disk', attrib={'type': 'file', 'device': 'disk'})
    el_driver = Element('driver', attrib={'name': 'qemu', 'type': 'raw'})
    el_source = Element('source', attrib={'file': disk})
    el_target = Element('target', attrib={'dev': 'vda', 'bus': 'virtio'})
    el_disk.append(el_driver)
    el_disk.append(el_source)
    el_disk.append(el_target)
    el_devices.append(el_disk)

    # <disk type="file" device="cdrom">
    #   <driver name="qemu" type="raw"/>
    #   <source file="${cdrom}"/>
    #   <target dev="hda" bus="ide"/>
    #   <readonly/>
    # </disk>
    el_disk = Element('disk', attrib={'type': 'file', 'device': 'cdrom'})
    el_driver = Element('driver', attrib={'name': 'qemu', 'type': 'raw'})
    el_source = Element('source', attrib={'file': cdrom})
    el_target = Element('target', attrib={'dev': 'hda', 'bus': 'ide'})
    el_readonly = Element('readonly')
    el_disk.append(el_driver)
    el_disk.append(el_source)
    el_disk.append(el_target)
    el_disk.append(el_readonly)
    el_devices.append(el_disk)

    # <interface type="${type}">
    #   <source bridge="${source}"/>
    #   <mac address="${mac_addr}"/>
    #   <model type="${model}"/>
    # </interface>
    el_interface = Element('interface', attrib={'type': nic['type']})
    el_source = Element('source', attrib={'bridge': nic['source']})
    el_model = Element('model', attrib={'type': nic['model']})
    el_interface.append(el_source)
    if nic['mac_addr']:
        el_mac = Element('mac', attrib={'address': nic['mac_addr']})
        el_interface.append(el_mac)
        el_interface.append(el_model)
        el_devices.append(el_interface)

    # <input type="mouse" bus="ps2"/>
    el_input = Element('input', attrib={'type': 'mouse', 'bus': 'ps2'})
    el_devices.append(el_input)
    # <graphics type="vnc" port="-1" listen="0.0.0.0" passwd="${vnc_password}"/>
    el_graphics = Element('graphics', attrib={'type': 'vnc', 'port': '-1',
                                              'listen': '0.0.0.0', 'passwd': vnc_password})
    el_devices.append(el_graphics)
    # <console type="pty"/>
    el_console = Element('console', attrib={'type': 'pty'})
    el_devices.append(el_console)

    root.append(el_devices)

    xml = ET.tostring(root).decode('utf-8').replace('\n', '')  # type: str

    return xml
