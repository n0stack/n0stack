package agent

import (
	"testing"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	"github.com/n0stack/n0core/target"
	uuid "github.com/satori/go.uuid"
)

// // TestProcessNotification test nominal scenarios.
// func TestProcessTaskMessage(t *testing.T) {
// 	id := uuid.NewV4().String()
// 	c, _ := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

// 	taskID := uuid.NewV4()
// 	mes := &message.Task{
// 		TaskID: taskID,
// 		Task:   "Test",
// 		Models: []model.AbstractModel{c},
// 	}

// 	// id := uuid.NewV4().String()
// 	// c, _ := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

// 	// specID := uuid.NewV4()
// 	// mes := &message.Notification{
// 	// 	SpecID:      specID,
// 	// 	Model:       c,
// 	// 	Event:       "APPLIED",
// 	// 	IsSucceeded: false,
// 	// 	Description: "foobar",
// 	// }

// 	f := func(in model.AbstractModel) (string, bool) {
// 		if !reflect.DeepEqual(c.Model, mes) { // modelの比較はstructの入れ子が比較できていない可能性がある
// 			t.Errorf("Got another message on MockRepository.StoreNotification:\ngot  %v\nwant %v", c.Model, mes)
// 		}

// 		return "", true
// 	}

// 	tg := &target.MockTarget{}
// 	tg.PatchApply = f

// 	a := &Agent{Targets: []target.Target{tg}, Notifier: }
// 	a.ProcessMessage(mes)
// }

func TestProcessNotificationMessage(t *testing.T) {
	id := uuid.NewV4().String()
	c, _ := model.NewCompute(id, "testing", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

	taskID := uuid.NewV4()
	mes := &message.Notification{
		TaskID:      taskID,
		Task:        "APPLIED",
		Operation:   "test",
		IsSucceeded: false,
		Description: "foobar",
		Model:       c,
	}

	a := &Agent{}
	err := a.ProcessMessage(mes)
	if !(err != nil && err.Error() == "Received notification message which is not supported, maybe there are stranger or distributor has bugs") {
		t.Errorf("Failed to get a correct error message '%v'", "Received notification message which is not supported, maybe there are stranger or distributor has bugs")
	}
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

func TestSupportingTypes(t *testing.T) {
	testType := "test/hoge"
	f := func() string {
		return testType
	}

	tg := &target.MockTarget{}
	tg.PatchManagingType = f

	a := &Agent{Targets: []target.Target{tg}}
	ts := a.SupportingTypes()

	if ts[0] != testType {
		t.Errorf("Got another supporting type: got  %v, want %v", ts[0], testType)
	}
}

func TestIsSupportModelType(t *testing.T) {
	testType := "test/hoge"
	f := func() string {
		return testType
	}

	tg1 := &target.MockTarget{}
	tg1.PatchManagingType = f

	a := &Agent{Targets: []target.Target{tg1}}
	tg2, ok := a.isSupportModelType(testType)

	if !ok {
		t.Errorf("Failed to isSupportModelType on %v", testType)
	}

	if tg1 != tg2 {
		t.Errorf("Got another target: got  %v, want %v", tg2, tg1)
	}
}
