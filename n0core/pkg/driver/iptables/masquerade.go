package iptables

import (
	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
)

func structMasqueradeRule(bridgeName, network string) (string, string, []string) {
	return "nat", "POSTROUTING", []string{
		"-o",
		bridgeName,
		"-s",
		network,
		"-j",
		"MASQUERADE",
	}
}

func CreateMasqueradeRule(bridgeName, network string) error {
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

func DeleteMasqueradeRule(bridgeName, network string) error {
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
