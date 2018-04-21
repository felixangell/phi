package lex

import "fmt"

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

func (t *Token) String() string {
	return fmt.Sprintf("lexeme: %s, type %s, at pos %d", t.Lexeme, t.Type, t.Start)
}