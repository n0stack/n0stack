package qcow2

import (
	"context"

	"github.com/n0stack/n0core/lib"
	"github.com/n0stack/proto/device/volume"

	n0stack "github.com/n0stack/proto"
)

const modelType = "device/volume/qcow2"

type Agent struct{}

func (a Agent) Get(ctx context.Context, model *n0stack.Model) (r *volume.Response, err error) {
	var (
		n *n0stack.Notification
		q qcow2
	)
	r = &volume.Response{
		Notification: n,
		Spec:         &q.Spec,
		Status:       &q.Status,
	}

	if n = q.getStatus(model); !n.Success {
		return
	}

	return
}

func (a Agent) Apply(ctx context.Context, spec *volume.Spec) (*volume.Response, error) {
	var (
		n *n0stack.Notification
		q qcow2
	)
	v := &volume.Response{
		Notification: n,
		Spec:         &q.Spec,
		Status:       &q.Status,
	}

	if n = q.getStatus(spec.Model); !n.Success {
		return v, nil
	}

	if q.Storage == nil {
		if n = q.createImage(spec.Storage.Bytes); !n.Success {
			return v, nil
		}

		if n = q.getStatus(spec.Model); !n.Success {
			return v, nil
		}
	} else {
		if n = q.resizeImage(spec.Storage.Bytes); !n.Success {
			return v, nil
		}
	}

	v.Notification = lib.MakeNotification("Apply", true, "")
	return v, nil
}

func (a Agent) Delete(ctx context.Context, model *n0stack.Model) (*volume.Response, error) {
	var (
		n *n0stack.Notification
		q qcow2
	)
	v := &volume.Response{
		Notification: n,
		Spec:         &q.Spec,
		Status:       &q.Status,
	}

	if n = q.getStatus(model); !n.Success {
		return v, nil
		// return v, grpc.Errorf(codes.InvalidArgument, v.Notification.Description)
	}

	if q.Storage != nil {
		if n = q.deleteImage(); !n.Success {
			return v, nil
		}
	}

	v.Notification = lib.MakeNotification("Delete", true, "")
	return v, nil
}
