

## Example

```go
package main

import (
	"fmt"
	"log"

	"code.cloudfoundry.org/bytefmt"
	n0stack "github.com/n0stack/go-proto"
	"github.com/n0stack/go-proto/device"
	"github.com/n0stack/go-proto/device/volume"
	"github.com/n0stack/go-proto/resource/storage"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	id, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d69ce")
	// id, _ := uuid.NewV4()
	size, _ := bytefmt.ToBytes("4G")
	s := &volume.Spec{
		Device: &device.Spec{
			Model: &n0stack.Model{
				Id:   id.Bytes(),
				Name: "test-volume",
			},
		},
		Storage: &storage.Spec{
			Bytes: size,
		},
	}

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := volume.NewAgentClient(conn)

	// n, err := cli.Apply(context.Background(), s)
	n, err := cli.Delete(context.Background(), s.Device.Model)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", n)
}
```