package iso

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type CloudConfig struct {
	isoPath string
	cfgPath string
}

func (c *CloudConfig) StructConfig(path, user string, key ssh.PublicKey) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open hostsfile: path='%s', err='%s'", path, err.Error())
	}
	defer f.Close()

	u := []struct {
		Name              string   `yaml:"name"`
		Groups            string   `yaml:"groups"`
		Shell             string   `yaml:"shell"`
		Sudo              []string `yaml:"sudo"`
		SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys"`
	}{
		{
			user,
			"sudo",
			"/bin/bash",
			[]string{
				"ALL=(ALL) NOPASSWD:ALL",
			},
			[]string{
				string(ssh.MarshalAuthorizedKey(key)),
			},
		},
	}

	buf, err := yaml.Marshal(struct {
		Users interface{} `yaml:"users"`
	}{
		u,
	})
	if err != nil {
		return fmt.Errorf("Failed to marshal for yaml: err='%s'", err.Error())
	}

	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf("Failed to write file: err='%s'", err.Error())
	}

	c.cfgPath = path
	return nil
}

func (c *CloudConfig) GenerateISO(path string) error {
	args := []string{
		"/usr/bin/cloud-localds",
		path,
		c.cfgPath,
	}

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to start process: args='%s', out='%s', err='%s'", args, string(out), err.Error())
	}

	c.isoPath = path
	return nil
}

func (c *CloudConfig) Delete() (err error) {
	if err := os.Remove(c.cfgPath); err != nil {
		return fmt.Errorf("Failed to delete cloud-config: path='%s', err='%s'", c.cfgPath, err.Error())
	}

	c.cfgPath = ""
	return nil
}
