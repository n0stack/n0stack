package agent

import (
	"bufio"
	"fmt"
	"os"

	"github.com/n0stack/n0core/gateway"
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	"github.com/n0stack/n0core/target"
	"github.com/satori/go.uuid"
)

// Agent is a processor which apply resources with targets.
//
// 1. Receive a message from gateway.
// 2. Apply resource with target.
// 3. Send a result message to gateway.
//
// Args:
// 	model_types:
// 	notification:
//
// Exapmle:
// 	>>> agent = Agent(notification, [flat_network_target])
type Agent struct {
	Targets  []target.Target
	Notifier gateway.Gateway
}

func NewAgent(t []target.Target, g gateway.Gateway, m map[string]string) (*Agent, error) {
	a := new(Agent)
	a.Targets = t
	a.Notifier = g

	for _, v := range t {
		d, ok := v.Initialize()
		if !ok {
			return nil, fmt.Errorf(d) // エラーについて考える
		}
	}

	id, err := a.getComputeUUID()
	if err != nil {
		return nil, err // エラーについて考える
	}
	// hostName, err := a.getHostName()
	c := model.NewCompute(id, "JOINED", "test_model", map[string]string{"hoge": "hoge"}, model.Dependencies{}, []string{"test/test"})

	n := &message.Notification{
		// SpecID: uuid.NewV4(),
		Model:       c,
		Event:       "APPLIED",
		IsSucceeded: true,
		Description: "Joined",
	}

	ok := a.Notifier.SendMessage(n)
	if !ok {
		return nil, fmt.Errorf("Failed to send spec message to initialize and join agent")
	}

	return a, nil
}

func (a Agent) ProcessMessage(am message.AbstractMessage) error {
	n, ok := am.(message.Notification)
	if !ok {
		return fmt.Errorf("Received notification message which is not supported, maybe there are stranger or distributor has bugs")
	}

	if n.IsSucceeded {
		return fmt.Errorf("Received notification message which is not succeeded, maybe there are stranger or distributor has bugs")
	}

	// c, ok := n.Model.(node.Compute)
	// if ok {
	// 	// joinの処理
	// }

	m := n.Model.ToModel()
	t, ok := a.isSupportModelType(m.Type)
	if !ok {
		return fmt.Errorf("Received model which is not supported, maybe there are stranger or distributor has bugs")
	}

	// TODO: check whether model is scheduled or not, and scheduling

	d, ok := t.Apply(n.Model)
	newN := &message.Notification{
		SpecID:      n.SpecID,
		Model:       n.Model,
		Event:       "APPLIED",
		IsSucceeded: ok,
		Description: d,
	}

	if !a.Notifier.SendMessage(newN) {
		return fmt.Errorf("Failed to send notification message")
	}

	return nil
}

// CollectMetrixes collect metrix to manage realized resources.
// func (a Agent) CollectMetrix() {}

// SupportingTypes return model types which is supported by agent.
func (a Agent) SupportingTypes() []string {
	t := make([]string, len(a.Targets))

	for i, v := range a.Targets {
		t[i] = v.ManagingType()
	}

	return t
}

func (a Agent) isSupportModelType(mt string) (target.Target, bool) {
	for _, t := range a.Targets {
		if t.ManagingType() == mt {
			return t, true
		}
	}

	return nil, false
}

func (a Agent) getComputeUUID() (uuid.UUID, error) {
	f, err := os.Open(`/sys/class/dmi/id/product_uuid`)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Failed to open /sys/class/dmi/id/product_uuid to read compute UUID")
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		i, err := uuid.FromString(s.Text())
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("Failed to parse compute UUID on /sys/class/dmi/id/product_uuid")
		}

		return i, nil
	}

	return uuid.UUID{}, fmt.Errorf("Failed to read compute UUID by /sys/class/dmi/id/product_uuid")
}
