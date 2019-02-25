package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func GetBlockStorage(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 0 {
		return listBlockStorage(conn)
	}

	for _, name := range c.Args() {
		err := getBlockStorage(name, conn)

		if err != nil {
			return err
		}
	}

	return nil
}

func listBlockStorage(conn *grpc.ClientConn) error {
	cl := pprovisioning.NewBlockStorageServiceClient(conn)
	res, err := cl.ListBlockStorages(context.Background(), &pprovisioning.ListBlockStoragesRequest{})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func getBlockStorage(name string, conn *grpc.ClientConn) error {
	cl := pprovisioning.NewBlockStorageServiceClient(conn)
	res, err := cl.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func DeleteBlockStorage(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 1 {
		name := c.Args().Get(0)
		return deleteBlockStorage(name, conn)
	}

	return fmt.Errorf("set valid arguments")
}

func deleteBlockStorage(name string, conn *grpc.ClientConn) error {
	cl := pprovisioning.NewBlockStorageServiceClient(conn)
	res, err := cl.DeleteBlockStorage(context.Background(), &pprovisioning.DeleteBlockStorageRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

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
		PrintGrpcError(err)
		return nil
	}

	fmt.Println(res.DownloadUrl)

	return nil
}
