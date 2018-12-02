package network

import (
	"fmt"
	"net"

	"github.com/n0stack/n0stack/n0core/pkg/util/net"
	"github.com/n0stack/n0stack/n0proto.go/budget/v0"
)

func CheckIPv4OnCIDR(request net.IP, cidr *net.IPNet) error {
	if !cidr.Contains(request) {
		return fmt.Errorf("Requested IPv4 '%s' is over from Network CIDR '%s'", request.String(), cidr.String())
	}

	return nil
}

func CheckConflictIPv4(request net.IP, reserved map[string]*pbudget.NetworkInterface) error {
	for k, v := range reserved {
		if request.String() == v.Ipv4Address {
			return fmt.Errorf("Network interface '%s' is already have IPv4 address '%s'", k, request.String())
		}
	}

	return nil
}

// O(len(reserved) ^ 2) なので要修正
func ScheduleNewIPv4(cidr *net.IPNet, reserved map[string]*pbudget.NetworkInterface) net.IP {
	// escape network address and broadcast address
	for ip := netutil.NextIP(cidr.IP); cidr.Contains(netutil.NextIP(ip)); ip = netutil.NextIP(ip) {
		if err := CheckConflictIPv4(ip, reserved); err == nil {
			return ip
		}
	}

	return nil
}
