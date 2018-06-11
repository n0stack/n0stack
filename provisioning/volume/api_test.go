package volume

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/n0stack/n0core/datastore/memory"
	"github.com/n0stack/proto.go/provisioning/v0"
	"github.com/n0stack/proto.go/v0"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestApplyVolume(t *testing.T) {
	a := VolumeAPI{
		dataStore: memory.NewMemoryDatastore(),
	}
	a.dataStore.Apply("version_test", &pprovisioning.Volume{
		Metadata: &pn0stack.Metadata{
			Name:    "version_test",
			Version: 3,
		},
		Spec: &pprovisioning.VolumeSpec{},
	})

	testCases := []struct {
		req    *pprovisioning.ApplyVolumeRequest
		volume *pprovisioning.Volume
		code   codes.Code
	}{
		{
			req: &pprovisioning.ApplyVolumeRequest{
				Metadata: &pn0stack.Metadata{
					Name:    "version_test",
					Version: 100,
				},
				Spec: &pprovisioning.VolumeSpec{},
			},
			volume: nil,
			code:   codes.InvalidArgument,
		},
		{
			req: &pprovisioning.ApplyVolumeRequest{
				Metadata: &pn0stack.Metadata{
					Name:    "version_0_test",
					Version: 100,
				},
				Spec: &pprovisioning.VolumeSpec{},
			},
			volume: nil,
			code:   codes.InvalidArgument,
		},
		{
			req: &pprovisioning.ApplyVolumeRequest{
				Metadata: &pn0stack.Metadata{
					Name: "lost_node_name",
				},
				Spec: &pprovisioning.VolumeSpec{},
			},
			volume: &pprovisioning.Volume{
				Metadata: &pn0stack.Metadata{
					Name: "lost_node_name",
				},
				Spec: &pprovisioning.VolumeSpec{},
				Status: &pprovisioning.VolumeStatus{
					State: pprovisioning.VolumeStatus_PENDING,
				},
			},
			code: codes.OK,
		},
	}

	for _, tc := range testCases {
		n, err := a.ApplyVolume(context.Background(), tc.req)

		if !reflect.DeepEqual(n, tc.volume) {
			t.Errorf("Wrong status value.\n\thave:%v\n\twant:%v", n, tc.volume)
		}

		if status.Code(err) != tc.code {
			t.Errorf("Wrong status code on %s.\n\thave:%v\n\twant:%v", tc.req.Metadata.Name, status.Code(err), tc.code)
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	a := VolumeAPI{
		dataStore: memory.NewMemoryDatastore(),
	}
	a.dataStore.Apply("lost_node_name", &pprovisioning.Volume{
		Metadata: &pn0stack.Metadata{
			Name:    "lost_node_name",
			Version: 1,
		},
		Spec: &pprovisioning.VolumeSpec{},
	})

	testCases := []struct {
		req    *pprovisioning.DeleteVolumeRequest
		volume *empty.Empty
		code   codes.Code
	}{
		{
			req: &pprovisioning.DeleteVolumeRequest{
				Name: "notfound",
			},
			volume: &empty.Empty{},
			code:   codes.NotFound,
		},
		{
			req: &pprovisioning.DeleteVolumeRequest{
				Name: "lost_node_name",
			},
			volume: &empty.Empty{},
			code:   codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		n, err := a.DeleteVolume(context.Background(), tc.req)

		if !reflect.DeepEqual(n, tc.volume) {
			t.Errorf("Wrong status value on %s.\n\thave:%v\n\twant:%v", tc.req.Name, n, tc.volume)
		}

		if status.Code(err) != tc.code {
			t.Errorf("Wrong status code on %s.\n\thave:%v\n\twant:%v", tc.req.Name, status.Code(err), tc.code)
		}
	}
}
