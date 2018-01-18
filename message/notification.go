package message

import (
	"fmt"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

type Notification struct {
	TaskID      uuid.UUID `yaml:"taskID"`
	Task        string
	Operation   string
	IsSucceeded bool `yaml:"isSucceeded" json:"is_succeeded"`
	Description string
	Model       model.AbstractModel
}

func (n *Notification) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]interface{})
	unmarshal(&m)

	n.TaskID = uuid.FromStringOrNil(m["taskID"].(string))
	n.Task = m["task"].(string)
	n.Operation = m["operation"].(string)
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
