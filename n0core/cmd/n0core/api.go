package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/n0stack/n0stack/n0core/pkg/api/deployment/flavor"
	"github.com/n0stack/n0stack/n0core/pkg/api/deployment/image"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/network"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/etcd"
	"github.com/n0stack/n0stack/n0proto.go/deployment/v0"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/rakyll/statik/fs"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	// "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	// "go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/urfave/cli"
)

func ServeAPI(ctx *cli.Context) error {
	etcdEndpoints := strings.Split(ctx.String("etcd-endpoints"), ",")

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctx.String("bind-address"), ctx.Int("bind-port")))
	if err != nil {
		return err
	}
	defer lis.Close()

	endpoint := fmt.Sprintf("localhost:%d", ctx.Int("bind-port"))
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dial:", err)
	}
	defer conn.Close()

	ds, err := etcd.NewEtcdDatastore(etcdEndpoints)
	if err != nil {
		return err
	}
	defer ds.Close()

	noa := node.CreateNodeAPI(ds)
	noc := ppool.NewNodeServiceClient(conn)

	nea := network.CreateNetworkAPI(ds)
	nec := ppool.NewNetworkServiceClient(conn)

	bsa := provisioning.CreateBlockStorageAPI(ds, noc)
	bsc := pprovisioning.NewBlockStorageServiceClient(conn)

	vma := provisioning.CreateVirtualMachineAPI(ds, noc, nec, bsc)
	vmc := pprovisioning.NewVirtualMachineServiceClient(conn)

	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	ia := image.CreateImageAPI(ds, bsc)
	ic := pdeployment.NewImageServiceClient(conn)

	fa := flavor.CreateFlavorAPI(ds, vmc, ic)

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
		)))
	ppool.RegisterNodeServiceServer(grpcServer, noa)
	ppool.RegisterNetworkServiceServer(grpcServer, nea)
	pprovisioning.RegisterBlockStorageServiceServer(grpcServer, bsa)
	pprovisioning.RegisterVirtualMachineServiceServer(grpcServer, vma)
	pdeployment.RegisterImageServiceServer(grpcServer, ia)
	pdeployment.RegisterFlavorServiceServer(grpcServer, fa)
	reflection.Register(grpcServer)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/v0/virtual_machines/:name/vncwebsocket", vma.ProxyWebsocket())
	e.GET("/static/virtual_machines/*", echo.WrapHandler(http.StripPrefix("/static/virtual_machines/", http.FileServer(statikFs))))

	log.Printf("[INFO] Started API: version=%s", version)
	// 本当は panic させる必要がある
	go e.Start("0.0.0.0:8080")
	return grpcServer.Serve(lis)
}
