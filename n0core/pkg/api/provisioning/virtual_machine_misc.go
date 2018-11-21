package provisioning

import (
	"context"
	"net"

	"github.com/n0stack/n0stack/n0core/pkg/driver/qemu"
	"github.com/n0stack/n0stack/n0core/pkg/util/net"
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
	"github.com/n0stack/n0stack/n0proto.go/provisioning/v0"
	"github.com/pkg/errors"
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

func (a *VirtualMachineAPI) addDefaultGateway(ctx context.Context, network *ppool.Network) (string, error) {
	_, ipn, err := net.ParseCIDR(network.Ipv4Cidr)
	if err != nil {
		return "", errors.Wrap(err, "Invalid CIDR in network")
	}

	ip := nettools.GetEndIP(ipn)

	a.networkAPI.ReserveNetworkInterface(ctx, &ppool.ReserveNetworkInterfaceRequest{
		NetworkName:          network.Name,
		NetworkInterfaceName: "default-gateway",
		Ipv4Address:          ip.String(),
		Annotations: map[string]string{
			AnnotationVirtualMachineVncWebSocketPort: "true",
		},
	})

	return ip.String(), nil
}
