package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/n0stack/n0stack/n0core/pkg/api/deployment/flavor"
	"github.com/n0stack/n0stack/n0core/pkg/api/deployment/image"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/network"
	"github.com/n0stack/n0stack/n0core/pkg/api/pool/node"
	"github.com/n0stack/n0stack/n0core/pkg/api/provisioning"
	"github.com/n0stack/n0stack/n0proto/deployment/v0"
	"github.com/n0stack/n0stack/n0proto/pool/v0"
	"github.com/n0stack/n0stack/n0proto/provisioning/v0"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/n0stack/n0stack/n0core/pkg/datastore/etcd"

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

	noe, err := etcd.NewEtcdDatastore("node", etcdEndpoints)
	if err != nil {
		return err
	}
	defer noe.Close()

	noa, err := node.CreateNodeAPI(noe)
	if err != nil {
		return err
	}
	noc := ppool.NewNodeServiceClient(conn)

	nee, err := etcd.NewEtcdDatastore("network", etcdEndpoints)
	if err != nil {
		return err
	}
	defer nee.Close()

	nea, err := network.CreateNetworkAPI(nee)
	if err != nil {
		return err
	}
	nec := ppool.NewNetworkServiceClient(conn)

	bse, err := etcd.NewEtcdDatastore("block_storage", etcdEndpoints)
	if err != nil {
		return err
	}
	defer bse.Close()

	bsa, err := provisioning.CreateBlockStorageAPI(bse, noc)
	if err != nil {
		return err
	}
	bsc := pprovisioning.NewBlockStorageServiceClient(conn)

	vme, err := etcd.NewEtcdDatastore("virtual_machine", etcdEndpoints)
	if err != nil {
		return err
	}
	defer vme.Close()

	vma, err := provisioning.CreateVirtualMachineAPI(vme, noc, nec, bsc)
	if err != nil {
		return err
	}
	vmc := pprovisioning.NewVirtualMachineServiceClient(conn)

	ie, err := etcd.NewEtcdDatastore("image", etcdEndpoints)
	if err != nil {
		return err
	}
	defer ie.Close()

	ia, err := image.CreateImageAPI(ie, bsc)
	if err != nil {
		return err
	}
	ic := pdeployment.NewImageServiceClient(conn)

	fe, err := etcd.NewEtcdDatastore("flavor", etcdEndpoints)
	if err != nil {
		return err
	}
	defer ie.Close()

	fa, err := flavor.CreateFlavorAPI(fe, vmc, ic)
	if err != nil {
		return err
	}

	// とりあえず log を表示するため利用する
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zapLogger),
			// grpc_auth.StreamServerInterceptor(auth),
			// grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapLogger),
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

	log.Printf("[INFO] Starting API")
	return grpcServer.Serve(lis)
}
