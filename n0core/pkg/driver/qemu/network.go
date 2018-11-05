package qemu

import (
	"encoding/json"
	"fmt"
	"net"
)

// (QEMU) netdev_add id=tap0 type=tap vhost=true ifname=tap0 script=no downscript=no
// (QEMU) device_add driver=virtio-net-pci netdev=tap0 id=test0 mac=52:54:00:df:89:29 bus=pci.0
// まだべき等ではない
// TODO:
//   - すでにアタッチされていた場合、エラー処理を文字列で判定する必要がある
//   - MACアドレスを変更する
func (q Qemu) AttachTap(label, tap string, mac net.HardwareAddr) error {
	netdevID := fmt.Sprintf("netdev-%s", label)
	devID := fmt.Sprintf("virtio-net-%s", label)

	// check to create netdev

	if err := q.tapNetdevAdd(netdevID, tap); err != nil {
		return fmt.Errorf("Failed to run netdev_add: err='%s'", err.Error())
	}

	if err := q.virtioNetPCIAdd(devID, netdevID, mac); err != nil {
		return fmt.Errorf("Failed to create virtio network device: err='%s'", err.Error())
	}

	return nil
}

func (q *Qemu) tapNetdevAdd(id, ifname string) error {
	cmd := struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Ifname     string `json:"ifname"`
		Vhost      bool   `json:"vhost"`
		Script     string `json:"script"`
		Downscript string `json:"downscript"`
	}{
		id,
		"tap",
		ifname,
		true,
		"no",
		"no",
	}

	bs, err := json.Marshal(map[string]interface{}{
		"execute":   "netdev_add",
		"arguments": cmd,
	})
	if err != nil {
		return err
	}

	_, err = q.qmp.Run(bs)
	if err != nil {
		return err
	}

	return err
}

func (q *Qemu) virtioNetPCIAdd(devID, netdevID string, mac net.HardwareAddr) error {
	cmd := struct {
		Driver string `json:"driver"`
		ID     string `json:"id"`
		Netdev string `json:"netdev"`
		Bus    string `json:"bus"`
		Mac    string `json:"mac"`
	}{
		"virtio-net-pci",
		devID,
		netdevID,
		"pci.0",
		mac.String(),
	}

	bs, err := json.Marshal(map[string]interface{}{
		"execute":   "device_add",
		"arguments": cmd,
	})
	if err != nil {
		return err
	}

	_, err = q.qmp.Run(bs)
	if err != nil {
		return err
	}

	return nil
}
