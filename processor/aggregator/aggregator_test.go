package aggregator

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	"github.com/n0stack/n0core/repository"
	"github.com/satori/go.uuid"
)

func TestProcessNotification(t *testing.T) {
	id := uuid.NewV4().String()
	c, _ := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

	specID := uuid.NewV4()
	mes := &message.Notification{
		SpecID:      specID,
		Model:       c,
		Event:       "APPLIED",
		IsSucceeded: true,
		Description: "foobar",
	}

	f := func(mod *message.Notification) bool {
		if !reflect.DeepEqual(mod, mes) {
			t.Errorf("Got another message on MockRepository.StoreNotification:\ngot  %v\nwant %v", mod, mes)
		}

		return true
	}

	r := &repository.MockRepository{}
	r.PatchStoreNotification = f

	a := &Aggregator{Repository: r}
	err := a.ProcessMessage(mes)
	if err != nil {
		t.Errorf("Failed to process message: %v", err.Error())
	}
}

func TestProcessSpec(t *testing.T) {
	id := uuid.NewV4().String()
	c, _ := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

	specID := uuid.NewV4()
	mes := &message.Spec{
		SpecID: specID,
		Models: []model.AbstractModel{c},
	}

	r := &repository.MockRepository{}
	a := &Aggregator{Repository: r}

	err := a.ProcessMessage(mes)
	if !(err != nil && err.Error() == "Received notification message which is not supported, maybe there are stranger or distributor has bugs") {
		t.Errorf("Could not specify notification message when got spec message:\nwant error message 'Received notification message which is not supported, maybe there are stranger or distributor has bugs'\ngot  error message '%v'", err.Error())
	}
}

func TestProcessMessageOnRepositoryFailure(t *testing.T) {
	id := uuid.NewV4().String()
	c, _ := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

	specID := uuid.NewV4()
	mes := &message.Notification{
		SpecID:      specID,
		Model:       c,
		Event:       "APPLIED",
		IsSucceeded: true,
		Description: "foobar",
	}

	f := func(mod *message.Notification) bool {
		return false
	}

	r := &repository.MockRepository{}
	r.PatchStoreNotification = f

	a := &Aggregator{Repository: r}
	err := a.ProcessMessage(mes)
	if !(err != nil && err.Error() == "Failed to store notification message") {
		t.Errorf("Could not handling failure to store notification message:\nwant error message 'Failed to store notification message'\ngot  error message '%v'", err.Error())
	}
}
