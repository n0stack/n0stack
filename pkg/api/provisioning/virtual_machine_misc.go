package provisioning

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/n0stack/n0core/pkg/driver/qemu"
	"github.com/n0stack/proto.go/provisioning/v0"
)

func GetAgentStateFromQemuState(s qemu.Status) VirtualMachineAgentState {
	switch s {
	case qemu.StatusRunning:
		return VirtualMachineAgentState_RUNNING

	case qemu.StatusShutdown, qemu.StatusGuestPanicked, qemu.StatusPreLaunch:
		return VirtualMachineAgentState_SHUTDOWN

	case qemu.StatusPaused, qemu.StatusSuspended:
		return VirtualMachineAgentState_PAUSED

	case qemu.StatusInternalError, qemu.StatusIOError:
		return VirtualMachineAgentState_FAILED

	case qemu.StatusInMigrate:
	case qemu.StatusFinishMigrate:
	case qemu.StatusPostMigrate:
	case qemu.StatusRestoreVM:
	case qemu.StatusSaveVM: // TODO: 多分PAUSED
	case qemu.StatusWatchdog:
	case qemu.StatusDebug:
	}

	return VirtualMachineAgentState_UNKNOWN
}

func TrimNetdevName(name string) string {
	return name[:16]
}

func StructNetDev(nics []*pprovisioning.VirtualMachineSpec_NIC, names []string) []*NetDev {
	nd := make([]*NetDev, 0, len(nics))
	for i, n := range nics {
		nd = append(nd, &NetDev{
			Name:            names[i],
			NetworkName:     n.NetworkName,
			HardwareAddress: n.HardwareAddress,
		})
	}

	return nd
}

func StructBlockDev(names []string, volumeAPI pprovisioning.VolumeServiceClient) ([]*BlockDev, error) {
	bd := make([]*BlockDev, 0, len(names))
	for i, n := range names {
		v, err := volumeAPI.GetVolume(context.Background(), &pprovisioning.GetVolumeRequest{Name: n})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return nil, grpc.Errorf(codes.InvalidArgument, "Volume '%s' is not found", n)
			}

			return nil, grpc.Errorf(codes.Internal, "Failed to get volume '%s' from API: %s", n, err.Error())
		}

		bd = append(bd, &BlockDev{
			Name:      names[i],
			Url:       v.Metadata.Annotations[AnnotationVolumePath],
			BootIndex: uint32(i),
		})
	}

	return bd, nil
}
