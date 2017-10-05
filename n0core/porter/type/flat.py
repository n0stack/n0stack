from n0core.porter.type import PorterType
from typing import Optional  # NOQA


class Flat(PorterType):
    """Flat type porter service.

    VM directly connect the bridge mastering the native interface.
    When you choise the native interface,
    this create a bridge named "br-flat" automatically.
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
                print("Bridge is already exists, ignore external interface option")  # WARN # NOQA
            self.bridge_index = self.get_interface_index(bridge_name)
        elif interface_name:
            self.bridge_index = self.create_bridge(self.BRIDGE_NAME, interface_name)[0]  # NOQA
        else:
            raise Exception("Need whether bridge name or interface name options") # NOQA
