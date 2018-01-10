package processor

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model/node"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

type MockRepository struct {
	f func(m *message.Notification)
}

func (mr MockRepository) StoreNotification(m *message.Notification) {
	mr.f(m)
}

func (mr MockRepository) DigModel(i *uuid.UUID, e string, d uint) (*model.Model, error) {
	return &model.Model{}, nil
}

func TestAggregatorProcessMessage(t *testing.T) {
	id, _ := uuid.NewV4()
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

	specID, _ := uuid.NewV4()
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

	r := &MockRepository{f: f}

	a := &Aggregator{Repository: r}
	a.ProcessMessage(mes)
}
