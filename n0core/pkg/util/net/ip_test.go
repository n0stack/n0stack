package nettools

import (
	"net"
	"testing"
)

func TestGetEndIP(t *testing.T) {
	_, ipn, _ := net.ParseCIDR("192.168.0.0/24")
	if ip := GetEndIP(ipn); ip.String() != "192.168.0.254" {
		t.Errorf("Failed to get end of ip: got='%s', want='%s'", ip.String(), "192.168.0.254")
	}
}

func TestIPv4Cidr(t *testing.T) {
	cidr := ParseCIDR("aa")
	if cidr != nil {
		t.Errorf("ParseCIDR() do not return nil when over: have=%v", cidr)
	}

	cidr = ParseCIDR("192.168.0.2/30")
	if cidr.String() != "192.168.0.2/30" {
		t.Errorf("String() is wrong: have='%s', want='%s'", cidr.String(), "192.168.0.2/30")
	}
	if cidr.Next().String() != "192.168.0.3/30" {
		t.Errorf("Next().String() is wrong: have='%s', want='%s'", cidr.Next().String(), "192.168.0.3/30")
	}
	if n := cidr.Next().Next(); n != nil {
		t.Errorf("Next() do not return nil when over: have=%v", n)
	}

	if cidr.IP().String() != "192.168.0.2" {
		t.Errorf("IP() is wrong: have='%s', want='%s'", cidr.IP().String(), "192.168.0.2/30")
	}
	if cidr.Network().String() != "192.168.0.0/30" {
		t.Errorf("Network() is wrong: have='%s', want='%s'", cidr.Network().String(), "192.168.0.0/30")
	}
	if cidr.SubnetMaskBits() != 30 {
		t.Errorf("SubnetMaskBits() is wrong: have='%d', want='%d'", cidr.SubnetMaskBits(), 30)
	}
	if cidr.SubnetMaskIP().String() != "255.255.255.252" {
		t.Errorf("SubnetMaskIP() is wrong: have='%s', want='%s'", cidr.SubnetMaskIP().String(), "192.168.0.2")
	}
}
