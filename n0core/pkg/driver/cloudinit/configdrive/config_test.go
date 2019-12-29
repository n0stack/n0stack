package configdrive

import (
	"io/ioutil"
	"net"
	"strings"
	"testing"

	img "n0st.ac/n0stack/n0core/pkg/driver/qemu_img"
	netutil "n0st.ac/n0stack/n0core/pkg/util/net"

	"golang.org/x/crypto/ssh"
)

func TestGenerateISO(t *testing.T) {
	key, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBITowPn2Ol1eCvXN5XV+Lb6jfXzgDbXyEdtayadDUJtFrcN2m2mjC1B20VBAoJcZtSYkmjrllS06Q26Te5sTYvE= testkey"))
	i := netutil.ParseCIDR("192.168.122.10/24")
	hw, _ := net.ParseMAC("52:54:00:78:71:f1")
	e := &CloudConfigEthernet{
		MacAddress: hw,
		Address4:   i,
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

func TestTRUNCFile(t *testing.T) {
	key, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBITowPn2Ol1eCvXN5XV+Lb6jfXzgDbXyEdtayadDUJtFrcN2m2mjC1B20VBAoJcZtSYkmjrllS06Q26Te5sTYvE= testkey"))
	i := netutil.ParseCIDR("192.168.122.10/24")
	hw, _ := net.ParseMAC("52:54:00:78:71:f1")
	e := &CloudConfigEthernet{
		MacAddress: hw,
		Address4:   i,
		Gateway4:   net.ParseIP("192.168.122.1"),
		NameServers: []net.IP{
			net.ParseIP("192.168.122.1"),
		},
	}
	c := StructConfig("user", "host", []ssh.PublicKey{key}, []*CloudConfigEthernet{e})

	padding := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	err := ioutil.WriteFile("meta-data", []byte(padding), 0664)
	if err != nil {
		t.Fatalf("Failed to write precondition: err='%s'", err.Error())
	}

	if err := c.GenerateMetadataFile("."); err != nil {
		t.Fatalf("Failed to generate metadata file: err='%s'", err.Error())
	}
	defer c.Delete()

	data, _ := ioutil.ReadFile("meta-data")
	if len(data) != 39 {
		t.Errorf("GenerateMetadataFile was wrong")
	}
}

func TestGenerateNetworkConfigFileAboutIPAddressIsMissing(t *testing.T) {
	key, _, _, _, _ := ssh.ParseAuthorizedKey([]byte("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBITowPn2Ol1eCvXN5XV+Lb6jfXzgDbXyEdtayadDUJtFrcN2m2mjC1B20VBAoJcZtSYkmjrllS06Q26Te5sTYvE= testkey"))
	hw, _ := net.ParseMAC("52:54:00:78:71:f1")
	e := &CloudConfigEthernet{
		MacAddress: hw,
	}
	c := StructConfig("user", "host", []ssh.PublicKey{key}, []*CloudConfigEthernet{e})

	if err := c.GenerateNetworkConfigFile("."); err != nil {
		t.Fatalf("Failed to generate network config file: err='%s'", err.Error())
	}
	defer c.Delete()

	data, _ := ioutil.ReadFile("network-config")
	if strings.Contains(string(data), "null") {
		t.Errorf("Failed to valid network file, subnets have null")
	}
}
