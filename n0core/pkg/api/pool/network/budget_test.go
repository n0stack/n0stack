package network

import (
	"net"
	"testing"

	"n0st.ac/n0stack/n0proto.go/budget/v0"
)

func TestScheduleNewIPv4(t *testing.T) {
	_, cidr, _ := net.ParseCIDR("192.168.0.0/30")

	cases := []struct {
		// cidr *net.IPNet
		reserved map[string]*pbudget.NetworkInterface
		result   net.IP
	}{
		{
			map[string]*pbudget.NetworkInterface{
				"hoge": {
					Ipv4Address: "192.168.0.1",
				},
			},
			net.ParseIP("192.168.0.2"),
		},
		{
			map[string]*pbudget.NetworkInterface{
				"foo": {
					Ipv4Address: "192.168.0.1",
				},
				"bar": {
					Ipv4Address: "192.168.0.2",
				},
			},
			nil,
		},
	}

	for _, c := range cases {
		ip := ScheduleNewIPv4(cidr, c.reserved)

		if c.result != nil && ip == nil {
			t.Errorf("Wrong generated IPv4 address\n\thave:nil\n\twant:%s", c.result)
		} else if c.result == nil && ip != nil {
			t.Errorf("Wrong generated IPv4 address\n\thave:%s\n\twant:nil", ip)
		} else if !ip.Equal(c.result) {
			t.Errorf("Wrong generated IPv4 address\n\thave:%s\n\twant:%s", ip, c.result)
		}
	}
}
