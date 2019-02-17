package main

import (
	"context"
	"fmt"
	"os"

	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func DownloadBlockStorage(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("set valid arguments")
	}
	resourceName := c.Args().Get(0)

	conn, err := ConnectAPI(c)
	if err != nil {
		return err
	}
	defer conn.Close()

	cl := pprovisioning.NewBlockStorageServiceClient(conn)
	res, err := cl.DownloadBlockStorage(context.Background(), &pprovisioning.DownloadBlockStorageRequest{Name: resourceName})
	if err != nil {
		fmt.Printf("Got error\n[%s] %s\n", grpc.Code(err).String(), grpc.ErrorDesc(err))
		return nil
	}

	fmt.Println(res.DownloadUrl)

	return nil
}

func ListBlockStorage(ctx context.Context, conn *grpc.ClientConn) error {
	cl := pprovisioning.NewBlockStorageServiceClient(conn)
	res, err := cl.ListBlockStorages(context.Background(), &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		fmt.Printf("Got error\n[%s] %s\n", grpc.Code(err).String(), grpc.ErrorDesc(err))
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println("")

	return nil
}

func GetBlockStorage(ctx context.Context, conn *grpc.ClientConn, resourceName string) error {
	cl := pprovisioning.NewBlockStorageServiceClient(conn)
	res, err := cl.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: resourceName})
	if err != nil {
		fmt.Printf("Got error\n[%s] %s\n", grpc.Code(err).String(), grpc.ErrorDesc(err))
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println("")

	return nil
}
