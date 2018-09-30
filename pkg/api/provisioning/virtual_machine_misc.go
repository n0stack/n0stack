package provisioning

import (
	"github.com/n0stack/proto.go/provisioning/v0"

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
