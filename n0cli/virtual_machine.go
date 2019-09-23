// +build ignore

package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli"

	pprovisioning "n0st.ac/n0stack/n0proto.go/provisioning/v0"
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
		PrintGrpcError(err)
		return nil
	}

	fmt.Println(res.ConsoleUrl)

	return nil
}
