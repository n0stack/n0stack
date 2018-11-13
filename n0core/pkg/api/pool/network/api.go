package network

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/n0stack/n0proto.go/budget/v0"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

type NetworkAPI struct {
	dataStore datastore.Datastore
}

func CreateNetworkAPI(ds datastore.Datastore) (*NetworkAPI, error) {
	a := &NetworkAPI{
		dataStore: ds,
	}
	a.dataStore.AddPrefix("network")

	return a, nil
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
	if _, _, err := net.ParseCIDR(req.Ipv4Cidr); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Field 'ipv4_cidr' is invalid : %s", err.Error())
	}

	res := &ppool.Network{}
	if err := a.dataStore.Get(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	// var err error
	res.Version, _ = datastore.CheckVersion(res.Version, req.Version)
	// if err != nil {
	// 	return nil, grpc.Errorf(codes.InvalidArgument, "Failed to check version: %s", err.Error())
	// }

	res.Name = req.Name
	res.Annotations = req.Annotations
	res.Ipv4Cidr = req.Ipv4Cidr
	res.Ipv6Cidr = req.Ipv6Cidr
	res.Domain = req.Domain

	res.State = ppool.Network_AVAILABLE
	if err := a.dataStore.Apply(req.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a NetworkAPI) DeleteNetwork(ctx context.Context, req *ppool.DeleteNetworkRequest) (*empty.Empty, error) {
	if err := a.dataStore.Delete(req.Name); err != nil {
		log.Printf("[WARNING] Failed to delete data from db: err='%s'", err.Error())
		return &empty.Empty{}, grpc.Errorf(codes.Internal, "Failed to delete '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return &empty.Empty{}, nil
}

// とりあえず IPv4 のスケジューリングのみに対応
func (a NetworkAPI) ReserveNetworkInterface(ctx context.Context, req *ppool.ReserveNetworkInterfaceRequest) (*ppool.Network, error) {
	if req.NetworkInterfaceName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Do not set field 'network_interface_name' as blank")
	}

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

	// 保存する際にパースするのでエラーは発生しない
	_, cidr, _ := net.ParseCIDR(res.Ipv4Cidr)
	var reqIPv4 net.IP
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

	var reqHW net.HardwareAddr
	if req.HardwareAddress == "" {
		reqHW = GenerateHardwareAddress(fmt.Sprintf("%s/%s", req.NetworkName, req.NetworkInterfaceName))
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
		Ipv4Address:     reqIPv4.String(),
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
