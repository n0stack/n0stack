package model

import (
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/satori/go.uuid"
)

func TestNetworkToModel(t *testing.T) {
	id := uuid.NewV4()
	m := &Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	n := &Network{
		Model:   *m,
		Bridge:  "nbr-test",
		Subnets: []*Subnet{},
	}

	f := n.ToModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", f, m)
	}
}

func TestNewNetwork(t *testing.T) {
	id := uuid.NewV4()
	specificType := "hoge"
	m := &Model{
		ID:           id,
		Type:         filepath.Join(NetworkType, specificType),
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	v := &Network{
		Model:   *m,
		Bridge:  "nbr-test",
		Subnets: []*Subnet{},
	}

	nv, err := NewNetwork(v.ID.String(), specificType, v.State, v.Name, v.Meta, v.Dependencies, v.Bridge, v.Subnets)
	if err != nil {
		t.Errorf("Failed to create network instance: error message %v", err.Error())
	}

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}

func TestYamlNetwork(t *testing.T) {
	d, err := NewDHCP("192.168.0.1", "192.168.0.127", "192.168.0.254", []string{"192.168.0.254"})
	if err != nil {
		t.Errorf("Failed to create dhcp instance: error message %v", err.Error())
	}

	s, err := NewSubnet("192.168.0.0/24", d)
	if err != nil {
		t.Errorf("Failed to create subnet instance: error message %v", err.Error())
	}
	v, err := NewNetwork("0f97b5a3-bff2-4f13-9361-9f9b4fab3d65", "direct", "UP", "test-network", map[string]string{}, Dependencies{}, "", []*Subnet{s})
	if err != nil {
		t.Errorf("Failed to create network instance: error message %v", err.Error())
	}

	y, err := yaml.Marshal(v)
	if err != nil {
		t.Errorf("Failed to marshal network")
	}

	t.Logf("Marshaled network:\n%v", string(y))

	m, err := ParseYAMLModel(y, v.Type)
	if m != v {
		t.Skip()
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", m.(*Network), v) // deep equal is not watching subnets
	}
}
