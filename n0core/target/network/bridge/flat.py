from typing import Dict, Optional  # NOQA

from n0library.logger import Logger
from n0core.target.network.bridge import Bridge


logger = Logger(__name__)


class Flat(Bridge):
    """Flat type porter service.

    VM directly connect the bridge mastering the native interface.
    When you choise the native interface,
    this class create a bridge named "br-flat" automatically.

    Args:
        interface_name: Set interface to create "br-flat" automatically.
        bridge_name: Set bridge_name like "br-flat".

    Notes:
        - When setted interface_name and bridge_name,
          priority of bridge_name is higher than interface_name.
    """

    TYPE = "flat"

    def __init__(self, interface):
        super().__init__(self.TYPE, interface)

    def apply_bridge(self, id, state="up", parameters={}):
        # type: (...) -> str
        """Create bridge mastering interface selected on args

        1. Create a bridge, when the bridge does not exist.
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
        ii = self._get_index(self.interface)
        if not ii:
            logger.error("Failed to get interface index of {}".format(self.interface))

        bn = self.get_bridge_name(id)
        bi = self._get_index(bn)

        if not bi:
            self.ip.link('add', ifname=bn, kind='bridge')
            bi = self._get_index(bn)

            if bi:
                logger.info("Bridge {}({}) is created, mastering {}({})".format(bn, bi, self.interface, ii))
            else:
                logger.critical("Failed to get created interface's index of {}, just after create interface.".format(bn))  # NOQA
                raise Exception()

        self.ip.link("set", index=ii, master=bi)
        self.ip.link('set', index=bi, state=state)
        self.ip.link('set', index=ii, state='up')

        return bn

    def delete_bridge(self, id):
        bn = self.get_bridge_name(id)
        bi = self._get_index(bn)
        if not bi:
            logger.error("Failed to get interface index of {}, when called delete_bridge.".format(bn))
            return

        self.ip.link('delete', index=bi)
