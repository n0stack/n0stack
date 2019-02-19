package virtualmachine

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"

	"github.com/n0stack/n0stack/n0core/pkg/driver/cloudinit/configdrive"
	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu"
	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu_img"
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
	p := filepath.Join(a.baseDirectory, name)

	if _, err := os.Stat(p); err != nil {
		if err := os.RemoveAll(p); err != nil { // TODO: check permission
			return errors.Wrapf(err, "Failed to rm '%s'", p)
		}
	}

	return nil
}

func SetPrefix(name string) string {
	return fmt.Sprintf("n0stack/%s", name)
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
	defer tx.RollbackWithLog()

	wd, err := a.GetWorkDirectory(name)
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get working directory '%s'", wd)
	}

	q, err := qemu.OpenQemu(SetPrefix(name))
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if q.IsRunning() {
		return nil, grpcutil.WrapGrpcErrorf(codes.AlreadyExists, "Qemu process is already running")
	}

	if err := q.Start(id, filepath.Join(wd, QmpMonitorSocketFile), websocket, vcpus, mem); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to start qemu process: err=%s", err.Error())
	}
	defer q.Close()
	tx.PushRollback("delete Qemu", func() error {
		return q.Delete()
	})

	eth := make([]*configdrive.CloudConfigEthernet, len(req.Netdevs))
	{
		for i, nd := range req.Netdevs {
			b, err := iproute2.NewBridge(netutil.StructLinuxNetdevName(nd.NetworkName))
			if err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to create bridge '%s': err='%s'", nd.NetworkName, err.Error())
			}
			tx.PushRollback("delete created bridge", func() error {
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

			t, err := iproute2.NewTap(netutil.StructLinuxNetdevName(nd.Name))
			if err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to create tap '%s': err='%s'", nd.Name, err.Error())
			}
			tx.PushRollback("delete created tap", func() error {
				if err := t.Delete(); err != nil {
					return fmt.Errorf("Failed to delete tap '%s': err='%s'", nd.Name, err.Error())
				}

				return nil
			})
			if err := t.SetMaster(b); err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to set master of tap '%s' as '%s': err='%s'", t.Name(), b.Name(), err.Error())
			}

			hw, err := net.ParseMAC(nd.HardwareAddress)
			if err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Hardware address '%s' is invalid on netdev '%s'", nd.HardwareAddress, nd.Name)
			}
			if err := q.AttachTap(nd.Name, t.Name(), hw); err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to attach tap: err='%s'", err.Error())
			}

			// Cloudinit settings
			eth[i] = &configdrive.CloudConfigEthernet{
				MacAddress: hw,
			}

			if nd.Ipv4AddressCidr != "" {
				ip := netutil.ParseCIDR(nd.Ipv4AddressCidr)
				if ip == nil {
					return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set valid ipv4_address_cidr: value='%s'", nd.Ipv4AddressCidr)
				}
				nameservers := make([]net.IP, len(nd.Nameservers))
				for i, n := range nd.Nameservers {
					nameservers[i] = net.ParseIP(n)
				}

				eth[i].Address4 = ip
				eth[i].Gateway4 = net.ParseIP(nd.Ipv4Gateway)
				eth[i].NameServers = nameservers

				// Gateway settings
				if nd.Ipv4Gateway != "" {
					mask := ip.SubnetMaskBits()
					gatewayIP := fmt.Sprintf("%s/%d", nd.Ipv4Gateway, mask)
					if err := b.SetAddress(gatewayIP); err != nil {
						return nil, grpcutil.WrapGrpcErrorf(codes.Internal, errors.Wrapf(err, "Failed to set gateway IP to bridge: value=%s", gatewayIP).Error())
					}
				}
			}
		}
	}

	{
		parsedKeys := make([]ssh.PublicKey, len(req.SshAuthorizedKeys))
		for i, k := range req.SshAuthorizedKeys {
			parsedKeys[i], _, _, _, err = ssh.ParseAuthorizedKey([]byte(k))
			if err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "ssh_authorized_keys is invalid: value='%s', err='%s'", k, err.Error())
			}
		}

		c := configdrive.StructConfig(req.LoginUsername, req.Name, parsedKeys, eth)
		p, err := c.Generate(wd)
		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to generate cloudinit configdrive:  err='%s'", err.Error())
		}
		req.Blockdevs = append(req.Blockdevs, &BlockDev{
			Name: "configdrive",
			Url: (&url.URL{
				Scheme: "file",
				Path:   p,
			}).String(),
			BootIndex: 50, // MEMO: 適当
		})
	}

	{
		for _, bd := range req.Blockdevs {
			u, err := url.Parse(bd.Url)
			if err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "url '%s' is invalid url: '%s'", bd.Url, err.Error())
			}

			i, err := img.OpenQemuImg(u.Path)
			if err != nil {
				return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to open qemu image: err='%s'", err.Error())
			}

			// この条件は雑
			if i.Info.Format == "raw" {
				if bd.BootIndex < 3 {
					if err := q.AttachISO(bd.Name, u, uint(bd.BootIndex)); err != nil {
						return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to attach iso '%s': err='%s'", u.Path, err.Error())
					}
				} else {
					if err := q.AttachRaw(bd.Name, u, uint(bd.BootIndex)); err != nil {
						return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to attach raw '%s': err='%s'", u.Path, err.Error())
					}
				}
			} else {
				if err := q.AttachQcow2(bd.Name, u, uint(bd.BootIndex)); err != nil {
					return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to attach image '%s': err='%s'", u.String(), err.Error())
				}
			}
		}
	}

	if err := q.Boot(); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to boot qemu: err=%s", err.Error())
	}

	s, err := q.Status()
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get status: err=%s", err.Error())
	}

	tx.Commit()
	return &BootVirtualMachineResponse{
		State:         GetAgentStateFromQemuState(s),
		WebsocketPort: uint32(websocket),
	}, nil
}

func (a VirtualMachineAgent) RebootVirtualMachine(ctx context.Context, req *RebootVirtualMachineRequest) (*RebootVirtualMachineResponse, error) {
	return nil, grpcutil.WrapGrpcErrorf(codes.Unimplemented, "")
}

func (a VirtualMachineAgent) ShutdownVirtualMachine(ctx context.Context, req *ShutdownVirtualMachineRequest) (*ShutdownVirtualMachineResponse, error) {
	q, err := qemu.OpenQemu(SetPrefix(req.Name))
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	defer q.Close()

	if !q.IsRunning() {
		return &ShutdownVirtualMachineResponse{
			State: VirtualMachineState_SHUTDOWN,
		}, nil
	}

	if req.Hard {
		if err := q.HardShutdown(); err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to hard shutdown qemu: %s", err.Error())
		}

		// after hard shutdown, qemu process exists immediately, so agent can't get status
		return &ShutdownVirtualMachineResponse{
			State: VirtualMachineState_SHUTDOWN,
		}, nil
	}

	if err := q.Shutdown(); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to shutdown qemu: %s", err.Error())
	}

	s, err := q.Status()

	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get status: err=%s", err.Error())
	}

	return &ShutdownVirtualMachineResponse{
		State: GetAgentStateFromQemuState(s),
	}, nil
}

func (a VirtualMachineAgent) DeleteVirtualMachine(ctx context.Context, req *DeleteVirtualMachineRequest) (*empty.Empty, error) {
	q, err := qemu.OpenQemu(SetPrefix(req.Name))
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	if q.IsRunning() {
		if err := q.Delete(); err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to delete qemu: %s", err.Error())
		}
	}
	if err := a.DeleteWorkDirectory(req.Name); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to delete work directory: %s", err.Error())
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
			}
		}
	}

	return &empty.Empty{}, nil
}

func GetAgentStateFromQemuState(s qemu.Status) VirtualMachineState {
	switch s {
	case qemu.StatusRunning:
		return VirtualMachineState_RUNNING

	case qemu.StatusShutdown, qemu.StatusGuestPanicked, qemu.StatusPreLaunch:
		return VirtualMachineState_SHUTDOWN

	case qemu.StatusPaused, qemu.StatusSuspended:
		return VirtualMachineState_PAUSED

	case qemu.StatusInternalError, qemu.StatusIOError:
		return VirtualMachineState_FAILED

	case qemu.StatusInMigrate:
	case qemu.StatusFinishMigrate:
	case qemu.StatusPostMigrate:
	case qemu.StatusRestoreVM:
	case qemu.StatusSaveVM: // TODO: 多分PAUSED
	case qemu.StatusWatchdog:
	case qemu.StatusDebug:
	}

	return VirtualMachineState_UNKNOWN
}
