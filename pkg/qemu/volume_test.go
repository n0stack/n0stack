// +build medium

package qemu

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/n0stack/n0core/pkg/qemu_img"
	uuid "github.com/satori/go.uuid"
)

// 大体10秒くらいかかる
// http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img に依存する
func TestQcow2Volume(t *testing.T) {
	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := OpenQemu(&id)
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	defer q.Delete()

	if _, ok := os.LookupEnv("DISABLE_KVM"); ok {
		q.isKVM = false
	}

	b, _ := bytefmt.ToBytes("512M")
	if err := q.StartProcess("test", "monitor.sock", 10000, 1, b); err != nil {
		t.Fatalf("Failed to start process: err='%s'", err.Error())
	}

	f := "cirros.qcow2"
	i, err := qemu_img.OpenQemuImg(f)
	if err != nil {
		t.Fatalf("Failed to open qemu-img, do not relate to this package code: err='%s'", err.Error())
	}

	u, _ := url.Parse("http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img")
	if err := i.Download(u); err != nil {
		t.Fatalf("Failed to download image, do not relate to this package code: url='%s', err='%s'", u.String(), err.Error())
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
}

// 大体20秒くらいかかる
// http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso に依存する
func TestISOVolume(t *testing.T) {
	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := OpenQemu(&id)
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	defer q.Delete()

	if _, ok := os.LookupEnv("DISABLE_KVM"); ok {
		q.isKVM = false
	}

	b, _ := bytefmt.ToBytes("512M")
	if err := q.StartProcess("test", "monitor.sock", 10000, 1, b); err != nil {
		t.Fatalf("Failed to start process: err='%s'", err.Error())
	}

	f := "ubuntu_mini.iso"
	i, err := qemu_img.OpenQemuImg(f)
	if err != nil {
		t.Fatalf("Failed to open qemu-img, do not relate to this package code: err='%s'", err.Error())
	}

	u, _ := url.Parse("http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso")
	if err := i.Download(u); err != nil {
		t.Fatalf("Failed to download iso, do not relate to this package code: url='%s', err='%s'", u.String(), err.Error())
	}
	defer i.Delete()

	p, _ := filepath.Abs(f)
	u, _ = url.Parse("file://" + p)
	if err := q.AttachISO("ubuntu_mini", u, 0); err != nil {
		t.Errorf("Failed to attach iso: err='%s'", err.Error())
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
}
