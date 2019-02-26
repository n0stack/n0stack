package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gobwas/glob"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	ppool "github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

func GetNode(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 0 {
		return listNode(conn)
	}

	for _, name := range c.Args() {
		if hasMeta(name) {
			getNodeByPattern(name, conn)
			return nil
		}

		err := getNode(name, conn)

		if err != nil {
			return err
		}
	}

	return nil
}

func listNode(conn *grpc.ClientConn) error {
	cl := ppool.NewNodeServiceClient(conn)
	res, err := cl.ListNodes(context.Background(), &ppool.ListNodesRequest{})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func getNode(name string, conn *grpc.ClientConn) error {
	cl := ppool.NewNodeServiceClient(conn)
	res, err := cl.GetNode(context.Background(), &ppool.GetNodeRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func getNodeByPattern(pattern string, conn *grpc.ClientConn) error {
	g, err := glob.Compile(pattern)
	if err != nil {
		return err
	}

	cl := ppool.NewNodeServiceClient(conn)
	res, err := cl.ListNodes(context.Background(), &ppool.ListNodesRequest{})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	for _, node := range res.GetNodes() {
		if g.Match(node.Name) {
			marshaler.Marshal(os.Stdout, node)
			fmt.Println()
		}
	}

	return nil
}

func DeleteNode(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 1 {
		name := c.Args().Get(0)
		return deleteNode(name, conn)
	}

	return fmt.Errorf("set valid arguments")
}

func deleteNode(name string, conn *grpc.ClientConn) error {
	cl := ppool.NewNodeServiceClient(conn)
	res, err := cl.DeleteNode(context.Background(), &ppool.DeleteNodeRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}
