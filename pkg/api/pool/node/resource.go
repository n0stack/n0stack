package node

import (
	"fmt"

	"github.com/n0stack/proto.go/resource/v0"
)

// TODO: 多分パッケージを切り出す

func CheckCompute(requestVcpus uint32, requestMemory uint64, total *presource.Compute, reserved map[string]*presource.Compute) error {
	usedCPU := uint32(0)
	usedMemory := uint64(0)

	for _, c := range reserved {
		usedCPU += c.Vcpus
		usedMemory += c.MemoryBytes
	}

	if total.Vcpus < usedCPU+requestVcpus {
		return fmt.Errorf("parameter='vcpus', total='%d', used='%d', requested='%d'", total.Vcpus, usedCPU, requestVcpus)
	}
	if total.MemoryBytes < usedMemory+requestMemory {
		return fmt.Errorf("parameter='memory', total='%d', used='%d', requested='%d'", total.MemoryBytes, usedMemory, requestMemory)
	}

	return nil
}

func CheckStorage(request uint64, total *presource.Storage, reserved map[string]*presource.Storage) error {
	usedStorage := uint64(0)

	for _, c := range reserved {
		usedStorage += c.Bytes
	}

	if total.Bytes < usedStorage+request {
		return fmt.Errorf("total='%d', used='%d', requested='%d'", total.Bytes, usedStorage, request)
	}

	return nil
}
