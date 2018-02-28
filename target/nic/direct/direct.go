package direct

import (
	"path/filepath"

	"github.com/n0stack/n0core/model"
)

const DirectType = "direct"

type Direct struct{}

func (d Direct) ManagingType() string {
	return filepath.Join(model.NICType, DirectType)
}

func (d *Direct) Operations(state, task string) ([]func(n model.AbstractModel) (string, bool, string), bool) {
	if !model.NetworkStateMachine[state][task] {
		return nil, false
	}

	return nil, false
}
