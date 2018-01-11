package processor

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/repository"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model/node"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

func TestAggregatorProcessMessage(t *testing.T) {
	id := uuid.NewV4()
	m := &model.Model{
		ID:           id,
		Type:         "test/test",
		State:        "testing",
		Name:         "test_model",
		Meta:         map[string]string{"hoge": "hoge"},
		Dependencies: model.Dependencies{},
	}

	c := &node.Compute{
		Model:           *m,
		SupportingTypes: []string{"test/test"},
	}

	specID := uuid.NewV4()
	mes := &message.Notification{
		SpecID:      specID,
		Model:       c,
		Event:       "APPLIED",
		IsSucceeded: true,
		Description: "foobar",
	}

	f := func(m *message.Notification) {
		if !reflect.DeepEqual(m, mes) {
			t.Errorf("Got another message on MockRepository.StoreNotification:\ngot  %v\nwant %v", m, mes)
		}
	}

	r := &repository.MockRepository{}
	r.PatchStoreNotification = f

	a := &Aggregator{Repository: r}
	a.ProcessMessage(mes)
}
