package virtualmachine

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"

	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"
	"github.com/n0stack/n0stack/n0core/pkg/driver/cloudinit/configdrive"
	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu"
	img "github.com/n0stack/n0stack/n0core/pkg/driver/qemu_img"
	grpcutil "github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	netutil "github.com/n0stack/n0stack/n0core/pkg/util/net"
	"github.com/n0stack/n0stack/n0proto.go/pkg/transaction"
)

const (
	QmpMonitorSocketFile   = "monitor.sock"
	VNCWebSocketPortOffset = 6900
)

type VirtualMachineAgent struct {
	baseDirectory string
	bridgeMutex   lock.MutexTable
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
		bridgeMutex:   lock.NewMemoryMutexTable(100),
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
	defer q.Close()

	if !q.IsRunning() {
		if err := q.Start(id, filepath.Join(wd, QmpMonitorSocketFile), vcpus, mem); err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to start qemu process: err=%s", err.Error())
		}
		tx.PushRollback("delete Qemu", func() error {
			return q.Delete()
		})

		eth := make([]*configdrive.CloudConfigEthernet, len(req.Netdevs))
		{
			for i, nd := range req.Netdevs {
				err := func() error {
					bn := netutil.StructLinuxNetdevName(nd.NetworkName)
					if !lock.WaitUntilLock(a.bridgeMutex, bn, 5*time.Second, 10*time.Millisecond) {
						return errors.Wrapf(datastore.LockError(), "Failed to lock bridge '%s'", bn)
					}
					defer a.bridgeMutex.Unlock(bn)

					b, err := iproute2.NewBridge(bn)
					if err != nil {
						return errors.Wrapf(err, "Failed to create bridge '%s'", nd.NetworkName)
					}
					tx.PushRollback("delete created bridge", func() error {
						if !lock.WaitUntilLock(a.bridgeMutex, b.Name(), 5*time.Second, 10*time.Millisecond) {
							return fmt.Errorf("Failed to lock bridge '%s': err='%s'", b.Name(), datastore.LockError().Error())
						}
						defer a.bridgeMutex.Unlock(b.Name())

						if _, err = b.DeleteIfNoSlave(); err != nil {
							return fmt.Errorf("Failed to delete bridge '%s': err='%s'", b.Name(), err.Error())
						}

						return nil
					})

					t, err := iproute2.NewTap(netutil.StructLinuxNetdevName(nd.Name))
					if err != nil {
						return grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to create tap '%s': err='%s'", nd.Name, err.Error())
					}
					tx.PushRollback("delete created tap", func() error {
						if err := t.Delete(); err != nil {
							return fmt.Errorf("Failed to delete tap '%s': err='%s'", nd.Name, err.Error())
						}

						return nil
					})
					if err := t.SetMaster(b); err != nil {
						return errors.Wrapf(err, "Failed to set master of tap '%s' as '%s'", t.Name(), b.Name())
					}

					hw, err := net.ParseMAC(nd.HardwareAddress)
					if err != nil {
						return errors.Wrapf(err, "Hardware address '%s' is invalid on netdev '%s'", nd.HardwareAddress, nd.Name)
					}
					if err := q.AttachTap(nd.Name, t.Name(), hw); err != nil {
						return errors.Wrapf(err, "Failed to attach tap")
					}

					// Cloudinit settings
					eth[i] = &configdrive.CloudConfigEthernet{
						MacAddress: hw,
					}

					if nd.Ipv4AddressCidr != "" {
						ip := netutil.ParseCIDR(nd.Ipv4AddressCidr)
						if ip == nil {
							return fmt.Errorf("Set valid ipv4_address_cidr: value='%s'", nd.Ipv4AddressCidr)
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
								return errors.Wrapf(err, "Failed to set gateway IP to bridge: value=%s", gatewayIP)
							}
						}
					}

					return nil
				}()

				if err != nil {
					return nil, grpcutil.WrapGrpcErrorf(codes.Internal, err.Error())
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

				if !i.IsExists() {
					return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "blockdev is not exists: blockdev=%s", bd.Name)
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
		WebsocketPort: uint32(q.GetVNCWebsocketPort()),
	}, nil
}

func (a VirtualMachineAgent) RebootVirtualMachine(ctx context.Context, req *RebootVirtualMachineRequest) (*RebootVirtualMachineResponse, error) {
	return nil, grpcutil.WrapGrpcErrorf(codes.Unimplemented, "")
}

func (a VirtualMachineAgent) ShutdownVirtualMachine(ctx context.Context, req *ShutdownVirtualMachineRequest) (*ShutdownVirtualMachineResponse, error) {
	return nil, grpcutil.WrapGrpcErrorf(codes.Unimplemented, "")
}

func (a VirtualMachineAgent) DeleteVirtualMachine(ctx context.Context, req *DeleteVirtualMachineRequest) (*empty.Empty, error) {
	q, err := qemu.OpenQemu(SetPrefix(req.Name))
	if err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to open qemu process: %s", err.Error())
	}
	defer q.Close()

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

		err = func() error {
			bn := netutil.StructLinuxNetdevName(nd.NetworkName)
			if !lock.WaitUntilLock(a.bridgeMutex, bn, 5*time.Second, 10*time.Millisecond) {
				return errors.Wrapf(datastore.LockError(), "Failed to lock bridge '%s'", bn)
			}
			defer a.bridgeMutex.Unlock(bn)

			b, err := iproute2.NewBridge(bn)
			if err != nil {
				return errors.Wrapf(err, "Failed to create bridge '%s'", nd.NetworkName)
			}

			if _, err := b.DeleteIfNoSlave(); err != nil {
				return errors.Wrapf(err, "Failed to delete bridge '%s'", b.Name())
			}

			return nil
		}()

		if err != nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.Internal, err.Error())
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
