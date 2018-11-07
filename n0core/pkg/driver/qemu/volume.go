package qemu

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// (QEMU) blockdev-add options={"driver":"qcow2","id":"drive-virtio-disk0","file":{"driver":"file","filename":"/home/h-otter/wk/test-qemu/ubuntu16.04.qcow2"}}
// (QEMU) device_add driver=virtio-blk-pci bus=pci.0 scsi=off drive=drive-virtio-disk0 id=virtio-disk0 bootindex=1
// まだべき等ではない
// TODO:
//   - すでにアタッチされていた場合、エラー処理を文字列で判定する必要がある
//   - bootindexがどうやって更新されるのかがわからない
func (q Qemu) AttachQcow2(label string, filepath *url.URL, index uint) error {
	driveID := fmt.Sprintf("drive-%s", label)
	devID := fmt.Sprintf("virtio-blk-%s", label)

	// check to create drive

	if err := q.qcow2BlockdevAdd(driveID, filepath); err != nil {
		return fmt.Errorf("Failed to run qcow2BlockdevAdd: err='%s'", err.Error())
	}

	if err := q.virtioBlkPCIAdd(devID, driveID, index); err != nil {
		return fmt.Errorf("Failed to create virtio block device: err='%s'", err.Error())
	}

	return nil
}

func (q Qemu) AttachRaw(label string, filepath *url.URL, index uint) error {
	driveID := fmt.Sprintf("drive-%s", label)
	devID := fmt.Sprintf("virtio-blk-%s", label)

	// check to create drive

	if err := q.rawBlockdevAdd(driveID, filepath, false); err != nil {
		return fmt.Errorf("Failed to run rawBlockdevAdd: err='%s'", err.Error())
	}

	if err := q.virtioBlkPCIAdd(devID, driveID, index); err != nil {
		return fmt.Errorf("Failed to create virtio block device: err='%s'", err.Error())
	}

	return nil
}

// TODO: 動作が怪しい
func (q Qemu) AttachISO(label string, filepath *url.URL, index uint) error {
	driveID := fmt.Sprintf("drive-scsi0-cd-%s", label)
	devID := fmt.Sprintf("scsi0-cd-%s", label)

	// check to create drive

	if err := q.rawBlockdevAdd(driveID, filepath, true); err != nil {
		return fmt.Errorf("Failed to run rawBlockdevAdd: err='%s'", err.Error())
	}

	if err := q.scsiCDAdd(devID, driveID, index); err != nil {
		return fmt.Errorf("Failed to create scsi cd device: err='%s'", err.Error())
	}

	return nil
}

func (q *Qemu) qcow2BlockdevAdd(nodeName string, filepath *url.URL) error {
	f := struct {
		Driver   string `json:"driver"`
		Filename string `json:"filename"`
	}{
		filepath.Scheme,
		filepath.Path,
	}
	cmd := struct {
		NodeName string      `json:"node-name"`
		Driver   string      `json:"driver"`
		File     interface{} `json:"file"`
	}{
		nodeName,
		"qcow2",
		f,
	}

	bs, err := json.Marshal(map[string]interface{}{
		"execute":   "blockdev-add",
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

func (q *Qemu) virtioBlkPCIAdd(devID, driveID string, bootIndex uint) error {
	cmd := struct {
		Driver    string `json:"driver"`
		ID        string `json:"id"`
		Drive     string `json:"drive"`
		Bus       string `json:"bus"`
		Scsi      string `json:"scsi"`
		BootIndex string `json:"bootindex"`
	}{
		"virtio-blk-pci",
		devID,
		driveID,
		"pci.0",
		"off",
		strconv.FormatUint(uint64(bootIndex), 10),
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

func (q *Qemu) rawBlockdevAdd(nodeName string, filepath *url.URL, readonly bool) error {
	f := struct {
		Driver   string `json:"driver"`
		Filename string `json:"filename"`
	}{
		filepath.Scheme,
		filepath.Path,
	}
	cmd := struct {
		NodeName string      `json:"node-name"`
		Driver   string      `json:"driver"`
		ReadOnly bool        `json:"read-only"`
		File     interface{} `json:"file"`
	}{
		nodeName,
		"raw",
		readonly,
		f,
	}

	bs, err := json.Marshal(map[string]interface{}{
		"execute":   "blockdev-add",
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

func (q *Qemu) scsiCDAdd(devID, driveID string, bootIndex uint) error {
	cmd := struct {
		Driver    string `json:"driver"`
		ID        string `json:"id"`
		Drive     string `json:"drive"`
		Bus       string `json:"bus"`
		BootIndex string `json:"bootindex"`
	}{
		"scsi-cd",
		devID,
		driveID,
		"scsi0.0",
		strconv.FormatUint(uint64(bootIndex), 10),
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
