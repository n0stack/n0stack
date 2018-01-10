package processor

import (
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/repository"
)

// Aggregator is a processor which store messages.
//
// 1. Receive a message from gateway.
// 2. Store messages to repository to provide repository functions.
//
// Args:
// 	repository: Data store to store result.
//
// Exaples:
type Aggregator struct {
	Repository repository.Repository
}

func (a Aggregator) ProcessMessage(m message.AbstractMessage) {
	n, ok := m.(message.Notification)
	if !ok {
		return
	}

	a.Repository.StoreNotification(&n)
}
