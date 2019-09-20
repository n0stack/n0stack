package n0provisioning

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	pprovisioning "n0st.ac/n0stack/n0proto.go/provisioning/v0"
	"google.golang.org/grpc"
)

func VirutalMachine_CreateVirtualMachine(req *pprovisioning.CreateVirtualMachineRequest) *GPRCUnit {
	Run := func(ctx context.Context, conn *grpc.ClientConn) (proto.Message, error) {
		cli := pprovisioning.NewVirtualMachineServiceClient(conn)
		return cli.CreateVirtualMachine(ctx, req)
	}

	String := func() string {
		r, _ := PbMarshaler.MarshalToString(req)
		return fmt.Sprintf("VirtualMachineService.CreateVirtualMachine(%s)", r)
	}

	return &GPRCUnit{
		runMethod:    Run,
		stringMethod: String,
	}
}
