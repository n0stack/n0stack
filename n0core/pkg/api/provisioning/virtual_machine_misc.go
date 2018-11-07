package provisioning

import (
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"

	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu"
)

func GetAPIStateFromAgentState(s VirtualMachineAgentState) pprovisioning.VirtualMachine_VirtualMachineState {
	switch s {
	case VirtualMachineAgentState_SHUTDOWN:
		return pprovisioning.VirtualMachine_SHUTDOWN

	case VirtualMachineAgentState_RUNNING:
		return pprovisioning.VirtualMachine_RUNNING

	case VirtualMachineAgentState_PAUSED:
		return pprovisioning.VirtualMachine_PAUSED
	}

	return pprovisioning.VirtualMachine_UNKNOWN
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

// TrimNetdevName trim network device name because Linux network device can use 15 characters.
func TrimNetdevName(name string) string {
	if len(name) <= 15 {
		return name
	}

	return name[:15]
}
