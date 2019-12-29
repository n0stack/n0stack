package n0deploy

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type SubcommandType int

const (
	SubcommandType_BOOTSTRAP SubcommandType = iota
	SubcommandType_DEPLOY
)

func Deploy(user, host, keyPath, n0deployFile string, bootstrap, deploy bool) error {
	n0depDir := filepath.Dir(n0deployFile)
	n0depFile := filepath.Base(n0deployFile)
	os.Chdir(n0depDir)

	buf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return errors.Wrap(err, "Failed to read key file")
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return errors.Wrap(err, "Failed to parse key")
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.SetDefaults()
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return errors.Wrap(err, "Failed to dial ssh")
	}
	defer conn.Close()

	// sudo permission

	b, err := ioutil.ReadFile(n0depFile)
	if err != nil {
		return errors.Wrapf(err, "Failed to read n0deploy file")
	}

	parser := NewSshParser(conn)
	n0dep, err := parser.Parse(string(b))
	if err != nil {
		return errors.Wrapf(err, "Failed to parse n0deploy file")
	}

	out := os.Stdout
	ctx := context.Background()
	inss := []Instruction{}
	if bootstrap {
		inss = append(inss, n0dep.Bootstrap...)
	}
	if deploy {
		inss = append(inss, n0dep.Deploy...)
	}

	for i, ins := range inss {
		fmt.Fprintf(out, "  [ Step %d/%d ] %s\n", i+1, len(inss), ins.String())

		if err := ins.Do(ctx, out); err != nil {
			return errors.Wrapf(err, "Failed to do instruction")
		}

		fmt.Fprintln(out, "")
	}

	return nil
}
