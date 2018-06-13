package iproute2

import (
	fmt "fmt"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/vishvananda/netlink"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Iproute2Agent struct {
	uplink netlink.Link
}

func NewIproute2Agent(uplink string) (*Iproute2Agent, error) {
	l, err := netlink.LinkByName(uplink)
	if err != nil {
		// return nil, err // エラーの場合はLinkがない場合のみだと考えておく
		return nil, err
	}

	i := &Iproute2Agent{
		uplink: l,
	}

	return i, nil
}

// ip link list $bridge_name
func (a Iproute2Agent) getBridge(name string) (*netlink.Bridge, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		// return nil, err // エラーの場合はLinkがない場合のみだと考えておく
		return nil, nil
	}

	b, ok := l.(*netlink.Bridge)
	if !ok {
		return nil, fmt.Errorf("Failed to get bridge because link is not bridge, LinkName:%s", l.Attrs().Name)
	}

	return b, nil
}

// ip link add type bridge name $name
// ip link set dev $bridge_name master $uplink promisc on up
func (a Iproute2Agent) createBridge(name string, uplink netlink.Link) error {
	l := netlink.NewLinkAttrs()
	l.Name = name
	b := &netlink.Bridge{LinkAttrs: l}

	if err := netlink.LinkAdd(b); err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(uplink, b); err != nil {
		return err
	}

	if err := netlink.SetPromiscOn(uplink); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(b); err != nil {
		return err
	}

	return nil
}

// ip link list $tap_name
func (a Iproute2Agent) getTap(name string) (netlink.Link, error) {
	l, err := netlink.LinkByName(name)
	if err != nil {
		// return nil, err // エラーの場合はLinkがない場合のみだと考えておく
		return nil, nil
	}

	// 何故かできない
	// t, ok := l.(*netlink.Tuntap)
	// if !ok {
	// 	return nil, fmt.Errorf("Failed to get tuntap because link is not tuntap, LinkName:%s", l.Attrs().Name)
	// }

	// if t.Mode != netlink.TUNTAP_MODE_TAP {
	// 	return nil, fmt.Errorf("Failed to get tuntap because link is not tap, LinkName:%s, Mode:%v", l.Attrs().Name, t.Mode)
	// }

	if l.Type() != "tun" { // なんでtapではなくtunになっているは謎
		return nil, fmt.Errorf("Failed to get tuntap because link is not tuntap, LinkName:%s", l.Type())
	}

	return l, nil
}

// ip tuntap add name $tap_name mode tap
// ip link set dev $tap_name master $bridge_name
func (a Iproute2Agent) createTap(name string, master *netlink.Bridge) error {
	l := netlink.NewLinkAttrs()
	l.Name = name
	t := &netlink.Tuntap{
		LinkAttrs: l,
		Mode:      netlink.TUNTAP_MODE_TAP,
	}

	if err := netlink.LinkAdd(t); err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(t, master); err != nil {
		return err
	}

	return nil
}

func (a Iproute2Agent) ApplyTap(ctx context.Context, req *ApplyTapRequest) (*Tap, error) {
	var b *netlink.Bridge

	switch req.Tap.Type {
	case Tap_FLAT:
		var err error
		b, err = a.getBridge(req.Tap.BridgeName)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to get bridge, bridge_name:%s err:%s", req.Tap.Name, err.Error())
		}
		if b == nil {
			if err := a.createBridge(req.Tap.BridgeName, a.uplink); err != nil {
				return nil, grpc.Errorf(codes.Internal, "Failed to create bridge, bridge_name:%s err:%s", req.Tap.Name, err.Error())
			}

			b, _ = a.getBridge(req.Tap.BridgeName)
		}

	case Tap_VLAN, Tap_VXLAN:
		return nil, grpc.Errorf(codes.Unimplemented, "")
	}

	t, err := a.getTap(req.Tap.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get tap, tap_name:%s, err:%s", req.Tap.Name, err.Error())
	}
	if t == nil {
		if err := a.createTap(req.Tap.Name, b); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to create tap, tap_name:%s err:%s", req.Tap.Name, err.Error())
		}
	}

	return req.Tap, nil
}

func (a Iproute2Agent) DeleteTap(ctx context.Context, req *DeleteTapRequest) (*google_protobuf.Empty, error) {
	t, err := a.getTap(req.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Failed to get tap, tap_name:%s, err:%s", req.Name, err.Error())
	}
	if t != nil {
		if err := netlink.LinkDel(t); err != nil {
			return nil, grpc.Errorf(codes.Internal, "Failed to create tap, tap_name:%s err:%s", req.Name, err.Error())
		}
	}

	// TODO: bridgeなどの後始末もする必要がある

	return &google_protobuf.Empty{}, nil
}
