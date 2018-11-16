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
