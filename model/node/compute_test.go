package node

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestComputeGetModel(t *testing.T) {
	id := uuid.NewV4()
	m := &model.Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	c := &Compute{
		Model:           *m,
		SupportingTypes: []string{"test/test"},
	}

	f := c.GetModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on GetModel:\ngot  %v\nwant %v", f, m)
	}
}
