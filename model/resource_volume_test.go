package model

import (
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
	yaml "gopkg.in/yaml.v2"
)

func TestVolumeToModel(t *testing.T) {
	id := uuid.NewV4()
	m := &Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	u, _ := url.Parse("file:///opt/n0core")
	v := &Volume{
		Model: *m,
		Size:  100000000000000,
		URL:   u,
	}

	f := v.ToModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on ToModel:\ngot  %v\nwant %v", f, m)
	}
}

func TestNewVolume(t *testing.T) {
	id := uuid.NewV4()
	specificType := "hoge"
	m := &Model{
		ID:           id,
		Type:         filepath.Join(VolumeType, specificType),
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: Dependencies{},
	}

	u, _ := url.Parse("file:///opt/n0core")
	v := &Volume{
		Model: *m,
		Size:  100000000000000,
		URL:   u,
	}

	nv, err := NewVolume(v.ID.String(), specificType, v.State, v.Name, v.Meta, v.Dependencies, v.Size, v.URL.String())
	if err != nil {
		t.Errorf("Failed to create volume instance: error message %v", err.Error())
	}

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}

func TestNewVolumeFailOnParseID(t *testing.T) {
	i := "hogehoge"
	c := fmt.Sprintf("Failed to parse uuid of id:\ngot %v", i)

	_, err := NewVolume(i, "test", "ALLOCATED", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, 1*1024*1024*1024, "file:///opt/n0core")
	if !(err != nil && err.Error() == c) {
		t.Errorf("Failed to issue error on parse id:\ngot  error message %v\nwant error message %v", err.Error(), c)
	}
}

func TestNewVolumeFailOnParseURL(t *testing.T) {
	u := "::hogehoge"
	c := fmt.Sprintf("Failed to parse url of path:\ngot %v", u)

	_, err := NewVolume(uuid.NewV4().String(), "test", "ALLOCATED", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, 1*1024*1024*1024, u)
	if !(err != nil && err.Error() == c) {
		t.Errorf("Failed to issue error on parse url:\ngot  error message %v\nwant error message %v", err.Error(), c)
	}
}

func TestYamlVolume(t *testing.T) {
	v, err := NewVolume(uuid.NewV4().String(), "test", "ALLOCATED", "test_model", map[string]string{"hoge": "hoge"}, Dependencies{}, 1*1024*1024*1024, "file:///opt/n0core")
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
