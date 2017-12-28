from netaddr.ip import IPAddress  # NOQA
import os
from shutil import rmtree
from typing import Any, Dict, List, Optional, Tuple, Union  # NOQA

from netns import NetNS as nsscope

from iptc import Chain, Match, Policy, Rule, Table, Target

from pyroute2 import IPRoute
from pyroute2 import NetlinkError
from pyroute2 import NetNS
from pyroute2 import NSPopen

from n0library.logger import Logger


logger = Logger(name=__name__)  # type: Logger


class Dnsmasq(object):
    """
    Manage namespace, veth pair, directory and dnsmasq process.

    Example:
        Create DHCP server.
        >>> masq = Dnsmasq("subnet01")
        >>> ip_address = ipaddress.ip_interface('192.168.1.1/24')
        >>> pool = ['192.168.1.2', '192.168.1.254']
        >>> bridge = 'br-subnet01'
        >>> masq.create_dhcp_server(ip_address, bridge, pool)

        Enable VM having the MAC address to communicate with dnsmasq process by DHCP.
        >>> masq.add_allowed_host('9e:e8:1d:d0:2c:c9')

        Delete allowed host.
        >>> masq.delete_allowed_host('9e:e8:1d:d0:2c:c9')

        Add static MAC-IP map entry. VM having the MAC address always receives specified IP address.
        >>> masq.add_host_entry('9e:e8:1d:d0:2c:c9', '192.168.1.11')

        Delete MAC-IP entry.
        >>> masq.delete_host_entry('9e:e8:1d:d0:2c:c9')

        Stop dnsmasq process.
        >>> masq.stop_process()

        Start dnsmasq process.
        >>> masq.start_process(pool)

        Respawn dnsmasq process to change its DHCP allocation pool.
        >>> pool = ['192.168.1.2', '192.168.1.128']
        >>> masq.respawn_process(pool)

        Delete DHCP server.
        >>> masq.delete_dhcp_server()
    """
    ip = IPRoute()  # type: IPRoute

    def __init__(self, network_id, bridge_name, log_dir=None):  # logging is disabled
        # type: (str, str, Optional[str]) -> None
        """
        Set names in order to create or delete resources.

        Args:
            subnet_id: Subnet ID.
            bridge_name: Name of bridge linked to Dnsmasq server.
            log_dir: Directory path to save log file.
        """
        id = network_id.split("-")[0]
        self.netns_name = 'nndhcp{}'.format(network_id)
        self.tap_name = 'ntndhcp{}'.format(id)
        self.peer_name = 'nedhcp{}'.format(id)
        self.bridge_name = bridge_name

        self.dirname = os.path.join('/var/lib/n0stack/', self.netns_name)  # type: str
        self.pid_filename = os.path.join(self.dirname, 'pid')  # type: str
        self.dhcp_hostsfilename = os.path.join(self.dirname, 'hosts')  # type: str
        self.dhcp_leasefilename = os.path.join(self.dirname, 'lease')  # type: str
        self.dhcp_optsfilename = os.path.join(self.dirname, 'opts')  # type: str
        # self.log_path = os.path.join(log_dir, 'dnsmasq-{}.log'.format(self.netns_name))  # type: str

        self.__hw_to_ip = {}  # type: Dict[str, str]
        self.__ip_to_hw = {}  # type: Dict[str, str]
        self.__load_hostsfile()

    def __load_hostsfile(self):
        # type: () -> None
        hostsfile = self.dhcp_hostsfilename  # type: str

        if not os.path.exists(hostsfile):
            return

        with open(hostsfile, 'r') as f:
            lines = f.readlines()

        for line in lines:
            pair = line.strip().split(',')  # type: List[str]

            if len(pair) != 2:
                continue

            hw_addr = pair[0]  # type: str
            ip_addr = pair[1]  # type: str
            self.__hw_to_ip[hw_addr] = ip_addr
            self.__ip_to_hw[ip_addr] = hw_addr

    def get_pid(self):
        # type: () -> Optional[int]
        """
        Get pid of running dnsmasq process on netns.

        Returns:
            pid.
            If pid file or process does not exist, return None.
        """
        if not os.path.exists(self.pid_filename):
            return None

        with open(self.pid_filename, 'r') as f:
            pid = int(f.read())  # type: int

        try:
            os.kill(pid, 0)
        except OSError:
            return None
        else:
            return pid

    def start_process(self, pool):
        # type: (Tuple[IPAddress, IPAddress]) -> None
        """
        Start dnsmasq process on netns.

        1. Create directory to save dnsmasq files.
        2. Set args and start process.

        Args:
            pool: Dnsmasq allocation pool. Allocate pool[0]-pool[1].

        Raises:
            Exception: If interface to bind does not exist, raise Exception.
        """
        if not os.path.exists(self.dirname):
            os.makedirs(self.dirname)
            open(self.dhcp_hostsfilename, 'w').close()
            open(self.dhcp_optsfilename, 'w').close()

        if self.get_pid() is not None:
            logger.warning("dnsmasq process in {} is already running".format(self.netns_name))
            return

        ns = NetNS(self.netns_name)  # type: NetNS

        if not ns.link_lookup(ifname=self.peer_name):
            raise Exception("Interface {} does not exist".format(self.peer_name))

        pid_file = '--pid-file={}'.format(self.pid_filename)  # type: str
        dhcp_hostsfile = '--dhcp-hostsfile={}'.format(self.dhcp_hostsfilename)  # type: str
        dhcp_optsfile = '--dhcp-optsfile={}'.format(self.dhcp_optsfilename)  # type: str
        dhcp_leasefile = '--dhcp-leasefile={}'.format(self.dhcp_leasefilename)  # type: str
        interface = '--interface={}'.format(self.peer_name)  # type: str
        dhcp_range = '--dhcp-range={},{},12h'.format(pool[0], pool[1])  # type: str
        # log_facility = '--log-facility={}'.format(self.log_path)
        cmd = ['/usr/sbin/dnsmasq',
               '--no-resolv',
               '--no-hosts',
               '--bind-interfaces',
               '--except-interface=lo',
               pid_file,
               dhcp_hostsfile,
               dhcp_optsfile,
               dhcp_leasefile,
               interface,
               dhcp_range]
            #    log_facility]  # type: List[str]
        NSPopen(self.netns_name, cmd)

    def stop_process(self):
        # type: () -> None
        """
        Stop dnsmasq process on netns.
        """
        pid = self.get_pid()  # type: Optional[int]

        if pid is not None:
            os.kill(pid, 9)
        else:
            logger.warning("dnsmasq process is not running in {}".format(self.netns_name))

    def respawn_process(self, pool):
        # type: (Tuple[IPAddress, IPAddress]) -> None
        """
        Respawn dnsmasq process on netns.

        Args:
            pool: Dnsmasq allocation pool. Allocate pool[0]-pool[1].
        """
        self.stop_process()
        self.start_process(pool)

    def init_iptables(self):
        # type: () -> None
        """
        Insert rules for dhcp server to iptabls in netns.

        1. Set policy of each chain DROP.
           in command: `iptables -P DROP $chain`
        2. Allow ICMP in/out.
           in command: `iptables -I $chain -p icmp -j ACCEPT`
        3. Allow DHCP out.
           in command: `iptables -I $chain -p udp --sport 67 --dport 68 -j ACCEPT`

        If rule already exists, skip insertion.
        """
        with nsscope(nsname=self.netns_name):
            table = Table(Table.FILTER)  # type: Table

            for chain in table.chains:
                chain.flush()
                chain.set_policy(Policy('DROP'))

            ping_rule = Rule()  # type: Rule
            ping_rule.protocol = 'icmp'
            ping_rule.target = Target(ping_rule, 'ACCEPT')

            input_chain = Chain(table, 'INPUT')  # type: Chain

            if all([ping_rule != rule for rule in input_chain.rules]):
                input_chain.insert_rule(ping_rule)

            output_chain = Chain(table, 'OUTPUT')  # type: Chain

            if all([ping_rule != rule for rule in output_chain.rules]):
                output_chain.insert_rule(ping_rule)

            dhcp_rule = Rule()  # type: Rule
            dhcp_rule.protocol = 'udp'
            match = Match(dhcp_rule, 'udp')  # type: Match
            match.sport = '67'
            match.dport = '68'
            dhcp_rule.add_match(match)
            dhcp_rule.target = Target(dhcp_rule, 'ACCEPT')

            if all([dhcp_rule != rule for rule in output_chain.rules]):
                output_chain.insert_rule(dhcp_rule)

    def create_dhcp_server(self, interface_addr, pool):  # interface_addrを文字列のcidrを受け取ってインターフェイスをいいかんじに作成する
        # type: (IPAddress, Tuple[IPAddress, IPAddress]) -> None
        """
        Create Dnsmasq server on specified subnet.

        1. Create netns if not exists.
           in command: `ip netns add $netns_name`
        2. Create veth pair.
           in command: `ip link add $tap_name type veth peer name $peer_name`
        3. Move one of the veth pair to netns.
           in command: `ip link set $peer_name netns $netns_name`
        4. Add ip address to the veth.
           in command: `ip netns exec $netns_name  ip addr add $address/$prefixlen dev $peer`
        5. Link the other to bridge.
           in command: `ip link set dev $tap_name master $bridge_name`
        6. Set up veths.
           in command: `ip link set $name up`
        7. Start dnsmasq process.

        Args:
            interface_addr: IP address of Dnsmasq server.
            pool: Dnsmasq allocation pool. Allocate pool[0]-pool[1].

        Raises:
            Exception: If specified bridge does not exist, raise Exception.
            Exception: If one of the veth pair exists and the other not, raise Exception.
        """
        bri_list = self.ip.link_lookup(ifname=self.bridge_name)  # type: List[Any]
        bri = None  # type: Optional[int]

        if bri_list:
            bri = bri_list[0]
        else:
            raise Exception("Specified bridge {} does not exist".format(self.bridge_name))

        ns = NetNS(self.netns_name, flags=os.O_CREAT)  # type: NetNS

        self.init_iptables()

        tap_name = self.tap_name  # type: str
        peer_name = self.peer_name  # type: str
        peer = None  # type: Optional[int]

        try:
            self.ip.link('add', ifname=tap_name, peer=peer_name, kind='veth')
        except NetlinkError as e:
            if e.code == 17:
                logger.warning("veth {} existing, ignore and continue".format(tap_name))
                peer_list = ns.link_lookup(ifname=peer_name)  # type: List[Any]

                if peer_list:
                    peer = peer_list[0]
                else:
                    raise Exception("One of the veth pair {} exists, but the other {} does not exists".format(tap_name, peer_name))  # NOQA
            else:
                raise e
        else:
            peer = self.ip.link_lookup(ifname=peer_name)[0]
            self.ip.link('set', index=peer, net_ns_fd=self.netns_name)

        interface_addr = interface_addr + 1  # rangeと被らないようにする必要がある
        address = str(interface_addr.ip)
        prefixlen = int(interface_addr.network.prefixlen)

        try:
            ns.addr('add', index=peer, address=address, prefixlen=prefixlen)
        except NetlinkError as e:
            if e.code == 17:
                logger.warning("IP address is already assinged to {}, ignore and continue".format(peer_name))
            else:
                raise e

        tap = self.ip.link_lookup(ifname=tap_name)[0]  # type: int
        self.ip.link('set', index=tap, master=bri)

        self.ip.link('set', index=tap, state='up')
        ns.link('set', index=peer, state='up')
        ns.close()

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
        self.stop_process()

        if os.path.exists(self.dirname):
            rmtree(self.dirname)
        else:
            logger.warning("dnsmasq directory {} does not exist".format(self.dirname))

        tap_list = self.ip.link_lookup(ifname=self.tap_name)  # type: List[Any]

        if tap_list:
            self.ip.link('del', index=tap_list[0])
        else:
            logger.warning("veth {} does not exist".format(self.tap_name))

        ns = NetNS(self.netns_name)  # type: NetNS
        ns.close()
        ns.remove()

    def _get_dhcp_allow_rule(self, hw_addr):
        # type: (str) -> Rule
        rule = Rule()  # type: Rule
        rule.protocol = 'udp'
        rule.target = Target(rule, 'ACCEPT')

        proto_match = Match(rule, 'udp')  # type: Match
        proto_match.sport = '68'
        proto_match.dport = '67'
        rule.add_match(proto_match)

        mac_match = Match(rule, 'mac')  # type: Match
        mac_match.mac_source = hw_addr
        rule.add_match(mac_match)

        return rule

    def add_allowed_host(self, hw_addr):
        # type: (str) -> None
        """
        Allow DHCP input from specified host.
        in command: `iptables -I INPUT -p udp --sport 68 --dport 67 -m --mac-source $hw_address`
        If rule already exists, skip insertion.

        Args:
            hw_addr: MAC address of host.
        """
        with nsscope(nsname=self.netns_name):
            chain = Chain(Table(Table.FILTER), 'INPUT')
            dhcp_rule = self._get_dhcp_allow_rule(hw_addr)  # type: Rule

            if all([dhcp_rule != rule for rule in chain.rules]):
                chain.insert_rule(dhcp_rule)

    def delete_allowed_host(self, hw_addr):
        # type: (str) -> None
        """
        Delete rule allowing DHCP input from specified host.

        Args:
            hw_addr: MAC address of host.
        """
        with nsscope(nsname=self.netns_name):
            chain = Chain(Table(Table.FILTER), 'INPUT')
            dhcp_rule = self._get_dhcp_allow_rule(hw_addr)  # type: Rule

            for rule in chain.rules:
                if rule == dhcp_rule:
                    chain.delete_rule(rule)

    def add_host_entry(self, hw_addr, ip_addr):
        # type: (str, str) -> None
        """
        Add MAC:IP mapping entry in order to assign IP address statically.

        1. Update entries.
        2. Write entries to dhcp-hostsfile.
        3. Send SIGHUP to dnsmasq process.

        Args:
            hw_address: MAC address of interface.
            ip_address: IP address of interface.

        Raise:
            Exception: If dnsmasq process is not running, raise Exeception.
            Exception: If specified IP address is already in entries, raise Exception.
        """
        pid = self.get_pid()  # type: Optional[int]

        if pid is None:
            raise Exception("dnsmasq process is not running in {}".format(self.netns_name))

        if ip_addr in self.__ip_to_hw:
            raise Exception("Specified IP address {} is already used".format(ip_addr))

        if hw_addr in self.__hw_to_ip:
            logger.warning("Spacified MAC address {} being used, update entry".format(hw_addr))
            self.__hw_to_ip[hw_addr] = ip_addr
            self.__ip_to_hw[ip_addr] = hw_addr

            with open(self.dhcp_hostsfilename, 'w') as f:
                for k, v in self.__hw_to_ip.items():
                    f.write('{},{}\n'.format(k, v))

        else:
            self.__hw_to_ip[hw_addr] = ip_addr
            self.__ip_to_hw[ip_addr] = hw_addr

            with open(self.dhcp_hostsfilename, 'a') as f:
                f.write('{},{}\n'.format(hw_addr, ip_addr))

        os.kill(pid, 1)

    def delete_host_entry(self, hw_addr):
        # type: (str) -> None
        """
        Delete MAC:IP mapping entry.
        If specified MAC address is not in entry, do nothing.

        Args:
            hw_addr: MAC address of entry to delete.

        Raises:
            Exception: If dnsmasq process is not running, raise Exeception.
        """
        pid = self.get_pid()  # type: Optional[int]

        if pid is None:
            raise Exception("dnsmasq process is not running in {}".format(self.netns_name))

        if hw_addr not in self.__hw_to_ip:
            logger.warning("Specified MAC address {} does not exist in entries".format(hw_addr))
            return

        ip_addr = self.__hw_to_ip[hw_addr]  # type: str
        del self.__ip_to_hw[ip_addr]
        del self.__hw_to_ip[hw_addr]

        with open(self.dhcp_hostsfilename, 'w') as f:
            for k, v in self.__hw_to_ip.items():
                f.write('{},{}\n'.format(k, v))

        os.kill(pid, 1)
