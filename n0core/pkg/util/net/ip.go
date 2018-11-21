package nettools

import (
	"fmt"
	"net"
)

func NextIP(ip net.IP) net.IP {
	res := net.ParseIP(ip.String())
	for i := len(res) - 1; i >= 0; i-- {
		res[i]++

		if res[i] > 0 {
			break
		}
	}

	return res
}

func GetEndIP(ipn *net.IPNet) net.IP {
	// next network address
	ip := NextIP(ipn.IP)

	for {
		// before broadcast address
		if !ipn.Contains(NextIP(NextIP(ip))) {
			return ip
		}

		ip = NextIP(ip)
	}
}

type IPv4Cidr struct {
	ip      net.IP
	network *net.IPNet
}

func ParseCIDR(cidr string) *IPv4Cidr {
	ip, ipn, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	return &IPv4Cidr{
		ip:      ip,
		network: ipn,
	}
}

func (c IPv4Cidr) String() string {
	return fmt.Sprintf("%s/%d", c.ip.String(), c.SubnetMaskBits())
}

func (c IPv4Cidr) IP() net.IP {
	return c.ip
}

func (c IPv4Cidr) Next() *IPv4Cidr {
	next := NextIP(c.ip)
	if !c.network.Contains(next) {
		return nil
	}

	return &IPv4Cidr{
		ip:      next,
		network: c.network,
	}
}

func (c IPv4Cidr) Network() *net.IPNet {
	net := *c.network

	return &net
}

func (c IPv4Cidr) SubnetMaskBits() int {
	m, _ := c.network.Mask.Size()

	return m
}

func (c IPv4Cidr) SubnetMaskIP() net.IP {
	return net.IPv4(c.network.Mask[0], c.network.Mask[1], c.network.Mask[2], c.network.Mask[3])
}
