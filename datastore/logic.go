package datastore

import (
	"fmt"

	n0stack "github.com/n0stack/proto.go/v0" // 本当は依存させたくない
)

// CheckVersion check whether version is valid or not
// return next version and error if version is invalid
//
// TODO: protoファイルに依存するのは綺麗ではないため、変えたい
func CheckVersion(previousMetadata, newMetadata *n0stack.Metadata) (uint64, error) {
	if previousMetadata == nil && newMetadata.Version != 0 {
		return 0, fmt.Errorf("Set 0 when create new object, have:%d, want:0", newMetadata.Version)
	}
	if previousMetadata != nil && newMetadata.Version != previousMetadata.Version {
		return 0, fmt.Errorf("Set the same version as stored in database, have:%d, want:%d", newMetadata.Version, previousMetadata.Version)
	}

	return newMetadata.Version + 1, nil
}
