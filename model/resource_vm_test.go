package model

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
	yaml "gopkg.in/yaml.v2"
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
		State:        "RUNNING",
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

	nv, err := NewVM(v.ID.String(), specificType, v.State, v.Name, v.Meta, v.Dependencies, v.Arch, v.VCPUs, v.Memory, v.VNCPassword)
	if err != nil {
		t.Errorf("Failed to create vm instance: error message %v", err.Error())
	}

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}

func TestNewVMFailOnParseID(t *testing.T) {
	i := "hogehoge"
	c := fmt.Sprintf("Failed to parse uuid of id:\ngot %v", i)

	_, err := NewVM(i, "test", "RUNNING", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, "x86/64", 1, 128*1024*1024*1024, "foobar")
	if !(err != nil && err.Error() == c) {
		t.Errorf("Failed to issue error on parse id:\ngot  error message %v\nwant error message %v", err.Error(), c)
	}
}

func TestYamlVM(t *testing.T) {
	v, err := NewVM(uuid.NewV4().String(), "test", "RUNNING", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, "x86/64", 1, 128*1024*1024*1024, "foobar")
	if err != nil {
		t.Errorf("Failed to create nic instance: error message %v", err.Error())
	}

	y, err := yaml.Marshal(v)
	if err != nil {
		t.Errorf("Failed to marshal nic")
	}

	t.Logf("%v", string(y))

	m, err := ParseYAMLModel(y, v.Type)
	if !reflect.DeepEqual(m, v) {
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", m, v) // deep equal is not watching subnets
	}
}
