package provisioning

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/driver/cloudinit/configdrive"
	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/n0stack/n0stack/n0core/pkg/driver/iptables"
	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu"
	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu_img"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
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
	cloudinit     []*configdrive.CloudConfig
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

func (a VirtualMachineAgentAPI) CreateVirtualMachineAgent(ctx context.Context, req *CreateVirtualMachineAgentRequest) (*VirtualMachineAgent, error) {
	var id uuid.UUID
	if req.Uuid == "" {
		id = uuid.NewV5(N0coreVirtualMachineNamespace, req.Name)
	} else {
		var err error
		id, err = uuid.FromString(req.Uuid)
		if err != nil {
			return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid uuid: %s", err.Error())
		}
	}

	tx := transaction.Begin()
	websocket := qemu.GetNewListenPort(VNCWebSocketPortOffset)

	wd, err := a.GetWorkDirectory(req.Name)
	if err != nil {
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to get working directory '%s'", wd)
	}

	q, err := qemu.OpenQemu(&id)
	if err != nil {
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if q.IsRunning() {
		return nil, WrapGrpcErrorf(codes.AlreadyExists, "")
	}

	if err := q.Start(req.Name, filepath.Join(wd, QMPMonitorSocketFile), websocket, req.Vcpus, req.MemoryBytes); err != nil {
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to start qemu process: err=%s", err.Error())
	}
	defer q.Close()
	tx.PushRollback("delete Qemu", func() error {
		return q.Delete()
	})

	createdNetdev := []*NetDev{}
	eth := make([]*configdrive.CloudConfigEthernet, len(req.Netdev))
	for i, nd := range req.Netdev {
		b, err := iproute2.NewBridge(TrimNetdevName(nd.NetworkName))
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.Internal, "Failed to create bridge '%s': err='%s'", nd.NetworkName, err.Error())
		}

		t, err := iproute2.NewTap(TrimNetdevName(nd.Name))
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.Internal, "Failed to create tap '%s': err='%s'", nd.Name, err.Error())
		}
		createdNetdev = append(createdNetdev, nd)

		if err := t.SetMaster(b); err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.Internal, "Failed to set master of tap '%s' as '%s': err='%s'", t.Name(), b.Name(), err.Error())
		}

		hw, err := net.ParseMAC(nd.HardwareAddress)
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.Internal, "Hardware address '%s' is invalid on netdev '%s'", nd.HardwareAddress, nd.Name)
		}

		if err := q.AttachTap(nd.Name, t.Name(), hw); err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.Internal, "Failed to attach tap: err='%s'", err.Error())
		}

		tx.PushRollback("delete created netdev", func() error {
			if err := t.Delete(); err != nil {
				return fmt.Errorf("Failed to delete tap '%s': err='%s'", nd.Name, err.Error())
			}

			links, err := b.ListSlaves()
			if err != nil {
				return fmt.Errorf("Failed to list links of bridge '%s': err='%s'", nd.NetworkName, err.Error())
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
					return fmt.Errorf("Failed to delete bridge '%s': err='%s'", b.Name(), err.Error())
				}
			}

			return nil
		})

		// Cloudinit settings
		ip, ipn, err := net.ParseCIDR(nd.Ipv4AddressCidr)
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid ipv4_address_cidr: value='%s', err='%s'", nd.Ipv4AddressCidr, err.Error())
		}
		hwaddr, err := net.ParseMAC(nd.HardwareAddress)
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid hardware_address: value='%s', err='%s'", nd.HardwareAddress, err.Error())
		}
		nameservers := make([]net.IP, len(nd.Nameservers))
		for i, n := range nd.Nameservers {
			nameservers[i] = net.ParseIP(n)
		}
		eth[i] = &configdrive.CloudConfigEthernet{
			MacAddress:  hwaddr,
			Address4:    ip,
			Network4:    ipn,
			Gateway4:    net.ParseIP(nd.Ipv4Gateway),
			NameServers: nameservers,
		}

		// Gateway settings
		if nd.Ipv4Gateway != "" {
			mask, _ := ipn.Mask.Size()
			gatewayIP := fmt.Sprintf("%s/%d", nd.Ipv4Gateway, mask)
			if err := b.SetAddress(gatewayIP); err != nil {
				WrapRollbackError(tx.Rollback())
				return nil, WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to set gateway IP to bridge: value=%s", gatewayIP).Error())
			}

			if iptables.CreateMasqueradeRule(b.Name(), ipn.String()); err != nil {
				WrapRollbackError(tx.Rollback())
				return nil, WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to create masquerade rule").Error())
			}
			tx.PushRollback("delete masquerade rule", func() error {
				return iptables.DeleteMasqueradeRule(b.Name(), ipn.String())
			})
		}
	}

	parsedKeys := make([]ssh.PublicKey, len(req.SshAuthorizedKeys))
	for i, k := range req.SshAuthorizedKeys {
		parsedKeys[i], _, _, _, err = ssh.ParseAuthorizedKey([]byte(k))
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.InvalidArgument, "ssh_authorized_keys is invalid: value='%s', err='%s'", k, err.Error())
		}
	}

	c := configdrive.StructConfig(req.LoginUsername, req.Name, parsedKeys, eth)
	p, err := c.Generate(wd)
	if err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to generate cloudinit configdrive:  err='%s'", err.Error())
	}
	req.Blockdev = append(req.Blockdev, &BlockDev{
		Name: "configdrive",
		Url: (&url.URL{
			Scheme: "file",
			Path:   p,
		}).String(),
		BootIndex: 50, // MEMO: 適当
	})

	for _, bd := range req.Blockdev {
		u, err := url.Parse(bd.Url)
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.InvalidArgument, "url '%s' is invalid url: '%s'", bd.Url, err.Error())
		}

		i, err := img.OpenQemuImg(u.Path)
		if err != nil {
			WrapRollbackError(tx.Rollback())
			return nil, WrapGrpcErrorf(codes.Internal, "Failed to open qemu image: err='%s'", err.Error())
		}

		// この条件は雑
		if i.Info.Format == "raw" {
			if bd.BootIndex < 3 {
				if err := q.AttachISO(bd.Name, u, uint(bd.BootIndex)); err != nil {
					WrapRollbackError(tx.Rollback())
					return nil, WrapGrpcErrorf(codes.Internal, "Failed to attach iso '%s': err='%s'", u.Path, err.Error())
				}
			} else {
				if err := q.AttachRaw(bd.Name, u, uint(bd.BootIndex)); err != nil {
					WrapRollbackError(tx.Rollback())
					return nil, WrapGrpcErrorf(codes.Internal, "Failed to attach raw '%s': err='%s'", u.Path, err.Error())
				}
			}
		} else {
			if err := q.AttachQcow2(bd.Name, u, uint(bd.BootIndex)); err != nil {
				WrapRollbackError(tx.Rollback())
				return nil, WrapGrpcErrorf(codes.Internal, "Failed to attach image '%s': err='%s'", u.String(), err.Error())
			}
		}
	}

	if err := q.Boot(); err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to boot qemu: err=%s", err.Error())

	}

	res := &VirtualMachineAgent{
		Name:          req.Name,
		Uuid:          id.String(),
		Vcpus:         req.Vcpus,
		MemoryBytes:   req.MemoryBytes,
		Blockdev:      req.Blockdev,
		Netdev:        req.Netdev,
		WebsocketPort: uint32(websocket),
	}
	if s, err := q.Status(); err != nil {
		WrapRollbackError(tx.Rollback())
		return nil, WrapGrpcErrorf(codes.Internal, "Failed to get status")
	} else {
		res.State = GetAgentStateFromQemuState(s)
	}

	return res, nil
}

func (a VirtualMachineAgentAPI) DeleteVirtualMachineAgent(ctx context.Context, req *DeleteVirtualMachineAgentRequest) (*empty.Empty, error) {
	id, err := uuid.FromString(req.Uuid)
	if err != nil {
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid uuid: value='%s', err='%s'", req.Uuid, err.Error())
	}

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

			// gateway settings
			if nd.Ipv4Gateway != "" {
				_, ipn, err := net.ParseCIDR(nd.Ipv4AddressCidr)
				if err != nil {
					return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid ipv4_address_cidr: value='%s', err='%s'", nd.Ipv4AddressCidr, err.Error())
				}
				if err := iptables.DeleteMasqueradeRule(b.Name(), ipn.String()); err != nil {
					return nil, WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to delete masquerade rule").Error())
				}
			}
		}
	}

	return &empty.Empty{}, nil
}

func (a VirtualMachineAgentAPI) BootVirtualMachineAgent(ctx context.Context, req *BootVirtualMachineAgentRequest) (*BootVirtualMachineAgentResponse, error) {
	id, err := uuid.FromString(req.Uuid)
	if err != nil {
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid uuid: value='%s', err='%s'", req.Uuid, err.Error())
	}

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
	id, err := uuid.FromString(req.Uuid)
	if err != nil {
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid uuid: value='%s', err='%s'", req.Uuid, err.Error())
	}

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
	id, err := uuid.FromString(req.Uuid)
	if err != nil {
		return nil, WrapGrpcErrorf(codes.InvalidArgument, "Set valid uuid: value='%s', err='%s'", req.Uuid, err.Error())
	}

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
