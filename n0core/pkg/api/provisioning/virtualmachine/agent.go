package virtualmachine

import (
	"context"
	"log"
	"os"
	"path/filepath"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"

	"github.com/n0stack/n0stack/bazel-n0stack/n0core/pkg/driver/iproute2"
	"github.com/n0stack/n0stack/bazel-n0stack/n0core/pkg/driver/iptables"
	"github.com/n0stack/n0stack/bazel-n0stack/n0core/pkg/driver/qemu"
	"github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	"github.com/n0stack/n0stack/n0core/pkg/util/net"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
)

const (
	QmpMonitorSocketFile   = "monitor.sock"
	VNCWebSocketPortOffset = 6900
)

type VirtualMachineAgent struct {
	baseDirectory string
}

func CreateVirtualMachineAgent(basedir string) (*VirtualMachineAgent, error) {
	b, err := filepath.Abs(basedir)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get absolute path")
	}

	if _, err := os.Stat(b); os.IsNotExist(err) {
		if err := os.MkdirAll(b, 0644); err != nil { // TODO: check permission
			return nil, errors.Wrapf(err, "Failed to mkdir '%s'", b)
		}
	}

	return &VirtualMachineAgent{
		baseDirectory: b,
	}, nil
}

func (a VirtualMachineAgent) GetWorkDirectory(name string) (string, error) {
	p := filepath.Join(a.baseDirectory, name)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0644); err != nil { // TODO: check permission
			return p, errors.Wrapf(err, "Failed to mkdir '%s'", p)
		}
	}

	return p, nil
}
func (a VirtualMachineAgent) DeleteWorkDirectory(name string) error {
	// p := filepath.Join(a.baseDirectory, name)

	// if _, err := os.Stat(p); os.IsNotExist(err) {
	// 	if err := os.MkdirAll(p, 0644); err != nil { // TODO: check permission
	// 		return p, errors.Wrapf(err, "Failed to mkdir '%s'", p)
	// 	}
	// }

	return nil
}

func (a VirtualMachineAgent) BootVirtualMachine(ctx context.Context, req *BootVirtualMachineRequest) (*BootVirtualMachineResponse, error) {
	name := req.Name
	id, err := uuid.FromString(req.Uuid)
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set valid uuid: %s", err.Error())
	}
	websocket := qemu.GetNewListenPort(VNCWebSocketPortOffset)
	vcpus := req.Vcpus
	mem := req.MemoryBytes

	tx := transaction.Begin()
	defer transaction.WrapRollbackError(tx.Rollback())

	wd, err := a.GetWorkDirectory(name)
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get working directory '%s'", wd)
	}

	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if q.IsRunning() {
		return nil, grpcutil.WrapGrpcErrorf(codes.AlreadyExists, "Qemu process is already running")
	}

	if err := q.Start(name, filepath.Join(wd, QmpMonitorSocketFile), websocket, vcpus, mem); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to start qemu process: err=%s", err.Error())
	}
	defer q.Close()
	tx.PushRollback("delete Qemu", func() error {
		return q.Delete()
	})

	tx.Commit()
	return &BootVirtualMachineResponse{
		// State: ,
		WebsocketPort: uint32(websocket),
	}, nil
}

func (a VirtualMachineAgent) RebootVirtualMachine(ctx context.Context, req *RebootVirtualMachineRequest) (*RebootVirtualMachineResponse, error) {
	return nil, grpcutil.WrapGrpcErrorf(codes.Unimplemented, "")
}

func (a VirtualMachineAgent) ShutdownVirtualMachine(ctx context.Context, req *ShutdownVirtualMachineRequest) (*ShutdownVirtualMachineResponse, error) {
	return nil, grpcutil.WrapGrpcErrorf(codes.Unimplemented, "")
}

func (a VirtualMachineAgent) DeleteVirtualMachine(ctx context.Context, req *DeleteVirtualMachineRequest) (*empty.Empty, error) {
	id := uuid.NewV4()

	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if !q.IsRunning() {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}
	if err := q.Delete(); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to delete qemu: %s", err.Error())
	}

	for _, nd := range req.Netdevs {
		t, err := iproute2.NewTap(netutil.StructLinuxNetdevName(nd.Name))
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to create tap '%s'", nd.Name).Error())
		}

		if err := t.Delete(); err != nil {
			log.Printf("Failed to delete tap '%s': err='%s'", nd.Name, err.Error())
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "") // TODO #89
		}

		b, err := iproute2.NewBridge(netutil.StructLinuxNetdevName(nd.NetworkName))
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to create bridge '%s'", nd.NetworkName).Error())
		}

		links, err := b.ListSlaves()
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to list links of bridge '%s'", nd.NetworkName).Error())
		}

		// TODO: 以下遅い気がする
		i := 0
		for _, l := range links {
			if _, err := iproute2.NewTap(l); err == nil {
				i++
			}
		}
		if i == 0 {
			if err := b.Delete(); err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to delete bridge '%s'", b.Name()).Error())
			}

			// gateway settings
			if nd.Ipv4Gateway != "" {
				ip := netutil.ParseCIDR(nd.Ipv4AddressCidr)
				if ip == nil {
					return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set valid ipv4_address_cidr: value='%s'", nd.Ipv4AddressCidr)
				}
				if err := iptables.DeleteMasqueradeRule(b.Name(), ip.Network()); err != nil {
					return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to delete masquerade rule").Error())
				}
			}
		}
	}

	return &empty.Empty{}, nil
}

// func GetAgentStateFromQemuState(s qemu.Status) VirtualMachineAgentState {
// 	switch s {
// 	case qemu.StatusRunning:
// 		return VirtualMachineAgentState_RUNNING

// 	case qemu.StatusShutdown, qemu.StatusGuestPanicked, qemu.StatusPreLaunch:
// 		return VirtualMachineAgentState_SHUTDOWN

// 	case qemu.StatusPaused, qemu.StatusSuspended:
// 		return VirtualMachineAgentState_PAUSED

// 	case qemu.StatusInternalError, qemu.StatusIOError:
// 		return VirtualMachineAgentState_FAILED

// 	case qemu.StatusInMigrate:
// 	case qemu.StatusFinishMigrate:
// 	case qemu.StatusPostMigrate:
// 	case qemu.StatusRestoreVM:
// 	case qemu.StatusSaveVM: // TODO: 多分PAUSED
// 	case qemu.StatusWatchdog:
// 	case qemu.StatusDebug:
// 	}

// 	return VirtualMachineAgentState_UNKNOWN
// }
