package gateway

import (
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model/node"
)

// Gateway provide methods of incoming or outgoing Messages with other services.
type Gateway interface {
	// StartReceiveMessage start receive message from something.
	StartReceiveMessage()

	// SendMessage send message to default destination like aggregator.
	SendMessage(m *message.AbstractMessage) bool

	// SendNotificationToCompute send message to destination compute
	SendNotificationToCompute(m *message.Notification, d *node.Compute) bool
}
