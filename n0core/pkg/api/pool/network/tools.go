package network

import (
	"github.com/n0stack/n0stack/n0proto.go/pool/v0"
)

func IsLockedForDeletion(network *ppool.Network) bool {
	for _, ni := range network.ReservedNetworkInterfaces {
		if ni.Annotations != nil {
			if _, ok := ni.Annotations[AnnotationNetworkInterfaceDisableDeletionLock]; ok {
				continue
			}
		}

		return true
	}

	return false
}
