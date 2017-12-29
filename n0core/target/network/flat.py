from ipaddress import ip_interface
from typing import Tuple, Dict, Optional, Union, cast  # NOQA

from n0library.logger import Logger
from n0core.target.network import Network as NetworkTarget
from n0core.model import Model  # NOQA
from n0core.model.resource.network import Network
from n0core.model.resource.nic import NIC
from n0core.target.network.dhcp.dnsmasq import Dnsmasq


logger = Logger(__name__)


class Flat(NetworkTarget):
    TYPE = "flat"

    def __init__(self, interface):
        # type: (str) -> None
        super().__init__(self.TYPE, interface)

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        resource_type = model.type.split("/")[1]

        if resource_type == "network":
            model = cast(Network, model)
            dnsmasq = Dnsmasq(model.id, self.bridge_name(model.id))

            if model.state == Network.STATES.UP:
                _ = self.apply_bridge(model.id)
                model.bridge = self.bridge_name(model.id)

                subnet = model.subnets[0]  # first version support only one subnet

                if dnsmasq.get_pid():  # dhcp is running
                    dnsmasq.respawn_process(subnet.dhcp.range)
                else:
                    dhcp_addr = ip_interface(subnet.cidr) + 1
                    dnsmasq.create_dhcp_server(dhcp_addr, subnet.dhcp.range)  # rangeが変わったらip addressも変えないといけないのでは？

            elif model.state == Network.STATES.DOWN:
                self.apply_bridge(model.id, state="down")
                dnsmasq.stop_process()

            elif model.state == Network.STATES.DELETED:
                dnsmasq.delete_dhcp_server()
                self.delete_bridge(model.id)

        elif resource_type == "nic":
            model = cast(NIC, model)
            network = cast(Network, model.depend_on("n0stack/n0core/resource/nic/network")[0].model)
            dnsmasq = Dnsmasq(network.id, network.bridge)

            if model.state == NIC.STATES.ATTACHED:
                dnsmasq.add_allowed_host(model.hw_addr)

                for i in model.ip_addrs:
                    dnsmasq.add_host_entry(model.hw_addr, i)

            elif model.state == NIC.STATES.DELETED:
                dnsmasq.delete_host_entry(model.hw_addr)
                dnsmasq.delete_allowed_host(model.hw_addr)

        return model, True, "Succeeded"

    def apply_bridge(self, id, state="up", parameters={}):
        # type: (str, str, Dict[str, str]) -> bool
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

        return True

    def delete_bridge(self, id):
        # type: (str) -> bool
        bn = self.bridge_name(id)
        bi = self._get_index(bn)
        if not bi:
            logger.error("Failed to get interface index of {}, when called delete_bridge.".format(bn))
            return False

        self.ip.link('delete', index=bi)

        return True
