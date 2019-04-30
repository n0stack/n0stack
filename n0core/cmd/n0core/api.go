package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/embed"

	"github.com/n0stack/n0stack/n0core/pkg/api/deployment/image"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/network"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning/blockstorage"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning/virtualmachine"
	pdeployment "github.com/n0stack/n0stack/n0proto.go/deployment/v0"
	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
	pprovisioning "github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/rakyll/statik/fs"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	// "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	// "go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/urfave/cli"
)

func OutputRecoveryLog(p interface{}) (err error) {
	log.Printf("[CRITICAL] panic happened: %v", p)
	log.Print(string(debug.Stack()))

	return nil
}

func ServeAPI(ctx *cli.Context) error {
	baseDir := ctx.String("base-directory")
	dbDir := filepath.Join(baseDir, "db")

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

	ds, err := embed.NewEmbedDatastore(dbDir)
	if err != nil {
		return err
	}
	defer ds.Close()

	noa := node.CreateNodeAPI(ds)
	noc := ppool.NewNodeServiceClient(conn)

	nea := network.CreateNetworkAPI(ds)
	nec := ppool.NewNetworkServiceClient(conn)

	bsa := blockstorage.CreateBlockStorageAPI(ds, noc)
	bsc := pprovisioning.NewBlockStorageServiceClient(conn)

	vma := virtualmachine.CreateVirtualMachineAPI(ds, noc, nec, bsc)
	// vmc := pprovisioning.NewVirtualMachineServiceClient(conn)

	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	ia := image.CreateImageAPI(ds, bsc)
	// ic := pdeployment.NewImageServiceClient(conn)

	// とりあえず log を表示するため利用する
	// zapLogger, err := zap.NewProduction()
	// if err != nil {
	// 	return err
	// }

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(OutputRecoveryLog)),
			// grpc_zap.StreamServerInterceptor(zapLogger),
			// grpc_auth.StreamServerInterceptor(auth),
			// grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(OutputRecoveryLog)),
			// grpc_zap.UnaryServerInterceptor(zapLogger),
			// grpc_auth.UnaryServerInterceptor(auth),
			// grpc_prometheus.UnaryServerInterceptor,
		)),
	)
	ppool.RegisterNodeServiceServer(grpcServer, noa)
	ppool.RegisterNetworkServiceServer(grpcServer, nea)
	pprovisioning.RegisterBlockStorageServiceServer(grpcServer, bsa)
	pprovisioning.RegisterVirtualMachineServiceServer(grpcServer, vma)
	pdeployment.RegisterImageServiceServer(grpcServer, ia)
	reflection.Register(grpcServer)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/n0core/api/v0/virtual_machines/:name/vncwebsocket", vma.ProxyWebsocket())
	e.GET("/n0core/api/v0/block_storage/download/:name", bsa.ProxyDownloadBlockStorage(8081, "/n0core/api/v0/block_storage/download/")) // ダウンロードしたときのファイル名を指定するため、このように指定した
	e.GET("/n0core/static/virtual_machines/*", echo.WrapHandler(http.StripPrefix("/n0core/static/virtual_machines/", http.FileServer(statikFs))))

	log.Printf("[INFO] Started API: version=%s", version)

	// 本当は panic させる必要がある
	go func() {
		if err := e.Start("0.0.0.0:8080"); err != nil {
			panic(err)
		}
	}()

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("SIGINT or SIGTERM received, stopping gracefully...")
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(lis)
}
