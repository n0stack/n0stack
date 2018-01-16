package model

import (
	"fmt"
	"net"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
	yaml "gopkg.in/yaml.v2"
)

func TestNICToModel(t *testing.T) {
	id := uuid.NewV4()
	m := &Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	h, _ := net.ParseMAC("01:23:45:67:89:ab")

	n := &NIC{
		Model:   *m,
		HWAddr:  h,
		IPAddrs: []net.IP{net.ParseIP("192.168.0.1")},
	}

	f := n.ToModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", f, m)
	}
}

func TestNewNIC(t *testing.T) {
	specificType := "test"
	m := &Model{
		ID:           uuid.NewV4(),
		Type:         filepath.Join(NICType, specificType),
		State:        "UP",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	h, _ := net.ParseMAC("01:23:45:67:89:ab")

	v := &NIC{
		Model:   *m,
		HWAddr:  h,
		IPAddrs: []net.IP{net.ParseIP("192.168.0.1")},
	}

	nv, err := NewNIC(v.ID.String(), specificType, v.State, v.Name, v.Meta, v.Dependencies, v.HWAddr.String(), []string{v.IPAddrs[0].String()})
	if err != nil {
		t.Errorf("Failed to create nic instance: error message %v", err.Error())
	}

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}

func TestNewNICFailOnParseID(t *testing.T) {
	i := "hogehoge"
	c := fmt.Sprintf("Failed to parse uuid of id:\ngot %v", i)

	_, err := NewNIC(i, "test", "UP", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, "01:23:45:67:89:ab", []string{"192.168.0.1"})
	if !(err != nil && err.Error() == c) {
		t.Errorf("Failed to issue error on parse id:\ngot  error message %v\nwant error message %v", err.Error(), c)
	}
}

func TestNewNICFailOnParseMAC(t *testing.T) {
	m := "hogehoge"
	c := fmt.Sprintf("Failed to parse mac address of hwAddr:\ngot %v", m)

	_, err := NewNIC(uuid.NewV4().String(), "test", "UP", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, m, []string{"192.168.0.1"})
	if !(err != nil && err.Error() == c) {
		t.Errorf("Failed to issue error on parse mac:\ngot  error message %v\nwant error message %v", err.Error(), c)
	}
}

func TestNewNICFailOnParseIP(t *testing.T) {
	i := "hogehoge"
	c := fmt.Sprintf("Failed to parse IP address:\ngot %v", i)

	_, err := NewNIC(uuid.NewV4().String(), "test", "UP", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, "01:23:45:67:89:ab", []string{i})
	if !(err != nil && err.Error() == c) {
		t.Errorf("Failed to issue error on parse mac:\ngot  error message %v\nwant error message %v", err.Error(), c)
	}
}

func TestYamlNIC(t *testing.T) {
	v, err := NewNIC(uuid.NewV4().String(), "test", "UP", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, "01:23:45:67:89:ab", []string{"192.168.0.1"})
	if err != nil {
		t.Errorf("Failed to create nic instance: error message %v", err.Error())
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
