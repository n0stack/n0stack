package iproute2

import (
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
