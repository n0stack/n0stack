// +build medium
// +build !without_external

package network

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/proto.go/budget/v0"
	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/v0"
)

func TestEmptyNetwork(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, err := CreateNetworkAPI(m)
	if err != nil {
		t.Fatalf("Failed to create Network API: err='%s'", err.Error())
	}

	listRes, err := na.ListNetworks(context.Background(), &ppool.ListNetworksRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("ListNetworks got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListNetworks do not return nil: res='%s'", listRes)
	}

	getRes, err := na.GetNetwork(context.Background(), &ppool.GetNetworkRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("GetNetwork got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetNetwork do not return nil: res='%s'", listRes)
	}
}

func TestApplyNetwork(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, err := CreateNetworkAPI(m)
	if err != nil {
		t.Fatalf("Failed to create Network API: err='%s'", err.Error())
	}

	n := &ppool.Network{
		Metadata: &pn0stack.Metadata{
			Name:    "test-network",
			Version: 1,
		},
		Spec: &ppool.NetworkSpec{
			Ipv4Cidr: "192.168.0.0/30",
			Domain:   "test.local",
		},
		Status: &ppool.NetworkStatus{
			State: ppool.NetworkStatus_AVAILABLE,
		},
	}

	applyRes, err := na.ApplyNetwork(context.Background(), &ppool.ApplyNetworkRequest{
		Metadata: &pn0stack.Metadata{
			Name: n.Metadata.Name,
		},
		Spec: n.Spec,
	})
	if err != nil {
		t.Fatalf("Failed to apply network: err='%s'", err.Error())
	}

	// diffが取れないので
	applyRes.XXX_sizecache = 0
	applyRes.Metadata.XXX_sizecache = 0
	applyRes.Spec.XXX_sizecache = 0
	applyRes.Status.XXX_sizecache = 0
	if diff := cmp.Diff(n, applyRes); diff != "" {
		t.Fatalf("ApplyNetwork response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := na.ListNetworks(context.Background(), &ppool.ListNetworksRequest{})
	if err != nil {
		t.Errorf("ListNetworks got error: err='%s'", err.Error())
	}
	if len(listRes.Networks) != 1 {
		t.Errorf("ListNetworks response is wrong: have='%d', want='%d'", len(listRes.Networks), 1)
	}

	getRes, err := na.GetNetwork(context.Background(), &ppool.GetNetworkRequest{Name: n.Metadata.Name})
	if err != nil {
		t.Errorf("GetNetwork got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(n, getRes); diff != "" {
		t.Errorf("GetNetwork response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := na.DeleteNetwork(context.Background(), &ppool.DeleteNetworkRequest{Name: n.Metadata.Name}); err != nil {
		t.Errorf("DeleteNetwork got error: err='%s'", err.Error())
	}
}

func TestNetworkAboutNetworkInterface(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, err := CreateNetworkAPI(m)
	if err != nil {
		t.Fatalf("Failed to create Network API: err='%s'", err.Error())
	}

	n := &ppool.Network{
		Metadata: &pn0stack.Metadata{
			Name:    "test-network",
			Version: 1,
		},
		Spec: &ppool.NetworkSpec{
			Ipv4Cidr: "192.168.0.0/30",
			Domain:   "test.local",
		},
		Status: &ppool.NetworkStatus{
			State: ppool.NetworkStatus_AVAILABLE,
		},
	}

	_, err = na.ApplyNetwork(context.Background(), &ppool.ApplyNetworkRequest{
		Metadata: &pn0stack.Metadata{
			Name: n.Metadata.Name,
		},
		Spec: n.Spec,
	})
	if err != nil {
		t.Fatalf("Failed to apply network: err='%s'", err.Error())
	}

	_, err = na.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
		Name:                 n.Metadata.Name,
		NetworkInterfaceName: "hogehoge",
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("[Invalid: no reserved network interface on Network -> NotFound] ReleaseNetworkInterface got error: err='%s'", err.Error())
	}

	reserveReq := &ppool.ReserveNetworkInterfaceRequest{
		Name:                 n.Metadata.Name,
		NetworkInterfaceName: "test-network-interface",
		NetworkInterface: &pbudget.NetworkInterface{
			Ipv4Address:     "192.168.0.1",
			HardwareAddress: "00:00:00:00:00:00",
		},
	}
	reserveRes, err := na.ReserveNetworkInterface(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("[Valid: no HardwareAddress] ReserveNetworkInterface got error: err='%s'", err.Error())
	}
	reserveRes.NetworkInterface.XXX_sizecache = 0
	if diff := cmp.Diff(reserveReq.Name, reserveRes.Name); diff != "" {
		t.Errorf("[Valid: no HardwareAddress] ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(reserveReq.NetworkInterfaceName, reserveRes.NetworkInterfaceName); diff != "" {
		t.Errorf("[Valid: no HardwareAddress] ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(reserveReq.NetworkInterface.Ipv4Address, reserveRes.NetworkInterface.Ipv4Address); diff != "" {
		t.Errorf("[Valid: no HardwareAddress] ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}
	if reserveRes.NetworkInterface.HardwareAddress == "" {
		t.Errorf("[Valid: no HardwareAddress] ReserveStorage response has blank hardware address, struct hardware address when getting blank request")
	}

	reserveReq = &ppool.ReserveNetworkInterfaceRequest{
		Name:                 n.Metadata.Name,
		NetworkInterfaceName: "test-network-interface2",
	}
	reserveRes, err = na.ReserveNetworkInterface(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("[Valid: no NetworkInterface] ReserveNetworkInterface got error: err='%s'", err.Error())
	}
	reserveRes.NetworkInterface.XXX_sizecache = 0
	if diff := cmp.Diff(reserveReq.Name, reserveRes.Name); diff != "" {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(reserveReq.NetworkInterfaceName, reserveRes.NetworkInterfaceName); diff != "" {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}
	if reserveRes.NetworkInterface.Ipv4Address != "192.168.0.2" {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response is wrong: ipv4_address_have=%s, ipv4_address_want=%s", reserveRes.NetworkInterface.Ipv4Address, "192.168.0.2")
	}
	if reserveRes.NetworkInterface.HardwareAddress == "" {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response has blank hardware address, struct hardware address when getting blank request")
	}

	// errors
	cases := []struct {
		name       string
		req        *ppool.ReserveNetworkInterfaceRequest
		res        *ppool.ReserveNetworkInterfaceResponse
		statusCode codes.Code
	}{
		{
			"Invalid: already exists",
			reserveReq,
			nil,
			codes.AlreadyExists,
		},
		{
			"Invalid: no NetworkInterfaceName -> InvalidArgument",
			&ppool.ReserveNetworkInterfaceRequest{
				Name: "invalid_argument",
				NetworkInterface: &pbudget.NetworkInterface{
					Ipv4Address: "192.168.0.1",
				},
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: no Name -> NotFound",
			&ppool.ReserveNetworkInterfaceRequest{
				NetworkInterfaceName: "not_found",
				NetworkInterface: &pbudget.NetworkInterface{
					Ipv4Address: "192.168.0.1",
				},
			},
			nil,
			codes.NotFound,
		},
		{
			"Invalid: no all -> InvalidArgument",
			&ppool.ReserveNetworkInterfaceRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: ipv4 address is over -> ResourceExhausted",
			&ppool.ReserveNetworkInterfaceRequest{
				Name:                 n.Metadata.Name,
				NetworkInterfaceName: "resource_exhausted",
				NetworkInterface: &pbudget.NetworkInterface{
					Ipv4Address: "192.168.0.1",
				},
			},
			nil,
			codes.ResourceExhausted,
		},
		{
			"Invalid: free ipv4 address is none -> ResourceExhausted",
			&ppool.ReserveNetworkInterfaceRequest{
				Name:                 n.Metadata.Name,
				NetworkInterfaceName: "resource_exhausted",
			},
			nil,
			codes.ResourceExhausted,
		},
		{
			"Invalid: Ipv4Address=aa -> InvalidArgument",
			&ppool.ReserveNetworkInterfaceRequest{
				Name:                 n.Metadata.Name,
				NetworkInterfaceName: "invalid_argument",
				NetworkInterface: &pbudget.NetworkInterface{
					Ipv4Address:     "aa",
					HardwareAddress: "00:00:00:00:00:00",
				},
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: Ipv4Address=aa -> InvalidArgument",
			&ppool.ReserveNetworkInterfaceRequest{
				Name:                 n.Metadata.Name,
				NetworkInterfaceName: "invalid_argument",
				NetworkInterface: &pbudget.NetworkInterface{
					Ipv4Address:     "192.168.10.1",
					HardwareAddress: "00:00:00:00:00:00",
				},
			},
			nil,
			codes.InvalidArgument,
		},
		// {
		// 	"Invalid: HardwareAddress=aa -> InvalidArgument",
		// 	&ppool.ReserveNetworkInterfaceRequest{
		// 		Name:                 n.Metadata.Name,
		// 		NetworkInterfaceName: "invalid_argument",
		// 		NetworkInterface: &pbudget.NetworkInterface{
		// 			Ipv4Address:     "192.168.0.1",
		// 			HardwareAddress: "aa",
		// 		},
		// 	},
		// 	nil,
		// 	codes.InvalidArgument,
		// },
	}

	for _, c := range cases {
		res, err := na.ReserveNetworkInterface(context.Background(), c.req)
		if err != nil && status.Code(err) != c.statusCode {
			t.Errorf("[%s] ReserveNetworkInterface got error: err='%s'", c.name, err.Error())
		}
		if res != c.res {
			t.Errorf("[%s] ReserveNetworkInterface response is not nil: res=%+v", c.name, reserveRes)
		}
	}

	releaseCases := []struct {
		name       string
		req        *ppool.ReleaseNetworkInterfaceRequest
		statusCode codes.Code
	}{
		{
			"no Name -> NotFound",
			&ppool.ReleaseNetworkInterfaceRequest{
				NetworkInterfaceName: reserveReq.NetworkInterfaceName,
			},
			codes.NotFound,
		},
		{
			"no StorageName -> NotFound",
			&ppool.ReleaseNetworkInterfaceRequest{
				Name: reserveReq.Name,
			},
			codes.NotFound,
		},
	}

	for _, c := range releaseCases {
		_, err := na.ReleaseNetworkInterface(context.Background(), c.req)
		if err != nil && status.Code(err) != c.statusCode {
			t.Errorf("[%s] ReleaseNetworkInterface got error: err='%s'", c.name, err.Error())
		}
	}

	_, err = na.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
		Name:                 reserveReq.Name,
		NetworkInterfaceName: reserveReq.NetworkInterfaceName,
	})
	if err != nil {
		t.Errorf("ReleaseNetworkInterface got error: err='%s'", err.Error())
	}
}
