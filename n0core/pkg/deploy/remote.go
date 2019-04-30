package deploy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type RemoteDeployer struct {
	ssh             *ssh.Client
	useSudo         bool
	password        string
	targetDirectory string
}

func NewRemoteDeployer(s *ssh.Client, target string) (*RemoteDeployer, error) {
	d := &RemoteDeployer{
		ssh:             s,
		targetDirectory: target,
	}

	err := d.CheckPriv()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to check privilege")
	}

	client, err := sftp.NewClient(d.ssh)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create new sftp client")
	}
	defer client.Close()

	if err := client.MkdirAll(d.targetDirectory); err != nil {
		return nil, errors.Wrap(err, "Failed to make target directory")
	}

	return d, nil
}

// targetDirectory にファイルを転送し、 path にシンボリックリンクを貼る
func (d RemoteDeployer) SendFile(body []byte, path string, permission os.FileMode) error {
	client, err := sftp.NewClient(d.ssh)
	if err != nil {
		return errors.Wrap(err, "Failed to create new sftp client")
	}
	defer client.Close()

	var target string
	if !filepath.IsAbs(path) {
		target = filepath.Join(d.targetDirectory, path)
	} else {
		filename := filepath.Base(path)
		target = filepath.Join(d.targetDirectory, filename)

		if _, err := client.Lstat(path); err == nil {
			if err := client.Remove(path); err != nil {
				return errors.Wrap(err, "Failed to remove old symbolic link")
			}
		}

		if err := client.Symlink(target, path); err != nil {
			return errors.Wrap(err, "Failed to create symbolic link")
		}
	}

	if _, err := client.Stat(target); err == nil {
		if err := client.Remove(target); err != nil {
			return errors.Wrap(err, "Failed to delete target")
		}
	}

	file, err := client.Create(target)
	if err != nil {
		return errors.Wrap(err, "Failed to create new file")
	}
	defer file.Close()

	if err := file.Chmod(permission); err != nil {
		return errors.Wrap(err, "Failed to change permission")
	}

	if _, err := file.Write(body); err != nil {
		return errors.Wrap(err, "Failed to write to file")
	}

	return nil
}

func (d RemoteDeployer) ReadSelf() ([]byte, error) {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get absolute path")
	}

	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("Use binary with absolute path to read self")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open self")
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	return buf.Bytes(), nil
}

func (d RemoteDeployer) Command(command string, stdout, stderr io.Writer) error {
	sess, err := d.ssh.NewSession()
	if err != nil {
		return errors.Wrap(err, "Failed to create new session")
	}
	defer sess.Close()

	sess.Stdout = stdout
	sess.Stderr = stderr

	if d.useSudo {
		if len(d.password) > 0 {
			command = "echo " + d.password + " | sudo -S " + command
		} else {
			command = "sudo " + command
		}
	}

	if err := sess.Run(command); err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			return fmt.Errorf("'%s' exit status is not 0: code=%d", strings.Replace(command, d.password, "**CENSORED**", -1), ee.ExitStatus())
		}

		return errors.Wrapf(err, "Failed to command '%s'", command)
	}

	return nil
}

func (d *RemoteDeployer) CheckPriv() error {
	isRemoteRoot, err := d.checkRemoteUserIsRoot()
	if err != nil {
		return err
	}

	if isRemoteRoot {
		return nil
	}

	d.useSudo = true

	ok, err := d.trySudoWithoutPassword()
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	for cnt := 0; cnt < 3; cnt++ {
		sess, err := d.ssh.NewSession()
		if err != nil {
			return fmt.Errorf("Failed to create new session")
		}
		defer sess.Close()

		fmt.Print("[sudo] input password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("Failed to read password")
		}
		fmt.Println("checking password...")

		err = sess.Run("echo " + string(password) + " | sudo -S id -u")
		if err != nil {
			continue
		}
		d.password = string(password)
		return nil
	}

	return fmt.Errorf("Wrong password")
}

func (d *RemoteDeployer) checkRemoteUserIsRoot() (bool, error) {
	sess, err := d.ssh.NewSession()
	if err != nil {
		return false, fmt.Errorf("Failed to create new session")
	}
	defer sess.Close()

	out, err := sess.Output("id -u")
	if err != nil {
		return false, fmt.Errorf("Failed to execute `id -u`")
	}
	uid, err := strconv.Atoi(strings.TrimRight(string(out), "\n"))
	if err != nil {
		return false, fmt.Errorf("Failed to convert to uid from: %s", out)
	}

	return (uid == 0), nil
}

func (d *RemoteDeployer) trySudoWithoutPassword() (bool, error) {
	sess, err := d.ssh.NewSession()
	if err != nil {
		return false, fmt.Errorf("Failed to create new session")
	}
	defer sess.Close()

	out, err := sess.Output("sudo -n echo -n ok")
	if err != nil {
		return false, err
	}

	return string(out) == "ok", nil
}
