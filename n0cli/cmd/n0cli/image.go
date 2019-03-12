package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gobwas/glob"
	"github.com/urfave/cli"
	"google.golang.org/grpc"

	pdeployment "github.com/n0stack/n0stack/n0proto.go/deployment/v0"
)

func GetImage(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 0 {
		return listImage(conn)
	}

	for _, name := range c.Args() {
		if hasMeta(name) {
			getImageByPattern(name, conn)
			return nil
		}

		err := getImage(name, conn)

		if err != nil {
			return err
		}
	}

	return nil
}

func listImage(conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.ListImages(context.Background(), &pdeployment.ListImagesRequest{})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func getImage(name string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.GetImage(context.Background(), &pdeployment.GetImageRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

func getImageByPattern(pattern string, conn *grpc.ClientConn) error {
	g, err := glob.Compile(pattern)
	if err != nil {
		return err
	}

	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.ListImages(context.Background(), &pdeployment.ListImagesRequest{})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	for _, image := range res.GetImages() {
		if g.Match(image.Name) {
			marshaler.Marshal(os.Stdout, image)
			fmt.Println()
		}
	}

	return nil
}

func DeleteImage(c *cli.Context) error {
	endpoint := c.GlobalString("api-endpoint")
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[DEBUG] Connected to '%s'\n", endpoint)

	if c.NArg() == 1 {
		name := c.Args().Get(0)
		return deleteImage(name, conn)
	}

	return fmt.Errorf("set valid arguments")
}

func deleteImage(name string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.DeleteImage(context.Background(), &pdeployment.DeleteImageRequest{Name: name})
	if err != nil {
		PrintGrpcError(err)
		return nil
	}

	marshaler.Marshal(os.Stdout, res)
	fmt.Println()

	return nil
}

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
		tags := c.StringSlice("t")
		a := c.Bool("apply-image")
		if a == true {
	        cl := pdeployment.NewImageServiceClient(conn)
	        _, err = cl.GetImage(context.Background(), &pdeployment.GetImageRequest{Name: img})
			if err != nil {
				err = registerApplyImage(img, conn)
				if err != nil {
					return err
				}
			}
		}
		return registerBlockStorage(tags, img, bs, conn)
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

func registerBlockStorage(tags []string, img, bs string, conn *grpc.ClientConn) error {
	cl := pdeployment.NewImageServiceClient(conn)
	res, err := cl.RegisterBlockStorage(context.Background(), &pdeployment.RegisterBlockStorageRequest{ImageName: img, BlockStorageName: bs, Tags: tags})
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

	if c.NArg() == 3 {
		img := c.Args().Get(0)
		t := c.Args().Get(1)
		bs := c.Args().Get(2)
		return tag(img, t, bs, conn)
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

	if c.NArg() == 2 {
		name := c.Args().Get(0)
		tag := c.Args().Get(1)
		return untag(name, tag, conn)
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