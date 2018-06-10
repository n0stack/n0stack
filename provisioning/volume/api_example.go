package volume

import (
	"context"
	"fmt"
	"log"

	"code.cloudfoundry.org/bytefmt"

	"github.com/n0stack/proto.go/provisioning/v0"
	"github.com/n0stack/proto.go/v0"

	"google.golang.org/grpc"
)

func apply() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "127.0.0.1", 20180), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	cli := pprovisioning.NewVolumeServiceClient(conn)

	b, _ := bytefmt.ToBytes("1G")
	v, err := cli.ApplyVolume(context.Background(), &pprovisioning.ApplyVolumeRequest{
		Metadata: &pn0stack.Metadata{
			Name: "test-volume",
			Annotations: map[string]string{
				"n0core/url":       "file:///tmp/test.qcow2",
				"n0core/node_name": "test-node",
			},
		},
		Spec: &pprovisioning.VolumeSpec{
			Bytes: b,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("applied volume:%v", v)
}

func delete() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "127.0.0.1", 20180), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	cli := pprovisioning.NewVolumeServiceClient(conn)

	v, err := cli.DeleteVolume(context.Background(), &pprovisioning.DeleteVolumeRequest{
		Name: "test-volume",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("delete volume:%v", v)
}
