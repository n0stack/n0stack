package message

import (
	"fmt"

	"github.com/n0stack/n0core/model"
	uuid "github.com/satori/go.uuid"
)

type Task struct {
	TaskID      uuid.UUID `yaml:"taskID"`
	Task        string
	Models      []model.AbstractModel // sliceにする必要があるかわからない
	Annotations map[string]string     // 必要かどうかわからない
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

	mi, ok := m["models"]
	if !ok {
		return nil
	}

	mms, ok := mi.([]interface{})
	if !ok {
		return fmt.Errorf("Failed to parse model")
	}

	t.Models = make([]model.AbstractModel, len(mms))
	for i, mm := range mms {
		var err error
		t.Models[i], err = model.MapToAbstractModel(mm.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
	}

	return nil
}
