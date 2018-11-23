package main

import (
	"fmt"
	"log"
	"os"

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

	app.Commands = []cli.Command{
		{
			Name:   "api",
			Usage:  "",
			Action: ServeAPI,
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
					Value: 20180,
				},
			},
		},
		{
			Name:   "agent",
			Usage:  "",
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
			},
		},
		{
			Name:  "deploy",
			Usage: "",
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
							Name: "identity-file, i",
						},
					},
				},
			},
		},
		{
			Name:  "install",
			Usage: "",
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
		{
			Name:   "mock-agent",
			Usage:  "",
			Action: ServeMockAgent,
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
					Value: 20181,
				},
			},
		},
	}

	log.SetFlags(log.Lshortfile | log.Ltime | log.Lmicroseconds)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start process, err:%s\n", err.Error())
		os.Exit(1)
	}
}
