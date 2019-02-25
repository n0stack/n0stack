package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"google.golang.org/grpc"

	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

func GetNetwork(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 0 {
		return listNetwork(conn)
	}

	for _, name := range c.Args() {
		err := getNetwork(name, conn)

		if err != nil {
			return err
		}
	}

	return nil
}

func listNetwork(conn *grpc.ClientConn) error {
	cl := ppool.NewNetworkServiceClient(conn)
	res, err := cl.ListNetworks(context.Background(), &ppool.ListNetworksRequest{})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func getNetwork(name string, conn *grpc.ClientConn) error {
	cl := ppool.NewNetworkServiceClient(conn)
	res, err := cl.GetNetwork(context.Background(), &ppool.GetNetworkRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func DeleteNetwork(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 1 {
		name := c.Args().Get(0)
		return deleteNetwork(name, conn)
	}

	return fmt.Errorf("set valid arguments")
}

func deleteNetwork(name string, conn *grpc.ClientConn) error {
	cl := ppool.NewNetworkServiceClient(conn)
	res, err := cl.DeleteNetwork(context.Background(), &ppool.DeleteNetworkRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}
