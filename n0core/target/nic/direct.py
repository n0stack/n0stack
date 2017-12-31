from typing import Tuple, Dict, Optional, Union, cast  # NOQA

from n0library.logger import Logger
from n0core.target import Target
from n0core.model import Model  # NOQA
from n0core.model.resource.network import Network
from n0core.model.resource.nic import NIC
from n0core.lib.dhcp.dnsmasq import Dnsmasq


logger = Logger(__name__)


class Direct(Target):
    TYPE = "direct"

    def __init__(self, interface):
        # type: (str) -> None
        super().__init__()

    def apply(self, model):
        # type: (Model) -> Tuple[Model, bool, str]
        model = cast(NIC, model)
        network = cast(Network, model.depend_on("n0stack/n0core/resource/nic/network")[0].model)

        if network.subnets[0].dhcp is None:
            return model, True, "Succeeded"

        dnsmasq = Dnsmasq(network.id, network.bridge)

        if model.state == NIC.STATES.ATTACHED:
            dnsmasq.add_allowed_host(model.hw_addr)

            for i in model.ip_addrs:
                dnsmasq.add_host_entry(model.hw_addr, i)

        elif model.state == NIC.STATES.DELETED:
            dnsmasq.delete_host_entry(model.hw_addr)
            dnsmasq.delete_allowed_host(model.hw_addr)

        return model, True, "Succeeded"
