package main

import (
	"context"
	"log"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	pauth "n0st.ac/n0stack/auth/v1alpha"
	piam "n0st.ac/n0stack/iam/v1alpha"
)

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
	if err := piam.RegisterProjectServiceHandlerFromEndpoint(ctx, mux, api.Host, opts); err != nil {
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
	e.POST("/api/*", echo.WrapHandler(mux))
	e.DELETE("/api/*", echo.WrapHandler(mux))
	e.PATCH("/api/*", echo.WrapHandler(mux))
	// websocket proxy ができてない

	log.Printf("[INFO] Started BFF: version=%s", version)
	return e.Start(c.String("listen-address"))
}
