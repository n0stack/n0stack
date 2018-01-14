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

func (d Direct) Apply(m model.AbstractModel) (string, bool) {
	return "", true
}
