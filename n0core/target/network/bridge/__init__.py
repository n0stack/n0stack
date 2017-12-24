from typing import Any, Callable, Tuple, List, Optional, cast  # NOQA
from pyroute2 import IPRoute


class Bridge:
    """
    Example:
        >>> b = Bridge(external_interface)
        >>> b.apply_beidge(id, state="up")
    """
    BRIDGE_PREFIX_FORMAT = "nbr-{}-{}"        # BRIDGE_PREFIX.format(bridge_type, network_id)
    META_PREFIX_FORMAT = "n0core/resource/network/{}"  # META_PREFIX.format(bridge_type)

    ip = IPRoute()

    def __init__(self, type, interface):
        self.__type = type
        self.__interface = interface

        self.__meta_prefix = self.META_PREFIX_FORMAT.format(type)

    @property
    def type(self):
        return self.__type

    @property
    def interface(self):
        return self.__interface

    @property
    def meta_prefix(self):
        return self.__meta_prefix

    def apply_bridge(self, state="up", parameters={}):
        # type: (...) -> str
        raise NotImplementedError

    def delete_bridge(self):
        raise NotImplementedError

    def get_bridge_name(self, id):
        return self.BRIDGE_PREFIX_FORMAT.format(self.type, id)

    @classmethod
    def _get_index(cls, name):
        # type: (str) -> Optional[int]
        """Translate interface name to iproute2 index.

        Args:
            name: Linux interface name like "eth0".

        Returns:
            iproute2 index.
            None when the interface do not exists.
        """
        ret = cls.ip.link_lookup(ifname=name)  # type: List[int]
        if ret:
            return ret[0]
        else:
            return None
