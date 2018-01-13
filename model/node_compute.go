package model

import (
	uuid "github.com/satori/go.uuid"
)

const ComputeType = "node/compute"

type Compute struct {
	Model `yaml:",inline"`

	SupportingTypes []string `yaml:"supportingTypes"`
}

func (c Compute) ToModel() *Model {
	return &c.Model
}

func NewCompute(id uuid.UUID, state, name string, meta map[string]string, dependencies Dependencies, supportingTypes []string) *Compute {
	return &Compute{
		Model: Model{
			ID:           id,
			Type:         ComputeType,
			State:        state,
			Name:         name,
			Meta:         meta,
			Dependencies: Dependencies{},
		},
		SupportingTypes: supportingTypes,
	}
}
