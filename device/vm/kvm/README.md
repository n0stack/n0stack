```go
package main

import (
	"fmt"
	"log"

	"github.com/n0stack/go-proto/resource/hwaddr"
	networkid "github.com/n0stack/go-proto/resource/networkid/v0"

	tap "github.com/n0stack/go-proto/device/tap/v0"
	"github.com/n0stack/go-proto/device/volume"

	"code.cloudfoundry.org/bytefmt"
	n0stack "github.com/n0stack/go-proto"
	"github.com/n0stack/go-proto/device"
	"github.com/n0stack/go-proto/device/vm"
	"github.com/n0stack/go-proto/resource/cpu"
	"github.com/n0stack/go-proto/resource/memory"
	"github.com/n0stack/go-proto/resource/storage"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	vmID, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d69ce")
	volumeID, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d79ce")
	tapID, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d89ce")

	// id, _ := uuid.NewV4()
	memorySize, _ := bytefmt.ToBytes("1G")
	volumeSize, _ := bytefmt.ToBytes("8G")
	s := &vm.Spec{
		Device: &device.Spec{
			Model: &n0stack.Model{
				Id:   vmID.Bytes(),
				Name: "test-vm",
			},
		},
		Cpu: &cpu.Spec{
			Vcpus:        1,
			Architecture: cpu.Architecture_x86_64,
		},
		Memory: &memory.Spec{
			Bytes: memorySize,
		},
		Volumes: []*volume.Spec{
			&volume.Spec{
				Model: &n0stack.Model{
					Id:   volumeID.Bytes(),
					Name: "test-volume",
				},
				Storage: &storage.Spec{
					Bytes: volumeSize,
				},
			},
		},
		Nics: []*vm.Spec_NIC{
			&vm.Spec_NIC{
				Model: &n0stack.Model{
					Id:   tapID.Bytes(),
					Name: "test-network",
				},
				Tap: &tap.Spec{
					NetworkID: &networkid.Spec{
						Type: networkid.Spec_FLAT,
					},
				},
				HwAddr: &hwaddr.Spec{
					Address: "52:54:00:00:12:34",
				},
			},
		},
	}

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := vm.NewStandardClient(conn)

	// n, err := cli.Apply(context.Background(), s)
	n, err := cli.Delete(context.Background(), s.Device.Model)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", n)
}
```