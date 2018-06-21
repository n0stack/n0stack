package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/n0stack/n0core/provisioning/compute"

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
			Action: func(ctx *cli.Context) error {
				lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
				if err != nil {
					return err
				}

				ne, err := etcd.NewEtcdDatastore("node", strings.Split(ctx.String("etcd-endpoints"), ","))
				if err != nil {
					return err
				}
				defer ne.Close()
				// starterをsliceでとったほうがいいかもしれない
				n, err := node.CreateNodeAPI(ne, ctx.String("memberlist-starter"))
				if err != nil {
					return err
				}

				nc, err := node.NewNodeConnections(fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
				if err != nil {
					return err
				}

				ve, err := etcd.NewEtcdDatastore("volume", strings.Split(ctx.String("etcd-endpoints"), ","))
				if err != nil {
					return err
				}
				defer ve.Close()
				v, err := volume.CreateVolumeAPI(ve, nc, ctx.String("volume-default-base-directory"))
				if err != nil {
					return err
				}

				ce, err := etcd.NewEtcdDatastore("compute", strings.Split(ctx.String("etcd-endpoints"), ","))
				if err != nil {
					return err
				}
				defer ce.Close()
				c, err := compute.CreateComputeAPI(ce, nc, ctx.String("compute-default-base-directory"))
				if err != nil {
					return err
				}

				s := grpc.NewServer()
				pprovisioning.RegisterNodeServiceServer(s, n)
				pprovisioning.RegisterVolumeServiceServer(s, v)
				pprovisioning.RegisterComputeServiceServer(s, c)
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
				cli.StringFlag{
					Name:  "compute-default-base-directory",
					Value: "/var/lib/n0core/kvm",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start process, err:%v", err.Error())
	}
}
