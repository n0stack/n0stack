package gateway

import (
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
)

// Gateway provide methods of incoming or outgoing Messages with other services.
type Gateway interface {
	// // StartReceiveMessage start receive message from something. 任意のコントローラを利用する
	// StartReceiveMessage()

	// SendNotification send message to distributor
	SendNotification(m *message.Notification) bool

	// SendNotificationToCompute send message to destination compute
	SendNotificationToCompute(m *message.Notification, d *model.Compute) bool
}
