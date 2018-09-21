package datastore

import (
	"fmt"
	"reflect"

	pn0stack "github.com/n0stack/proto.go/v0"
)

type HavingMetadata interface {
	GetMetadata() *pn0stack.Metadata
}

// CheckVersion check whether version is valid or not
// return next version and error if version is invalid
//
// TODO: protoファイルに依存するのは綺麗ではないため、変えたい
func CheckVersion(previous, new HavingMetadata) (uint64, error) {
	if reflect.ValueOf(previous).IsNil() || reflect.ValueOf(previous.GetMetadata()).IsNil() {
		if new.GetMetadata().Version != 0 {
			return 0, fmt.Errorf("Set 0 when create new object: have=%d, want=0", new.GetMetadata().Version)
		}

		return 1, nil
	}

	if new.GetMetadata().Version != previous.GetMetadata().Version {
		return 0, fmt.Errorf("Set the same version as stored in database: have=%d, want=%d", new.GetMetadata().Version, previous.GetMetadata().Version)
	}

	return new.GetMetadata().Version + 1, nil
}
