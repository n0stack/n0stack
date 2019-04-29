package main

import (
	"fmt"
	"log"
	"os"

	"code.cloudfoundry.org/bytefmt"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/urfave/cli"
)

var version = "undefined"

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname")
	}

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
							Name:  "base-directory",
							Value: "/var/lib/n0core",
						},
					},
				},
				{
					Name:   "agent",
					Usage:  "Daemon which administrate n0stack cluster node",
					Action: ServeAgent,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Value: hostname,
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
							Value: 20181,
						},
						cli.StringFlag{
							Name:  "base-directory",
							Value: "/var/lib/n0core",
						},
						cli.StringFlag{
							Name:  "location",
							Usage: "<Datacenter>/<AvailavilityZone>/<Cell>/<Rack>/<Unit(int)>",
						},
						cli.UintFlag{
							Name:  "cpu-capacity-milli-cores",
							Value: uint(node.GetTotalCPUMilliCores()) * 1000,
						},
						cli.Uint64Flag{
							Name:  "memory-capacity-bytes",
							Value: node.GetTotalMemory(),
						},
						cli.Uint64Flag{
							Name:  "storage-capacity-bytes",
							Value: uint64(100 * bytefmt.GIGABYTE),
						},
					},
				},
			},
		},
		{
			Name:  "deploy",
			Usage: "Deploy n0core to remote host with ssh",
			Subcommands: []cli.Command{
				{
					Name:      "agent",
					Action:    DeployAgent,
					ArgsUsage: "[user]@[hostname] [agent options]",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "base-directory",
							Value: "/var/lib/n0core",
						},
						cli.StringFlag{
							Name:  "send-to",
							Value: "/tmp",
						},
						cli.StringFlag{
							Name: "identity-file, i",
						},
					},
				},
			},
		},
		{
			Name:  "install",
			Usage: "Install n0core on localhost",
			Subcommands: []cli.Command{
				{
					Name:   "agent",
					Action: InstallAgent,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "base-directory",
							Value: "/var/lib/n0core",
						},
						cli.StringFlag{
							Name: "arguments",
						},
					},
				},
			},
		},
	}

	log.SetFlags(log.Llongfile | log.Ltime | log.Lmicroseconds)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start process, err:%s\n", err.Error())
		os.Exit(1)
	}
}
