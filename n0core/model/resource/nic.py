from os.path import join
from enum import Enum
from netaddr import EUI, mac_unix
from netaddr.ip import IPAddress
from typing import Dict, List, Union  # NOQA

from n0core.model import Model
from n0core.model import _Dependency # NOQA


class NIC(Model):
    """NIC manage IP address resource.

    Example:
        ```yaml
        id: 0a0615bf-8d26-4e9f-bfbc-bbd0890fcd4f
        type: resource/nic
        name: port
        state: attached
        hw_addr: ffffffffffff
        ip_addrs:
        - 192.168.0.1
        - fe08::1
        dependencies:
        - model:
            id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
            type: resource/network/vlan
            state: up
            name: hogehoge
            meta:
              n0stack/n0stack/resource/network/vlan_id: 100
            bridge: nvlan0f97b5a3
            subnets:
              - cidr: 192.168.0.0/24
                dhcp:
                  range: 192.168.0.1-192.168.0.127
                  nameservers:
                    - 192.168.0.254
                  gateway: 192.168.0.254
            parameters:
          label: n0stack/n0core/resource/nic/network
        ```

    States:
        attached: Attached NIC.
        deleted: Deleted NIC.

    Meta:

    Labels:
        n0stack/n0core/resource/nic/network: Network to be attached.

    Property:

    Args:
        id: UUID
        type:
        state:
        name: Name of resource.
        hw_addr: Hardware address.
        ip_addrs: IP addresses to assign.
        meta:
        dependencies: List of dependency to
    """

    STATES = Enum("STATES", ["ATTACHED", "DELETED"])

    def __init__(self,
                 id,              # type: str
                 type,            # type: str
                 state,           # type: Enum
                 name,            # type: str
                 hw_addr,         # type: str
                 ip_addrs,        # type: List[str]
                 meta={},         # type: Dict[str, str]
                 dependencies=[]  # type: List[_Dependency]
                 ):
        # type: (...) -> None
        super().__init__(id=id,
                         type=join("resource/nic", type),
                         state=state,
                         name=name,
                         meta=meta,
                         dependencies=dependencies)

        self.hw_addr = EUI(hw_addr)  # hw_addrを誰が生成するかを決定したあとに要修正
        self.hw_addr.dialect = mac_unix
        self.__ip_addrs = [map(lambda i: IPAddress(i), ip_addrs)]

        self.meta = meta
        self.dependencies = dependencies

    @property
    def ip_addrs(self):
        # type: () -> List[IPAddress]
        return self.__ip_addrs
