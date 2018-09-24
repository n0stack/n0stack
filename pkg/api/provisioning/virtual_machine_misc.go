package provisioning

import (
	"context"
	"log"
	"strconv"

	"github.com/n0stack/proto.go/budget/v0"
	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/provisioning/v0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/n0stack/n0core/pkg/driver/qemu"
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

func (a VirtualMachineAPI) reserveCompute(name string, annotations map[string]string, reqCpu, limitCpu uint32, reqMem, limitMem uint64) (string, string, error) {
	var rcr *ppool.ReserveComputeResponse
	var err error
	if node, ok := annotations[AnnotationRequestNodeName]; !ok {
		rcr, err = a.nodeAPI.ScheduleCompute(context.Background(), &ppool.ScheduleComputeRequest{
			ComputeName: name,
			Compute: &pbudget.Compute{
				RequestCpuMilliCore: reqCpu,
				LimitCpuMilliCore:   limitCpu,
				RequestMemoryBytes:  reqMem,
				LimitMemoryBytes:    limitMem,
			},
		})
	} else {
		rcr, err = a.nodeAPI.ReserveCompute(context.Background(), &ppool.ReserveComputeRequest{
			Name:        node,
			ComputeName: name,
			Compute: &pbudget.Compute{
				RequestCpuMilliCore: reqCpu,
				LimitCpuMilliCore:   limitCpu,
				RequestMemoryBytes:  reqMem,
				LimitMemoryBytes:    limitMem,
			},
		})
	}
	if err != nil {
		return "", "", err // TODO: #89
	}

	return rcr.Name, rcr.ComputeName, nil
}

func (a VirtualMachineAPI) releaseCompute(node, compute string) error {
	_, err := a.nodeAPI.ReleaseCompute(context.Background(), &ppool.ReleaseComputeRequest{
		Name:        node,
		ComputeName: compute,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to release compute '%s': %s", compute, err.Error())

		// Notfound でもとりあえず問題ないため、処理を続ける
		if status.Code(err) != codes.NotFound {
			return grpc.Errorf(codes.Internal, "Failed to release compute '%s': please retry", compute)
		}
	}

	return nil
}

func (a VirtualMachineAPI) reserveNics(name string, nics []*pprovisioning.VirtualMachineSpec_NIC) ([]*pprovisioning.VirtualMachineSpec_NIC, []string, error) {
	// res.Status.NetworkInterfaceNames = make([]string, 0, len(req.Spec.Nics))
	networkInterfaceNames := make([]string, 0, len(nics))

	for i, n := range nics {
		ni, err := a.networkAPI.ReserveNetworkInterface(context.Background(), &ppool.ReserveNetworkInterfaceRequest{
			Name:                 n.NetworkName,
			NetworkInterfaceName: name + strconv.Itoa(i),
			NetworkInterface: &pbudget.NetworkInterface{
				HardwareAddress: nics[i].HardwareAddress,
				Ipv4Address:     nics[i].Ipv4Address,
				Ipv6Address:     nics[i].Ipv6Address,
			},
		})
		if err != nil {
			return nil, nil, err // TODO: #89
		}

		nics[i].HardwareAddress = ni.NetworkInterface.HardwareAddress
		nics[i].Ipv4Address = ni.NetworkInterface.Ipv4Address
		nics[i].Ipv6Address = ni.NetworkInterface.Ipv6Address
		networkInterfaceNames = append(networkInterfaceNames, ni.NetworkInterfaceName)
	}

	return nics, networkInterfaceNames, nil
}

func (a VirtualMachineAPI) releaseNics(nics []*pprovisioning.VirtualMachineSpec_NIC, networkInterfaces []string) error {
	for i, n := range nics {
		_, err := a.networkAPI.ReleaseNetworkInterface(context.Background(), &ppool.ReleaseNetworkInterfaceRequest{
			Name:                 n.NetworkName,
			NetworkInterfaceName: networkInterfaces[i],
		})
		if err != nil {
			log.Printf("[ERROR] Failed to release network interface '%s': %s", networkInterfaces[i], err.Error())

			// Notfound でもとりあえず問題ないため、処理を続ける
			if status.Code(err) != codes.NotFound {
				return grpc.Errorf(codes.Internal, "Failed to release network interface '%s': please check network interface on your own", networkInterfaces[i])
			}
		}
	}

	return nil
}

func (a VirtualMachineAPI) reserveVolume(names []string) ([]*BlockDev, error) {
	bd := make([]*BlockDev, 0, len(names))
	for i, n := range names {
		v, err := a.volumeAPI.SetInuseVolume(context.Background(), &pprovisioning.SetInuseVolumeRequest{Name: n})
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

func (a VirtualMachineAPI) relaseVolumes(names []string) error {
	for _, n := range names {
		_, err := a.volumeAPI.SetAvailableVolume(context.Background(), &pprovisioning.SetAvailableVolumeRequest{Name: n})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				return grpc.Errorf(codes.InvalidArgument, "Volume '%s' is not found", n)
			}

			return grpc.Errorf(codes.Internal, "Failed to get volume '%s' from API: %s", n, err.Error())
		}
	}

	return nil
}
