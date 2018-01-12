package aggregator

import (
	"fmt"

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

func (a Aggregator) ProcessMessage(m message.AbstractMessage) error {
	n, ok := m.(*message.Notification)
	if !ok {
		return fmt.Errorf("Received notification message which is not supported, maybe there are stranger or distributor has bugs")
	}

	if !a.Repository.StoreNotification(n) {
		return fmt.Errorf("Failed to store notification message")
	}

	return nil
}
