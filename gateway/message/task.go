package message

import (
	"fmt"

	"github.com/n0stack/n0core/model"
	uuid "github.com/satori/go.uuid"
)

type Task struct {
	TaskID      uuid.UUID `yaml:"taskID"`
	Task        string
	Model       model.AbstractModel
	Annotations map[string]string // 必要かどうかわからない
}

func (t *Task) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]interface{})
	unmarshal(&m)

	var ok bool
	t.TaskID = uuid.FromStringOrNil(m["taskID"].(string))
	t.Task = m["task"].(string)
	t.Annotations, ok = m["annotations"].(map[string]string)
	if !ok {
		t.Annotations = map[string]string{}
	}

	mi, ok := m["model"]
	if !ok {
		return nil
	}

	mm, ok := mi.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("Failed to parse model")
	}

	var err error
	t.Model, err = model.MapToAbstractModel(mm)
	if err != nil {
		return err
	}

	return nil
}
