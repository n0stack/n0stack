package node

import (
	"fmt"

	"github.com/n0stack/n0stack/n0proto.go/budget/v0"
)

// TODO: 多分パッケージを切り出す

func CheckCompute(requestCpus, totalCpus uint32, requestMemory, totalMemory uint64, reserved map[string]*pbudget.Compute) error {
	usedCPU := uint32(0)
	usedMemory := uint64(0)

	for _, c := range reserved {
		usedCPU += c.RequestCpuMilliCore
		usedMemory += c.RequestMemoryBytes
	}

	if totalCpus < usedCPU+requestCpus {
		return fmt.Errorf("parameter='cpu_milli_core', total='%d', used='%d', requested='%d'", totalCpus, usedCPU, requestCpus)
	}
	if totalMemory < usedMemory+requestMemory {
		return fmt.Errorf("parameter='memory_bytes', total='%d', used='%d', requested='%d'", totalMemory, usedMemory, requestMemory)
	}

	return nil
}

func CheckStorage(request, total uint64, reserved map[string]*pbudget.Storage) error {
	usedStorage := uint64(0)

	for _, c := range reserved {
		usedStorage += c.RequestBytes
	}

	if total < usedStorage+request {
		return fmt.Errorf("total='%d', used='%d', requested='%d'", total, usedStorage, request)
	}

	return nil
}
