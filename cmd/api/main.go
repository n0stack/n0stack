package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/n0stack/n0core/datastore/etcd"
	"github.com/n0stack/n0core/provisioning/node"
	"github.com/n0stack/n0core/provisioning/volume"
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

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Join to API and serve some daemons.",
			Action: func(c *cli.Context) error {
				lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.String("bind-address"), c.Int("bind-port")))
				if err != nil {
					return err
				}

				ne, err := etcd.NewEtcdDatastore("node", strings.Split(c.String("etcd-endpoints"), ","))
				if err != nil {
					return err
				}
				defer ne.Close()
				// starterをsliceでとったほうがいいかもしれない
				n, err := node.CreateNodeAPI(ne, c.String("memberlist-starter"))
				if err != nil {
					return err
				}

				ve, err := etcd.NewEtcdDatastore("volume", strings.Split(c.String("etcd-endpoints"), ","))
				if err != nil {
					return err
				}
				defer ve.Close()
				v, err := volume.CreateVolumeAPI(ve, n, c.String("volume-default-base-directory"))
				if err != nil {
					return err
				}

				s := grpc.NewServer()
				pprovisioning.RegisterNodeServiceServer(s, n)
				pprovisioning.RegisterVolumeServiceServer(s, v)
				reflection.Register(s)

				log.Printf("[INFO] Starting API")
				return s.Serve(lis)
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
				cli.StringFlag{
					Name:  "volume-default-base-directory",
					Value: "/var/lib/n0core/qcow2",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start process, err:%v", err.Error())
	}
}
