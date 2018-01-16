package model

import (
	"fmt"

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

func NewCompute(id, state, name string, meta map[string]string, dependencies Dependencies, supportingTypes []string) (*Compute, error) {
	i, err := uuid.FromString(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse uuid of id:\ngot %v", id)
	}

	return &Compute{
		Model: Model{
			ID:           i,
			Type:         ComputeType,
			State:        state,
			Name:         name,
			Meta:         meta,
			Dependencies: Dependencies{},
		},
		SupportingTypes: supportingTypes,
	}, nil
}
