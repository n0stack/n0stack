from typing import Any, Callable, Tuple, List, Optional, cast  # NOQA
from pyroute2 import IPRoute

from n0library.logger import Logger

import n0core.lib.proto
from n0core.porter.exceptions import ReceivedUnsupportedMessage


logger = Logger(__name__)


class PorterType(object):
    """
    PorterType class is network type abstract class.
    For example, flat and vlan.
    """

    ip = IPRoute()

    @staticmethod
    def _default_proto_method(message):
        # type: (Any) -> None
        """Only raise a Exception meaning recieved unsupported message.

        Exceptions:
            ReceivedUnsupportedMessage: Received unsupported message on porter class.
        """
        raise ReceivedUnsupportedMessage(message.__class__.__name__)

    def __getattr__(self, name):
        # type: (str) -> Optional[Callable[[Any], None]]
        """Catch undefined methods handling to messages.

        Returns:
            The return value of _default_proto_method, when getting not defined proto message.
            None when the message undefined on protobuf.
        """
        if hasattr(n0core.lib.proto, name):
            return self._default_proto_method
        raise AttributeError

    @classmethod
    def get_interface_index(cls, interface_name):
        # type: (str) -> Optional[int]
        """Translate interface name to iproute2 index.

        Args:
            interface_name: Linux interface name like "eth0".

        Returns:
            iproute2 index.
            None when the interface do not exists.
        """
        ret = cls.ip.link_lookup(ifname=interface_name)  # type: List[int]
        if ret:
            return ret[0]
        else:
            return None

    @classmethod
    def create_bridge(cls, bridge_name, interface_name):
        # type: (str, str) -> Tuple[int, int]
        """Create bridge mastering interface selected on args

        1. Create a bridge.
           When the bridge already exists, ignoring.
           in command: `ip link add dev $bridge_name type bridge`
        2. Set interface master to the bridge.
           in command: `ip link set dev $interface_name master $bridge_name`
        3. Set up the interface and the bridge.
           in command: `ip link set up dev $names`

        Args:
            bridge_name: Name of creating bridge.
            interface_name: Name of bridge slave interface.

        Returns:
            Tuple
            - Index of bridge interface.
            - Index of interface interface.
        """
        ini = cls.get_interface_index(interface_name)
        if not ini:
            logger.error("Failed to get interface index of {}".format(interface_name))

        bri = cls.get_interface_index(bridge_name)
        if bri:
            logger.error("Already exists {}({}); keep to continue...".format(bridge_name, bri))
        else:
            cls.ip.link('add', ifname=bridge_name, kind='bridge')
            bri = cls.get_interface_index(bridge_name)
            if bri:
                logger.info("Bridge {}({}) is created, mastering {}({})".format(bridge_name, bri, interface_name, ini))
            else:
                logger.critical("Failed to get created interface's index of {}, just after create interface.".format(bridge_name))  # NOQA
        cls.ip.link("set", index=ini, master=bri)
        cls.ip.link('set', index=bri, state='up')
        cls.ip.link('set', index=ini, state='up')

        return (cast(int, bri), cast(int, ini))
