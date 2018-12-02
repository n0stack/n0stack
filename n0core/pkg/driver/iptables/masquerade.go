package iptables

import (
	"net"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
)

func structMasqueradeRule(bridgeName string, network *net.IPNet) (string, string, []string) {
	return "nat", "POSTROUTING", []string{
		"!",
		"-o",
		bridgeName,
		"-s",
		network.String(),
		"-j",
		"MASQUERADE",
	}
}

func CreateMasqueradeRule(bridgeName string, network *net.IPNet) error {
	ipt, err := iptables.New()
	if err != nil {
		return errors.Wrapf(err, "Failed to create iptables instance")
	}

	// iptables -t nat -A POSTROUTING -o br0 -s 192.168.0.0/24 -j MASQUERADE
	table, chain, rule := structMasqueradeRule(bridgeName, network)

	if exists, err := ipt.Exists(table, chain, rule...); err != nil {
		return errors.Wrapf(err, "Failed to check rule")
	} else if exists {
		return nil
	}

	if err := ipt.Insert(table, chain, 1, rule...); err != nil {
		return errors.Wrapf(err, "Failed to create iptables instance: rule='%s'", rule)
	}

	return nil
}

func DeleteMasqueradeRule(bridgeName string, network *net.IPNet) error {
	ipt, err := iptables.New()
	if err != nil {
		return errors.Wrapf(err, "Failed to create iptables instance")
	}

	// iptables -t nat -A POSTROUTING -o br0 -s 192.168.0.0/24 -j MASQUERADE
	table, chain, rule := structMasqueradeRule(bridgeName, network)

	if exists, err := ipt.Exists(table, chain, rule...); err != nil {
		return errors.Wrapf(err, "Failed to check rule")
	} else if !exists {
		return nil
	}

	if err := ipt.Delete(table, chain, rule...); err != nil {
		return errors.Wrapf(err, "Failed to delete iptables instance: rule='%s'", rule)
	}

	return nil
}
