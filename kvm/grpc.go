package kvm

import (
	"context"

	"github.com/n0stack/n0core/notification"
	uuid "github.com/satori/go.uuid"

	"github.com/golang/protobuf/ptypes/empty"
	pkvm "github.com/n0stack/proto.go/kvm/v0"
	pqcow2 "github.com/n0stack/proto.go/qcow2/v0"
	ptap "github.com/n0stack/proto.go/tap/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Agent struct {
	VNCPortMax uint
	VNCPortMin uint
}

const (
	modelType = "kvm"
)

func (a *Agent) Get(ctx context.Context, req *pkvm.GetRequest) (r *pkvm.KVM, errRes error) {
	return
}

func (a *Agent) Apply(ctx context.Context, req *pkvm.ApplyRequest) (r *pkvm.KVM, errRes error) {
	var k kvm
	r = &pkvm.KVM{
		Id:     req.Id,
		Spec:   &k.Spec,
		Status: &k.Status,
	}
	starting := false

	var err error
	k.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(k.getKVM())

	// if vm is not running
	if k.args == nil {
		// check CPU usage
		// check Memory usage

		notification.Notify(k.getVNCPort(a.VNCPortMin, a.VNCPortMax))
		notification.Notify(k.runVM(req.Spec.Vcpus, req.Spec.MemoryBytes))
		starting = true
	}

	notification.Notify(k.connectQMP())
	// 気持ち悪いが、どこでQMPをCloseするかがわかりにくいのでこうしている。要リファクタリング
	k.qmp.Connect()
	defer k.qmp.Disconnect()

	conn, err := grpc.Dial("localhost:20180", grpc.WithInsecure())
	if err != nil {
		return nil, nil
	}
	defer conn.Close()

	for i, v := range req.Spec.Volumes {
		cli := pqcow2.NewQcow2ServiceClient(conn)

		r, err := cli.Get(context.Background(), &pqcow2.GetRequest{Id: v})
		if err != nil {
			return nil, nil
		}

		notification.Notify(k.attachVolume(r.Id, r.Status.Url, i+1))
	}

	for _, n := range req.Spec.Nics {
		cli := ptap.NewTapServiceClient(conn)

		r, err := cli.Get(context.Background(), &ptap.GetRequest{Id: n.Tap})
		if err != nil {
			return nil, nil
		}

		notification.Notify(k.attachNIC(r.Id, r.Status.Tap, n.Hwaddr))
	}

	if starting {
		notification.Notify(k.bootVM())
	}

	return
}

func (a *Agent) Delete(ctx context.Context, req *pkvm.DeleteRequest) (e *empty.Empty, errRes error) {
	var k kvm
	e = &empty.Empty{}

	var err error
	k.id, err = uuid.FromString(req.Id)
	if err != nil {
		errRes = grpc.Errorf(codes.InvalidArgument, "message:Failed to validate uuid\tgot:%v", req.Id)
		return
	}

	notification.Notify(k.getKVM())

	// if vm is running
	if k.args != nil {
		// vmのシャットダウン

		notification.Notify(k.kill())
	}

	// remove networkはする

	// remove workdir

	return
}

func (a *Agent) Boot(ctx context.Context, req *pkvm.ActionRequest) (r *pkvm.KVM, errRes error) {
	return
}

func (a *Agent) Reboot(ctx context.Context, req *pkvm.ActionRequest) (r *pkvm.KVM, errRes error) {
	return
}

func (a *Agent) HardReboot(ctx context.Context, req *pkvm.ActionRequest) (r *pkvm.KVM, errRes error) {
	return
}

func (a *Agent) Shutdown(ctx context.Context, req *pkvm.ActionRequest) (r *pkvm.KVM, errRes error) {
	return
}

func (a *Agent) HardShutdown(ctx context.Context, req *pkvm.ActionRequest) (r *pkvm.KVM, errRes error) {
	return
}

func (a *Agent) Save(ctx context.Context, req *pkvm.ActionRequest) (r *pkvm.KVM, errRes error) {
	return
}
