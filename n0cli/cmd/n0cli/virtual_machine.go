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

func GetVirtualMachine(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 0 {
		return listVirtualMachine(conn)
	}

	for _, name := range c.Args() {
		err := getVirtualMachine(name, conn)

		if err != nil {
			return err
		}
	}

	return nil
}

func listVirtualMachine(conn *grpc.ClientConn) error {
	cl := pprovisioning.NewVirtualMachineServiceClient(conn)
	res, err := cl.ListVirtualMachines(context.Background(), &pprovisioning.ListVirtualMachinesRequest{})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func getVirtualMachine(name string, conn *grpc.ClientConn) error {
	cl := pprovisioning.NewVirtualMachineServiceClient(conn)
	res, err := cl.GetVirtualMachine(context.Background(), &pprovisioning.GetVirtualMachineRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func DeleteVirtualMachine(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 1 {
		name := c.Args().Get(0)
		return deleteVirtualMachine(name, conn)
	}

	return fmt.Errorf("set valid arguments")
}

func deleteVirtualMachine(name string, conn *grpc.ClientConn) error {
	cl := pprovisioning.NewVirtualMachineServiceClient(conn)
	res, err := cl.DeleteVirtualMachine(context.Background(), &pprovisioning.DeleteVirtualMachineRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func BootVirtualMachine(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 1 {
		name := c.Args().Get(0)
		return bootVirtualMachine(name, conn)
	}

	return fmt.Errorf("set valid arguments")
}

func bootVirtualMachine(name string, conn *grpc.ClientConn) error {
	cl := pprovisioning.NewVirtualMachineServiceClient(conn)
	res, err := cl.BootVirtualMachine(context.Background(), &pprovisioning.BootVirtualMachineRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

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
