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

    def __init__(self, subnet_id):
        # type: (str) -> None
        """
        Set names in order to create or delete resources.

        Args:
            subnet_id: Subnet ID.
        """
        self.netns_name = 'dhcp-' + subnet_id
        self.tap_name = 'tap-dhcp-' + subnet_id
        self.peer_name = 'eth-dhcp-' + subnet_id
        self.pid_dirname = os.path.join('/var/run/', self.netns_name)
        self.pid_filename = os.path.join(self.pid_dirname, 'dnsmasq.pid')

    def _start_dnsmasq_process(self, pool):
        # type: (Tuple[str, str]) -> None
        """
        Start dnsmasq process on netns.

        1. Create directory where to save pid file.
        2. Start dnsmasq process.

        Args:
            pool: DHCP allocation pool. Allocate pool[0]-pool[1].
        """
        if not os.path.exists(self.pid_dirname):
            os.mkdir(self.pid_dirname)

        interface = '--interface=' + self.peer_name
        dhcp_range = '--dhcp-range=' + pool[0] + ',' + pool[1] + ',' + '12h'
        pid_file = '--pid-file=' + self.pid_filename
        cmd = ['/usr/sbin/dnsmasq',
               '--no-resolv',
               '--no-hosts',
               interface,
               dhcp_range,
               pid_file]
        NSPopen(self.netns_name, cmd)

    def create_dhcp_server(self, interface_addr, bridge_name, pool):
        # (IPv4Interface, str, Tuple[str, str]) -> None
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
            interface_addr: IP address of DHCP server.
            bridge_name: Name of bridge linked to DHCP server.
            pool: DHCP allocation pool. Allocate pool[0]-pool[1].
        """
        netns = NetNS(self.netns_name, flags=os.O_CREAT)

        tap_name = self.tap_name
        peer_name = self.peer_name
        DHCP.ip.link('add', ifname=tap_name, peer=peer_name, kind='veth')

        tap = DHCP.ip.link_lookup(ifname=tap_name)[0]
        bri = DHCP.ip.link_lookup(ifname=bridge_name)[0]
        DHCP.ip.link('set', index=tap, master=bri)

        peer = DHCP.ip.link_lookup(ifname=peer_name)[0]
        DHCP.ip.link('set', index=peer, net_ns_fd=self.netns_name)

        address = str(interface_addr.ip)
        prefixlen = int(interface_addr.network.prefixlen)
        netns.addr('add', index=peer, address=address, prefixlen=prefixlen)

        DHCP.ip.link('set', index=tap, state='up')
        netns.link('set', index=peer, state='up')
        netns.close()

        self._start_dnsmasq_process(pool)

    def delete_dhcp_server(self):
        # type : () -> None
        """
        Delete DHCP server on specified subnet.

        1. Kill dnsmasq process.
        2. Delete veth pairs.
           in command: `ip link del $tap_name`
        2. Delete related netns.
           in command: `ip netns del $netns_name`
        """
        with open(self.pid_filename, 'r') as f:
            pid = int(f.read())
        os.kill(pid, 9)
        os.remove(self.pid_filename)
        os.rmdir(self.pid_dirname)

        tap = DHCP.ip.link_lookup(ifname=self.tap_name)[0]
        DHCP.ip.link('del', index=tap)

        netns = NetNS(self.netns_name)
        netns.close()
        netns.remove()
