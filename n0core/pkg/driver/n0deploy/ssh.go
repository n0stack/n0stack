package n0deploy

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func NewSshParser(s *ssh.Client) *Parser {
	return &Parser{
		NewCopyInstruction: NewSshCopyInstruction(s),
		NewRunInstruction:  NewSshRunInstruction(s),
	}
}

type SshRunInstruction struct {
	command   string
	sshClinet *ssh.Client
}

func NewSshRunInstruction(s *ssh.Client) func(string) (Instruction, error) {
	return func(line string) (Instruction, error) {
		return &SshRunInstruction{
			command:   strings.TrimPrefix(line, "RUN "),
			sshClinet: s,
		}, nil
	}
}

func (i SshRunInstruction) Do(ctx context.Context, out io.Writer) error {
	sess, err := i.sshClinet.NewSession()
	if err != nil {
		return errors.Wrap(err, "Failed to create new session")
	}
	defer sess.Close()

	sess.Stdout = out
	sess.Stderr = out

	if err := sess.Run(i.command); err != nil { // sh -c が必要かどうかわからない
		if ee, ok := err.(*ssh.ExitError); ok {
			return fmt.Errorf("'%s' exit status is not 0: code=%d", i.command, ee.ExitStatus())
		}

		return errors.Wrapf(err, "Failed to command '%s'", i.command)
	}

	return nil
}

func (i SshRunInstruction) String() string {
	return fmt.Sprintf("RUN %s", i.command)
}

type SshCopyInstruction struct {
	src       string
	dst       string
	sshClinet *ssh.Client
}

func NewSshCopyInstruction(s *ssh.Client) func(string, string) (Instruction, error) {
	return func(src, dst string) (Instruction, error) {
		return &SshCopyInstruction{
			src:       src,
			dst:       dst,
			sshClinet: s,
		}, nil
	}
}

func (i SshCopyInstruction) Do(ctx context.Context, out io.Writer) error {
	client, err := sftp.NewClient(i.sshClinet)
	if err != nil {
		return errors.Wrap(err, "Failed to create new sftp client")
	}
	defer client.Close()

	srcFile, err := os.Open(i.src)
	if err != nil {
		return errors.Wrap(err, "Failed to open localfile")
	}
	defer srcFile.Close()

	// ディレクトリの処理がうまくできていない
	dstFile, err := client.Create(i.dst)
	if err != nil {
		return errors.Wrap(err, "Failed to create remote file")
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Wrap(err, "Failed to copy file")
	}

	return nil
}

func (i SshCopyInstruction) String() string {
	return fmt.Sprintf("COPY %s %s", i.src, i.dst)
}
