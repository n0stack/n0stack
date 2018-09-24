package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"reflect"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/pkg/datastore"
	"github.com/n0stack/proto.go/budget/v0"
	"github.com/n0stack/proto.go/pool/v0"
)

type NetworkAPI struct {
	dataStore datastore.Datastore
}

func CreateNetworkAPI(ds datastore.Datastore) (*NetworkAPI, error) {
	a := &NetworkAPI{
		dataStore: ds,
	}

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
	if reflect.ValueOf(res.Metadata).IsNil() {
		return nil, grpc.Errorf(codes.NotFound, "")
	}

	return res, nil
}

func (a NetworkAPI) ApplyNetwork(ctx context.Context, req *ppool.ApplyNetworkRequest) (*ppool.Network, error) {
	res := &ppool.Network{
		Metadata: req.Metadata,
		Spec:     req.Spec,
		Status:   &ppool.NetworkStatus{},
	}

	if _, _, err := net.ParseCIDR(req.Spec.Ipv4Cidr); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Field 'ipv4_cidr' is invalid : %s", err.Error())
	}

	prev := &ppool.Network{}
	if err := a.dataStore.Get(req.Metadata.Name, prev); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
	}
	var err error
	res.Metadata.Version, err = datastore.CheckVersion(prev, req)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "Failed to check version: %s", err.Error())
	}

	res.Status.State = ppool.NetworkStatus_AVAILABLE
	if err := a.dataStore.Apply(req.Metadata.Name, res); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' for db, please retry or contact for the administrator of this cluster", req.Metadata.Name)
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
func (a NetworkAPI) ReserveNetworkInterface(ctx context.Context, req *ppool.ReserveNetworkInterfaceRequest) (*ppool.ReserveNetworkInterfaceResponse, error) {
	if req.NetworkInterfaceName == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Do not set field 'network_interface_name' as blank")
	}

	n := &ppool.Network{}
	if err := a.dataStore.Get(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if n == nil {
		return nil, grpc.Errorf(codes.NotFound, "Node '%s' is not found", req.Name)
	}
	if n.Status.ReservedNetworkInterfaces == nil {
		n.Status.ReservedNetworkInterfaces = make(map[string]*pbudget.NetworkInterface)
	}
	if _, ok := n.Status.ReservedNetworkInterfaces[req.NetworkInterfaceName]; ok {
		return nil, grpc.Errorf(codes.AlreadyExists, "Network interface '%s' is already exists on Node '%s'", req.NetworkInterfaceName, req.Name)
	}

	// 保存する際にパースするのでエラーは発生しない
	_, cidr, _ := net.ParseCIDR(n.Spec.Ipv4Cidr)

	var reqIPv4 net.IP
	if req.NetworkInterface == nil || req.NetworkInterface.Ipv4Address == "" {
		if reqIPv4 = ScheduleNewIPv4(cidr, n.Status.ReservedNetworkInterfaces); reqIPv4 == nil {
			return nil, grpc.Errorf(codes.ResourceExhausted, "ipv4_address is full on Network '%s'", req.Name)
		}
	} else {
		reqIPv4 = net.ParseIP(req.NetworkInterface.Ipv4Address)
		if reqIPv4 == nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "ipv4_address field is invalid")
		}

		if err := CheckIPv4OnCIDR(reqIPv4, cidr); err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "ipv4_address field is invalid: %s", err.Error())
		}
		if err := CheckConflictIPv4(reqIPv4, n.Status.ReservedNetworkInterfaces); err != nil {
			return nil, grpc.Errorf(codes.AlreadyExists, "ipv4_address field is invalid: %s", err.Error())
		}
	}

	var reqHW net.HardwareAddr
	if req.NetworkInterface == nil || req.NetworkInterface.HardwareAddress == "" {
		reqHW = GenerateHardwareAddress(fmt.Sprintf("%s/%s", req.Name, req.NetworkInterfaceName))
	} else {
		var err error
		reqHW, err = net.ParseMAC(req.NetworkInterface.HardwareAddress)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "hardware_address field is invalid")
		}
	}

	res := &ppool.ReserveNetworkInterfaceResponse{
		Name:                 req.Name,
		NetworkInterfaceName: req.NetworkInterfaceName,
		NetworkInterface: &pbudget.NetworkInterface{
			Annotations:     req.NetworkInterface.Annotations,
			HardwareAddress: reqHW.String(),
			Ipv4Address:     reqIPv4.String(),
		},
	}
	n.Status.ReservedNetworkInterfaces[req.NetworkInterfaceName] = res.NetworkInterface
	if err := a.dataStore.Apply(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to store data on db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return res, nil
}

func (a NetworkAPI) ReleaseNetworkInterface(ctx context.Context, req *ppool.ReleaseNetworkInterfaceRequest) (*empty.Empty, error) {
	n := &ppool.Network{}
	if err := a.dataStore.Get(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to get data from db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to get '%s' from db, please retry or contact for the administrator of this cluster", req.Name)
	}
	if n == nil {
		return nil, grpc.Errorf(codes.NotFound, "Do not exists network '%s'", req.Name)
	}

	if _, ok := n.Status.ReservedNetworkInterfaces[req.NetworkInterfaceName]; !ok {
		return nil, grpc.Errorf(codes.NotFound, "Do not exists network interface '%s' on network '%s'", req.NetworkInterfaceName, req.Name)
	}
	delete(n.Status.ReservedNetworkInterfaces, req.NetworkInterfaceName)

	if err := a.dataStore.Apply(req.Name, n); err != nil {
		log.Printf("[WARNING] Failed to apply data for db: err='%s'", err.Error())
		return nil, grpc.Errorf(codes.Internal, "Failed to store '%s' on db, please retry or contact for the administrator of this cluster", req.Name)
	}

	return &empty.Empty{}, nil
}
