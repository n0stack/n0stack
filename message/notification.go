package message

import (
	"fmt"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

type Notification struct {
	SpecID      uuid.UUID `yaml:"specID"      json:"spec_id"`
	Model       model.AbstractModel
	Event       string // enum的なのにしたい
	IsSucceeded bool   `yaml:"isSucceeded" json:"is_succeeded"`
	Description string
}

// func (n *Notification) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	type Alias Notification
// 	a := &struct {
// 		Alias `yaml:",inline"`
// 		Model model.Model `yaml:"model"`
// 	}{}

// 	unmarshal(&a)

// 	*n = (Notification)(a.Alias)

// 	// var err error
// 	// n.Model, err = model.MapToAbstractModel(a.Model)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

func (n *Notification) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]interface{})
	unmarshal(&m)

	n.SpecID = uuid.FromStringOrNil(m["specID"].(string))
	n.Event = m["event"].(string)
	n.IsSucceeded = m["isSucceeded"].(bool)
	n.Description = m["description"].(string)

	mi, ok := m["model"]
	if !ok {
		return nil
	}

	mm, ok := mi.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("Failed to parse model")
	}

	var err error
	n.Model, err = model.MapToAbstractModel(mm)
	if err != nil {
		return err
	}

	return nil
}
