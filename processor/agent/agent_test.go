package agent

import (
	"reflect"
	"testing"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	"github.com/n0stack/n0core/target"
	uuid "github.com/satori/go.uuid"
)

// TestProcessNotification test nominal scenarios.
func TestProcessNotificationResource(t *testing.T) {
	t.SkipNow()

	id := uuid.NewV4()
	c := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

	specID := uuid.NewV4()
	mes := &message.Notification{
		SpecID:      specID,
		Model:       c,
		Event:       "APPLIED",
		IsSucceeded: false,
		Description: "foobar",
	}

	f := func(in model.AbstractModel) (string, bool) {
		if !reflect.DeepEqual(c.Model, mes) { // modelの比較はstructの入れ子が比較できていない可能性がある
			t.Errorf("Got another message on MockRepository.StoreNotification:\ngot  %v\nwant %v", c.Model, mes)
		}

		return "", true
	}

	tg := &target.MockTarget{}
	tg.PatchApply = f

	// a := &Agent{Targets: []target.Target{tg}, Notifier: }
	// a.ProcessMessage(mes)
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
