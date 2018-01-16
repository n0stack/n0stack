package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
	yaml "gopkg.in/yaml.v2"
)

func TestComputeToModel(t *testing.T) {
	id := uuid.NewV4()
	m := &Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
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

func TestNewCompute(t *testing.T) {
	id := uuid.NewV4()
	m := &Model{
		ID:           id,
		Type:         ComputeType,
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	v := &Compute{
		Model:           *m,
		SupportingTypes: []string{"test/test"},
	}

	nv, err := NewCompute(v.ID.String(), v.State, v.Name, v.Meta, v.Dependencies, v.SupportingTypes)
	if err != nil {
		t.Errorf("Failed to create compute instance: error message %v", err.Error())
	}

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewCompute:\ngot  %v\nwant %v", v, nv)
	}
}

func TestNewComputeFailOnParseID(t *testing.T) {
	i := "hogehoge"
	c := fmt.Sprintf("Failed to parse uuid of id:\ngot %v", i)

	_, err := NewCompute(i, "UP", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, []string{"test/test"})
	if !(err != nil && err.Error() == c) {
		t.Errorf("Failed to issue error on parse id:\ngot  error message %v\nwant error message %v", err.Error(), c)
	}
}

func TestYamlCompute(t *testing.T) {
	v, err := NewCompute(uuid.NewV4().String(), "UP", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, []string{"test/test"})
	if err != nil {
		t.Errorf("Failed to create compute instance: error message %v", err.Error())
	}

	y, err := yaml.Marshal(v)
	if err != nil {
		t.Errorf("Failed to marshal nic")
	}

	m, err := ParseYAMLModel(y, v.Type)
	if !reflect.DeepEqual(m, v) {
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", m, v) // deep equal is not watching subnets
	}
}
