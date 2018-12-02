// +build medium

package iptables

import (
	"net"
	"testing"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
)

func TestMasquerade(t *testing.T) {
	b, err := iproute2.NewBridge("masq-br")
	if err != nil {
		t.Fatalf("failed to create bridge: err='%s'", err.Error())
	}
	defer b.Delete()

	if err := b.SetAddress("172.31.255.254/24"); err != nil {
		t.Fatalf("failed to set address to bridge: err='%s'", err.Error())
	}

	_, n, _ := net.ParseCIDR("172.31.255.0/24")
	if err := CreateMasqueradeRule(b.Name(), n); err != nil {
		t.Errorf("Failed to create masquerade rule: err='%s'", err.Error())
	}

	if err := DeleteMasqueradeRule(b.Name(), n); err != nil {
		t.Errorf("Failed to delete masquerade rule: err='%s'", err.Error())
	}
}
