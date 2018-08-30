package parser

type Program struct {
	Sentences []*Sentence
}

type Sentence struct {
	Line       []byte
	LineNumber int
	Opecode    *Token
	Operands   []*Token
	// Comment    *Token
}

type TokenKind int

const (
	Opecode        TokenKind = iota
	Comment        TokenKind = iota
	Operand_String TokenKind = iota
	Operand_Json   TokenKind = iota
)

type Token struct {
	Kind  TokenKind
	Value []byte
}
