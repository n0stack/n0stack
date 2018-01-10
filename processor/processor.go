package processor

import (
	"github.com/n0stack/n0core/message"
)

type Processor interface {
	ProcessMessage(m message.AbstractMessage)
}
