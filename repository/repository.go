package repository

import (
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	"github.com/satori/go.uuid"
)

// Repository is application service to store messages on RDBMS, KVS, GraphDB and etc.
type Repository interface {
	// DigModel dig model for directed graph to specified depth on specified event.
	//
	// Args:
	// 	id: Model ID such as uuid.
	// 	event: Notification event such as "APPLIED" and "SCHEDULED".
	// 	depth: Depth of model dependency.
	// 		   For example, "VM -> Volume" is 1, "VM" is 0, and "VM -> Volume -> Volume agent" is 2.
	//
	// Return:
	// 	Model on event which is setted models until depth.
	//
	// Example:
	// 	>>> m = r.read("...", event="APPLIED", depth=1)
	// 	>>> m.dependencies -> not None
	// 	>>> m.dependencies.model.dependencies -> None
	DigModel(i *uuid.UUID, e string, d uint) (*model.Model, error)

	// Schedule is needed to implement *after v0.0.3*.
	//
	// Args:
	// 	model: Model of necessary to schedule models.
	// 	ids: List of necessary to create models.
	//
	// Return:
	// 	Model which is attached scheduled agent model.
	// Schedule(m *model.AbstractModel, i []string) *model.AbstractModel

	// Store store message to provide query methods like Repository.read and Repository.schedule.
	//
	// Args:
	// 	message: Message to store.
	// 			 Model on the top is only stored.
	StoreNotification(m *message.Notification)
}
