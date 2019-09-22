package grpcutil

import (
	"net"
	"testing"

	"google.golang.org/grpc"
)

// port番号を乱数にしても良さそう
func PrepareMockedGRPC(t *testing.T, register func(*grpc.Server)) (*grpc.ClientConn, func() error) {
	listen := "localhost:20190"

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		t.Fatalf("net.Listen() returns err=%v", err)
	}

	grpcServer := grpc.NewServer()
	register(grpcServer)
	go grpcServer.Serve(lis)

	conn, err := grpc.Dial(listen, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("grpc.Dial() returns err=%+v", err)
	}

	return conn, func() error {
		if err := lis.Close(); err != nil {
			return err
		}

		return nil
	}
}
