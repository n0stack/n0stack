package tap

// ip tuntap add dev $tap_name mode tap # user group
// ip link set master $bridge_name dev $tap_name
// ip link set up dev $tap_name
// iptables -P FORWARD ACCEPT

const NICType = "tap"
