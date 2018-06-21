package compute

import (
	"context"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"

	"github.com/n0stack/n0core/provisioning/node"
	"github.com/n0stack/n0core/provisioning/node/iproute2"
	"github.com/n0stack/n0core/provisioning/node/kvm"
	uuid "github.com/satori/go.uuid"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/datastore"
	"github.com/n0stack/proto.go/provisioning/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	VMNamespace                = "6565cf73-9845-49d7-bdd7-76970e5ebe4f"
	AnnotationNodeName         = "n0core/node_name"
	AnnotationQmpPath          = "n0core/qmp_path"
	AnnotationVncWebsocketPort = "n0core/vnc_websocket_port"

	ApplyIproute2Deadline = time.Second * 10
	ApplyKVMDeadline      = time.Second * 10
)

type ComputeAPI struct {
	dataStore       datastore.Datastore
	nodeConnections *node.NodeConnections
	qmpBaseDir      string
}

func CreateComputeAPI(ds datastore.Datastore, nc *node.NodeConnections, qmpBaseDir string) (*ComputeAPI, error) {
	return &ComputeAPI{
		dataStore:       ds,
		nodeConnections: nc,
		qmpBaseDir:      qmpBaseDir,
	}, nil
}

func (a ComputeAPI) ListComputes(ctx context.Context, req *pprovisioning.ListComputesRequest) (*pprovisioning.ListComputesResponse, error) {
	res := &pprovisioning.ListComputesResponse{}
	f := func(s int) []proto.Message {
		res.Computes = make([]*pprovisioning.Compute, s)
		for i := range res.Computes {
			res.Computes[i] = &pprovisioning.Compute{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Computes {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}
	if len(res.Computes) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a ComputeAPI) GetCompute(ctx context.Context, req *pprovisioning.GetComputeRequest) (*pprovisioning.Compute, error) {
	res := &pprovisioning.Compute{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "message:Failed to get from db\tgot:%v", err.Error())
	}

	if res.Metadata == nil {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a ComputeAPI) genQmpPath(name string) string {
	return filepath.Join(a.qmpBaseDir, name, "monitor.sock")
}

// well known portと /proc/sys/net/ipv4/ip_local_port_range を避ける
// TODO: overflowしないか心配
func (a ComputeAPI) genVncWebsocketPort() uint32 {
	const LinuxIPLocalPortMin = 32768
	const LinuxIPLocalPortMax = 60999

	return rand.Uint32()%(LinuxIPLocalPortMax-LinuxIPLocalPortMin) + LinuxIPLocalPortMin
}

func (a ComputeAPI) ApplyCompute(ctx context.Context, req *pprovisioning.ApplyComputeRequest) (*pprovisioning.Compute, error) {
	prev := &pprovisioning.Compute{}
	err := a.dataStore.Get(req.Metadata.Name, prev)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get db, got:%v.", err.Error())
	}
	if prev.Metadata == nil && req.Metadata.Version != 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set the same version as Get result, have:%d, want:0.", req.Metadata.Version)
	}
	if prev.Metadata != nil && req.Metadata.Version != prev.Metadata.Version {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set the same version as Get result, have:%d, want:%d.", req.Metadata.Version, prev.Metadata.Version)
	}

	res := &pprovisioning.Compute{
		Metadata: req.Metadata,
		Spec:     req.Spec,
	}
	if prev.Status == nil {
		res.Status = &pprovisioning.ComputeStatus{}
	} else {
		res.Status = prev.Status
	}

	nn, ok := req.Metadata.Annotations[AnnotationNodeName]
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set '%s'", AnnotationNodeName)
	}

	conn, err := a.nodeConnections.GetConnection(nn)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Fail to dial to node, err:%v node_name:%s", err.Error(), nn)
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set valid '%s'", AnnotationNodeName) // TODO: エラーメッセージを何とかする
	}
	defer conn.Close()

	// volume
	// get volume
	// use url
	// set component as used

	// iproute2
	ipCli := iproute2.NewIproute2ServiceClient(conn)
	for l, n := range req.Spec.Nics {
		ctx, cancel := context.WithTimeout(context.Background(), ApplyIproute2Deadline)
		defer cancel()

		reqTap := &iproute2.ApplyTapRequest{
			Tap: &iproute2.Tap{
				Name:       l,
				BridgeName: n.NetworkName, // TODO: 多分長さ制限で引っかかったりするので何か関数などでラップする必要がある
				Type:       iproute2.Tap_FLAT,
			},
		}

		if resTap, err := ipCli.ApplyTap(ctx, reqTap); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to ApplyTap, err:'%s', req:'%s', res'%s'", err.Error(), reqTap, resTap)
		}
	}

	// kvm
	ns, _ := uuid.FromString(VMNamespace)
	vmID := uuid.NewV5(ns, req.Metadata.Name)

	qp, ok := res.Metadata.Annotations[AnnotationQmpPath]
	if !ok {
		qp = a.genQmpPath(req.Metadata.Name)
	}

	svwp, ok := res.Metadata.Annotations[AnnotationVncWebsocketPort]
	var vwp uint32
	if ok {
		vwp64, err := strconv.ParseInt(svwp, 10, 32)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Set unsigned integer to '%s', err:'%s', have:'%s'", AnnotationVncWebsocketPort, err.Error(), svwp)
		}

		vwp = uint32(vwp64)
	} else {
		vwp = a.genVncWebsocketPort()
	}

	reqKVM := &kvm.KVM{
		Name: req.Metadata.Name,
		// generate from name
		Uuid:             vmID.String(),
		CpuCores:         req.Spec.Vcpus,
		MemoryBytes:      req.Spec.MemoryBytes,
		Nics:             map[string]*kvm.KVM_NIC{},
		Storages:         map[string]*kvm.KVM_Storage{},
		QmpPath:          qp,
		VncWebsocketPort: vwp,
	}
	for l, n := range req.Spec.Nics {
		reqKVM.Nics[l] = &kvm.KVM_NIC{
			TapName: l,
			HwAddr:  n.HardwareAddress,
		}
	}
	for l, v := range req.Spec.Volumes {
		reqKVM.Storages[l] = &kvm.KVM_Storage{
			Url: v.VolumeName,

			// とりあえずvolumeがひとつの場合のみをサポート
			BootIndex: 1,
		}
	}

	kvmCli := kvm.NewKVMServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ApplyKVMDeadline)
	defer cancel()

	resKVM, err := kvmCli.ApplyKVM(ctx, &kvm.ApplyKVMRequest{Kvm: reqKVM})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to apply kvm, err:'%s', request: '%s', response: '%s'", err.Error(), reqKVM, resKVM)
	}

	res.Metadata.Annotations[AnnotationQmpPath] = resKVM.QmpPath
	res.Metadata.Annotations[AnnotationVncWebsocketPort] = strconv.FormatUint(uint64(resKVM.VncWebsocketPort), 10)

	res.Metadata.Version++
	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to apply for db, got:%v.", err.Error())
	}
	log.Printf("[INFO] On Applly, applied Network:%v", res)

	return res, nil
}

func (a ComputeAPI) DeleteCompute(ctx context.Context, req *pprovisioning.DeleteComputeRequest) (*empty.Empty, error) {
	prev := &pprovisioning.Compute{}
	err := a.dataStore.Get(req.Name, prev)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get db, got:%v.", err.Error())
	}

	nn, ok := prev.Metadata.Annotations[AnnotationNodeName]
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set '%s'", AnnotationNodeName)
	}

	conn, err := a.nodeConnections.GetConnection(nn)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Fail to dial to node, err:%v node_name:%s", err.Error(), nn)
	}
	if conn == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set valid '%s'", AnnotationNodeName) // TODO: エラーメッセージを何とかする
	}
	defer conn.Close()

	// kvm
	kvmCli := kvm.NewKVMServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), ApplyKVMDeadline)
	defer cancel()

	if _, err := kvmCli.DeleteKVM(ctx, &kvm.DeleteKVMRequest{Name: req.Name}); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to apply kvm, err:'%s', name: '%s'", err.Error(), req.Name)
	}

	// iproute2
	ipCli := iproute2.NewIproute2ServiceClient(conn)
	for l := range prev.Spec.Nics {
		ctx, cancel := context.WithTimeout(context.Background(), ApplyIproute2Deadline)
		defer cancel()

		if _, err := ipCli.DeleteTap(ctx, &iproute2.DeleteTapRequest{Name: l}); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to DeleteTap, err:'%s', name:'%s'", err.Error(), l)
		}
	}

	d, err := a.dataStore.Delete(req.Name)
	if err != nil {
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "message:Failed to delete from db.\tgot:%v", err.Error())
	}
	if d < 1 {
		return &empty.Empty{}, grpc.Errorf(codes.NotFound, "")
	}

	return &empty.Empty{}, nil
}

func (a ComputeAPI) WatchCompute(req *pprovisioning.WatchComputesRequest, res pprovisioning.ComputeService_WatchComputeServer) error {
	return grpc.Errorf(codes.Unimplemented, "")
}

func (a ComputeAPI) Boot(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a ComputeAPI) Reboot(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a ComputeAPI) HardReboot(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a ComputeAPI) Shutdown(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a ComputeAPI) HardShutdown(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}

func (a ComputeAPI) Save(ctx context.Context, req *pprovisioning.ActionComputeRequest) (*pprovisioning.Compute, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}
