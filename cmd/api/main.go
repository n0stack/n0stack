package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "n0core API"
	// app.Usage = ""
	app.Version = "0.1.0" // CIで取るようにする

	app.Commands = []cli.Command{
		{
			Name:   "node",
			Usage:  "",
			Action: ServeNodeAPI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "etcd-endpoints",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name:  "bind-address",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "bind-port",
					Value: 20181,
				},
			},
		},
		{
			Name:   "network",
			Usage:  "",
			Action: ServeNetworkAPI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "etcd-endpoints",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name:  "bind-address",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "bind-port",
					Value: 20182,
				},
			},
		},
		{
			Name:   "volume",
			Usage:  "",
			Action: ServeVolumeAPI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "etcd-endpoints",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name:  "bind-address",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "bind-port",
					Value: 20183,
				},
				cli.StringFlag{
					Name: "node-api-endpoint",
				},
			},
		},
		{
			Name:   "virtual_machine",
			Usage:  "",
			Action: ServeNetworkAPI,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "etcd-endpoints",
				},
				cli.StringFlag{
					// interfaceからも取れるようにしたい
					Name:  "bind-address",
					Value: "0.0.0.0",
				},
				cli.IntFlag{
					Name:  "bind-port",
					Value: 20184,
				},
				cli.StringFlag{
					Name: "node-api-endpoint",
				},
				cli.StringFlag{
					Name: "network-api-endpoint",
				},
				cli.StringFlag{
					Name: "volume-api-endpoint",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start process, err:%v", err.Error())
	}
}
