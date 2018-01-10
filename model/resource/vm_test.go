package resource

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestVMGetModel(t *testing.T) {
	id, _ := uuid.NewV4()
	m := &model.Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	v := &VM{
		Model:       *m,
		Arch:        "x86/64",
		VCPUs:       1,
		Memory:      128 * 1024 * 1024 * 1024,
		VNCPassword: "foobar",
	}

	f := v.GetModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on GetModel:\ngot  %v\nwant %v", f, m)
	}
}
