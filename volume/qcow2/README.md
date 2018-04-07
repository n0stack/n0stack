

## Example

```go
package main

import (
	"fmt"
	"log"

	"code.cloudfoundry.org/bytefmt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pkvm "github.com/n0stack/go.proto/kvm/v0"
	pqcow2 "github.com/n0stack/go.proto/qcow2/v0"
	ptap "github.com/n0stack/go.proto/tap/v0"
)

func main() {
	volumeID, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d79ce")
	// volumeID, _ := uuid.NewV4()
	size, _ := bytefmt.ToBytes("4G")
	req := &pqcow2.ApplyRequest{
		Id: volumeID.String(),
		Spec: &pqcow2.Spec{
			Bytes: size,
		},
	}

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := pqcow2.NewQcow2ServiceClient(conn)

	// r, err := cli.Get(context.Background(), s.Model)
	r, err := cli.Apply(context.Background(), req)
	// r, err := cli.Delete(context.Background(), &qcow2.DeleteRequest{Id: req.Id})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", req)
	fmt.Printf("%v\n", r)

	tapID, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d89ce")
	s := &ptap.ApplyRequest{
		Id: tapID.String(),
		Spec: &ptap.Spec{
			Type: ptap.Spec_FLAT,
		},
	}
	fmt.Printf("%v\n", s)

	tapCli := ptap.NewTapServiceClient(conn)

	// r, err := cli.Get(context.Background(), s.Model)
	t, err := tapCli.Apply(context.Background(), s)
	// r, err := cli.Delete(context.Background(), &tap.DeleteRequest{Id: s.Id})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", t)
	fmt.Printf("%v\n", err)

	vmID, _ := uuid.FromString("614de68e-8e7c-429e-830f-76593a1d69ce")
	// id, _ := uuid.NewV4()
	memorySize, _ := bytefmt.ToBytes("1G")
	vreq := &pkvm.ApplyRequest{
		Id: vmID.String(),
		Spec: &pkvm.Spec{
			Vcpus:       1,
			MemoryBytes: memorySize,
			Volumes: []string{
				volumeID.String(),
			},
			Nics: []*pkvm.Spec_NIC{
				&pkvm.Spec_NIC{
					Tap:    tapID.String(),
					Hwaddr: "52:54:00:00:12:34",
				},
			},
		},
	}

	vmCli := pkvm.NewKVMServiceClient(conn)

	n, err := vmCli.Apply(context.Background(), vreq)
	// n, err := vmCli.Delete(context.Background(), s.Device.Model)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", vreq)
	fmt.Printf("%v\n", n)
}

```