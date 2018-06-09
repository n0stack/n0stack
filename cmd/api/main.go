package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/n0stack/n0core/datastore/etcd"
	"github.com/n0stack/n0core/provisioning/node"
	pprovisioning "github.com/n0stack/proto.go/provisioning/v0"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	app := cli.NewApp()
	app.Name = "n0core API"
	// app.Usage = ""
	app.Version = "0.1.0" // CIで取るようにする

	// command action
	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Join to API and serve some daemons.",
			Action: func(c *cli.Context) error {
				var err error
				e, err := etcd.NewEtcdDatastore("node", strings.Split(c.String("etcd-endpoints"), ","))
				if err != nil {
					return err
				}

				lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.String("bind-address"), c.Int("bind-port")))
				if err != nil {
					return err
				}

				s := grpc.NewServer()

				// starterをsliceでとったほうがいいかもしれない
				n, err := node.CreateNodeAPI(e, c.String("memberlist-starter"))
				if err != nil {
					return err
					fmt.Printf("%v", err.Error())

				}

				pprovisioning.RegisterNodeServiceServer(s, n)
				reflection.Register(s)

				log.Printf("[INFO] Starting API")
				if err := s.Serve(lis); err != nil {
					return err
				}

				return nil
			},
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("%v", err.Error())
	}
}
