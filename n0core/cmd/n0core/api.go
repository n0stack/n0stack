package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	pauth "n0st.ac/n0stack/auth/v1alpha"
	piam "n0st.ac/n0stack/iam/v1alpha"
	authn "n0st.ac/n0stack/n0core/pkg/api/auth/authentication"
	user "n0st.ac/n0stack/n0core/pkg/api/iam/user"
	"n0st.ac/n0stack/n0core/pkg/datastore/etcd"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/urfave/cli"
)

func OutputRecoveryLog(p interface{}) (err error) {
	log.Printf("[CRITICAL] panic happened: %v", p)
	log.Print(string(debug.Stack()))

	return nil
}

func ServeAPI(cctx *cli.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cctx.String("bind-address"), cctx.Int("bind-port")))
	if err != nil {
		return err
	}
	defer lis.Close()

	endpoint := fmt.Sprintf("localhost:%d", cctx.Int("bind-port"))
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	ds, err := etcd.NewEtcdDatastore(strings.Split(cctx.String("etcd-endpoints"), ","))
	if err != nil {
		return err
	}
	defer ds.Close()

	userapi := user.CreateUserAPI(ds)
	userClient := piam.NewUserServiceClient(conn)

	secret := cctx.String("token-secret")
	authapi := authn.CreateAuthenticationAPI(userClient, []byte(secret))

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(OutputRecoveryLog)),
		// grpc_auth.StreamServerInterceptor(auth),
		// grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(OutputRecoveryLog)),
			// grpc_auth.UnaryServerInterceptor(auth),
			// grpc_prometheus.UnaryServerInterceptor,
		)),
	)
	piam.RegisterUserServiceServer(grpcServer, userapi)
	pauth.RegisterAuthenticationServiceServer(grpcServer, authapi)
	reflection.Register(grpcServer)

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("SIGINT or SIGTERM received, stopping gracefully...")
		grpcServer.GracefulStop()
	}()

	log.Printf("[INFO] Started API: version=%s", version)
	return grpcServer.Serve(lis)
}
