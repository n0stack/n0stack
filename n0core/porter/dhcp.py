from ipaddress import IPv4Interface # noqa
import os
from typing import Tuple # noqa

from pyroute2 import IPRoute
from pyroute2 import NetNS
from pyroute2 import NSPopen


class DHCP(object):
    """
    Manage namespaces, veth pairs and dnsmasq processes.
    """
    ip = IPRoute()

    @classmethod
    def _get_veth_names(cls, subnet_id):
        # type: (str) -> Tuple[str, str]
        """
        Get names of veths linked to DHCP server on specified subnet.

        Args:
            subnet_id: Subnet ID.

        Returns:
            Name of one of the veth pair.
            Name of the other.
        """
        return 'tap-dhcp-' + subnet_id, 'eth-dhcp-' + subnet_id

    @classmethod
    def _get_netns_name(cls, subnet_id):
        # type: (str) -> str
        """
        Gets netns name by subnet.

        Args:
            subnet_id: Uuid of subnet.

        Returns:
            netns name.
        """
        return 'dhcp-' + subnet_id

    @classmethod
    def _get_pid_filename(cls, netns_name):
        # type: (str) -> Tuple[str, str]
        """
        Get dnsmasq pid filename by netns.

        Args:
            netns_name: netns name.

        Returns:
            Path to directory of pid file.
            Path to dnsmasq pid file.
        """
        dirname = os.path.join('/var/run/', netns_name)
        return dirname, os.path.join(dirname, 'dnsmasq.pid')

    @classmethod
    def _start_dnsmasq_process(cls, netns_name, interface_name, pool):
        # type: (str, str, Tuple[str, str]) -> None
        """
        Start dnsmasq process on netns.

        1. Create directory where to save pid file.
        2. Start dnsmasq process.

        Args:
            netns_name: netns name.
            interface_name: Name of interface used by dnsmasq.
            pool: DHCP allocation pool. Allocate pool[0]-pool[1].
        """
        dirname, pid_filename = cls._get_pid_filename(netns_name)
        if not os.path.exists(dirname):
            os.mkdir(dirname)

        interface = '--interface=' + interface_name
        dhcp_range = '--dhcp-range=' + pool[0] + ',' + pool[1] + ',' + '12h'
        pid_file = '--pid-file=' + pid_filename
        cmd = ['/usr/sbin/dnsmasq',
               '--no-resolv',
               '--no-hosts',
               interface,
               dhcp_range,
               pid_file]
        NSPopen(netns_name, cmd)

    @classmethod
    def create_dhcp_server(cls, subnet_id, interface_addr, bridge_name, pool):
        # type: (str, IPv4Interface, str, Tuple[str, str]) -> None
        """
        Create DHCP server on specified subnet.

        1. Create netns if not exists.
           in command: `ip netns add $netns_name`
        2. Create veth pair.
           in command: `ip link add $tap_name type veth peer name $peer_name`
        3. Link one of the veth pair to bridge.
           in command: `ip link set dev $tap_name master $bridge_name`
        4. Move the other veth to netns.
           in command: `ip link set $peer_name netns $netns_name`
        5. Add ip address to the veth.
           in command: `ip netns exec $netns_name \
                        ip addr add $address/$prefixlen dev $peer`
        6. Set up veths.
           in command: `ip link set $name up`
        7. Start dnsmasq process.

        Args:
            subnet_id: Subnet id.
            interface_addr: IP address of DHCP server.
            bridge_name: Name of bridge linked to DHCP server.
            pool: DHCP allocation pool. Allocate pool[0]-pool[1].
        """
        netns_name = cls._get_netns_name(subnet_id)
        netns = NetNS(netns_name, flags=os.O_CREAT)

        tap_name, peer_name = cls._get_veth_names(subnet_id)
        cls.ip.link('add', ifname=tap_name, peer=peer_name, kind='veth')

        tap = cls.ip.link_lookup(ifname=tap_name)[0]
        bri = cls.ip.link_lookup(ifname=bridge_name)[0]
        cls.ip.link('set', index=tap, master=bri)

        peer = cls.ip.link_lookup(ifname=peer_name)[0]
        cls.ip.link('set', index=peer, net_ns_fd=netns_name)
        address = str(interface_addr.ip)
        prefixlen = int(interface_addr.network.prefixlen)
        netns.addr('add', index=peer, address=address, prefixlen=prefixlen)

        cls.ip.link('set', index=tap, state='up')
        netns.link('set', index=peer, state='up')
        netns.close()

        cls._start_dnsmasq_process(netns_name, peer_name, pool)

    @classmethod
    def delete_dhcp_server(cls, subnet_id):
        # type : (str) -> None
        """
        Delete DHCP server on specified subnet.

        1. Kill dnsmasq process.
        2. Delete veth pairs.
           in command: `ip link del $tap_name`
        2. Delete related netns.
           in command: `ip netns del $netns_name`

        Args:
            subnet_id: Subnet id.
        """
        netns_name = cls._get_netns_name(subnet_id)

        dirname, pid_filename = cls._get_pid_filename(netns_name)
        with open(pid_filename, 'r') as f:
            pid = int(f.read())
        os.kill(pid, 9)
        os.remove(pid_filename)
        os.rmdir(dirname)

        tap_name, _ = cls._get_veth_names(subnet_id)
        tap = cls.ip.link_lookup(ifname=tap_name)[0]
        cls.ip.link('del', index=tap)

        netns = NetNS(netns_name)
        netns.close()
        netns.remove()
