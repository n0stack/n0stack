package flat

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/satori/go.uuid"

	"google.golang.org/grpc/codes"

	n0stack "github.com/n0stack/go-proto"
	tap "github.com/n0stack/go-proto/device/tap/v0"
	"github.com/n0stack/go-proto/resource/networkid/v0"
	"google.golang.org/grpc"
)

type Agent struct {
	// DB *gorm.DB
	InterfaceName string
}

const (
	modelType = "device/tap/flat"
)

func notify(n *n0stack.Notification) {
	if !n.Success {
		panic(n)
	}
}

// Apply スペックを元にステートレスに適用する
func (a Agent) Apply(ctx context.Context, req *tap.ApplyRequest) (r *tap.Resource, errRes error) {
	var f flat
	r = &tap.Resource{
		Model:  req.Model,
		Spec:   &f.Spec,
		Status: &f.Status,
	}

	// Validation
	if req.Spec.NetworkID.Type != networkid.Spec_FLAT {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Supporting type is FLAT\tgot:%v", req.Spec.NetworkID.Type)
		return
	}

	var err error
	f.id, err = uuid.FromBytes(req.Model.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Model.Id)
		return
	}

	notify(f.getInterface(a.InterfaceName))
	notify(f.getBridge())

	if f.Bridge == "" {
		notify(f.createBridge())
		notify(f.getBridge())
	}

	notify(f.applyBridge())

	notify(f.getTap())
	if f.Tap == "" {
		notify(f.createTap())
		notify(f.getTap())
	}

	notify(f.applyTap())

	return
}

func (a Agent) Delete(ctx context.Context, req *tap.DeleteRequest) (e *empty.Empty, errRes error) {
	var f flat
	e = &empty.Empty{}

	// Validation
	var err error
	f.id, err = uuid.FromBytes(req.Model.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Model.Id)
		return
	}

	notify(f.getInterface(a.InterfaceName))
	notify(f.getBridge())
	notify(f.getTap())

	if f.Tap == "" { // 条件分岐かなり適当 no tap + bridge = fail, no tap + no bridge = fail
		errRes = grpc.Errorf(codes.NotFound, "message:Tap was already deleted")
	} else {
		f.deleteTap()
	}

	if f.Bridge == "" {
		errRes = grpc.Errorf(codes.NotFound, "message:Bridge was already deleted")
	} else {
		f.deleteBridge() // 接続しているtapが0個になった時に削除
	}

	return
}
