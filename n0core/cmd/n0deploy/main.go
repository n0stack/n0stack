package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/n0stack/n0stack/n0core/pkg/driver/n0deploy"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
)

var version = "undefined"

const n0deployFilename = "n0deployfile"

type SubcommandType int

const (
	SubcommandType_BOOTSTRAP SubcommandType = iota
	SubcommandType_DEPLOY
)

func WrapAction(t SubcommandType) func(cctx *cli.Context) error {
	return func(cctx *cli.Context) error {
		if cctx.NArg() != 1 {
			return fmt.Errorf("Argument usage: %s", cctx.Command.ArgsUsage)
		}
		h := strings.Split(cctx.Args()[0], "@")
		if len(h) != 2 {
			return fmt.Errorf("Argument usage: %s", cctx.Command.ArgsUsage)
		}
		user := h[0]
		host := h[1]

		keyPath := cctx.GlobalString("identity-file")

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
		conn, err := ssh.Dial("tcp", host+":22", config)
		if err != nil {
			return errors.Wrap(err, "Failed to dial ssh")
		}
		defer conn.Close()

		// sudo permission

		b, err := ioutil.ReadFile(n0deployFilename)
		if err != nil {
			return errors.Wrapf(err, "Failed to read n0deployfile")
		}

		parser := n0deploy.NewSshParser(conn)
		n0dep, err := parser.Parse(string(b))
		if err != nil {
			return errors.Wrapf(err, "Failed to parse n0deployfile")
		}

		out := os.Stdout
		ctx := context.Background()
		inss := []n0deploy.Instruction{}
		switch t {
		case SubcommandType_BOOTSTRAP:
			inss = n0dep.Bootstrap
		case SubcommandType_DEPLOY:
			inss = n0dep.Deploy
		}

		for i, ins := range inss {
			fmt.Fprintf(out, ">>> Step %d/%d: %s\n", i+1, len(n0dep.Bootstrap), ins.String())

			if err := ins.Do(ctx, out); err != nil {
				return errors.Wrapf(err, "Failed to do instruction")
			}

			fmt.Fprintln(out, "")
		}

		return nil
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "n0deploy"
	app.Version = version
	app.Usage = "The n0stack deployment tool"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "identity-file, i",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "bootstrap",
			Usage:     "Bootstrap",
			ArgsUsage: "[user]@[hostname]",
			Action:    WrapAction(SubcommandType_BOOTSTRAP),
		},
		{
			Name:      "deploy",
			Usage:     "Deploy",
			ArgsUsage: "[user]@[hostname]",
			Action:    WrapAction(SubcommandType_DEPLOY),
		},
	}

	log.SetFlags(log.Llongfile | log.Ltime | log.Lmicroseconds)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
