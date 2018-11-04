package network

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"net"

	"github.com/n0stack/n0stack/n0proto.go/budget/v0"
)

func GenerateHardwareAddress(id string) net.HardwareAddr {
	cs := crc32.Checksum([]byte(id), crc32.IEEETable)
	b, _ := hex.DecodeString(fmt.Sprintf("5254%08x", cs))
	return net.HardwareAddr(b)
}

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
	for ip := NextIP(cidr.IP); cidr.Contains(NextIP(ip)); ip = NextIP(ip) {
		if err := CheckConflictIPv4(ip, reserved); err == nil {
			return ip
		}
	}

	return nil
}

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
