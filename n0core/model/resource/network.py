from os.path import join
from enum import Enum
from ipaddress import IPv4Network, IPv6Network  # NOQA
from ipaddress import IPv4Address, IPv6Address  # NOQA
from typing import Dict, List, Optional, Tuple, Union  # NOQA

from n0core.model import Model
from n0core.model import _Dependency # NOQA


class Network(Model):
    """Network manage network range resource.

    Example:
        ```yaml
        id: 0f97b5a3-bff2-4f13-9361-9f9b4fab3d65
        type: resource/network/vlan
        name: hogehoge
        state: up
        bridge: br-flat
        subnets:
          - cidr: 192.168.0.0/24
            dhcp:
              range: 192.168.0.1-192.168.0.127
              nameservers:
                - 192.168.0.254
              gateway: 192.168.0.254
        meta:
          n0stack/n0core/resource/network/vlan/id: 100
        ```

    States:
        up: Up network.
        down: Down network.
        deleted: Delete network.

    Meta:
        n0stack/n0core/resource/network/vlan/id: VLAN ID on vlan network type.
        n0stack/n0core/resource/network/vxlan/id: VXLAN ID on vxlan network type.

    Labels:

    Property:

    Args:
        id: UUID
        type:
        state:
        name: Name of resource.
        bridge: Bridge which provide service network in Linux.
        subnets: Subnets which manage IP range.
        meta:
        dependencies: List of dependency to
    """

    STATES = Enum("STATES", ["ATTACHED", "DELETED"])

    def __init__(self,
                 id,              # type: str
                 type,            # type: str
                 state,           # type: Enum
                 name,            # type: str
                 bridge,          # type: str
                 subnets=[],      # type: List[_Subnet]
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

        self.__id = id
        self.__type = type
        self.state = state

        self.bridge = bridge
        self.__subnets = subnets

        self.meta = meta
        self.dependencies = dependencies

    @property
    def subnets(self):
        # type: () -> List[_Subnet]
        return self.__subnets

    def add_subnet(self,
                   cidr,         # type: Union[IPv4Network, IPv6Network]
                   range,        # type: Union[Tuple[IPv4Address, IPv4Address], Tuple[IPv6Address, IPv6Address]]
                   nameservers,  # type: List[Union[IPv4Address, IPv6Address]]
                   gateway       # type: Union[IPv4Address, IPv6Address]
                   ):
        # type: (...) -> None
        for s in self.subnets:
            if cidr in s.cidr or s.cidr in cidr:
                raise Exception  # 例外を飛ばす already exists

        d = _DHCP(range, nameservers, gateway)
        s = _Subnet(cidr, d)

        self.__subnets.append(s)


class _Subnet:
    """Subnet for inner class

    Example:
        ```yaml
        cidr: 192.168.0.0/24
        dhcp:
          range: 192.168.0.1-192.168.0.127
          nameservers:
            - 192.168.0.254
          gateway: 192.168.0.254
        ```

    Args:
        cidr: Network range to use network.
        dhcp: DHCP options.
    """

    def __init__(self,
                 cidr,  # type: Union[IPv4Network, IPv6Network]
                 dhcp,  # type: _DHCP
                 ):
        # type: (...) -> None
        self.__cidr = cidr
        self.__dhcp = dhcp

    @property
    def cidr(self):
        # type: () -> Union[IPv4Network, IPv6Network]
        return self.__cidr

    @property
    def dhcp(self):
        # type: () -> _DHCP
        return self.__dhcp


class _DHCP:
    """DHCP options.

    Example:
        ```yaml
        range: 192.168.0.1-192.168.0.127
        nameservers:
          - 192.168.0.254
        gateway: 192.168.0.254
        ```

    Args:
        range: Network range to allocate with DHCP.
        nameservers: DNS addresses to publish.
        gateway: Gateway address for default route.
    """

    def __init__(self,
                 range,        # type: Union[Tuple[IPv4Address, IPv4Address], Tuple[IPv6Address, IPv6Address]]
                 nameservers,  # type: List[Union[IPv4Address, IPv6Address]]
                 gateway       # type: Union[IPv4Address, IPv6Address]
                 ):
        # type: (...) -> None
        self.__range = range
        self.__nameservers = nameservers
        self.__gateway = gateway

    @property
    def range(self):
        # type: () -> Union[Tuple[IPv4Address, IPv4Address], Tuple[IPv6Address, IPv6Address]]
        return self.__range

    @property
    def nameservers(self):
        # type: () -> List[Union[IPv4Address, IPv6Address]]
        return self.__nameservers

    @property
    def gateway(self):
        # type: () -> Union[IPv4Address, IPv6Address]
        return self.__gateway
