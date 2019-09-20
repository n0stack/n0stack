package node

import (
	"n0st.ac/n0stack/n0proto.go/pool/v0"
)

func IsLockedForDeletion(node *ppool.Node) bool {
	for _, c := range node.ReservedComputes {
		if _, ok := c.Annotations[AnnotationComputeDisableDeletionLock]; ok {
			continue
		}

		return true
	}

	for _, s := range node.ReservedStorages {
		if _, ok := s.Annotations[AnnotationStorageDisableDeletionLock]; ok {
			continue
		}

		return true
	}

	return false
}
