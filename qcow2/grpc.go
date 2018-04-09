package qcow2

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/n0stack/n0core/notification"
	qcow2 "github.com/n0stack/proto.go/qcow2/v0"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const modelType = "device/volume/qcow2"

type Agent struct{}

func (a Agent) Get(ctx context.Context, req *qcow2.GetRequest) (r *qcow2.Qcow2, errRes error) {
	var v volume
	r = &qcow2.Qcow2{
		Id:     req.Id,
		Spec:   &v.Spec,
		Status: &v.Status,
	}

	var err error
	v.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(v.getQcow2())

	return
}

func (a Agent) Apply(ctx context.Context, req *qcow2.ApplyRequest) (q *qcow2.Qcow2, errRes error) {
	var v volume
	q = &qcow2.Qcow2{
		Id:     req.Id,
		Spec:   &v.Spec,
		Status: &v.Status,
	}

	var err error
	v.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(v.getQcow2())

	if v.Bytes == 0 {
		notification.Notify(v.createImage(req.Spec.Bytes))
		notification.Notify(v.getQcow2(req.Spec.SoueceUrl))
	} else {
		notification.Notify(v.resizeImage(req.Spec.Bytes))
	}

	return
}

func (a Agent) Delete(ctx context.Context, req *qcow2.DeleteRequest) (e *empty.Empty, errRes error) {
	var v volume
	e = &empty.Empty{}

	var err error
	v.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(v.getQcow2())

	if v.Bytes != 0 {
		notification.Notify(v.deleteImage()) // -> deleteWorkdir にする必要がある
	}

	return
}
