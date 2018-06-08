package main

import (
	"os"
	"time"

	"github.com/n0stack/n0core/provisioning/node"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "n0core Agent"
	// app.Usage = ""
	app.Version = "0.1.0" // CIで取るようにする

	// command action
	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Join to API and serve some daemons.",
			Action: func(c *cli.Context) error {
				if err := node.JoinNode(
					c.String("name"),
					c.String("advertise-address"),
					c.String("api-address"),
				); err != nil {
					return err
				}

				time.Sleep(600 * time.Second)

				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "name",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name: "advertise-address",
				},
				cli.StringFlag{
					Name: "api-address",
				},
			},
		},
	}
	app.Run(os.Args)
}
