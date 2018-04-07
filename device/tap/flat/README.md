```go
package main

import (
	"fmt"
	"log"

	tap "github.com/n0stack/go-proto/device/tap/v0"
	"github.com/n0stack/go-proto/resource/networkid/v0"

	n0stack "github.com/n0stack/go-proto"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	id, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d69ee")
	// id, _ := uuid.NewV4()
	s := &tap.ApplyRequest{
		Model: &n0stack.Model{
			Id:   id.Bytes(),
			Name: "test-network",
		},
		Spec: &tap.Spec{
			NetworkID: &networkid.Spec{
				Type: networkid.Spec_FLAT,
			},
		},
	}
	fmt.Printf("%v\n", s)

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := tap.NewStandardClient(conn)

	// r, err := cli.Get(context.Background(), s.Model)
	// r, err := cli.Apply(context.Background(), s)
	r, err := cli.Delete(context.Background(), &tap.DeleteRequest{
		Model: s.Model,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", err)
}
```