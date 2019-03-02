package lex

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	Word      TokenType = "word"
	Symbol              = "sym"
	Character           = "char"
	String              = "string"
	Number              = "num"
)

type Token struct {
	Lexeme string
	Type   TokenType
	Start  int
}

func (t *Token) Equals(str string) bool {
	return strings.Compare(str, t.Lexeme) == 0
}

func (t *Token) IsType(typ TokenType) bool {
	return t.Type == typ
}

func NewToken(lexeme string, kind TokenType, start int) *Token {
	return &Token{lexeme, kind, start}
}

func (t *Token) String() string {
	return fmt.Sprintf("{ %s = %s }", t.Lexeme, string(t.Type))
}
