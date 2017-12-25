from ipaddress import ip_interface
from typing import Tuple, Dict, Optional  # NOQA

from n0library.logger import Logger
from n0core.target.network import Network
from n0core.model import Model  # NOQA
from n0core.target.network.dhcp.dnsmasq import Dnsmasq


logger = Logger(__name__)


class Flat(Network):
    TYPE = "flat"

    def __init__(self, interface):
        super().__init__(self.TYPE, interface)

    def apply(self, model):
        # type: (Model) -> Tuple[bool, str]
        resource_type = model.type.split("/")[1]
        if resource_type == "network":
            d = Dnsmasq(model.id, self.bridge_name(model.id))

            if model.state == "up":
                model["bridge"] = self.apply_bridge(model.id)

                s = model["subnets"][0]  # first version support only one subnet

                if d.get_pid():  # dhcp is running
                    d.respawn_process(s["dhcp"]["range"].split("-"))
                else:
                    d.create_dhcp_server(s["cidr"], s["dhcp"]["range"].split("-"))  # rangeが変わったらip addressも変えないといけないのでは？

            elif model.state == "down":
                self.apply_bridge(model.id, state="down")
                d.stop_process()

            elif model.state == "deleted":
                d.delete_dhcp_server()
                self.delete_bridge(model.id)

        elif resource_type == "port":
            nid = model.depend_on("n0core/port/network")[0].model.id
            nb = model.depend_on("n0core/port/network")[0].model.bridge
            d = Dnsmasq(nid, nb)

            if model.state == "up":
                d.add_allowed_host(model["hw_addr"])

                for i in model["ip_addrs"]:
                    d.add_host_entry(model["hw_addr"], i)

            elif model.state == "delete":
                d.delete_host_entry(model["hw_addr"])
                d.delete_allowed_host(model["hw_addr"])

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

        bn = self.bridge_name(id)
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

        return bn

    def delete_bridge(self, id):
        bn = self.bridge_name(id)
        bi = self._get_index(bn)
        if not bi:
            logger.error("Failed to get interface index of {}, when called delete_bridge.".format(bn))
            return

        self.ip.link('delete', index=bi)
