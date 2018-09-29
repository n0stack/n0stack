package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "n0core Agent"
	// app.Usage = ""
	app.Version = "0.1.0" // CIで取るようにする

	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "",
			Action: Serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "name",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name: "advertise-address",
				},
				cli.StringFlag{
					Name: "node-api-endpoint",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name:  "bind-address",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "bind-port",
					Value: 20180,
				},
				cli.StringFlag{
					Name:  "base-directory",
					Value: "/var/lib/n0core",
				},
			},
		},
		{
			Name:   "mock",
			Usage:  "",
			Action: Mock,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "name",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name: "advertise-address",
				},
				cli.StringFlag{
					Name: "node-api-endpoint",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name:  "bind-address",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "bind-port",
					Value: 20180,
				},
			},
		},
	}

	log.SetFlags(log.Lshortfile)

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start process, err:%v", err.Error())
	}
}
