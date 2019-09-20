// +build ignore
// +build medium

package configdrive

import (
	"net/url"
	"path/filepath"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"n0st.ac/n0stack/n0core/pkg/driver/qemu"
	img "n0st.ac/n0stack/n0core/pkg/driver/qemu_img"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
)

func TestBootUbuntu(t *testing.T) {
	id, _ := uuid.FromString("5fd6c569-172f-4b25-84cd-b76cc336cfdd")
	q, err := qemu.OpenQemu(&id)
	if err != nil {
		t.Fatalf("Failed to open qemu: err='%s'", err.Error())
	}
	defer q.Delete()

	b, _ := bytefmt.ToBytes("512M")
	if err := q.Start("test", "monitor.sock", 10000, 1, b); err != nil {
		t.Fatalf("Failed to start process: err='%s'", err.Error())
	}

	f := "../../../../.cache/ubuntu.qcow2"
	i, err := img.OpenQemuImg(f)
	if err != nil {
		t.Fatalf("Failed to open qemu-img, do not relate to this package code: err='%s'", err.Error())
	}

	u, _ := url.Parse("https://cloud-images.ubuntu.com/xenial/current/xenial-server-cloudimg-amd64-disk1.img")
	if err := i.Download(u); err != nil {
		t.Fatalf("Failed to download image, do not relate to this package code: url='%s', err='%s'", u.String(), err.Error())
	}
	defer i.Delete()

	p, _ := filepath.Abs(f)
	u, _ = url.Parse("file://" + p)
	if err := q.AttachQcow2("ubuntu", u, 0); err != nil {
		t.Errorf("Failed to attach qcow2: err='%s'", err.Error())
	}

	c := &CloudConfig{}

	key, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBITowPn2Ol1eCvXN5XV+Lb6jfXzgDbXyEdtayadDUJtFrcN2m2mjC1B20VBAoJcZtSYkmjrllS06Q26Te5sTYvE= testkey"))
	if err := c.StructConfig("test.cfg", "test", key); err != nil {
		t.Fatalf("Failed to struct config: err='%s'", err.Error())
	}
	defer c.Delete()

	isof := "test.iso"
	iso, err := c.GenerateISO(isof)
	if err != nil {
		t.Errorf("Failed to generate iso: err='%s'", err.Error())
	}
	defer iso.Delete()

	p, _ = filepath.Abs(isof)
	u, _ = url.Parse("file://" + p)
	if err := q.AttachISO("cloudcfg", u, 1); err != nil {
		t.Errorf("Failed to attach iso: err='%s'", err.Error())
	}

	if err := q.Boot(); err != nil {
		t.Errorf("Failed to boot: err='%s'", err.Error())
	}

	s, err := q.Status()
	if err != nil {
		t.Errorf("Failed to get status: err='%s'", err.Error())
	}
	if s != qemu.StatusRunning {
		t.Errorf("Status is mismatch: want='%v', have='%v'", qemu.StatusRunning, s)
	}
}
