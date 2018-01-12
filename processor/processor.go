package processor

import (
	"github.com/n0stack/n0core/message"
)

// Processor is enterprise service to provide abstract process which is shown on overall architecture.
//
// "n0core" is based on onion architecture.
// Application service, which is target, repository and gateway, is depending for Processor,
// and enterprise service is depending for nothing,
// so life cycle of Processor must be long.
type Processor interface {
	ProcessMessage(m message.AbstractMessage) error
}
