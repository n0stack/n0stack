package message

import (
	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

type Notification struct {
	SpecID      uuid.UUID
	Model       model.AbstractModel
	Event       string // enum的なのにしたい
	IsSucceeded bool
	Description string
}
