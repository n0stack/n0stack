package datastore

import (
	"fmt"
)

// CheckVersion check whether version is valid or not
// return next version and error if version is invalid
//
// TODO: protoファイルに依存するのは綺麗ではないため、変えたい
func CheckVersion(previous, new uint64) (uint64, error) {
	if new != previous {
		return 0, fmt.Errorf("Set the same version as stored in database: have=%d, want=%d", new, previous)
	}

	return new + 1, nil
}
