package lex

type TokenType uint

const (
	Word TokenType = iota
)

type Token struct {
	Lexeme string
	Type TokenType
	Start int
}

func NewToken(lexeme string, kind TokenType, start int) *Token {
	return &Token {lexeme, kind, start}
}