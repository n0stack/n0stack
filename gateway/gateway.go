package gateway

import (
	"time"

	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

// Gateway provide methods of incoming or outgoing Messages with other services.
type Gateway interface {
	// // StartReceiveMessage start receive message from something. 任意のコントローラを利用する
	// StartReceiveMessage()

	// SendNotification send message to distributor
	SendNotification(
		id uuid.UUID,
		task string,
		operation string,
		notifiedAt time.Time,
		isSucceeded bool,
		description string,
		model model.AbstractModel,
	) bool

	// SendNotificationToAgent send message to destination compute
	SendNotificationToAgent(
		id uuid.UUID,
		task string,
		operation string,
		notifiedAt time.Time,
		isSucceeded bool,
		description string,
		model model.AbstractModel,
	) bool
}
