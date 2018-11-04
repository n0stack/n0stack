package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/n0stack/n0stack/n0proto.go/deployment/v0"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func Get(ctx *cli.Context) error {
	if ctx.NArg() == 1 {
		return list(ctx)
	} else if ctx.NArg() == 2 {
		return get(ctx)
	}

	return fmt.Errorf("set valid arguments")
}

func get(ctx *cli.Context) error {
	resourceType := ctx.Args().Get(0)
	resourceName := ctx.Args().Get(1)

	endpoint := ctx.String("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	switch resourceType {
	case "node":
		cl := ppool.NewNodeServiceClient(conn)
		res, err := cl.GetNode(context.Background(), &ppool.GetNodeRequest{Name: resourceName})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "network":
		cl := ppool.NewNetworkServiceClient(conn)
		res, err := cl.GetNetwork(context.Background(), &ppool.GetNetworkRequest{Name: resourceName})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "block_storage":
		cl := pprovisioning.NewBlockStorageServiceClient(conn)
		res, err := cl.GetBlockStorage(context.Background(), &pprovisioning.GetBlockStorageRequest{Name: resourceName})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "virtual_machine":
		cl := pprovisioning.NewVirtualMachineServiceClient(conn)
		res, err := cl.GetVirtualMachine(context.Background(), &pprovisioning.GetVirtualMachineRequest{Name: resourceName})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "image":
		cl := pdeployment.NewImageServiceClient(conn)
		res, err := cl.GetImage(context.Background(), &pdeployment.GetImageRequest{Name: resourceName})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "flavor":
		cl := pdeployment.NewFlavorServiceClient(conn)
		res, err := cl.GetFlavor(context.Background(), &pdeployment.GetFlavorRequest{Name: resourceName})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	default:
		return fmt.Errorf("resource type '%s' is not existing\n", resourceType)
	}

	return nil
}

func list(ctx *cli.Context) error {
	resourceType := ctx.Args().Get(0)

	endpoint := ctx.String("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	switch resourceType {
	case "node":
		cl := ppool.NewNodeServiceClient(conn)
		res, err := cl.ListNodes(context.Background(), &ppool.ListNodesRequest{})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "network":
		cl := ppool.NewNetworkServiceClient(conn)
		res, err := cl.ListNetworks(context.Background(), &ppool.ListNetworksRequest{})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "block_storage":
		cl := pprovisioning.NewBlockStorageServiceClient(conn)
		res, err := cl.ListBlockStorages(context.Background(), &pprovisioning.ListBlockStoragesRequest{})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "virtual_machine":
		cl := pprovisioning.NewVirtualMachineServiceClient(conn)
		res, err := cl.ListVirtualMachines(context.Background(), &pprovisioning.ListVirtualMachinesRequest{})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "image":
		cl := pdeployment.NewImageServiceClient(conn)
		res, err := cl.ListImages(context.Background(), &pdeployment.ListImagesRequest{})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	case "flavor":
		cl := pdeployment.NewFlavorServiceClient(conn)
		res, err := cl.ListFlavors(context.Background(), &pdeployment.ListFlavorsRequest{})
		if err != nil {
			fmt.Printf("got error response: %s\n", err.Error())
			return nil
		}

		marshaler.Marshal(os.Stdout, res)

	default:
		return fmt.Errorf("resource type '%s' is not existing\n", resourceType)
	}

	fmt.Println("")

	return nil
}
