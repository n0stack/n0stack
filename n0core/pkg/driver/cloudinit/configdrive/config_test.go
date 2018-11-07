package configdrive

import (
	"net"
	"testing"

	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu_img"

	"golang.org/x/crypto/ssh"
)

func TestGenerateISO(t *testing.T) {
	key, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBITowPn2Ol1eCvXN5XV+Lb6jfXzgDbXyEdtayadDUJtFrcN2m2mjC1B20VBAoJcZtSYkmjrllS06Q26Te5sTYvE= testkey"))
	i, n, _ := net.ParseCIDR("192.168.122.10/24")
	hw, _ := net.ParseMAC("52:54:00:78:71:f1")
	e := &CloudConfigEthernet{
		MacAddress: hw,
		Address4:   i,
		Network4:   n,
		Gateway4:   net.ParseIP("192.168.122.1"),
		NameServers: []net.IP{
			net.ParseIP("192.168.122.1"),
		},
	}
	c := StructConfig("user", "host", []ssh.PublicKey{key}, []*CloudConfigEthernet{e})

	f, err := c.Generate(".")
	if err != nil {
		t.Errorf("Failed to generate iso: err='%s'", err.Error())
	}

	_, err = img.OpenQemuImg(f)
	if err != nil {
		t.Errorf("Failed to open qemu image, maybe this is not valid image: path='%s', err='%s'", f, err.Error())
	}

	if err := c.Delete(); err != nil {
		t.Errorf("Failed to delete: err='%s'", err.Error())
	}
}
