package node

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestComputeToModel(t *testing.T) {
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

	f := c.ToModel()

	if !reflect.DeepEqual(f, m) { // これ本当に正しいか怪しい
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", f, m)
	}
}
