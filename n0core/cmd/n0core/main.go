package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var version = "undefined"

func main() {
	app := cli.NewApp()
	app.Name = "n0core"
	app.Version = version
	app.Usage = "The n0stack cluster manager"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Serve daemons",
			Subcommands: []cli.Command{
				{
					Name:   "api",
					Usage:  "Daemon which provide n0stack cluster API",
					Action: ServeAPI,
					Flags: []cli.Flag{
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
							Name: "etcd-endpoints",
						},
						cli.StringFlag{
							Name: "token-secret",
						},
						cli.StringFlag{
							Name: "listen-url",
						},
					},
				},
				{
					Name:   "bff",
					Usage:  "Daemon which provide bff for n0stack API",
					Action: ServeBFF,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "api-url",
						},
						cli.StringFlag{
							Name:  "listen-address",
							Value: "0.0.0.0:8080",
						},
					},
				},
				// {
				// 	Name:   "agent",
				// 	Usage:  "Daemon which administrate n0stack cluster node",
				// 	Action: ServeAgent,
				// 	Flags: []cli.Flag{
				// 		cli.StringFlag{
				// 			Name:  "name",
				// 			Value: hostname,
				// 		},
				// 		cli.StringFlag{
				// 			// interfaceからも取れるようにしたい
				// 			Name: "advertise-address",
				// 		},
				// 		cli.StringFlag{
				// 			Name: "node-api-endpoint",
				// 		},
				// 		cli.StringFlag{
				// 			// interfaceからも取れるようにしたい
				// 			Name:  "bind-address",
				// 			Value: "0.0.0.0",
				// 		},
				// 		cli.IntFlag{
				// 			Name:  "bind-port",
				// 			Value: 20181,
				// 		},
				// 		cli.StringFlag{
				// 			Name:  "base-directory",
				// 			Value: "/var/lib/n0core",
				// 		},
				// 		cli.StringFlag{
				// 			Name:  "location",
				// 			Usage: "<Datacenter>/<AvailavilityZone>/<Cell>/<Rack>/<Unit(int)>",
				// 		},
				// 		cli.UintFlag{
				// 			Name:  "cpu-capacity-milli-cores",
				// 			Value: uint(node.GetTotalCPUMilliCores()) * 1000,
				// 		},
				// 		cli.Uint64Flag{
				// 			Name:  "memory-capacity-bytes",
				// 			Value: node.GetTotalMemory(),
				// 		},
				// 		cli.Uint64Flag{
				// 			Name:  "storage-capacity-bytes",
				// 			Value: uint64(100 * bytefmt.GIGABYTE),
				// 		},
				// 	},
				// },
			},
		},
	}

	log.SetFlags(log.Llongfile | log.Ltime | log.Lmicroseconds)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start process, err:%s\n", err.Error())
		os.Exit(1)
	}
}
