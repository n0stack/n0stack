from typing import Tuple  # NOQA
from pyroute2 import IPRoute


class PorterType(object):
    """
    PorterType class is network type abstract class.
    For example, flat and vlan.
    """
    ip = IPRoute()

    @classmethod
    def get_interface_index(cls, interface_name):
        # type: (str) -> int
        """Translate interface name to iproute2 index.

        Args:
            interface_name: Linux interface name like "eth0".

        Returns:
            iproute2 index.
            If the interface do not exists, return None.
        """
        ret = cls.ip.link_lookup(ifname=interface_name)
        if ret:
            return ret[0]
        else:
            return None

    @classmethod
    def create_bridge(cls, bridge_name, interface_name):
        # type: (str, str) -> Tuple[int, int]
        """Create bridge mastering interface selected on args

        1. Create a bridge.
           in command: `ip link add dev $bridge_name type bridge`
        2. Set interface master to the bridge.
           in command: `ip link set dev $interface_name master $bridge_name`
        3. Set up the interface and the bridge.
           in command: `ip link set up dev $names`

        When the bridge already exists, ignoring.

        Args:
            bridge_name: Name of creating bridge.
            interface_name: Name of bridge slave interface.

        Returns:
            Tuple
            - Index of bridge interface.
            - Index of interface interface.
        """
        ini = cls.get_interface_index(interface_name)  # type: int # NOQA
        bri = cls.get_interface_index(bridge_name)  # type: int # NOQA
        if bri:
            print("Already exists %s(%d); keep to continue..." % (bridge_name, bri))  # ERROR # NOQA
        else:
            cls.ip.link('add', ifname=bridge_name, kind='bridge')
            bri = cls.get_interface_index(bridge_name)
            print("Bridge %s(%d) is created, mastering %s(%d)" % (bridge_name, bri, interface_name, ini))  # NOQA
        cls.ip.link("set", index=ini, master=bri)
        cls.ip.link('set', index=bri, state='up')
        cls.ip.link('set', index=ini, state='up')

        return (bri, ini)
