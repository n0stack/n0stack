```go
package main

import (
	"fmt"
	"log"

	tap "github.com/n0stack/go.proto/tap/v0"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	id, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d69ee")
	// id, _ := uuid.NewV4()
	s := &tap.ApplyRequest{
		Id: id.String(),
		Spec: &tap.Spec{
			Type: tap.Spec_FLAT,
		},
	}
	fmt.Printf("%v\n", s)

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := tap.NewTapServiceClient(conn)

	// r, err := cli.Get(context.Background(), s.Model)
	// r, err := cli.Apply(context.Background(), s)
	r, err := cli.Delete(context.Background(), &tap.DeleteRequest{
		Id: id.String(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", err)
}

```