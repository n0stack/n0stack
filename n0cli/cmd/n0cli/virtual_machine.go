package main

import (
	"context"
	"fmt"

	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func OpenConsoleOfVirtualMachine(c *cli.Context) error {
	if c.NArg() != 1 {
		return fmt.Errorf("set valid arguments")
	}
	resourceName := c.Args().Get(0)

	conn, err := ConnectAPI(c)
	if err != nil {
		return err
	}
	defer conn.Close()

	cl := pprovisioning.NewVirtualMachineServiceClient(conn)
	res, err := cl.OpenConsole(context.Background(), &pprovisioning.OpenConsoleRequest{Name: resourceName})
	if err != nil {
		fmt.Printf("Got error\n[%s] %s\n", grpc.Code(err).String(), grpc.ErrorDesc(err))
		return nil
	}

	fmt.Println(res.ConsoleUrl)

	return nil
}
