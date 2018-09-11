package iso

import (
	"testing"

	"github.com/n0stack/n0core/pkg/driver/qemu_img"

	"golang.org/x/crypto/ssh"
)

func TestGenerateISO(t *testing.T) {
	c := &CloudConfig{}

	p, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBMwfns/Woug7RBC3YjRbXecM/4LGVb5S1u8muo2v3ealA5mX2fnAKUssQY84XXE08m18nnWClVWU2goUsVtlgp0= testkey"))
	if err := c.StructConfig("test.cfg", "hoge", p); err != nil {
		t.Fatalf("Failed to struct config: err='%s'", err.Error())
	}

	f := "test.iso"
	err := c.GenerateISO(f)
	if err != nil {
		t.Errorf("Failed to generate iso: err='%s'", err.Error())
	}

	i, err := img.OpenQemuImg(f)
	if err != nil {
		t.Errorf("Failed to open qemu image, maybe this is not valid image: path='%s', err='%s'", f, err.Error())
	}
	defer i.Delete()

	if err := c.Delete(); err != nil {
		t.Errorf("Failed to delete: path='%s', err='%s'", c.cfgPath, err.Error())
	}
}
