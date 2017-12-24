from typing import Dict, Optional  # NOQA

from n0library.logger import Logger
from n0core.porter.type import PorterType
from n0core.porter.dhcp.dnsmasq import Dnsmasq
import n0core.lib.proto  # NOQA


logger = Logger(__name__)


class Flat(PorterType):
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

    BRIDGE_NAME = "br-flat"

    def __init__(self,
                 interface_name=None,   # type: Optional[str]
                 bridge_name=None       # type: Optional[str]
                 ):
        # type: (...) -> None
        super().__init__()

        if bridge_name:
            if interface_name:
                logger.warning("Bridge is already exists, ignoring external interface option")

            self.bridge_index = self.get_interface_index(bridge_name)

        elif interface_name:
            self.bridge_index = self.create_bridge(self.BRIDGE_NAME, interface_name)[0]

        else:
            raise Exception("Need either bridge name or interface name options")

        self.dhcp = {}  # type: Dict[str, Dnsmasq]

    def AttachNetworkInterfaceRequest(self, message):
        # type: (n0core.lib.proto.AttachNetworkInterfaceRequest) -> None
        pass

    def DetachNetworkInterfaceRequest(self, message):
        # type: (n0core.lib.proto.DetachNetworkInterfaceRequest) -> None
        pass

    def UpdateNetworkInterfaceRequest(self, message):
        # type: (n0core.lib.proto.UpdateNetworkInterfaceRequest) -> None
        pass

    def CreateIPv4SubnetRequest(self, message):
        # type: (n0core.lib.proto.CreateIPv4SubnetRequest) -> None
        if message.subnet_id not in self.dhcp:
            self.dhcp[message.subnet_id] = Dnsmasq(message.subnet_id, "")

    def DeleteIPv4SubnetRequest(self, message):
        # type: (n0core.lib.proto.DeleteIPv4SubnetRequest) -> None
        pass

    def UpdateIPv4SubnetRequest(self, message):
        # type: (n0core.lib.proto.UpdateIPv4SubnetRequest) -> None
        pass

    def CreateNetworkRequest(self, message):
        # type: (n0core.lib.proto.CreateNetworkRequest) -> None
        pass

    def DeleteNetworkRequest(self, message):
        # type: (n0core.lib.proto.DeleteNetworkRequest) -> None
        pass

    def UpdateNetworkRequest(self, message):
        # type: (n0core.lib.proto.UpdateNetworkRequest) -> None
        pass
