package deploy

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

type LocalDeployer struct {
	targetDirectory string
}

func NewLocalDeployer(target string) (*LocalDeployer, error) {
	t, err := filepath.Abs(target)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get absolute path")
	}

	d := &LocalDeployer{
		targetDirectory: t,
	}

	if os.MkdirAll(t, 0755); err != nil {
		return nil, errors.Wrap(err, "Failed to create target directory")
	}

	return d, nil
}

// targetDirectory にファイルを転送し、 path にシンボリックリンクを貼る
func (d LocalDeployer) SaveFile(body []byte, path string, permission os.FileMode) error {
	var target string
	if !filepath.IsAbs(path) {
		target = filepath.Join(d.targetDirectory, path)
	} else {
		filename := filepath.Base(path)
		target = filepath.Join(d.targetDirectory, filename)

		if _, err := os.Lstat(path); err == nil {
			if err := os.Remove(path); err != nil {
				return errors.Wrap(err, "Failed to remove old symbolic link")
			}
		}

		if err := os.Symlink(target, path); err != nil {
			return errors.Wrap(err, "Failed to create symbolic link")
		}
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, permission)
	if err != nil {
		return errors.Wrap(err, "Failed to create new file")
	}
	defer file.Close()

	if _, err := file.Write(body); err != nil {
		return errors.Wrap(err, "Failed to write to file")
	}

	return nil
}

func (d LocalDeployer) InstallBinary(path string) (string, error) {
	binpath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", errors.Wrap(err, "Failed to get absolute path")
	}

	installPath := filepath.Join(path, filepath.Base(binpath))
	if err := os.Rename(binpath, installPath); err != nil {
		return "", errors.Wrapf(err, "Failed to move file: %s", installPath)
	}

	return installPath, nil
}

func (d LocalDeployer) Link(srcpath, dstpath string) error {
	self, err := filepath.Abs(srcpath)
	if err != nil {
		return errors.Wrap(err, "Failed to get self absolute path")
	}

	dst, err := filepath.Abs(dstpath)
	if err != nil {
		return errors.Wrap(err, "Failed to get destination absolute path")
	}

	if _, err := os.Lstat(dst); err == nil {
		if err := os.Remove(dst); err != nil {
			return errors.Wrap(err, "Failed to remove old symbolic link")
		}
	}

	if err := os.Symlink(self, dst); err != nil {
		return errors.Wrap(err, "Failed to create symbolic link")
	}

	return nil
}

func (d LocalDeployer) RestartDaemon(daemon string, stdout, stderr io.Writer) error {
	cmd := []string{
		"systemctl",
		"daemon-reload",
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		return errors.Wrapf(err, "Failed to command '%s'", strings.Join(cmd, " "))
	}

	cmd = []string{
		"systemctl",
		"restart",
		daemon,
	}
	c = exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		return errors.Wrapf(err, "Failed to command '%s'", strings.Join(cmd, " "))
	}

	cmd = []string{
		"systemctl",
		"enable",
		daemon,
	}
	c = exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		return errors.Wrapf(err, "Failed to command '%s'", strings.Join(cmd, " "))
	}

	return nil
}

func (d LocalDeployer) DaemonStatus(daemon string, stdout, stderr io.Writer) error {
	cmd := []string{
		"systemctl",
		"status",
		daemon,
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		return errors.Wrapf(err, "Failed to command '%s'", strings.Join(cmd, " "))
	}

	return nil
}

func (d LocalDeployer) StopDaemon(daemon string, stdout, stderr io.Writer) error {
	cmd := []string{
		"systemctl",
		"stop",
		daemon,
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			if status, ok := ee.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 5 { // Failed to stop n0core-agent.service: Unit n0core-agent.service not loaded.
					return errors.Wrapf(err, "Failed to command '%s'", strings.Join(cmd, " "))
				}
			}
		} else {
			return errors.Wrapf(err, "Failed to command '%s'", strings.Join(cmd, " "))
		}
	}

	return nil
}

func (d LocalDeployer) InstallPackages(packages []string, stdout, stderr io.Writer) error {
	{
		cmd := []string{
			"apt",
			"update",
		}

		c := exec.Command(cmd[0], cmd[1:]...)
		c.Stdout = stdout
		c.Stderr = stderr
		if err := c.Run(); err != nil {
			return errors.Wrap(err, "Failed to command 'apt'")
		}
	}

	{
		cmd := []string{
			"apt",
			"install",
			"-y",
		}
		cmd = append(cmd, packages...)

		c := exec.Command(cmd[0], cmd[1:]...)
		c.Stdout = stdout
		c.Stderr = stderr
		if err := c.Run(); err != nil {
			return errors.Wrap(err, "Failed to command 'apt'")
		}
	}

	return nil
}

func (d LocalDeployer) SetSysctl(key string, value []byte) error {
	p := filepath.Join("/proc/sys/", strings.Replace(key, ".", "/", -1))
	return ioutil.WriteFile(p, value, 0644)
}
