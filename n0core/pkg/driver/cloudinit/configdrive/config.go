package configdrive

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type cloudConfigUser struct {
	Name              string   `yaml:"name"`
	Groups            string   `yaml:"groups"`
	Shell             string   `yaml:"shell"`
	Sudo              []string `yaml:"sudo"`
	SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys,omitempty"`
	LockPasswd        bool     `yaml:"lock_passwd"`
	// Passwd            string   `yaml:"passwd"` // mkpasswd --method=SHA-512 --rounds=4096
}

type cloudConfigNameServer struct {
	Addressees []net.IP `yaml:"addresses"`
}

type cloudNetworkSubnet struct {
	Type       string   `yaml:"type"`
	Address    net.IP   `yaml:"address"`
	NetMask    net.IP   `yaml:"netmask"`
	Gateway    net.IP   `yaml:"gateway,omitempty"`
	DnsServers []net.IP `yaml:"dns_nameservers,omitempty"`
	DnsSearch  []string `yaml:"dns_search,omitempty"`
}

type cloudNetworkConfig struct {
	Type       string                `yaml:"type"`
	Name       string                `yaml:"name"`
	MacAddress string                `yaml:"mac_address"`
	Subnets    []*cloudNetworkSubnet `yaml:"subnets,omitempty"`
}

type CloudConfig struct {
	Metadata struct {
		InstanceID    string `yaml:"instance-id"`
		LocalHostname string `yaml:"local-hostname"`
	}
	Userdata struct {
		Users []cloudConfigUser `yaml:"users"`
	}
	NetworkConfig struct {
		Version int                   `yaml:"version"`
		Config  []*cloudNetworkConfig `yaml:"config"`
	}

	isoPath           string
	metadataPath      string
	userdataPath      string
	networkConfigPath string
}

type CloudConfigEthernet struct {
	MacAddress net.HardwareAddr
	Address4   net.IP
	Network4   *net.IPNet
	Gateway4   net.IP
	// Address6    []net.IP
	// Mask6 int
	// Gateway6    net.IP
	NameServers []net.IP
	// NameSearch  []string
}

func StructConfig(user, hostname string, keys []ssh.PublicKey, eth []*CloudConfigEthernet) *CloudConfig {
	c := &CloudConfig{}

	ks := make([]string, len(keys))
	for i, k := range keys {
		ks[i] = string(ssh.MarshalAuthorizedKey(k))
	}
	c.Userdata.Users = []cloudConfigUser{
		{
			user,
			"sudo",
			"/bin/bash",
			[]string{
				"ALL=(ALL) NOPASSWD:ALL",
			},
			ks,
			true,
		},
	}

	c.Metadata.InstanceID = hostname
	c.Metadata.LocalHostname = hostname

	c.NetworkConfig.Version = 1
	c.NetworkConfig.Config = make([]*cloudNetworkConfig, len(eth))
	for i, e := range eth {
		c.NetworkConfig.Config[i] = &cloudNetworkConfig{
			Type:       "physical",
			Name:       fmt.Sprintf("eth%d", i),
			MacAddress: e.MacAddress.String(),
			Subnets:    make([]*cloudNetworkSubnet, 0),
		}

		c.NetworkConfig.Config[i].Subnets = append(c.NetworkConfig.Config[i].Subnets, &cloudNetworkSubnet{
			Type:       "static",
			Address:    e.Address4,
			NetMask:    net.IPv4(e.Network4.Mask[0], e.Network4.Mask[1], e.Network4.Mask[2], e.Network4.Mask[3]),
			Gateway:    e.Gateway4,
			DnsServers: e.NameServers,
		})
	}

	return c
}

func (c *CloudConfig) GenerateUserdataFile(basedir string) error {
	c.userdataPath = filepath.Join(basedir, "user-data")

	f, err := os.OpenFile(c.userdataPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open hostsfile: path='%s', err='%s'", c.userdataPath, err.Error())
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintln("#cloud-config")); err != nil {
		return fmt.Errorf("Failed to write file: err='%s'", err.Error())
	}

	buf, err := yaml.Marshal(c.Userdata)
	if err != nil {
		return fmt.Errorf("Failed to marshal for yaml: err='%s'", err.Error())
	}

	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf("Failed to write file: err='%s'", err.Error())
	}

	return nil
}

func (c *CloudConfig) GenerateMetadataFile(basedir string) error {
	c.metadataPath = filepath.Join(basedir, "meta-data")

	f, err := os.OpenFile(c.metadataPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open hostsfile: path='%s', err='%s'", c.metadataPath, err.Error())
	}
	defer f.Close()

	buf, err := yaml.Marshal(c.Metadata)
	if err != nil {
		return fmt.Errorf("Failed to marshal for yaml: err='%s'", err.Error())
	}

	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf("Failed to write file: err='%s'", err.Error())
	}

	return nil
}

func (c *CloudConfig) GenerateNetworkConfigFile(basedir string) error {
	c.networkConfigPath = filepath.Join(basedir, "network-config")

	f, err := os.OpenFile(c.networkConfigPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open hostsfile: path='%s', err='%s'", c.networkConfigPath, err.Error())
	}
	defer f.Close()

	buf, err := yaml.Marshal(c.NetworkConfig)
	if err != nil {
		return fmt.Errorf("Failed to marshal for yaml: err='%s'", err.Error())
	}

	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf("Failed to write file: err='%s'", err.Error())
	}

	return nil
}

// Generate ISO and yaml files
func (c *CloudConfig) Generate(basedir string) (string, error) {
	c.isoPath = filepath.Join(basedir, "cloudconfig.iso")

	if err := c.GenerateMetadataFile(basedir); err != nil {
		return "", err
	}

	if err := c.GenerateUserdataFile(basedir); err != nil {
		return "", err
	}

	if err := c.GenerateNetworkConfigFile(basedir); err != nil {
		return "", err
	}

	args := []string{
		"/usr/bin/genisoimage",
		"-output",
		c.isoPath,
		"-volid",
		"cidata",
		"-joliet",
		"-rock",
		c.userdataPath,
		c.metadataPath,
		c.networkConfigPath,
	}

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed to start process: args='%s', out='%s', err='%s'", args, string(out), err.Error())
	}

	return c.isoPath, nil
}

// Clean generated yaml files
func (c *CloudConfig) Clean() error {
	if c.metadataPath != "" {
		if err := os.Remove(c.metadataPath); err != nil {
			return fmt.Errorf("Failed to delete generated metadata file: path='%s', err='%s'", c.metadataPath, err.Error())
		}
		c.metadataPath = ""
	}

	if c.userdataPath != "" {
		if err := os.Remove(c.userdataPath); err != nil {
			return fmt.Errorf("Failed to delete generated userdata file: path='%s', err='%s'", c.userdataPath, err.Error())
		}
		c.userdataPath = ""
	}

	if c.networkConfigPath != "" {
		if err := os.Remove(c.networkConfigPath); err != nil {
			return fmt.Errorf("Failed to delete generated userdata file: path='%s', err='%s'", c.networkConfigPath, err.Error())
		}
		c.networkConfigPath = ""
	}

	return nil
}

// Delete generated the iso file
func (c *CloudConfig) Delete() error {
	if err := c.Clean(); err != nil {
		return err
	}

	if c.isoPath != "" {
		if err := os.Remove(c.isoPath); err != nil {
			return fmt.Errorf("Failed to delete generated iso: path='%s', err='%s'", c.isoPath, err.Error())
		}
		c.isoPath = ""
	}

	return nil
}
