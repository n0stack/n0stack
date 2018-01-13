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

	// m.Meta = map[string]string{"hoge": "hogehoge"}
	f := c.ToModel()

	if !reflect.DeepEqual(f, m) { // これ本当に正しいか怪しい、というか原理的に壊れることはないと思うんだが :thinking_face:
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", f, m)
	}
}

func TestNewVM(t *testing.T) {
	id := uuid.NewV4()
	m := &model.Model{
		ID:           id,
		Type:         ComputeType,
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	v := &Compute{
		Model:           *m,
		SupportingTypes: []string{"test/test"},
	}

	nv := NewCompute(v.ID, v.State, v.Name, v.Meta, v.Dependencies, v.SupportingTypes)

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}
