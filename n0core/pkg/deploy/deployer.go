package deploy

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Deployer struct {
	ssh             *ssh.Client
	targetDirectory string
}

func NewDeployer(s *ssh.Client, target string) (*Deployer, error) {
	d := &Deployer{
		ssh:             s,
		targetDirectory: target,
	}

	client, err := sftp.NewClient(d.ssh)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create new sftp client")
	}
	defer client.Close()

	if err := client.MkdirAll(d.targetDirectory); err != nil {
		return nil, errors.Wrap(err, "Failed to create target directory")
	}

	return d, nil
}

// targetDirectory にファイルを転送し、 path にシンボリックリンクを貼る
func (d Deployer) SendFile(body []byte, path string, permission os.FileMode) error {
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

		if err := client.Symlink(target, path); err != nil {
			return errors.Wrap(err, "Failed to create symbolic link")
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

func (d Deployer) ReadSelf() ([]byte, error) {
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get absolute path")
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

func (d Deployer) RestartDaemon(daemon string) error {
	sess, err := d.ssh.NewSession()
	if err != nil {
		return errors.Wrap(err, "Failed to create new session")
	}
	defer sess.Close()

	reload := "systemctl daemon-reload"
	if out, err := sess.CombinedOutput(reload); err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			return fmt.Errorf("'daemon-reload' exit status is not 0: code=%d, output=\n%s", ee.ExitStatus(), string(out))
		}

		return errors.Wrapf(err, "Failed to command '%s'", reload)
	}

	restart := fmt.Sprintf("systemctl restart %s", daemon)
	if out, err := sess.CombinedOutput(restart); err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			return fmt.Errorf("'restart' exit status is not 0: code=%d, output=\n%s", ee.ExitStatus(), string(out))
		}

		return errors.Wrapf(err, "Failed to command '%s'", restart)
	}

	return nil
}

// func (d Deployer) Close() (error) {}
