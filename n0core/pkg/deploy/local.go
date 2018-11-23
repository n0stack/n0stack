package deploy

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

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

		if err := os.Symlink(target, path); err != nil {
			return errors.Wrap(err, "Failed to create symbolic link")
		}
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, permission)
	if err != nil {
		return errors.Wrap(err, "Failed to create new file")
	}
	defer file.Close()

	if _, err := file.Write(body); err != nil {
		return errors.Wrap(err, "Failed to write to file")
	}

	return nil
}

func (d LocalDeployer) RestartDaemon(daemon string) error {
	// sess, err := d.ssh.NewSession()
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to create new session")
	// }
	// defer sess.Close()

	// reload := "systemctl daemon-reload"
	// if out, err := sess.CombinedOutput(reload); err != nil {
	// 	if ee, ok := err.(*ssh.ExitError); ok {
	// 		return fmt.Errorf("'daemon-reload' exit status is not 0: code=%d, output=\n%s", ee.ExitStatus(), string(out))
	// 	}

	// 	return errors.Wrapf(err, "Failed to command '%s'", reload)
	// }

	// restart := fmt.Sprintf("systemctl restart %s", daemon)
	// if out, err := sess.CombinedOutput(restart); err != nil {
	// 	if ee, ok := err.(*ssh.ExitError); ok {
	// 		return fmt.Errorf("'restart' exit status is not 0: code=%d, output=\n%s", ee.ExitStatus(), string(out))
	// 	}

	// 	return errors.Wrapf(err, "Failed to command '%s'", restart)
	// }

	return nil
}

func (d LocalDeployer) InstallPackages(packages []string, stdout, stderr io.Writer) error {
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

	return nil
}
