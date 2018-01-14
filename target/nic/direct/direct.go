package none

import (
	"github.com/n0stack/n0core/model"
)

type Direct struct{}

func (f Direct) Apply(m model.AbstractModel) (string, bool) {
	return "", true
}
