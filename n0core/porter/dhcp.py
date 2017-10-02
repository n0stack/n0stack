from ipaddress import IPv4Interface # NOQA
import os
from shutil import rmtree
from typing import Tuple # NOQA
from warnings import warn

from pyroute2 import IPRoute
from pyroute2 import NetNS
from pyroute2 import NSPopen


class Dnsmasq(object):
    """
    Manage namespace, veth pair, directory and dnsmasq process.
    """
    ip = IPRoute()

    def __init__(self, subnet_id):
        # type: (str) -> None
        """
        Set names in order to create or delete resources.

        Args:
            subnet_id: Subnet ID.
        """
        self.netns_name = 'dhcp-{}'.format(subnet_id)
        self.tap_name = 'tap-dhcp-{}'.format(subnet_id)
        self.peer_name = 'eth-dhcp-{}'.format(subnet_id)
        self.dirname = os.path.join('/var/lib/n0stack/', self.netns_name)
        self.pid_filename = os.path.join(self.dirname, 'pid')
        self.dhcp_hostsfilename = os.path.join(self.dirname, 'hosts')
        self.dhcp_leasefilename = os.path.join(self.dirname, 'lease')
        self.dhcp_optsfilename = os.path.join(self.dirname, 'opts')

    def get_pid(self):
        # type: () -> int
        """
        Get pid of running dnsmasq process on netns.

        Returns:
            pid.
            If pid file or process does not exist, return None.
        """
        if not os.path.exists(self.pid_filename):
            return None
        with open(self.pid_filename, 'r') as f:
            pid = int(f.read())
        try:
            os.kill(pid, 0)
        except OSError:
            return None
        else:
            return pid

    def start_process(self, pool):
        # type: (Tuple[str, str]) -> None
        """
        Start dnsmasq process on netns.

        1. Create directory to save dnsmasq files.
        2. Set args and start process.

        Args:
            pool: Dnsmasq allocation pool. Allocate pool[0]-pool[1].

        Raises:
            Exception: If dnsmasq process is already running, raise Exception.
        """
        if not os.path.exists(self.dirname):
            os.makedirs(self.dirname)

        if self.get_pid() is not None:
            raise Exception("dnsmasq process in {} is already running".format(self.netns_name)) # NOQA

        pid_file = '--pid-file={}'.format(self.pid_filename)
        dhcp_hostsfile = '--dhcp-hostsfile={}'.format(self.dhcp_hostsfilename)
        dhcp_optsfile = '--dhcp-optsfile={}'.format(self.dhcp_optsfilename)
        dhcp_leasefile = '--dhcp-leasefile={}'.format(self.dhcp_leasefilename)
        interface = '--interface={}'.format(self.peer_name)
        dhcp_range = '--dhcp-range={},{},12h'.format(pool[0], pool[1])
        cmd = ['/usr/sbin/dnsmasq',
               '--no-resolv',
               '--no-hosts',
               '--except-interface=lo',
               pid_file,
               dhcp_hostsfile,
               dhcp_optsfile,
               dhcp_leasefile,
               interface,
               dhcp_range]
        NSPopen(self.netns_name, cmd)

    def stop_process(self):
        # type: () -> None
        """
        Stop dnsmasq process on netns.
        """
        pid = self.get_pid()
        if pid is not None:
            os.kill(pid, 9)
        else:
            warn("dnsmasq process is not running in {}".format(self.netns_name)) # NOQA

    def respawn_process(self, pool):
        # type: (Tuple[str, str]) -> None
        """
        Respawn dnsmasq process on netns.

        Args:
            pool: Dnsmasq allocation pool. Allocate pool[0]-pool[1].
        """
        self.stop_process()
        self.start_process(pool)

    def create_dhcp_server(self, interface_addr, bridge_name, pool):
        # type: (IPv4Interface, str, Tuple[str, str]) -> None
        """
        Create Dnsmasq server on specified subnet.

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
            interface_addr: IP address of Dnsmasq server.
            bridge_name: Name of bridge linked to Dnsmasq server.
            pool: Dnsmasq allocation pool. Allocate pool[0]-pool[1].

        Raises:
            Exception: If spcified bridge does not exist, raise Exception.
        """
        bri = Dnsmasq.ip.link_lookup(ifname=bridge_name)
        if bri:
            bri = bri[0]
        else:
            raise Exception("Specified bridge {} does not exist".format(bridge_name)) # NOQA

        netns = NetNS(self.netns_name, flags=os.O_CREAT)

        tap_name = self.tap_name
        peer_name = self.peer_name
        Dnsmasq.ip.link('add', ifname=tap_name, peer=peer_name, kind='veth')

        tap = Dnsmasq.ip.link_lookup(ifname=tap_name)[0]
        Dnsmasq.ip.link('set', index=tap, master=bri)

        peer = Dnsmasq.ip.link_lookup(ifname=peer_name)[0]
        Dnsmasq.ip.link('set', index=peer, net_ns_fd=self.netns_name)

        address = str(interface_addr.ip)
        prefixlen = int(interface_addr.network.prefixlen)
        netns.addr('add', index=peer, address=address, prefixlen=prefixlen)

        Dnsmasq.ip.link('set', index=tap, state='up')
        netns.link('set', index=peer, state='up')
        netns.close()

        self.start_process(pool)

    def delete_dhcp_server(self):
        # type: () -> None
        """
        Delete Dnsmasq server on specified subnet.

        1. Kill dnsmasq process.
        2. Delete directory for dnsmasq files.
        3. Delete veth pairs.
           in command: `ip link del $tap_name`
        4. Delete related netns.
           in command: `ip netns del $netns_name`

        Even if some resources don't exist, go on to delete existing resources.
        """
        self.stop_process();

        if os.path.exists(self.dirname):
            rmtree(self.dirname)
        else:
            warn("dnsmasq directory {} does not exist".format(self.dirname)) # NOQA

        tap = Dnsmasq.ip.link_lookup(ifname=self.tap_name)
        if tap:
            Dnsmasq.ip.link('del', index=tap[0])
        else:
            warn("veth {} does not exist".format(self.tap_name))

        netns = NetNS(self.netns_name)
        netns.close()
        netns.remove()

    def add_host_entry(self, hw_addr, ip_addr):
        # type: (str, str) -> None
        """
        Add MAC:IP mapping in order to assign IP address statically.

        1. Write mapping to dhcp-hostsfile.
        2. Send SIGHUP to dnsmasq process.

        Args:
            hw_address: MAC address of interface.
            ip_address: IP address of interface.

        Raise:
            Exception: If dnsmasq process is not running, raise Exeception.
        """
        pid = self.get_pid()
        if pid is None:
            raise Exception("dnsmasq process is not running in {}".format(self.netns_name)) # NOQA

        with open(self.dhcp_hostsfilename, 'a') as f:
            f.write('{},{}\n'.format(hw_addr, ip_addr))

        os.kill(pid, 1)
