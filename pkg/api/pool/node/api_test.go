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

	// errors
	// TODO: memory, CPU
	reserveRes, err = na.ReserveCompute(context.Background(), reserveReq)
	if err != nil && status.Code(err) != codes.AlreadyExists {
		t.Errorf("[AlreadyExists] ReserveCompute got error, not AlreadyExists: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[AlreadyExists] ReserveCompute response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
		Name:        n.Metadata.Name,
		ComputeName: "test-compute2",
		Compute: &pbudget.Compute{
			LimitCpuMilliCore:   1000,
			RequestCpuMilliCore: 1000,
			LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
			RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		},
	})
	if err != nil && status.Code(err) != codes.ResourceExhausted {
		t.Errorf("[ResourceExhausted] ReserveCompute got error, not ResourceExhausted: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[ResourceExhausted] ReserveCompute response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
		Name: "not_found",
		Compute: &pbudget.Compute{
			LimitCpuMilliCore:   1000,
			RequestCpuMilliCore: 1000,
			LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
			RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		},
	})
	if err != nil && status.Code(err) != codes.InvalidArgument {
		t.Errorf("[InvalidArgument] ReserveCompute got error, not InvalidArgument: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[InvalidArgument] ReserveCompute response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
		Name:        n.Metadata.Name,
		ComputeName: "test-compute2",
	})
	if err != nil && status.Code(err) != codes.InvalidArgument {
		t.Errorf("[InvalidArgument] ReserveCompute got error, not InvalidArgument: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[InvalidArgument] ReserveCompute response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
		ComputeName: "not_found",
		Compute: &pbudget.Compute{
			LimitCpuMilliCore:   1000,
			RequestCpuMilliCore: 1000,
			LimitMemoryBytes:    1 * bytefmt.GIGABYTE,
			RequestMemoryBytes:  1 * bytefmt.GIGABYTE,
		},
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("[NotFound] ReserveCompute got error, not NotFound: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[NotFound] ReserveCompute response is not nil: res=%+v", reserveRes)
	}

	_, err = na.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		ComputeName: reserveReq.ComputeName,
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("[NotFound] ReleaseCompute got error, not NotFound: err='%s'", err.Error())
	}
	_, err = na.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		Name: reserveReq.Name,
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("[NotFound] ReleaseCompute got error, not NotFound: err='%s'", err.Error())
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
	reserveRes, err = na.ReserveStorage(context.Background(), reserveReq)
	if err != nil && status.Code(err) != codes.AlreadyExists {
		t.Errorf("[AlreadyExists] ReserveStorage got error, not AlreadyExists: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[AlreadyExists] ReserveStorage response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveStorage(context.Background(), &ppool.ReserveStorageRequest{
		Name:        n.Metadata.Name,
		StorageName: "test-storage2",
		Storage: &pbudget.Storage{
			LimitBytes:   1 * bytefmt.GIGABYTE,
			RequestBytes: 1 * bytefmt.GIGABYTE,
		},
	})
	if err != nil && status.Code(err) != codes.ResourceExhausted {
		t.Errorf("[ResourceExhausted] ReserveStorage got error, not ResourceExhausted: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[ResourceExhausted] ReserveStorage response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveStorage(context.Background(), &ppool.ReserveStorageRequest{
		Name: "not_found",
		Storage: &pbudget.Storage{
			LimitBytes:   1 * bytefmt.GIGABYTE,
			RequestBytes: 1 * bytefmt.GIGABYTE,
		},
	})
	if err != nil && status.Code(err) != codes.InvalidArgument {
		t.Errorf("[InvalidArgument] ReserveStorage got error, not InvalidArgument: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[InvalidArgument] ReserveStorage response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveStorage(context.Background(), &ppool.ReserveStorageRequest{
		Name:        n.Metadata.Name,
		StorageName: "test-storage2",
	})
	if err != nil && status.Code(err) != codes.InvalidArgument {
		t.Errorf("[InvalidArgument] ReserveStorage got error, not InvalidArgument: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[InvalidArgument] ReserveStorage response is not nil: res=%+v", reserveRes)
	}
	reserveRes, err = na.ReserveStorage(context.Background(), &ppool.ReserveStorageRequest{
		StorageName: "not_found",
		Storage: &pbudget.Storage{
			LimitBytes:   1 * bytefmt.GIGABYTE,
			RequestBytes: 1 * bytefmt.GIGABYTE,
		},
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("[NotFound] ReserveStorage got error, not NotFound: err='%s'", err.Error())
	}
	if reserveRes != nil {
		t.Errorf("[NotFound] ReserveStorage response is not nil: res=%+v", reserveRes)
	}

	_, err = na.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		StorageName: reserveReq.StorageName,
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("[NotFound] ReleaseStorage got error, not NotFound: err='%s'", err.Error())
	}
	_, err = na.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		Name: reserveReq.Name,
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Errorf("[NotFound] ReleaseStorage got error, not NotFound: err='%s'", err.Error())
	}

	_, err = na.ReleaseStorage(context.Background(), &ppool.ReleaseStorageRequest{
		Name:        reserveReq.Name,
		StorageName: reserveReq.StorageName,
	})
	if err != nil {
		t.Errorf("ReleaseStorage got error: err='%s'", err.Error())
	}
}
