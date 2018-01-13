package model

import (
	"net/url"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/satori/go.uuid"
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

	nv := NewVolume(v.ID, specificType, v.State, v.Name, v.Meta, v.Dependencies, v.Size, v.URL)

	if !reflect.DeepEqual(v, nv) {
		t.Errorf("Got another model on NewVM:\ngot  %v\nwant %v", v, nv)
	}
}
