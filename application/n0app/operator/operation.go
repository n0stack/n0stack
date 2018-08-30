package operator

import (
	"github.com/n0stack/n0core/application/n0app/parser"
)

type OutputTypes int

const (
	Stdout OutputTypes = iota
	Stderr OutputTypes = iota
)

type OutputLine struct {
	OutputType OutputTypes
	Line       string
	Err        error
}

type Operations interface {
	FROM(program []*parser.Token) (chan OutputLine, error)
	RUN(program []*parser.Token) (chan OutputLine, error)
	COPY(program []*parser.Token) (chan OutputLine, error)
	DAEMON(program []*parser.Token) (chan OutputLine, error)
}
