package node

import "github.com/n0stack/n0core/model"

type Compute struct {
	model.Model

	SupportingTypes []string
}

func (c Compute) GetModel() *model.Model {
	return &c.Model
}
