package node

import (
	"github.com/n0stack/n0core/model"
	uuid "github.com/satori/go.uuid"
)

const ComputeType = "node/compute"

type Compute struct {
	model.Model `yaml:",inline"`

	SupportingTypes []string `yaml:"supportingTypes"`
}

func (c Compute) ToModel() *model.Model {
	return &c.Model
}

func NewCompute(id uuid.UUID, state, name string, meta map[string]string, dependencies model.Dependencies, supportingTypes []string) *Compute {
	return &Compute{
		Model: model.Model{
			ID:           id,
			Type:         ComputeType,
			State:        state,
			Name:         name,
			Meta:         meta,
			Dependencies: model.Dependencies{},
		},
		SupportingTypes: supportingTypes,
	}
}
