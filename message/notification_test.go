package message

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	uuid "github.com/satori/go.uuid"
	yaml "gopkg.in/yaml.v2"
)

func TestNotificationUnmarshalYAML(t *testing.T) {
	id, _ := uuid.FromString("1578ce2b-b845-41b2-9c73-7e05009785e6")
	c := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

	specID, _ := uuid.FromString("2efbfd8d-6136-4390-a513-033e7c5f2391")
	mes := &Notification{
		SpecID:      specID,
		Model:       c,
		Event:       "APPLIED",
		IsSucceeded: true,
		Description: "foobar",
	}

	y, err := yaml.Marshal(mes)
	if err != nil {
		t.Errorf("Failed to marshal message to yaml: error message %v", err.Error())
	}

	r := []byte(`specID: 2efbfd8d-6136-4390-a513-033e7c5f2391
model:
  id: 1578ce2b-b845-41b2-9c73-7e05009785e6
  type: node/compute
  state: testing
  name: test_model
  meta:
    hoge: hoge
  dependencies: []
  supportingTypes:
  - test/test
event: APPLIED
isSucceeded: true
description: foobar
`)
	if !reflect.DeepEqual(r, y) {
		t.Errorf("Failed to marshal to yaml:\ngot\n%v\nwant\n%v", string(y), string(r))
	}

	n := Notification{}
	err = yaml.Unmarshal(y, &n)
	if err != nil {
		t.Errorf("Failed to unmarshal message to yaml: error message %v", err.Error())
	}
	if !reflect.DeepEqual(*mes, n) {
		t.Errorf("Failed to unmarshal to yaml:\ngot  %v,\nwant %v", n, *mes)
	}
}
