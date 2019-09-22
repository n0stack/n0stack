package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	pauth "n0st.ac/n0stack/auth/v1alpha"
	piam "n0st.ac/n0stack/iam/v1alpha"
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
							Name: "api-url",
						},
						cli.StringFlag{
							Name:  "listen-address",
							Value: "0.0.0.0:8080",
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

	api, err := url.Parse(c.String("api-url"))
	if err != nil {
		return err
	}

	var opts []grpc.DialOption
	if api.Scheme == "http" {
		opts = append(opts, grpc.WithInsecure())
	}

	mux := runtime.NewServeMux()

	// とりあえず動くようにした。
	if err := piam.RegisterUserServiceHandlerFromEndpoint(ctx, mux, api.Host, opts); err != nil {
		return err
	}
	if err := pauth.RegisterAuthenticationServiceHandlerFromEndpoint(ctx, mux, api.Host, opts); err != nil {
		return err
	}

	// /n0core にプロキシ
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/*", echo.WrapHandler(mux))
	// websocket proxy ができてない

	log.Printf("[INFO] Started BFF: version=%s", version)
	return e.Start(c.String("listen-address"))
}
