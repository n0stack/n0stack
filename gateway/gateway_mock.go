package gateway

import (
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
)

type MockGateway struct {
	PatchSendMessage               func(m *message.AbstractMessage) bool
	PatchSendNotificationToCompute func(m *message.Notification, d *model.Compute) bool
}

func (mock MockGateway) SendMessage(m *message.AbstractMessage) bool {
	return mock.PatchSendMessage(m)
}

func (mock MockGateway) SendNotificationToCompute(m *message.Notification, d *model.Compute) bool {
	return mock.PatchSendNotificationToCompute(m, d)
}
