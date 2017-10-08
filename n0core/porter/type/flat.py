from typing import Optional  # NOQA

from n0library.logger import Logger
from n0core.porter.type import PorterType


logger = Logger(__name__)


class Flat(PorterType):
    """Flat type porter service.

    VM directly connect the bridge mastering the native interface.
    When you choise the native interface,
    this class create a bridge named "br-flat" automatically.
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
