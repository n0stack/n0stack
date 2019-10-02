package grpccmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli"
	piam "n0st.ac/n0stack/iam/v1alpha"
)

func TestGenerateFlag(t *testing.T) {
	flags := GenerateFlags(piam.UserServiceClient.CreateUser, []string{})
	if diff := cmp.Diff(flags, []cli.Flag{
		cli.StringFlag{
			Name: "name",
		},
	}); diff != "" {
		t.Errorf("GenerateFlag response is wrong: diff=(-want +got)\n%s", diff)
	}

	t.Logf("%v", flags)
}

// func TestGenerateGRPCGetter(t *testing.T) {
// 	cli := pprovisioning.NewVirtualMachineServiceClient(&grpc.ClientConn{})
// 	GenerateGRPCGetter(cli.GetVirtualMachine)
// }
