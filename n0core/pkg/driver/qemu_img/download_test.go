// +build medium

package img

import (
	"net/url"
	"testing"
)

func TestDownloadQcow2(t *testing.T) {
	p := "test.qcow2"

	i, err := OpenQemuImg(p)
	if err != nil {
		t.Fatalf("Cannot open '%s': err='%s'", p, err.Error())
	}
	if i.IsExists() {
		t.Errorf("Test environment is invalid, image is already existing: path='%s'", p)
	}

	u, err := url.Parse("http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img")
	if err := i.Download(u); err != nil {
		t.Errorf("Failed to download image: err='%s'", err.Error())
	}
	if !i.IsExists() {
		t.Errorf("Failed to download image: image is not existing yet")
	}

	if err := i.Delete(); err != nil {
		t.Errorf("Failed to delete image: err='%s'", err.Error())
	}
	if i.IsExists() {
		t.Errorf("Failed to delete image: image is existing yet")
	}
}

func TestDownloadISO(t *testing.T) {
	p := "mini.iso"

	i, err := OpenQemuImg(p)
	if err != nil {
		t.Fatalf("Cannot open '%s': err='%s'", p, err.Error())
	}
	if i.IsExists() {
		t.Errorf("Test environment is invalid, image is already existing: path='%s'", p)
	}

	u, err := url.Parse("http://archive.ubuntu.com/ubuntu/dists/bionic-updates/main/installer-amd64/current/images/netboot/mini.iso")
	if err := i.Download(u); err != nil {
		t.Errorf("Failed to download image: err='%s'", err.Error())
	}
	if !i.IsExists() {
		t.Errorf("Failed to download image: image is not existing yet")
	}

	if err := i.Delete(); err != nil {
		t.Errorf("Failed to delete image: err='%s'", err.Error())
	}
	if i.IsExists() {
		t.Errorf("Failed to delete image: image is existing yet")
	}
}
