package network

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/google/go-cmp/cmp"
	"n0st.ac/n0stack/n0core/pkg/datastore/memory"
	ppool "n0st.ac/n0stack/n0proto.go/pool/v0"
)

func TestListNetworkAboutEmpty(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNetworkAPI(m)

	listRes, err := na.ListNetworks(context.Background(), &ppool.ListNetworksRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("ListNetworks got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListNetworks do not return nil: res='%s'", listRes)
	}
}

func TestGetNetworkAboutError(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNetworkAPI(m)

	cases := []struct {
		name    string
		req     *ppool.GetNetworkRequest
		res     *ppool.Network
		errCode codes.Code
	}{
		{
			"empty",
			&ppool.GetNetworkRequest{
				Name: "",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"not found",
			&ppool.GetNetworkRequest{
				Name: "hogehoge",
			},
			nil,
			codes.NotFound,
		},
	}

	for _, c := range cases {
		res, err := na.GetNetwork(context.Background(), c.req)
		if err == nil {
			t.Errorf("[%s] GetNetwork do not get error", c.name)
		} else if grpc.Code(err) != c.errCode {
			t.Errorf("[%s] GetNetwork get wrong error: want='%v', have='%v'", c.name, c.errCode, grpc.Code(err))
		}

		if res != c.res {
			t.Errorf("[%s] GetNetwork is incorrect: want='%v', have='%v'", c.name, c.res, res)
		}
	}
}

func TestApplyNetwork(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNetworkAPI(m)

	applyRes, err := na.ApplyNetwork(context.Background(), &ppool.ApplyNetworkRequest{
		Name: "test-network",
		Annotations: map[string]string{
			"test-annotation": "testing",
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		Ipv4Cidr: "192.168.0.0/30",
		Domain:   "test.local",
	})
	if err != nil {
		t.Fatalf("Failed to apply network: err='%s'", err.Error())
	}

	expected := &ppool.Network{
		Name: "test-network",
		Annotations: map[string]string{
			"test-annotation": "testing",
		},
		Labels: map[string]string{
			"test-label": "testing",
		},

		Ipv4Cidr: "192.168.0.0/30",
		Domain:   "test.local",
		State:    ppool.Network_AVAILABLE,
	}
	// diffが取れないので
	applyRes.XXX_sizecache = 0
	if diff := cmp.Diff(expected, applyRes); diff != "" {
		t.Fatalf("ApplyNetwork response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := na.ListNetworks(context.Background(), &ppool.ListNetworksRequest{})
	if err != nil {
		t.Errorf("ListNetworks got error: err='%s'", err.Error())
	}
	if len(listRes.Networks) != 1 {
		t.Errorf("ListNetworks response is wrong: have='%d', want='%d'", len(listRes.Networks), 1)
	}

	getRes, err := na.GetNetwork(context.Background(), &ppool.GetNetworkRequest{Name: expected.Name})
	if err != nil {
		t.Errorf("GetNetwork got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(expected, getRes); diff != "" {
		t.Errorf("GetNetwork response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := na.DeleteNetwork(context.Background(), &ppool.DeleteNetworkRequest{Name: expected.Name}); err != nil {
		t.Errorf("DeleteNetwork got error: err='%s'", err.Error())
	}
}

func TestApplyNetworkTableDriven(t *testing.T) {
	cases := []struct {
		name string
		req  *ppool.ApplyNetworkRequest
		res  *ppool.Network
		code codes.Code
	}{
		{
			"full",
			&ppool.ApplyNetworkRequest{
				Name:     "test-network",
				Ipv4Cidr: "192.168.0.0/30",
				Ipv6Cidr: "2001:db8::/64",
				Domain:   "test.local",
			},
			&ppool.Network{
				Name:     "test-network",
				Ipv4Cidr: "192.168.0.0/30",
				Ipv6Cidr: "2001:db8::/64",
				Domain:   "test.local",
				State:    ppool.Network_AVAILABLE,
			},
			codes.OK,
		},
		{
			"only ipv4",
			&ppool.ApplyNetworkRequest{
				Name:     "test-network",
				Ipv4Cidr: "192.168.1.0/30",
				Domain:   "test.local",
			},
			&ppool.Network{
				Name:     "test-network",
				Ipv4Cidr: "192.168.1.0/30",
				Domain:   "test.local",
				State:    ppool.Network_AVAILABLE,
			},
			codes.OK,
		},
		{
			"only ipv6",
			&ppool.ApplyNetworkRequest{
				Name:     "test-network",
				Ipv6Cidr: "2001:db8::/64",
				Domain:   "test.local",
			},
			&ppool.Network{
				Name:     "test-network",
				Ipv6Cidr: "2001:db8::/64",
				Domain:   "test.local",
				State:    ppool.Network_AVAILABLE,
			},
			codes.OK,
		},
		{
			"no network",
			&ppool.ApplyNetworkRequest{
				Name:   "test-network",
				Domain: "test.local",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"no name",
			&ppool.ApplyNetworkRequest{
				Ipv4Cidr: "192.168.0.0/30",
				Ipv6Cidr: "2001:db8::/64",
				Domain:   "test.local",
			},
			nil,
			codes.InvalidArgument,
		},
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := memory.NewMemoryDatastore()
	na := NewMockNetworkAPI(m)
	originDatastore := na.api.dataStore
	for i, c := range cases {
		na.api.dataStore = originDatastore.AddPrefix(fmt.Sprintf("%d", i))
		res, err := na.ApplyNetwork(ctx, c.req)

		if res != nil {
			res.XXX_sizecache = 0
		}
		if diff := cmp.Diff(c.res, res); diff != "" {
			t.Errorf("[%s] ApplyNetwork response is wrong: diff=(-want +got)\n%s", c.name, diff)
		}

		if grpc.Code(err) != c.code {
			t.Errorf("[%s] Response code is invalid, want=%v, have=%v", c.name, c.code, grpc.Code(err))
		}
	}
}

func TestNetworkAboutNetworkInterface(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNetworkAPI(m)

	n := &ppool.Network{
		Name:     "test-network",
		Ipv4Cidr: "192.168.0.0/30",
		Domain:   "test.local",
		State:    ppool.Network_AVAILABLE,
	}

	_, err := na.ApplyNetwork(context.Background(), &ppool.ApplyNetworkRequest{
		Name:     "test-network",
		Ipv4Cidr: "192.168.0.0/30",
		Domain:   "test.local",
	})
	if err != nil {
		t.Fatalf("Failed to apply network: err='%s'", err.Error())
	}

	_, err = na.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
		NetworkName:          n.Name,
		NetworkInterfaceName: "hogehoge",
	})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("[Invalid: no reserved network interface on Network -> NotFound] ReleaseNetworkInterface got error: err='%s'", err.Error())
	}

	reserveReq := &ppool.ReserveNetworkInterfaceRequest{
		NetworkName:          n.Name,
		NetworkInterfaceName: "test-network-interface",
		Ipv4Address:          "192.168.0.1",
		HardwareAddress:      "00:00:00:00:00:00",
	}
	reserveRes, err := na.ReserveNetworkInterface(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("[Valid: no HardwareAddress] ReserveNetworkInterface got error: err='%s'", err.Error())
	}
	if _, ok := reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName]; !ok {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response do not have requested network interface")
	}
	if reserveReq.HardwareAddress != reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName].HardwareAddress {
		t.Errorf("[Valid: no HardwareAddress] ReserveStorage response about 'HardwareAddress' is wrong: want=%s, have=%s", reserveReq.HardwareAddress, reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName].HardwareAddress)
	}
	if reserveReq.Ipv4Address != reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName].Ipv4Address {
		t.Errorf("[Valid: no HardwareAddress] ReserveStorage response about 'Ipv4Address' is wrong: want=%s, have=%s", reserveReq.Ipv4Address, reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName].Ipv4Address)
	}
	if reserveReq.NetworkName != reserveRes.Name {
		t.Errorf("[Valid: no HardwareAddress] ReserveStorage response about 'Name' is wrong: want=%s, have=%s", reserveReq.NetworkName, reserveRes.Name)
	}

	reserveReq = &ppool.ReserveNetworkInterfaceRequest{
		NetworkName:          n.Name,
		NetworkInterfaceName: "test-network-interface2",
	}
	reserveRes, err = na.ReserveNetworkInterface(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("[Valid: no NetworkInterface] ReserveNetworkInterface got error: err='%s'", err.Error())
	}
	if _, ok := reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName]; !ok {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response do not have requested network interface")
	}
	if reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName].Ipv4Address != "" {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response is wrong: ipv4_address_have=%s, ipv4_address_want=%s", reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName].Ipv4Address, "")
	}
	if reserveRes.ReservedNetworkInterfaces[reserveReq.NetworkInterfaceName].HardwareAddress == "" {
		t.Errorf("[Valid: no NetworkInterface] ReserveStorage response has blank hardware address, struct hardware address when getting blank request")
	}

	// errors
	cases := []struct {
		name       string
		req        *ppool.ReserveNetworkInterfaceRequest
		res        *ppool.Network
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
				NetworkName: "invalid_argument",
				Ipv4Address: "192.168.0.1",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: no Name -> NotFound",
			&ppool.ReserveNetworkInterfaceRequest{
				NetworkInterfaceName: "not_found",
				Ipv4Address:          "192.168.0.1",
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
				NetworkName:          n.Name,
				NetworkInterfaceName: "resource_exhausted",
				Ipv4Address:          "192.168.0.1",
			},
			nil,
			codes.ResourceExhausted,
		},
		// {
		// 	"Invalid: free ipv4 address is none -> ResourceExhausted",
		// 	&ppool.ReserveNetworkInterfaceRequest{
		// 		NetworkName:          n.Name,
		// 		NetworkInterfaceName: "resource_exhausted",
		// 	},
		// 	nil,
		// 	codes.ResourceExhausted,
		// },
		{
			"Invalid: Ipv4Address=aa -> InvalidArgument",
			&ppool.ReserveNetworkInterfaceRequest{
				NetworkName:          n.Name,
				NetworkInterfaceName: "invalid_argument",
				Ipv4Address:          "aa",
				HardwareAddress:      "00:00:00:00:00:00",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: Ipv4Address=aa -> InvalidArgument",
			&ppool.ReserveNetworkInterfaceRequest{
				NetworkName:          n.Name,
				NetworkInterfaceName: "invalid_argument",
				Ipv4Address:          "192.168.10.1",
				HardwareAddress:      "00:00:00:00:00:00",
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
		if err != nil && grpc.Code(err) != c.statusCode {
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
				NetworkName: reserveReq.NetworkName,
			},
			codes.NotFound,
		},
	}

	for _, c := range releaseCases {
		_, err := na.ReleaseNetworkInterface(context.Background(), c.req)
		if err != nil && grpc.Code(err) != c.statusCode {
			t.Errorf("[%s] ReleaseNetworkInterface got error: err='%s'", c.name, err.Error())
		}
	}

	_, err = na.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
		NetworkName:          reserveReq.NetworkName,
		NetworkInterfaceName: reserveReq.NetworkInterfaceName,
	})
	if err != nil {
		t.Errorf("ReleaseNetworkInterface got error: err='%s'", err.Error())
	}
}

func TestDeletionLock(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNetworkAPI(m)

	n := &ppool.Network{
		Name:     "test-network",
		Ipv4Cidr: "192.168.0.0/30",
		Domain:   "test.local",
		State:    ppool.Network_AVAILABLE,
	}

	_, err := na.ApplyNetwork(context.Background(), &ppool.ApplyNetworkRequest{
		Name:     "test-network",
		Ipv4Cidr: "192.168.0.0/30",
		Domain:   "test.local",
	})
	if err != nil {
		t.Fatalf("Failed to apply network: err='%s'", err.Error())
	}

	reserveReq := &ppool.ReserveNetworkInterfaceRequest{
		NetworkName:          n.Name,
		NetworkInterfaceName: "test-network-interface",
		Ipv4Address:          "192.168.0.1",
		HardwareAddress:      "00:00:00:00:00:00",
	}
	_, err = na.ReserveNetworkInterface(context.Background(), reserveReq)
	if err != nil {
		t.Fatalf("ReserveNetworkInterface got error: err='%s'", err.Error())
	}

	_, err = na.DeleteNetwork(context.Background(), &ppool.DeleteNetworkRequest{
		Name: n.Name,
	})
	if grpc.Code(err) != codes.FailedPrecondition {
		t.Errorf("deletion lock do not worked")
	}
}
