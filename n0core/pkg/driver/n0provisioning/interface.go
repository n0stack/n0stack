package n0provisioning

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
)

var PbMarshaler = &jsonpb.Marshaler{
	EnumsAsInts:  true,
	EmitDefaults: false,
	OrigName:     true,
}

type ExecuteUnit interface {
	Run(ctx context.Context, conn *grpc.ClientConn) (proto.Message, error)
	String() string
}

type GPRCUnit struct {
	runMethod    func(ctx context.Context, conn *grpc.ClientConn) (proto.Message, error)
	stringMethod func() string
}

func (e GPRCUnit) Run(ctx context.Context, conn *grpc.ClientConn) (proto.Message, error) {
	return e.runMethod(ctx, conn)
}
func (e GPRCUnit) String() string {
	return e.stringMethod()
}

//go:generate gen_grpc_units -package=n0st.ac/n0stack/n0proto.go/provisioning/v0 -service=VirtualMachine
//go:generate gen_grpc_units -package=n0st.ac/n0stack/n0proto.go/provisioning/v0 -service=BlockStorage
