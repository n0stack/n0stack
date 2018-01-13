package agent

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	"github.com/n0stack/n0core/model/node"
	"github.com/n0stack/n0core/target"
	uuid "github.com/satori/go.uuid"
)

// TestProcessNotification test nominal scenarios.
func TestProcessNotificationResource(t *testing.T) {
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
		IsSucceeded: false,
		Description: "foobar",
	}

	f := func(in model.AbstractModel) (string, bool) {
		if !reflect.DeepEqual(m, mes) { // modelの比較はstructの入れ子が比較できていない可能性がある
			t.Errorf("Got another message on MockRepository.StoreNotification:\ngot  %v\nwant %v", m, mes)
		}

		return "", true
	}

	tg := &target.MockTarget{}
	tg.PatchApply = f

	a := &Aggregator{Repository: r}
	a.ProcessMessage(mes)
}

// TestMultipleTargets test to separate multiple targets.
// success := func(m *message.Notification) {
// 	if !reflect.DeepEqual(m, mes) {
// 		t.Errorf("Got another message on MockRepository.StoreNotification:\ngot  %v\nwant %v", m, mes)
// 	}
// }

// fail := func(m *message.Notification) {
// 	if !reflect.DeepEqual(m, mes) {
// 		t.Errorf("Got another message on MockRepository.StoreNotification:\ngot  %v\nwant %v", m, mes)
// 	}
// }

// func TestAgentGetComputeUUID(t *testing.T) {
// 	a := &Agent{}
// 	i, _ := a.getComputeUUID()

// 	t.Errorf(i.String())
// }
