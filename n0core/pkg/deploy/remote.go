package deploy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RemoteDeployer struct {
	ssh             *ssh.Client
	targetDirectory string
}

func NewRemoteDeployer(s *ssh.Client, target string) (*RemoteDeployer, error) {
	d := &RemoteDeployer{
		ssh:             s,
		targetDirectory: target,
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

	if err := sess.Run(command); err != nil {
		if ee, ok := err.(*ssh.ExitError); ok {
			return fmt.Errorf("'%s' exit status is not 0: code=%d", command, ee.ExitStatus())
		}

		return errors.Wrapf(err, "Failed to command '%s'", command)
	}

	return nil
}
