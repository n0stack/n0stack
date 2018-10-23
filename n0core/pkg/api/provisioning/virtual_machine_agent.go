package provisioning

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"

	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu_img"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

var N0coreVirtualMachineNamespace uuid.UUID

const (
	QMPMonitorSocketFile   = "monitor.sock"
	VNCWebSocketPortOffset = 6900
)

func init() {
	N0coreVirtualMachineNamespace, _ = uuid.FromString("a015d18d-b2c3-4181-8028-6f707ef31c95")
}

type VirtualMachineAgentAPI struct {
	baseDirectory string
}

func CreateVirtualMachineAgentAPI(basedir string) (*VirtualMachineAgentAPI, error) {
	b, err := filepath.Abs(basedir)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get absolute path")
	}

	if _, err := os.Stat(b); os.IsNotExist(err) {
		if err := os.MkdirAll(b, 0644); err != nil { // TODO: check permission
			return nil, errors.Wrapf(err, "Failed to mkdir '%s'", b)
		}
	}

	return &VirtualMachineAgentAPI{
		baseDirectory: b,
	}, nil
}

func (a VirtualMachineAgentAPI) GetWorkDirectory(name string) (string, error) {
	p := filepath.Join(a.baseDirectory, name)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0644); err != nil { // TODO: check permission
			return p, errors.Wrapf(err, "Failed to mkdir '%s'", p)
		}
	}

	return p, nil
}

func (a VirtualMachineAgentAPI) CreateVirtualMachineAgent(ctx context.Context, req *CreateVirtualMachineAgentRequest) (res *VirtualMachineAgent, errRes error) {
	id := uuid.NewV5(N0coreVirtualMachineNamespace, req.Name)
	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if q.IsRunning() {
		return nil, grpc.Errorf(codes.AlreadyExists, "")
	}

	wd, err := a.GetWorkDirectory(req.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get working directory '%s'", wd)
	}
	websocket := qemu.GetNewListenPort(VNCWebSocketPortOffset)

	if err := q.Start(req.Name, filepath.Join(wd, QMPMonitorSocketFile), websocket, req.Vcpus, req.MemoryBytes); err != nil {
		log.Printf("Failed to start qemu process: err=%s", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to start qemu process")
	}
	defer q.Close()

	createdNetdev := []*NetDev{}
	for _, nd := range req.Netdev {
		b, err := iproute2.NewBridge(TrimNetdevName(nd.NetworkName))
		if err != nil {
			log.Printf("Failed to create bridge '%s': err='%s'", nd.NetworkName, err.Error())
			errRes = grpc.Errorf(codes.Internal, "") // TODO #89
			goto DeleteNetDev
		}

		t, err := iproute2.NewTap(TrimNetdevName(nd.Name))
		if err != nil {
			log.Printf("Failed to create tap '%s': err='%s'", nd.Name, err.Error())
			errRes = grpc.Errorf(codes.Internal, "") // TODO #89
			goto DeleteNetDev
		}
		createdNetdev = append(createdNetdev, nd)

		if err := t.SetMaster(b); err != nil {
			log.Printf("Failed to set master of tap '%s' as '%s': err='%s'", t.Name(), b.Name(), err.Error())
			errRes = grpc.Errorf(codes.Internal, "") // TODO #89
			goto DeleteNetDev
		}

		hw, err := net.ParseMAC(nd.HardwareAddress)
		if err != nil {
			errRes = grpc.Errorf(codes.InvalidArgument, "Hardware address '%s' is invalid on netdev '%s'", nd.HardwareAddress, nd.Name)
			goto DeleteNetDev
		}

		if err := q.AttachTap(nd.Name, t.Name(), hw); err != nil {
			log.Printf("Failed to attach tap: err='%s'", err.Error())
			errRes = grpc.Errorf(codes.Internal, "Failed to attach tap")
			goto DeleteNetDev
		}
	}

	for _, bd := range req.Blockdev {
		u, err := url.Parse(bd.Url)
		if err != nil {
			errRes = grpc.Errorf(codes.InvalidArgument, "url '%s' is invalid url: '%s'", bd.Url, err.Error())
			goto DeleteNetDev
		}

		i, err := img.OpenQemuImg(u.Path)
		if err != nil {
			log.Printf("Failed to open qemu image: err='%s'", err.Error())
			errRes = grpc.Errorf(codes.Internal, "") // TODO #89
			goto DeleteNetDev
		}

		// この条件は雑
		if i.Info.Format == "raw" {
			if err := q.AttachISO(bd.Name, u, uint(bd.BootIndex)); err != nil {
				log.Printf("Failed to attach iso '%s': err='%s'", u.Path, err.Error())
				errRes = grpc.Errorf(codes.Internal, "") // TODO #89
				goto DeleteNetDev
			}
		} else {
			if err := q.AttachQcow2(bd.Name, u, uint(bd.BootIndex)); err != nil {
				log.Printf("Failed to attach image '%s': err='%s'", u.String(), err.Error())
				errRes = grpc.Errorf(codes.Internal, "") // TODO #89
				goto DeleteNetDev
			}
		}
	}

	if err := q.Boot(); err != nil {
		log.Printf("Failed to boot qemu: err=%s", err.Error())
		errRes = grpc.Errorf(codes.Internal, "Failed to boot qemu")
		goto DeleteNetDev

	}

	res = &VirtualMachineAgent{
		Name:          req.Name,
		Uuid:          id.String(),
		Vcpus:         req.Vcpus,
		MemoryBytes:   req.MemoryBytes,
		Blockdev:      req.Blockdev,
		Netdev:        req.Netdev,
		WebsocketPort: uint32(websocket),
	}
	if s, err := q.Status(); err != nil {
		errRes = grpc.Errorf(codes.Internal, "Failed to get status")
		goto DeleteNetDev
	} else {
		res.State = GetAgentStateFromQemuState(s)
	}

	return

DeleteNetDev:
	if err := q.Delete(); err != nil {
		log.Printf("Failed to delete qemu: err=%s", err.Error())
	}

	for _, nd := range createdNetdev {
		t, err := iproute2.NewTap(TrimNetdevName(nd.Name))
		if err != nil {
			log.Printf("Failed to create tap '%s': err='%s'", nd.Name, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
		}

		if err := t.Delete(); err != nil {
			log.Printf("Failed to delete tap '%s': err='%s'", nd.Name, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
		}

		b, err := iproute2.NewBridge(TrimNetdevName(nd.NetworkName))
		if err != nil {
			log.Printf("Failed to create bridge '%s': err='%s'", nd.NetworkName, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
		}

		links, err := b.ListSlaves()
		if err != nil {
			log.Printf("Failed to list links of bridge '%s': err='%s'", nd.NetworkName, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
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
				log.Printf("Failed to delete bridge '%s': err='%s'", b.Name(), err.Error())
				return nil, grpc.Errorf(codes.Internal, "") // TODO #89
			}
		}
	}

	return
}

func (a VirtualMachineAgentAPI) DeleteVirtualMachineAgent(ctx context.Context, req *DeleteVirtualMachineAgentRequest) (*empty.Empty, error) {
	id := uuid.NewV5(N0coreVirtualMachineNamespace, req.Name)
	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if !q.IsRunning() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	if err := q.Delete(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to delete qemu: %s", err.Error())
	}

	for _, nd := range req.Netdev {
		t, err := iproute2.NewTap(TrimNetdevName(nd.Name))
		if err != nil {
			log.Printf("Failed to create tap '%s': err='%s'", nd.Name, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
		}

		if err := t.Delete(); err != nil {
			log.Printf("Failed to delete tap '%s': err='%s'", nd.Name, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
		}

		b, err := iproute2.NewBridge(TrimNetdevName(nd.NetworkName))
		if err != nil {
			log.Printf("Failed to create bridge '%s': err='%s'", nd.NetworkName, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
		}

		links, err := b.ListSlaves()
		if err != nil {
			log.Printf("Failed to list links of bridge '%s': err='%s'", nd.NetworkName, err.Error())
			return nil, grpc.Errorf(codes.Internal, "") // TODO #89
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
				log.Printf("Failed to delete bridge '%s': err='%s'", b.Name(), err.Error())
				return nil, grpc.Errorf(codes.Internal, "") // TODO #89
			}
		}
	}

	return &empty.Empty{}, nil
}

func (a VirtualMachineAgentAPI) BootVirtualMachineAgent(ctx context.Context, req *BootVirtualMachineAgentRequest) (*BootVirtualMachineAgentResponse, error) {
	id := uuid.NewV5(N0coreVirtualMachineNamespace, req.Name)
	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if !q.IsRunning() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}
	defer q.Close()

	if err := q.Boot(); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to boot qemu: %s", err.Error())
	}

	s, err := q.Status()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get qemu status: %s", err.Error())
	}

	return &BootVirtualMachineAgentResponse{
		State: GetAgentStateFromQemuState(s),
	}, nil
}

func (a VirtualMachineAgentAPI) RebootVirtualMachineAgent(ctx context.Context, req *RebootVirtualMachineAgentRequest) (*RebootVirtualMachineAgentResponse, error) {
	id := uuid.NewV5(N0coreVirtualMachineNamespace, req.Name)
	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if !q.IsRunning() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}
	defer q.Close()

	if req.Hard {
		if err := q.HardReset(); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to hard reboot qemu: %s", err.Error())
		}
	} else {
		return nil, grpc.Errorf(codes.Unimplemented, "reboot is unimplemented")
	}

	s, err := q.Status()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get qemu status: %s", err.Error())
	}

	return &RebootVirtualMachineAgentResponse{
		State: GetAgentStateFromQemuState(s),
	}, nil
}

func (a VirtualMachineAgentAPI) ShutdownVirtualMachineAgent(ctx context.Context, req *ShutdownVirtualMachineAgentRequest) (*ShutdownVirtualMachineAgentResponse, error) {
	id := uuid.NewV5(N0coreVirtualMachineNamespace, req.Name)
	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if !q.IsRunning() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}
	defer q.Close()

	if req.Hard {
		if err := q.HardShutdown(); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to hard shutdown qemu: %s", err.Error())
		}
	} else {
		if err := q.Shutdown(); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to shutdown qemu: %s", err.Error())
		}
	}

	s, err := q.Status()
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get qemu status: %s", err.Error())
	}

	return &ShutdownVirtualMachineAgentResponse{
		State: GetAgentStateFromQemuState(s),
	}, nil
}
