package model

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
)

func TestVMToModel(t *testing.T) {
	id := uuid.NewV4()
	m := &Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	v := &VM{
		Model:       *m,
		Arch:        "x86/64",
		VCPUs:       1,
		Memory:      128 * 1024 * 1024 * 1024,
		VNCPassword: "foobar",
	}

	f := v.ToModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", f, m)
	}
}

func TestNewVM(t *testing.T) {
	id := uuid.NewV4()
	specificType := "hoge"
	m := &Model{
		ID:           id,
		Type:         filepath.Join(VMType, specificType),
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	v := &VM{
		Model:       *m,
		Arch:        "x86/64",
		VCPUs:       1,
		Memory:      128 * 1024 * 1024 * 1024,
		VNCPassword: "foobar",
	}

	nv := NewVM(v.ID, specificType, v.State, v.Name, v.Meta, v.Dependencies, v.Arch, v.VCPUs, v.Memory, v.VNCPassword)

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}
