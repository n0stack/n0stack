package flat

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/satori/go.uuid"

	"google.golang.org/grpc/codes"

	n0stack "github.com/n0stack/proto"
	"github.com/n0stack/proto/resource/networkid/v0"
	"google.golang.org/grpc"

	network "github.com/n0stack/proto/device/network/v0"
)

type Agent struct {
	// DB *gorm.DB
	InterfaceName string
}

const (
	modelType = "device/network/flat"
)

func notify(n *n0stack.Notification) {
	if !n.Success {
		panic(n)
	}
}

// Apply スペックを元にステートレスに適用する
func (a Agent) Apply(ctx context.Context, req *network.ApplyRequest) (r *network.Resource, errRes error) {
	var f flat
	r = &network.Resource{
		Model:  req.Model,
		Spec:   &f.Spec,
		Status: &f.Status,
	}

	// Validation
	if req.Spec.NetworkId.Type != networkid.Spec_FLAT {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Supporting type is FLAT\tgot:%v", req.Spec.NetworkId.Type)
		return
	}

	var err error
	f.id, err = uuid.FromBytes(req.Model.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Model.Id)
		return
	}

	notify(f.getBridge(a.InterfaceName))

	if f.Status.Bridge == "" {
		notify(f.createBridge())
		notify(f.getBridge(a.InterfaceName))
	}

	notify(f.applyBridge())

	return
}

func (a Agent) Delete(ctx context.Context, req *network.DeleteRequest) (e *empty.Empty, errRes error) {
	var f flat
	e = &empty.Empty{}

	// Validation
	var err error
	f.id, err = uuid.FromBytes(req.Model.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Model.Id)
		return
	}

	notify(f.getBridge(a.InterfaceName))

	if f.Status.Bridge == "" {
		errRes = grpc.Errorf(codes.NotFound, "message:Already deleted")
		return
	}

	f.deleteBridge()

	return
}
