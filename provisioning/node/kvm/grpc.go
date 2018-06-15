package kvm

import (
	"log"
	"net"
	"net/url"

	"github.com/digitalocean/go-qemu/qmp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	uuid "github.com/satori/go.uuid"
	context "golang.org/x/net/context"
)

type KVMAgent struct {
	qmp map[string]*qmp.SocketMonitor
}

func NewKVMAgent() (*KVMAgent, error) {
	return &KVMAgent{
		qmp: map[string]*qmp.SocketMonitor{},
	}, nil
}

func (a KVMAgent) ApplyKVM(ctx context.Context, req *ApplyKVMRequest) (*KVM, error) {
	// validation

	p, err := a.getProcess(req.Kvm.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to getProcess, err:'%s'", err.Error())
	}

	started := false
	if p == nil {
		u, err := uuid.FromString(req.Kvm.Uuid)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse uuid, err:'%s', uuid:'%s'", err.Error(), req.Kvm.Uuid)
		}

		err = a.startProcess(
			u,
			req.Kvm.Name,
			req.Kvm.QmpPath,
			req.Kvm.VncWebsocketPort,
			req.Kvm.CpuCores,
			req.Kvm.MemoryBytes,
		)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to startProcess, err:'%s'", err.Error())
		}

		started = true
	}

	log.Printf("[DEBUG] before connectQMP")

	q, err := a.connectQMP(req.Kvm.Name, req.Kvm.QmpPath)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to connectQMP, err:'%s'", err.Error())
	}

	log.Printf("[DEBUG] after connectQMP")

	// if started {
	// 	a.startCheckEvents()
	// }

	// Storage
	for label, v := range req.Kvm.Storages {
		index := v.BootIndex
		u, err := url.Parse(v.Url)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse url, err:'%s', url:'%s'", err.Error(), v.Url)
		}

		if err := a.attachStorage(q, label, u, index); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to attachStorage, err:'%s'", err.Error())
		}
	}

	log.Printf("[DEBUG] before NIC")

	// Network
	for label, v := range req.Kvm.Nics {
		m, err := net.ParseMAC(v.HwAddr)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Failed to parse hardware address, err:'%s', hwaddr:'%s'", err.Error(), v.HwAddr)
		}

		if err := a.attachNIC(q, label, v.TapName, m); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to attachNIC, err:'%s'", err.Error())
		}
	}

	log.Printf("[DEBUG] after NIC")

	if started {
		if err := a.boot(q); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to Boot, err:'%s'", err.Error())
		}
		log.Printf("[DEBUG] after Boot")
	}

	return req.Kvm, nil
}

func (a KVMAgent) DeleteKVM(ctx context.Context, req *DeleteKVMRequest) (*google_protobuf.Empty, error) {
	if v, ok := a.qmp[req.Name]; ok {
		defer v.Disconnect()
		delete(a.qmp, req.Name)
	}

	p, err := a.getProcess(req.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to getProcess, err:'%s'", err.Error())
	}

	if err := p.Kill(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to kill process, err:'%s', pid:'%d'", err.Error(), p.Pid)
	}

	return &google_protobuf.Empty{}, nil
}

// (QEMU) cont
func (a KVMAgent) Boot(ctx context.Context, req *ActionKVMRequest) (*google_protobuf.Empty, error) {
	q, err := a.connectQMP(req.Name, req.QmpPath)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to connectQMP, err:'%s'", err.Error())
	}

	if err := a.boot(q); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to boot, err:'%s'", err.Error())
	}

	return &google_protobuf.Empty{}, nil
}

func (a KVMAgent) Reboot(context.Context, *ActionKVMRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

func (a KVMAgent) HardReboot(context.Context, *ActionKVMRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

func (a KVMAgent) Shutdown(context.Context, *ActionKVMRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

func (a KVMAgent) HardShutdown(context.Context, *ActionKVMRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}

func (a KVMAgent) Save(context.Context, *ActionKVMRequest) (*google_protobuf.Empty, error) {
	return &google_protobuf.Empty{}, nil
}
