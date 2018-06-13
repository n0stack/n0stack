package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/n0stack/n0core/provisioning/node"
	"github.com/n0stack/n0core/provisioning/node/iproute2"
	"github.com/n0stack/n0core/provisioning/node/qcow2"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	app := cli.NewApp()
	app.Name = "n0core Agent"
	// app.Usage = ""
	app.Version = "0.1.0" // CIで取るようにする

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Join to API and serve some daemons.",
			Action: func(c *cli.Context) error {
				err := node.JoinNode(c.String("name"), c.String("advertise-address"), c.String("api-address"), c.Int("api-port"))
				if err != nil {
					return err
				}

				lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.String("bind-address"), c.Int("bind-port")))
				if err != nil {
					return err
				}

				s := grpc.NewServer()
				qcow2.RegisterQcow2ServiceServer(s, &qcow2.Qcow2Agent{})

				i, err := iproute2.NewIproute2Agent(c.String("uplink-interface"))
				if err != nil {
					return err
				}
				iproute2.RegisterIproute2ServiceServer(s, i)

				reflection.Register(s)

				log.Printf("[INFO] Starting Agent")
				return s.Serve(lis)
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
				cli.IntFlag{
					Name: "api-port",
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
					Name: "uplink-interface",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start process, err:%v", err.Error())
	}
}
