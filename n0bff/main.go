package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pdeployment "github.com/n0stack/n0stack/n0proto.go/deployment/v0"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
	pprovisioning "github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
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
					Name:   "bff",
					Usage:  "Daemon which provide bff for n0stack API",
					Action: ServeBFF,
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
			},
		},
	}

	log.SetFlags(log.Llongfile | log.Ltime | log.Lmicroseconds)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start process, err:%s\n", err.Error())
		os.Exit(1)
	}
}

func ServeBFF(c *cli.Context) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	api := "api:20180"

	// とりあえず動くようにした。
	if err := ppool.RegisterNetworkServiceHandlerFromEndpoint(ctx, mux, api, opts); err != nil {
		return err
	}
	if err := pprovisioning.RegisterBlockStorageServiceHandlerFromEndpoint(ctx, mux, api, opts); err != nil {
		return err
	}
	if err := pprovisioning.RegisterVirtualMachineServiceHandlerFromEndpoint(ctx, mux, api, opts); err != nil {
		return err
	}
	if err := pdeployment.RegisterImageServiceHandlerFromEndpoint(ctx, mux, api, opts); err != nil {
		return err
	}
	if err := ppool.RegisterNodeServiceHandlerFromEndpoint(ctx, mux, api, opts); err != nil {
		return err
	}

	u := &url.URL{
		Scheme: "http",
		Host:   "api:8080",
	}
	// /n0core にプロキシ
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/*", echo.WrapHandler(mux))
	// websocket proxy ができてない
	e.GET("/n0core/*", echo.WrapHandler(httputil.NewSingleHostReverseProxy(u)))

	log.Printf("[INFO] Started BFF: version=%s", version)
	return e.Start("0.0.0.0:8080")
}
