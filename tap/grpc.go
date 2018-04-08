package tap

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/notification"
	tap "github.com/n0stack/proto.go/tap/v0"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Agent struct {
	// DB *gorm.DB
	InterfaceName string
}

const (
	modelType = "tap"
)

func (a Agent) Get(ctx context.Context, req *tap.GetRequest) (r *tap.Tap, errRes error) {
	var f flat
	r = &tap.Tap{
		Id:     req.Id,
		Spec:   &f.Spec,
		Status: &f.Status,
	}

	// Validation
	var err error
	f.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(f.getInterface(a.InterfaceName))
	notification.Notify(f.getBridge())
	notification.Notify(f.getTap())

	return
}

// Apply スペックを元にステートレスに適用する
func (a Agent) Apply(ctx context.Context, req *tap.ApplyRequest) (r *tap.Tap, errRes error) {
	var f flat
	r = &tap.Tap{
		Id:     req.Id,
		Spec:   &f.Spec,
		Status: &f.Status,
	}

	// Validation
	if req.Spec.Type != tap.Spec_FLAT {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Supporting type is FLAT\tgot:%v", req.Spec.Type)
		return
	}

	var err error
	f.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(f.getInterface(a.InterfaceName))
	notification.Notify(f.getBridge())

	if f.Bridge == "" {
		notification.Notify(f.createBridge())
		notification.Notify(f.getBridge())
	}

	notification.Notify(f.applyBridge())

	notification.Notify(f.getTap())
	if f.Tap == "" {
		notification.Notify(f.createTap())
		notification.Notify(f.getTap())
	}

	notification.Notify(f.applyTap())

	return
}

func (a Agent) Delete(ctx context.Context, req *tap.DeleteRequest) (e *empty.Empty, errRes error) {
	var f flat
	e = &empty.Empty{}

	// Validation
	var err error
	f.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(f.getInterface(a.InterfaceName))
	notification.Notify(f.getBridge())
	notification.Notify(f.getTap())

	if f.Tap == "" { // 条件分岐かなり適当 no tap + bridge = fail, no tap + no bridge = fail
		errRes = grpc.Errorf(codes.NotFound, "message:Tap was already deleted")
	} else {
		notification.Notify(f.deleteTap())
	}

	if f.Bridge == "" {
		errRes = grpc.Errorf(codes.NotFound, "message:Bridge was already deleted")
	} else {
		notification.Notify(f.deleteBridge()) // TODO: 接続しているtapが0個になった時に削除
	}

	return
}
