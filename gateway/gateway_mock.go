package gateway

import (
	"time"

	"github.com/n0stack/n0core/model"
	uuid "github.com/satori/go.uuid"
)

type MockGateway struct {
	PatchSendNotification func(
		id uuid.UUID,
		task string,
		operation string,
		notifiedAt time.Time,
		isSucceeded bool,
		description string,
		model model.AbstractModel,
	) bool

	PatchSendNotificationToAgent func(
		id uuid.UUID,
		task string,
		operation string,
		notifiedAt time.Time,
		isSucceeded bool,
		description string,
		model model.AbstractModel,
	) bool
}

func (mock MockGateway) SendMessage(
	id uuid.UUID,
	task string,
	operation string,
	notifiedAt time.Time,
	isSucceeded bool,
	description string,
	model model.AbstractModel,
) bool {
	return mock.PatchSendNotification(id, task, operation, notifiedAt, isSucceeded, description, model)
}

func (mock MockGateway) SendNotificationToCompute(
	id uuid.UUID,
	task string,
	operation string,
	notifiedAt time.Time,
	isSucceeded bool,
	description string,
	model model.AbstractModel,
) bool {
	return mock.PatchSendNotificationToAgent(id, task, operation, notifiedAt, isSucceeded, description, model)
}
