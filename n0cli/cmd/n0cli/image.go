package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
	"google.golang.org/grpc"

	pdeployment "github.com/n0stack/n0stack/n0proto.go/deployment/v0"
)

func RegisterBlockStorage(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 2 {
		img := c.Args().Get(0)
		bs := c.Args().Get(1)
		a := c.Bool("apply-image")
		if a == true {
			fmt.Println("---> Try 'GetImageRequest ", img, "'")
			cl := pdeployment.NewImageServiceClient(conn)
			_, err = cl.GetImage(context.Background(), &pdeployment.GetImageRequest{Name: img})
			if err != nil {
				PrintGrpcError(err)
				fmt.Println("---> Send 'ApplyImageRequest ", img, "'")
				err = registerApplyImage(img, conn)
				if err != nil {
					return err
				}
				fmt.Println("---> Success 'ApplyImageRequest ", img, "'")
				fmt.Println("---> Send 'RegisterBlockStorageRequest'")
			}
		}
		return registerBlockStorage(img, bs, conn)
	}
	return fmt.Errorf("set valid arguments.")
}

func registerApplyImage(name string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.ApplyImage(context.Background(), &pdeployment.ApplyImageRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}
	marshaler.Marshal(os.Stdout, res)
	fmt.Println()
	return nil
}

func registerBlockStorage(img, bs string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.RegisterBlockStorage(context.Background(), &pdeployment.RegisterBlockStorageRequest{ImageName: img, BlockStorageName: bs})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}
	marshaler.Marshal(os.Stdout, res)
	fmt.Println()
	return nil
}

func UnregisterBlockStorage(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 2 {
		img := c.Args().Get(0)
		bs := c.Args().Get(1)
		return unregisterBlockStorage(img, bs, conn)
	}
	return fmt.Errorf("set valid arguments.")
}

func unregisterBlockStorage(img, bs string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.UnregisterBlockStorage(context.Background(), &pdeployment.UnregisterBlockStorageRequest{ImageName: img, BlockStorageName: bs})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}
	marshaler.Marshal(os.Stdout, res)
	fmt.Println()
	return nil
}

func Tag(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 2 {
		str := strings.Split(c.Args().Get(0), ":")
		if len(str) != 2 {
			return fmt.Errorf("set valid arguments.")
		}
		bs := c.Args().Get(1)
		return tag(str[0], str[1], bs, conn)
	}
	return fmt.Errorf("set valid arguments.")
}

func tag(name, tag, bs string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	fmt.Println(name, tag, bs)
	res, err := cl.TagImage(context.Background(), &pdeployment.TagImageRequest{Name: name, Tag: tag, RegisteredBlockStorageName: bs})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}
	marshaler.Marshal(os.Stdout, res)
	fmt.Println()
	return nil
}

func Untag(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 1 {
		str := strings.Split(c.Args().Get(0), ":")
		if len(str) != 2 {
			return fmt.Errorf("set valid arguments.")
		}
		return untag(str[0], str[1], conn)
	}
	return fmt.Errorf("set valid arguments.")
}

func untag(name, tag string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.UntagImage(context.Background(), &pdeployment.UntagImageRequest{Name: name, Tag: tag})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}
	marshaler.Marshal(os.Stdout, res)
	fmt.Println()
	return nil
}
