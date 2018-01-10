package resource

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestNetworkGetModel(t *testing.T) {
	id := uuid.NewV4()
	m := &model.Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	n := &Network{
		Model:   *m,
		Bridge:  "nbr-test",
		Subnets: []subnet{},
	}

	f := n.GetModel()

	if !reflect.DeepEqual(f, m) {
		t.Errorf("Got another model on GetModel:\ngot  %v\nwant %v", f, m)
	}
}
