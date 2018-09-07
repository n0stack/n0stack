package qemu

import (
	"net/url"
	"path/filepath"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/n0stack/n0core/module/qemu_img"
	uuid "github.com/satori/go.uuid"
)

// 大体30秒くらいかかる
// http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img に依存する
func TestQcow2Volume(t *testing.T) {
	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := OpenQemu(&id)
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	if q.IsRunning() {
		t.Fatalf("Test environment is invalid, process is already existing: uuid='%s'", id.String())
	}

	b, _ := bytefmt.ToBytes("512M")
	if err := q.StartProcess("test", "monitor.sock", 10000, 1, b); err != nil {
		t.Errorf("Failed to start process: err='%s'", err.Error())
	}
	if !q.IsRunning() {
		t.Errorf("Failed to start process, qemu is not running yet")
	}

	f := "cirros.qcow2"
	i, err := qemu_img.OpenQemuImg("cirros.qcow2")
	if err != nil {
		t.Fatalf("Failed to open qemu-img, do not relate to this package code: err='%s'", err.Error())
	}

	u, _ := url.Parse("http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img")
	if err := i.Download(u); err != nil {
		t.Fatalf("Failed to download image, do not relate to this package code: err='%s'", err.Error())
	}
	defer i.Delete()

	p, _ := filepath.Abs(f)
	u, _ = url.Parse("file://" + p)
	if err := q.AttachQcow2("cirros", u, 0); err != nil {
		t.Errorf("Failed to attach qcow2: err='%s'", err.Error())
	}

	if err := q.Boot(); err != nil {
		t.Errorf("Failed to boot: err='%s'", err.Error())
	}

	s, err := q.Status()
	if err != nil {
		t.Errorf("Failed to get status: err='%s'", err.Error())
	}
	if s != StatusRunning {
		t.Errorf("Status is mismatch: want='%v', have='%v'", StatusRunning, s)
	}

	if err := q.Kill(); err != nil {
		t.Errorf("Failed to kill process: err='%s'", err.Error())
	}
	if q.IsRunning() {
		t.Errorf("Failed to kill process, qemu is running yet")
	}
}

func TestISOVolume(t *testing.T) {}
