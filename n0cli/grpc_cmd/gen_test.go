package grpccmd

// func TestGenerateFlag(t *testing.T) {
// 	flags := GenerateFlags(pprovisioning.GetVirtualMachineRequest{})
// 	if diff := cmp.Diff(flags, []cli.Flag{
// 		cli.StringFlag{
// 			Name: "name",
// 		},
// 	}); diff != "" {
// 		t.Errorf("GenerateFlag response is wrong: diff=(-want +got)\n%s", diff)
// 	}

// 	t.Logf("%v", flags)
// }

// func TestGenerateGRPCGetter(t *testing.T) {
// 	cli := pprovisioning.NewVirtualMachineServiceClient(&grpc.ClientConn{})
// 	GenerateGRPCGetter(cli.GetVirtualMachine)
// }
