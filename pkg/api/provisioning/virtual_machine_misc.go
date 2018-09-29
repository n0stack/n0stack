package provisioning

import (
	"context"
	"log"
	"net/url"
	"strconv"

	"github.com/n0stack/proto.go/budget/v0"
	"github.com/n0stack/proto.go/pool/v0"
	"github.com/n0stack/proto.go/provisioning/v0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/n0stack/n0core/pkg/driver/qemu"
)

func GetAPIStateFromAgentState(s VirtualMachineAgentState) pprovisioning.VirtualMachineStatus_VirtualMachineState {
	switch s {
	case VirtualMachineAgentState_SHUTDOWN:
		return pprovisioning.VirtualMachineStatus_SHUTDOWN

	case VirtualMachineAgentState_RUNNING:
		return pprovisioning.VirtualMachineStatus_RUNNING

	case VirtualMachineAgentState_PAUSED:
		return pprovisioning.VirtualMachineStatus_PAUSED
	}

	return pprovisioning.VirtualMachineStatus_SHUTDOWN
}

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
	if len(name) <= 16 {
		return name
	}

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
		return "", "", grpc.Errorf(codes.Internal, "") // TODO: #89
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
			log.Printf("Failed to relserve network interface '%s' from API: %s", name+strconv.Itoa(i), err.Error())
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

func (a VirtualMachineAPI) reserveBlockStorage(names []string) ([]*BlockDev, error) {
	bd := make([]*BlockDev, 0, len(names))
	for i, n := range names {
		v, err := a.blockstorageAPI.SetInuseBlockStorage(context.Background(), &pprovisioning.SetInuseBlockStorageRequest{Name: n})
		if err != nil {
			log.Printf("Failed to get block storage '%s' from API: %s", n, err.Error())
			if status.Code(err) != codes.NotFound {
				return nil, grpc.Errorf(codes.Internal, "Failed to set block storage '%s' as in use from API", n)
			}

			return nil, grpc.Errorf(codes.InvalidArgument, "BlockStorage '%s' is not found", n)
		}

		u := url.URL{
			Scheme: "file",
			Path:   v.Metadata.Annotations[AnnotationBlockStoragePath],
		}
		bd = append(bd, &BlockDev{
			Name:      names[i],
			Url:       u.String(),
			BootIndex: uint32(i),
		})
	}

	return bd, nil
}

func (a VirtualMachineAPI) relaseBlockStorages(names []string) error {
	for _, n := range names {
		_, err := a.blockstorageAPI.SetAvailableBlockStorage(context.Background(), &pprovisioning.SetAvailableBlockStorageRequest{Name: n})
		if err != nil {
			log.Printf("Failed to get block storage '%s' from API: %s", n, err.Error())

			if status.Code(err) != codes.NotFound {
				return grpc.Errorf(codes.Internal, "Failed to get block storage '%s' as in use from API", n)
			}
		}
	}

	return nil
}
