from os.path import join
from enum import Enum
from netaddr import IPSet
from netaddr.ip import IPAddress, IPNetwork, IPRange
from typing import Dict, List, Optional, Tuple, Union, Any  # NOQA

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

    STATES = Enum("STATES", ["UP", "DOWN", "DELETED"])

    def __init__(self,
                 id,              # type: str
                 type,            # type: str
                 state,           # type: Enum
                 name,            # type: str
                 bridge="",       # type: str
                 subnets=[],      # type: List[_Subnet]
                 meta={},         # type: Dict[str, str]
                 dependencies=[]  # type: List[_Dependency]
                 ):
        # type: (...) -> None
        super().__init__(id=id,
                         type=join("resource/network", type),
                         state=state,
                         name=name,
                         meta=meta,
                         dependencies=dependencies)

        self.bridge = bridge
        self.__subnets = subnets

    @property
    def subnets(self):
        # type: () -> List[_Subnet]
        return self.__subnets

    def apply_subnet(self,
                     cidr,         # type: str
                     range,        # type: str
                     nameservers,  # type: List[str]
                     gateway       # type: str
                     ):
        # type: (...) -> None
        for i, s in enumerate(self.subnets):
            if IPSet([cidr]) & IPSet(s.cidr):
                self.subnets.pop(i)

        ip_range = IPRange(range.split("-")[0], range.split("-")[1])
        ns = list(map(lambda n: IPAddress(n), nameservers))
        dhcp = _DHCP(ip_range, ns, IPAddress(gateway))
        subnet = _Subnet(IPNetwork(cidr), dhcp)

        self.__subnets.append(subnet)


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
                 cidr,  # type: IPAddress
                 dhcp,  # type: _DHCP
                 ):
        # type: (...) -> None
        self.__cidr = cidr
        self.__dhcp = dhcp

    @property
    def cidr(self):
        # type: () -> IPAddress
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
                 range,        # type: Tuple[IPAddress, IPAddress]
                 nameservers,  # type: List[IPAddress]
                 gateway       # type: IPAddress
                 ):
        # type: (...) -> None
        self.__range = range
        self.__nameservers = nameservers
        self.__gateway = gateway

    @property
    def range(self):
        # type: () -> Tuple[IPAddress, IPAddress]
        return self.__range

    @property
    def nameservers(self):
        # type: () -> List[IPAddress]
        return self.__nameservers

    @property
    def gateway(self):
        # type: () -> IPAddress
        return self.__gateway
