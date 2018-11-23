package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// References:
// 	https://github.com/grpc/grpc-go/pull/1406/files#diff-34c6b408d72845d076d47126c29948d1R591
// 	https://qiita.com/ktateish/items/ff5c045e3cebc59bf119
func MixHandler(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

func ServeAgent(ctx *cli.Context) error {
	bind := fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port"))
	baseDirectory := ctx.String("base-directory")
	nodeName := ctx.String("name")
	advertiseAddress := ctx.String("advertise-address")
	nodeAPI := ctx.String("node-api-endpoint")
	vlanInterface := ctx.String("vlan-external-interface")

	bvm := filepath.Join(baseDirectory, "virtual_machine")
	vma, err := provisioning.CreateVirtualMachineAgentAPI(bvm, vlanInterface)
	if err != nil {
		return err
	}

	bbs := filepath.Join(baseDirectory, "block_storage")
	va, err := provisioning.CreateBlockStorageAgentAPI(bbs)
	if err != nil {
		return err
	}

	// とりあえず log を表示するため利用する
	// zapLogger, err := zap.NewProduction()
	// if err != nil {
	// 	return err
	// }

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			// grpc_zap.StreamServerInterceptor(zapLogger),
			// grpc_auth.StreamServerInterceptor(auth),
			// grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			// grpc_zap.UnaryServerInterceptor(zapLogger),
			// grpc_auth.UnaryServerInterceptor(auth),
			// grpc_prometheus.UnaryServerInterceptor,
		)),
	)
	provisioning.RegisterBlockStorageAgentServiceServer(grpcServer, va)
	provisioning.RegisterVirtualMachineAgentServiceServer(grpcServer, vma)
	reflection.Register(grpcServer)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static(provisioning.DownloadBlockStorageHTTPPrefix, bbs)

	if err := node.RegisterNodeToAPI(nodeName, advertiseAddress, nodeAPI); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Started API: version=%s", version)
	// TODO: セキュリティ的に問題あり 暫定処置
	go e.Start("0.0.0.0:8081")
	return grpcServer.Serve(lis)
}
