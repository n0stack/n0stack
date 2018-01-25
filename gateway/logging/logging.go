package file_gateway

import (
	"fmt"

	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	yaml "gopkg.in/yaml.v2"
)

type LoggingGateway struct{}

func (l LoggingGateway) SendNotification(m *message.Notification) bool {
	buf, _ := yaml.Marshal(m)
	fmt.Printf("Send message:\n%v\n", string(buf))

	return true
}

func (l LoggingGateway) SendNotificationToCompute(m *message.Notification, d *model.Compute) bool {
	fmt.Printf("Send message:\n%v\nTo:\n%v\n", m, d)
	println(m)
	println(d)

	return true
}

// // FileGateway support to parse yaml and json
// type FileGateway struct {
// 	Processor processor.Processor
// 	FilePath  string
// }

// type Gateway interface {
// 	// StartReceiveMessage start receive message from something.
// 	StartReceiveMessage()

// 	// SendNotification send spec message message to distributor
// 	SendSpec(m *message.Spec) bool

// 	// SendNotification send notificatioin message to aggregator
// 	SendNotification(m *message.Notification) bool

// 	// SendNotificationToCompute send message to destination compute
// 	SendNotificationToCompute(m *message.Notification, d *node.Compute) bool
// }

// func (fg FileGateway) StartReceiveMessage() {
// 	buf, err := ioutil.ReadFile(fg.FilePath)
// 	if err != nil {
// 		// logging
// 		return
// 	}

// 	m := &model.Model{}
// 	err = yaml.Unmarshal(buf, &m)
// 	if err != nil {
// 		// logging
// 		return
// 	}

// }

// func (fg FileGateway) SendSpec(m *message.Spec) bool {
// 	// logging
// 	return true
// }

// func (fg FileGateway) SendNotification(m *message.Notification) bool {
// 	// logging
// 	return true
// }

// func (fg FileGateway) SendNotificationToCompute(m *message.Notification, d *node.Compute) bool {
// 	// logging
// 	return true
// }
