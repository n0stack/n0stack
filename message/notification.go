package message

import (
	"fmt"
	"time"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

type Notification struct {
	TaskID      uuid.UUID `yaml:"taskID" json:"task_id"`
	Task        string
	NotifiedAt  time.Time `yaml:"notifiedAt" json:"notified_at"`
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

	t := m["notifiedAt"].(string)
	var err error
	n.NotifiedAt, err = time.Parse(time.RFC3339, t)
	if err != nil {
		return err
	}

	mi, ok := m["model"]
	if !ok {
		return nil
	}

	mm, ok := mi.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("Failed to parse model")
	}

	n.Model, err = model.MapToAbstractModel(mm)
	if err != nil {
		return err
	}

	return nil
}
