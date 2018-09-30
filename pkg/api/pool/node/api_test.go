// +build medium
// +build !without_external

package node

import (
	"context"
	"testing"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/go-cmp/cmp"
	"github.com/n0stack/n0core/pkg/datastore/memory"
	"github.com/n0stack/proto.go/budget/v0"
	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/v0"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEmptyNodes(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, err := CreateNodeAPI(m)
	if err != nil {
		t.Fatalf("Failed to create Node API: err='%s'", err.Error())
	}

	listRes, err := na.ListNodes(context.Background(), &ppool.ListNodesRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("ListNode got error, not NotFound: err='%s'", err.Error())
	}
	if listRes != nil {
		t.Errorf("ListNodes do not return nil: res='%s'", listRes)
	}

	getRes, err := na.GetNode(context.Background(), &ppool.GetNodeRequest{})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("GetNode got error, not NotFound: err='%s'", err.Error())
	}
	if getRes != nil {
		t.Errorf("GetNode do not return nil: res='%s'", listRes)
	}
}

func TestApplyNode(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, err := CreateNodeAPI(m)
	if err != nil {
		t.Fatalf("Failed to create Node API: err='%s'", err.Error())
	}

	n := &ppool.Node{
		Metadata: &pn0stack.Metadata{
			Name:    "test-node",
			Version: 1,
		},
		Spec: &ppool.NodeSpec{
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
		},
		Status: &ppool.NodeStatus{
			State: ppool.NodeStatus_Ready,
		},
	}

	applyRes, err := na.ApplyNode(context.Background(), &ppool.ApplyNodeRequest{
		Metadata: &pn0stack.Metadata{
			Name: n.Metadata.Name,
		},
		Spec: n.Spec,
	})
	if err != nil {
		t.Fatalf("Failed to apply node: err='%s'", err.Error())
	}

	applyRes.XXX_sizecache = 0
	applyRes.Metadata.XXX_sizecache = 0
	applyRes.Spec.XXX_sizecache = 0
	applyRes.Status.XXX_sizecache = 0
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

	getRes, err := na.GetNode(context.Background(), &ppool.GetNodeRequest{Name: n.Metadata.Name})
	if err != nil {
		t.Errorf("GetNode got error: err='%s'", err.Error())
	}
	if diff := cmp.Diff(n, getRes); diff != "" {
		t.Errorf("GetNode response is wrong: diff=(-want +got)\n%s", diff)
	}

	if _, err := na.DeleteNode(context.Background(), &ppool.DeleteNodeRequest{Name: n.Metadata.Name}); err != nil {
		t.Errorf("DeleteNode got error: err='%s'", err.Error())
	}
}

func TestNodeAboutCompute(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, err := CreateNodeAPI(m)
	if err != nil {
		t.Fatalf("Failed to create Node API: err='%s'", err.Error())
	}

	n := &ppool.Node{
		Metadata: &pn0stack.Metadata{
			Name:    "test-node",
			Version: 1,
		},
		Spec: &ppool.NodeSpec{
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
		},
		Status: &ppool.NodeStatus{
			State: ppool.NodeStatus_Ready,
		},
	}

	_, err = na.ApplyNode(context.Background(), &ppool.ApplyNodeRequest{
		Metadata: &pn0stack.Metadata{
			Name: n.Metadata.Name,
		},
		Spec: n.Spec,
	})
	if err != nil {
		t.Fatalf("Failed to apply node: err='%s'", err.Error())
	}

	reserveReq := &ppool.ReserveComputeRequest{
		Name:        n.Metadata.Name,
		ComputeName: "test-compute",
		Compute: &pbudget.Compute{
			LimitCpuMilliCore:   1000,
			RequestCpuMilliCore: 1000,
			LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
			RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		},
	}
	reserveRes, err := na.ReserveCompute(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("ReserveCompute got error: err='%s'", err.Error())
	}
	reserveRes.Compute.XXX_sizecache = 0
	if diff := cmp.Diff(reserveReq.Name, reserveRes.Name); diff != "" {
		t.Errorf("ReserveCompute response is wrong: diff=(-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(reserveReq.ComputeName, reserveRes.ComputeName); diff != "" {
		t.Errorf("ReserveCompute response is wrong: diff=(-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(reserveReq.Compute, reserveRes.Compute); diff != "" {
		t.Errorf("ReserveCompute response is wrong: diff=(-want +got)\n%s", diff)
	}

	cases := []struct {
		name       string
		req        *ppool.ReserveComputeRequest
		res        *ppool.ReserveComputeResponse
		statusCode codes.Code
	}{
		{
			"already exists",
			reserveReq,
			nil,
			codes.AlreadyExists,
		},
		{
			"no ComputeName -> InvalidArgument",
			&ppool.ReserveComputeRequest{
				Name: "invalid_argument",
				Compute: &pbudget.Compute{
					LimitCpuMilliCore:   1000,
					RequestCpuMilliCore: 1000,
					LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
					RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
				},
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"no Compute -> InvalidArgument",
			&ppool.ReserveComputeRequest{
				Name:        n.Metadata.Name,
				ComputeName: "test-compute2",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"no Name -> NotFound",
			&ppool.ReserveComputeRequest{
				ComputeName: "not_found",
				Compute: &pbudget.Compute{
					LimitCpuMilliCore:   1000,
					RequestCpuMilliCore: 1000,
					LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
					RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
				},
			},
			nil,
			codes.NotFound,
		},
		{
			"no all -> InvalidArgument",
			&ppool.ReserveComputeRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"request over -> ResourceExhausted",
			&ppool.ReserveComputeRequest{
				Name:        n.Metadata.Name,
				ComputeName: "test-compute2",
				Compute: &pbudget.Compute{
					RequestCpuMilliCore: 1000,
					RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
				},
			},
			nil,
			codes.ResourceExhausted,
		},
	}

	// TODO: memory, CPU
	for _, c := range cases {
		res, err := na.ReserveCompute(context.Background(), c.req)
		if err != nil && status.Code(err) != c.statusCode {
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
			"no Name -> NotFound",
			&ppool.ReleaseComputeRequest{
				ComputeName: reserveReq.ComputeName,
			},
			codes.NotFound,
		},
		{
			"no ComputeName -> NotFound",
			&ppool.ReleaseComputeRequest{
				Name: reserveReq.Name,
			},
			codes.NotFound,
		},
	}

	for _, c := range releaseCases {
		_, err := na.ReleaseCompute(context.Background(), c.req)
		if err != nil && status.Code(err) != c.statusCode {
			t.Errorf("[%s] ReleaseCompute got error: err='%s'", c.name, err.Error())
		}
	}

	_, err = na.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		Name:        reserveReq.Name,
		ComputeName: reserveReq.ComputeName,
	})
	if err != nil {
		t.Errorf("ReleaseCompute got error: err='%s'", err.Error())
	}
}

func TestNodeAboutStorage(t *testing.T) {
	m := memory.NewMemoryDatastore()
	na, err := CreateNodeAPI(m)
	if err != nil {
		t.Fatalf("Failed to create Node API: err='%s'", err.Error())
	}

	n := &ppool.Node{
		Metadata: &pn0stack.Metadata{
			Name:    "test-node",
			Version: 1,
		},
		Spec: &ppool.NodeSpec{
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
		},
		Status: &ppool.NodeStatus{
			State: ppool.NodeStatus_Ready,
		},
	}

	_, err = na.ApplyNode(context.Background(), &ppool.ApplyNodeRequest{
		Metadata: &pn0stack.Metadata{
			Name: n.Metadata.Name,
		},
		Spec: n.Spec,
	})
	if err != nil {
		t.Fatalf("Failed to apply node: err='%s'", err.Error())
	}

	reserveReq := &ppool.ReserveStorageRequest{
		Name:        n.Metadata.Name,
		StorageName: "test-storage",
		Storage: &pbudget.Storage{
			LimitBytes:   1 * bytefmt.GIGABYTE,
			RequestBytes: 1 * bytefmt.GIGABYTE,
		},
	}
	reserveRes, err := na.ReserveStorage(context.Background(), reserveReq)
	if err != nil {
		t.Errorf("ReserveStorage got error: err='%s'", err.Error())
	}
	reserveRes.Storage.XXX_sizecache = 0
	if diff := cmp.Diff(reserveReq.Name, reserveRes.Name); diff != "" {
		t.Errorf("ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(reserveReq.StorageName, reserveRes.StorageName); diff != "" {
		t.Errorf("ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(reserveReq.Storage, reserveRes.Storage); diff != "" {
		t.Errorf("ReserveStorage response is wrong: diff=(-want +got)\n%s", diff)
	}

	// errors
	cases := []struct {
		name       string
		req        *ppool.ReserveStorageRequest
		res        *ppool.ReserveStorageResponse
		statusCode codes.Code
	}{
		{
			"already exists",
			reserveReq,
			nil,
			codes.AlreadyExists,
		},
		{
			"no StorageName -> InvalidArgument",
			&ppool.ReserveStorageRequest{
				Name: "not_found",
				Storage: &pbudget.Storage{
					LimitBytes:   1 * bytefmt.GIGABYTE,
					RequestBytes: 1 * bytefmt.GIGABYTE,
				},
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"no Storage -> InvalidArgument",
			&ppool.ReserveStorageRequest{
				Name:        n.Metadata.Name,
				StorageName: "test-storage2",
			},
			nil,
			codes.InvalidArgument,
		},
		{
			"no Name -> NotFound",
			&ppool.ReserveStorageRequest{
				StorageName: "not_found",
				Storage: &pbudget.Storage{
					LimitBytes:   1 * bytefmt.GIGABYTE,
					RequestBytes: 1 * bytefmt.GIGABYTE,
				},
			},
			nil,
			codes.NotFound,
		},
		{
			"no all -> InvalidArgument",
			&ppool.ReserveStorageRequest{},
			nil,
			codes.InvalidArgument,
		},
		{
			"request over -> ResourceExhausted",
			&ppool.ReserveStorageRequest{
				Name:        n.Metadata.Name,
				StorageName: "test-storage2",
				Storage: &pbudget.Storage{
					LimitBytes:   1 * bytefmt.GIGABYTE,
					RequestBytes: 1 * bytefmt.GIGABYTE,
				},
			},
			nil,
			codes.ResourceExhausted,
		},
	}

	for _, c := range cases {
		res, err := na.ReserveStorage(context.Background(), c.req)
		if err != nil && status.Code(err) != c.statusCode {
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
			"no Name -> NotFound",
			&ppool.ReleaseStorageRequest{
				StorageName: reserveReq.StorageName,
			},
			codes.NotFound,
		},
		{
			"no StorageName -> NotFound",
			&ppool.ReleaseStorageRequest{
				Name: reserveReq.Name,
			},
			codes.NotFound,
		},
	}

	for _, c := range releaseCases {
		_, err := na.ReleaseStorage(context.Background(), c.req)
		if err != nil && status.Code(err) != c.statusCode {
			t.Errorf("[%s] ReleaseStorage got error: err='%s'", c.name, err.Error())
		}
	}

	_, err = na.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		Name:        reserveReq.Name,
		StorageName: reserveReq.StorageName,
	})
	if err != nil {
		t.Errorf("ReleaseStorage got error: err='%s'", err.Error())
	}
}
