package model

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func MapToAbstractModel(m map[interface{}]interface{}) (AbstractModel, error) {
	y, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}

	t, ok := m["type"].(string)
	if !ok {
		return nil, fmt.Errorf("Failed to parse type on model")
	}

	if strings.HasPrefix(t, ComputeType) {
		v := Compute{}

		err = yaml.Unmarshal(y, &v)
		if err != nil {
			return nil, err
		}

		return &v, nil
	} else if strings.HasPrefix(t, NetworkType) {
		v := Network{}

		err = yaml.Unmarshal(y, &v)
		if err != nil {
			return nil, err
		}

		return &v, nil
	} else if strings.HasPrefix(t, NICType) {
		v := NIC{}

		err = yaml.Unmarshal(y, &v)
		if err != nil {
			return nil, err
		}

		return &v, nil
	} else if strings.HasPrefix(t, VolumeType) {
		v := Volume{}

		err = yaml.Unmarshal(y, &v)
		if err != nil {
			return nil, err
		}

		return &v, nil
	} else if strings.HasPrefix(t, VMType) {
		v := VM{}

		err = yaml.Unmarshal(y, &v)
		if err != nil {
			return nil, err
		}

		return &v, nil
	}

	return nil, fmt.Errorf("Unsupported model type is setted")
}
