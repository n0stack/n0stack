package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/cenkalti/backoff"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning/blockstorage"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning/virtualmachine"
	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

const AnnotationNodeAgentVersion = "github.com/n0stack/n0stack/n0core/agent_version"

func RegisterNodeToAPI(ctx context.Context, api, name, address, version string, cpu uint32, mem, storage uint64, dc, az, cell, rack string, unit uint32) error {
	ip, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return err
	}

	ar := &ppool.ApplyNodeRequest{
		Name: name,
		Annotations: map[string]string{
			AnnotationNodeAgentVersion: version,
		},

		Address:     ip.String(),
		IpmiAddress: node.GetIpmiAddress(),
		Serial:      node.GetSerial(),

		CpuMilliCores: cpu,
		MemoryBytes:   mem,
		StorageBytes:  storage,

		Datacenter:       dc,
		AvailabilityZone: az,
		Cell:             cell,
		Rack:             rack,
		Unit:             unit,
	}

	conn, err := grpc.Dial(api, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	cli := ppool.NewNodeServiceClient(conn)

	n, err := cli.GetNode(ctx, &ppool.GetNodeRequest{Name: name})
	if err != nil {
		if grpc.Code(err) != codes.NotFound {
			return err
		}
	} else {
		log.Printf("[INFO] Get old Node: Node=%v", n)
	}

	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5)
	err = backoff.Retry(func() error {
		n, err = cli.ApplyNode(context.Background(), ar)
		if err != nil {
			return err
		}

		log.Printf("[INFO] Applied Node to APi: Node=%v", n)
		return nil
	}, b)
	if err != nil {
		return err
	}

	return nil
}

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

	cpu := uint32(ctx.Uint("cpu-capacity-milli-cores"))
	memory := ctx.Uint64("memory-capacity-bytes")
	storage := ctx.Uint64("storage-capacity-bytes")

	location := strings.Split(ctx.String("location"), "/")
	if len(location) != 5 {
		return fmt.Errorf("invalid argument 'location'")
	}
	dc := location[0]
	az := location[1]
	cell := location[2]
	rack := location[3]
	u, err := strconv.ParseUint(location[4], 10, 32)
	if err != nil {
		return err
	}
	unit := uint32(u)

	bvm := filepath.Join(baseDirectory, "virtual_machine")
	vma, err := virtualmachine.CreateVirtualMachineAgent(bvm)
	if err != nil {
		return err
	}

	bbs := filepath.Join(baseDirectory, "block_storage")
	va, err := blockstorage.CreateBlockStorageAgentAPI(bbs)
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
	blockstorage.RegisterBlockStorageAgentServiceServer(grpcServer, va)
	virtualmachine.RegisterVirtualMachineAgentServiceServer(grpcServer, vma)
	reflection.Register(grpcServer)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static(blockstorage.DownloadBlockStorageHTTPPrefix, bbs)

	c := context.Background()
	if err := RegisterNodeToAPI(c, nodeAPI, nodeName, advertiseAddress, version, cpu, memory, storage, dc, az, cell, rack, unit); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Started API: version=%s", version)
	// TODO: セキュリティ的に問題あり 暫定処置
	go e.Start("0.0.0.0:8081")

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("SIGINT or SIGTERM received, stopping gracefully...")
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(lis)
}
