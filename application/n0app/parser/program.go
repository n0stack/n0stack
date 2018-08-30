package parser

type Program struct {
	Sentences []Sentence
}

type Sentence struct {
	Tokens     []Token
	LineNumber int
}

type TokenKind int

const (
	TokenKind_Opecode        TokenKind = iota
	TokenKind_Comment        TokenKind = iota
	TokenKind_Operand_String TokenKind = iota
	TokenKind_Operand_Json   TokenKind = iota
)

type Token struct {
	Kind  TokenKind
	Value []byte
}
