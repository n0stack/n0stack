package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/lock"
	"github.com/n0stack/n0stack/n0core/pkg/util/grpc"
	"github.com/n0stack/n0stack/n0core/pkg/util/net"
	"github.com/n0stack/n0stack/n0proto.go/budget/v0"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

type NetworkAPI struct {
	dataStore datastore.Datastore
}

func CreateNetworkAPI(ds datastore.Datastore) *NetworkAPI {
	a := &NetworkAPI{
		dataStore: ds.AddPrefix("network"),
	}

	return a
}

func (a NetworkAPI) ListNetworks(ctx context.Context, req *ppool.ListNetworksRequest) (*ppool.ListNetworksResponse, error) {
	res := &ppool.ListNetworksResponse{}
	f := func(s int) []proto.Message {
		res.Networks = make([]*ppool.Network, s)
		for i := range res.Networks {
			res.Networks[i] = &ppool.Network{}
		}

		m := make([]proto.Message, s)
		for i, v := range res.Networks {
			m[i] = v
		}

		return m
	}

	if err := a.dataStore.List(f); err != nil {
		log.Printf("[WARNING] Failed to list data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to list from db, please retry or contact for the administrator of this cluster")
	}
	if len(res.Networks) == 0 {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NetworkAPI) GetNetwork(ctx context.Context, req *ppool.GetNetworkRequest) (*ppool.Network, error) {
	if req.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set name")
	}

	res := &ppool.Network{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NetworkAPI) ApplyNetwork(ctx context.Context, req *ppool.ApplyNetworkRequest) (*ppool.Network, error) {
	ipv4 := netutil.ParseCIDR(req.Ipv4Cidr)
	ipv6 := netutil.ParseCIDR(req.Ipv6Cidr)
	{
		if req.Name == "" {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Set any 'name'")
		}

		if ipv4 == nil && ipv6 == nil {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Field 'ipv4_cidr' and 'ipv6_cidr' are invalid")
		}
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	network := &ppool.Network{}
	if err := a.dataStore.Get(req.Name, network); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get data from db: err='%s'", err.Error())
	}

	{
		if network.Name != "" && ipv4.String() != network.Ipv4Cidr {
			return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "ipv4 cidr is different from previous ipv4 cidr")
		}

		res, err := a.ListNetworks(ctx, &ppool.ListNetworksRequest{})
		if err != nil {
			if grpc.Code(err) != codes.NotFound {
				return nil, grpcutil.WrapGrpcErrorf(grpc.Code(err), errors.Wrapf(err, "Failed to list networks").Error())
			}
		} else {
			if ipv4 != nil {
				for _, v := range res.Networks {
					if v.Name == req.Name {
						continue
					}

					existing := netutil.ParseCIDR(v.Ipv4Cidr)
					if existing != nil && netutil.IsConflicting(ipv4, existing) {
						return nil, grpcutil.WrapGrpcErrorf(codes.InvalidArgument, "Field 'ipv4_cidr' is conflicting with network=%s", v.Name)
					}
				}
			}

			// TODO: check IPv6 conflicting
		}
	}

	network.Name = req.Name
	network.Annotations = req.Annotations
	network.Ipv4Cidr = req.Ipv4Cidr
	network.Ipv6Cidr = req.Ipv6Cidr
	network.Domain = req.Domain

	network.State = ppool.Network_AVAILABLE
	if err := a.dataStore.Apply(req.Name, network); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to apply data for db: err='%s'", err.Error())
	}

	return network, nil
}

func (a NetworkAPI) DeleteNetwork(ctx context.Context, req *ppool.DeleteNetworkRequest) (*empty.Empty, error) {
	if req.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Set name")
	}

	if !a.dataStore.Lock(req.Name) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.Name)

	network := &ppool.Network{}
	if err := a.dataStore.Get(req.Name, network); err != nil {
		return nil, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to get data from db: err='%s'", err.Error())
	}
	if network.Name == "" {
		return nil, grpcutil.WrapGrpcErrorf(codes.NotFound, "")
	}

	if IsLockedForDeletion(network) {
		return nil, grpcutil.WrapGrpcErrorf(codes.FailedPrecondition, "Network has some network interfaces, so is locked for deletion")
	}

	if err := a.dataStore.Delete(req.Name); err != nil {
		return &empty.Empty{}, grpcutil.WrapGrpcErrorf(codes.Internal, "Failed to delete data from db: err='%s', name='%s'", err.Error(), req.Name)
	}

	return &empty.Empty{}, nil
}

// とりあえず IPv4 のスケジューリングのみに対応
func (a NetworkAPI) ReserveNetworkInterface(ctx context.Context, req *ppool.ReserveNetworkInterfaceRequest) (*ppool.Network, error) {
	if req.NetworkInterfaceName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Do not set field 'network_interface_name' as blank")
	}

	if !lock.WaitUntilLock(a.dataStore, req.NetworkName, 1*time.Second, 50*time.Millisecond) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.NetworkName)

	res := &ppool.Network{}
	if err := a.dataStore.Get(req.NetworkName, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.NetworkName)
	}
	if res.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Network '%s' is not found", req.NetworkName)
	}
	if res.ReservedNetworkInterfaces == nil {
		res.ReservedNetworkInterfaces = make(map[string]*pbudget.NetworkInterface)
	}
	if _, ok := res.ReservedNetworkInterfaces[req.NetworkInterfaceName]; ok {
		return nil, grpc.Errorf(codes.AlreadyExists, "Network interface '%s' is already exists on Network '%s'", req.NetworkInterfaceName, req.NetworkName)
	}

	// err != nil の場合は Ipv4Cidr がからのとき
	_, cidr, err := net.ParseCIDR(res.Ipv4Cidr)
	var reqIPv4 net.IP
	if err == nil {
		if req.Ipv4Address == "" {
			if reqIPv4 = ScheduleNewIPv4(cidr, res.ReservedNetworkInterfaces); reqIPv4 == nil {
				return nil, grpc.Errorf(codes.ResourceExhausted, "ipv4_address is full on Network '%s'", req.NetworkName)
			}
		} else {
			reqIPv4 = net.ParseIP(req.Ipv4Address)
			if reqIPv4 == nil {
				return nil, grpc.Errorf(codes.InvalidArgument, "ipv4_address field is invalid")
			}

			if err := CheckIPv4OnCIDR(reqIPv4, cidr); err != nil {
				return nil, grpc.Errorf(codes.InvalidArgument, "ipv4_address field is invalid: %s", err.Error())
			}
			if err := CheckConflictIPv4(reqIPv4, res.ReservedNetworkInterfaces); err != nil {
				return nil, grpc.Errorf(codes.ResourceExhausted, "ipv4_address field is invalid: %s", err.Error())
			}
		}
	}

	var reqHW net.HardwareAddr
	if req.HardwareAddress == "" {
		reqHW = netutil.GenerateHardwareAddress(fmt.Sprintf("%s/%s", req.NetworkName, req.NetworkInterfaceName))
	} else {
		var err error
		reqHW, err = net.ParseMAC(req.HardwareAddress)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "hardware_address field is invalid")
		}
	}

	res.ReservedNetworkInterfaces[req.NetworkInterfaceName] = &pbudget.NetworkInterface{
		Annotations:     req.Annotations,
		HardwareAddress: reqHW.String(),
	}
	if reqIPv4 != nil {
		res.ReservedNetworkInterfaces[req.NetworkInterfaceName].Ipv4Address = reqIPv4.String()
	}

	if err := a.dataStore.Apply(req.NetworkName, res); err != nil {
		log.Printf("[WARNING] Failed to store data on db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.NetworkName)
	}

	return res, nil
}

func (a NetworkAPI) ReleaseNetworkInterface(ctx context.Context, req *ppool.ReleaseNetworkInterfaceRequest) (*empty.Empty, error) {
	n := &ppool.Network{}
	if err := a.dataStore.Get(req.NetworkName, n); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.NetworkName)
	}
	if n.Name == "" {
		return nil, grpc.Errorf(codes.NotFound, "Do not exists network '%s'", req.NetworkName)
	}

	if !lock.WaitUntilLock(a.dataStore, req.NetworkName, 1*time.Second, 50*time.Millisecond) {
		return nil, datastore.LockError()
	}
	defer a.dataStore.Unlock(req.NetworkName)

	if n.ReservedNetworkInterfaces == nil {
		return nil, grpc.Errorf(codes.NotFound, "Do not exists network interface '%s' on network '%s'", req.NetworkInterfaceName, req.NetworkName)
	}
	if _, ok := n.ReservedNetworkInterfaces[req.NetworkInterfaceName]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Do not exists network interface '%s' on network '%s'", req.NetworkInterfaceName, req.NetworkName)
	}
	delete(n.ReservedNetworkInterfaces, req.NetworkInterfaceName)

	if err := a.dataStore.Apply(req.NetworkName, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.NetworkName)
	}

	return &empty.Empty{}, nil
}
