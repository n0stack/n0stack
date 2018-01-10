package resource

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestVolumeGetModel(t *testing.T) {
	id := uuid.NewV4()
	m := &model.Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	u, _ := url.Parse("file:///opt/n0core")
	v := &Volume{
		Model: *m,
		Size:  100000000000000,
		URL:   u,
	}

	f := v.GetModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on GetModel:\ngot  %v\nwant %v", f, m)
	}
}
