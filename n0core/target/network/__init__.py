from typing import Any, Callable, Tuple, List, Optional, cast  # NOQA
from pyroute2 import IPRoute


from typing import Tuple  # NOQA

from n0core.model import Model  # NOQA
from n0core.target import Target  # NOQA


class Network(Target):
    """
    Example:
        >>> b = Network(external_interface)
        >>> b.apply_beidge(id, state="up")
    """

    BRIDGE_FORMAT = "n{}{}"        # BRIDGE_PREFIX.format(bridge_type, id(like vlan_id))
    META_PREFIX_FORMAT = "n0core/resource/network/{}"  # META_PREFIX.format(bridge_type)

    ip = IPRoute()

    def __init__(self, type, interface):
        # type: (str, str) -> None
        super().__init__()

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

    def apply(self, model):
        # type: (Model) -> Tuple[bool, str]
        """
        Args:
            model: model is Model which you want to apply.

        Return:
            - Return succeeded bool
            - Return result description
        """
        pass

    def apply_bridge(self, state="up", parameters={}):
        # type: (...) -> bool
        raise NotImplementedError

    def delete_bridge(self, id):
        # type: (str) -> bool
        raise NotImplementedError

    def bridge_name(self, id):
        i = id.split("-")[0]
        return self.BRIDGE_FORMAT.format(self.type, i)

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

    def _decode_parameters(self, model):
        # type: (Model) -> Dict[str, str]
        d = {}

        for k, v in filter(lambda k, v: self.meta_prefix in k, model.meta.items()):
            d[k.split("/")[-1]] = v

        return d
