// +build medium

package node

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestEmptyNode(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNodeAPI(m)

	listRes, err := na.ListNodes(context.Background(), &ppool.ListNodesRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("ListNode got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListNodes do not return nil: res='%s'", listRes)
	}

	getRes, err := na.GetNode(context.Background(), &ppool.GetNodeRequest{})
	if err != nil && grpc.Code(err) != codes.NotFound {
		t.Errorf("GetNode got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetNode do not return nil: res='%s'", listRes)
	}
}

func TestApplyNode(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNodeAPI(m)

	n := &ppool.Node{
		Name:    "test-node",
		Version: 1,

		Address:     "10.0.0.1",
		IpmiAddress: "192.168.0.1",
		Serial:      "aa",

		CpuMilliCores: 1000,
		MemoryBytes:   1 * bytefmt.GIGABYTE,
		StorageBytes:  10 * bytefmt.GIGABYTE,

		Datacenter:       "test-dc",
		AvailavilityZone: "test-az",
		Cell:             "test-cell",
		Rack:             "test-rack",
		Unit:             1,

		State: ppool.Node_Ready,
	}

	applyRes, err := na.ApplyNode(context.Background(), &ppool.ApplyNodeRequest{
		Name: "test-node",

		Address:     "10.0.0.1",
		IpmiAddress: "192.168.0.1",
		Serial:      "aa",

		CpuMilliCores: 1000,
		MemoryBytes:   1 * bytefmt.GIGABYTE,
		StorageBytes:  10 * bytefmt.GIGABYTE,

		Datacenter:       "test-dc",
		AvailavilityZone: "test-az",
		Cell:             "test-cell",
		Rack:             "test-rack",
		Unit:             1,
	})
	if err != nil {
		t.Fatalf("Failed to apply node: err='%s'", err.Error())
	}

	// diffが取れないので
	applyRes.XXX_sizecache = 0
	if diff := cmp.Diff(n, applyRes); diff != "" {
		t.Fatalf("ApplyNode response is wrong: diff=(-want +got)\n%s", diff)
	}

	listRes, err := na.ListNodes(context.Background(), &ppool.ListNodesRequest{})
	if err != nil {
		t.Errorf("ListNode got error: err='%s'", err.Error())
	}
	if len(listRes.Nodes) != 1 {
		t.Errorf("ListNodes response is wrong: have='%d', want='%d'", len(listRes.Nodes), 1)
	}

	getRes, err := na.GetNode(context.Background(), &ppool.GetNodeRequest{Name: n.Name})
	if err != nil {
		t.Errorf("GetNode got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(n, getRes); diff != "" {
		t.Errorf("GetNode response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := na.DeleteNode(context.Background(), &ppool.DeleteNodeRequest{Name: n.Name}); err != nil {
		t.Errorf("DeleteNode got error: err='%s'", err.Error())
	}
}

func TestNodeAboutCompute(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNodeAPI(m)

	n := &ppool.Node{
		Name:    "test-node",
		Version: 1,

		Address:     "10.0.0.1",
		IpmiAddress: "192.168.0.1",
		Serial:      "aa",

		CpuMilliCores: 1000,
		MemoryBytes:   1 * bytefmt.GIGABYTE,
		StorageBytes:  10 * bytefmt.GIGABYTE,

		Datacenter:       "test-dc",
		AvailavilityZone: "test-az",
		Cell:             "test-cell",
		Rack:             "test-rack",
		Unit:             1,

		State: ppool.Node_Ready,
	}

	_, err := na.ApplyNode(context.Background(), &ppool.ApplyNodeRequest{
		Name: "test-node",

		Address:     "10.0.0.1",
		IpmiAddress: "192.168.0.1",
		Serial:      "aa",

		CpuMilliCores: 1000,
		MemoryBytes:   1 * bytefmt.GIGABYTE,
		StorageBytes:  10 * bytefmt.GIGABYTE,

		Datacenter:       "test-dc",
		AvailavilityZone: "test-az",
		Cell:             "test-cell",
		Rack:             "test-rack",
		Unit:             1,
	})
	if err != nil {
		t.Fatalf("Failed to apply node: err='%s'", err.Error())
	}

	reserveReq := &ppool.ReserveComputeRequest{
		NodeName:            n.Name,
		ComputeName:         "test-compute",
		LimitCpuMilliCore:   1000,
		RequestCpuMilliCore: 1000,
		LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
		RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
	}
	reserveRes, err := na.ReserveCompute(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("ReserveCompute got error: err='%s'", err.Error())
	}
	if reserveReq.NodeName != reserveRes.Name {
		t.Errorf("ReserveCompute response about 'Name' is wrong: want=%s, have=%s", reserveReq.NodeName, reserveRes.Name)
	}
	if _, ok := reserveRes.ReservedComputes[reserveReq.ComputeName]; !ok {
		t.Errorf("ReserveCompute response do not have requested compute")
	}

	cases := []struct {
		name       string
		req        *ppool.ReserveComputeRequest
		res        *ppool.Node
		statusCode codes.Code
	}{
		{
			"Invalid: already exists",
			reserveReq,
			nil,
			codes.AlreadyExists,
		},
		{
			"Invalid: no ComputeName -> InvalidArgument",
			&ppool.ReserveComputeRequest{
				NodeName:            "invalid_argument",
				LimitCpuMilliCore:   1000,
				RequestCpuMilliCore: 1000,
				LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
				RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: no Compute -> InvalidArgument",
			&ppool.ReserveComputeRequest{
				NodeName:    n.Name,
				ComputeName: "test-compute2",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: no NodeName -> NotFound",
			&ppool.ReserveComputeRequest{
				ComputeName:         "not_found",
				LimitCpuMilliCore:   1000,
				RequestCpuMilliCore: 1000,
				LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
				RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
			},
			nil,
			codes.NotFound,
		},
		{
			"Invalid: no all -> InvalidArgument",
			&ppool.ReserveComputeRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: request over -> ResourceExhausted",
			&ppool.ReserveComputeRequest{
				NodeName:            n.Name,
				ComputeName:         "test-compute2",
				RequestCpuMilliCore: 1000,
				RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
			},
			nil,
			codes.ResourceExhausted,
		},
	}

	// TODO: memory, CPU
	for _, c := range cases {
		res, err := na.ReserveCompute(context.Background(), c.req)
		if err != nil && grpc.Code(err) != c.statusCode {
			t.Errorf("[%s] ReserveCompute got error: err='%s'", c.name, err.Error())
		}
		if res != c.res {
			t.Errorf("[%s] ReserveCompute response is not nil: res=%+v", c.name, reserveRes)
		}
	}

	releaseCases := []struct {
		name       string
		req        *ppool.ReleaseComputeRequest
		statusCode codes.Code
	}{
		{
			"Invalid: no NodeName -> NotFound",
			&ppool.ReleaseComputeRequest{
				ComputeName: reserveReq.ComputeName,
			},
			codes.NotFound,
		},
		{
			"Invalid: no ComputeName -> NotFound",
			&ppool.ReleaseComputeRequest{
				NodeName: reserveReq.NodeName,
			},
			codes.NotFound,
		},
	}

	for _, c := range releaseCases {
		_, err := na.ReleaseCompute(context.Background(), c.req)
		if err != nil && grpc.Code(err) != c.statusCode {
			t.Errorf("[%s] ReleaseCompute got error: err='%s'", c.name, err.Error())
		}
	}

	_, err = na.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		NodeName:    reserveReq.NodeName,
		ComputeName: reserveReq.ComputeName,
	})
	if err != nil {
		t.Errorf("ReleaseCompute got error: err='%s'", err.Error())
	}
}

func TestNodeAboutStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na := NewMockNodeAPI(m)

	n := &ppool.Node{
		Name:    "test-node",
		Version: 1,

		Address:     "10.0.0.1",
		IpmiAddress: "192.168.0.1",
		Serial:      "aa",

		CpuMilliCores: 1000,
		MemoryBytes:   1 * bytefmt.GIGABYTE,
		StorageBytes:  10 * bytefmt.GIGABYTE,

		Datacenter:       "test-dc",
		AvailavilityZone: "test-az",
		Cell:             "test-cell",
		Rack:             "test-rack",
		Unit:             1,

		State: ppool.Node_Ready,
	}

	_, err := na.ApplyNode(context.Background(), &ppool.ApplyNodeRequest{
		Name: "test-node",

		Address:     "10.0.0.1",
		IpmiAddress: "192.168.0.1",
		Serial:      "aa",

		CpuMilliCores: 1000,
		MemoryBytes:   1 * bytefmt.GIGABYTE,
		StorageBytes:  1 * bytefmt.GIGABYTE,

		Datacenter:       "test-dc",
		AvailavilityZone: "test-az",
		Cell:             "test-cell",
		Rack:             "test-rack",
		Unit:             1,
	})
	if err != nil {
		t.Fatalf("Failed to apply node: err='%s'", err.Error())
	}

	reserveReq := &ppool.ReserveStorageRequest{
		NodeName:     n.Name,
		StorageName:  "test-storage",
		LimitBytes:   1 * bytefmt.GIGABYTE,
		RequestBytes: 1 * bytefmt.GIGABYTE,
	}
	reserveRes, err := na.ReserveStorage(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("ReserveStorage got error: err='%s'", err.Error())
	}
	if reserveReq.NodeName != reserveRes.Name {
		t.Errorf("ReserveStorage response about 'Name' is wrong: want=%s, have=%s", reserveReq.NodeName, reserveRes.Name)
	}
	if _, ok := reserveRes.ReservedStorages[reserveReq.StorageName]; !ok {
		t.Errorf("ReserveStorage response do not have requested compute")
	}

	// errors
	cases := []struct {
		name       string
		req        *ppool.ReserveStorageRequest
		res        *ppool.Node
		statusCode codes.Code
	}{
		{
			"Invalid: already exists",
			reserveReq,
			nil,
			codes.AlreadyExists,
		},
		{
			"Invalid: no StorageName -> InvalidArgument",
			&ppool.ReserveStorageRequest{
				NodeName:     "not_found",
				LimitBytes:   1 * bytefmt.GIGABYTE,
				RequestBytes: 1 * bytefmt.GIGABYTE,
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: no Storage -> InvalidArgument",
			&ppool.ReserveStorageRequest{
				NodeName:    n.Name,
				StorageName: "test-storage2",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: no NodeName -> NotFound",
			&ppool.ReserveStorageRequest{
				StorageName:  "not_found",
				LimitBytes:   1 * bytefmt.GIGABYTE,
				RequestBytes: 1 * bytefmt.GIGABYTE,
			},
			nil,
			codes.NotFound,
		},
		{
			"Invalid: no all -> InvalidArgument",
			&ppool.ReserveStorageRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"Invalid: request over -> ResourceExhausted",
			&ppool.ReserveStorageRequest{
				NodeName:     n.Name,
				StorageName:  "test-storage2",
				LimitBytes:   1 * bytefmt.GIGABYTE,
				RequestBytes: 1 * bytefmt.GIGABYTE,
			},
			nil,
			codes.ResourceExhausted,
		},
	}

	for _, c := range cases {
		res, err := na.ReserveStorage(context.Background(), c.req)
		if err != nil && grpc.Code(err) != c.statusCode {
			t.Errorf("[%s] ReserveStorage got error: err='%s'", c.name, err.Error())
		}
		if res != c.res {
			t.Errorf("[%s] ReserveStorage response is not nil: res=%+v", c.name, reserveRes)
		}
	}

	releaseCases := []struct {
		name       string
		req        *ppool.ReleaseStorageRequest
		statusCode codes.Code
	}{
		{
			"Invalid: no NodeName -> NotFound",
			&ppool.ReleaseStorageRequest{
				StorageName: reserveReq.StorageName,
			},
			codes.NotFound,
		},
		{
			"Invalid: no StorageName -> NotFound",
			&ppool.ReleaseStorageRequest{
				NodeName: reserveReq.NodeName,
			},
			codes.NotFound,
		},
	}

	for _, c := range releaseCases {
		_, err := na.ReleaseStorage(context.Background(), c.req)
		if err != nil && grpc.Code(err) != c.statusCode {
			t.Errorf("[%s] ReleaseStorage got error: err='%s'", c.name, err.Error())
		}
	}

	_, err = na.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		NodeName:    reserveReq.NodeName,
		StorageName: reserveReq.StorageName,
	})
	if err != nil {
		t.Errorf("ReleaseStorage got error: err='%s'", err.Error())
	}
}
