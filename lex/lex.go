package lex

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	startingPos int
	pos         int
	input       []rune
}

func New(input string) *Lexer {
	return &Lexer{
		startingPos: 0,
		pos:         0,
		input:       []rune(input),
	}
}

func (l *Lexer) consume() rune {
	consumed := l.peek()
	l.pos++
	return consumed
}

func (l *Lexer) expect(c rune) (rune, bool) {
	if l.hasNext() && l.peek() == c {
		return l.consume(), true
	}
	if !l.hasNext() {
		return rune(0), false
	}
	// TODO, fail?
	return l.consume(), true
}

func (l *Lexer) next(offs int) rune {
	return l.input[l.pos+offs]
}

func (l *Lexer) peek() rune {
	return l.input[l.pos]
}

func (l *Lexer) hasNext() bool {
	return l.pos < len(l.input)
}

func (l *Lexer) recognizeString() *Token {
	l.expect('"')
	for l.hasNext() && l.peek() != '"' {
		l.consume()
	}
	l.expect('"')
	return NewToken(l.captureLexeme(), String, l.startingPos)
}

func (l *Lexer) recognizeCharacter() *Token {
	l.expect('\'')
	for l.hasNext() && l.peek() != '\'' {
		l.consume()
	}
	l.expect('\'')
	return NewToken(l.captureLexeme(), Character, l.startingPos)
}

func (l *Lexer) recognizeNumber() *Token {
	for l.hasNext() && unicode.IsDigit(l.peek()) {
		l.consume()
	}
	if l.hasNext() && l.peek() == '.' {
		l.consume()
		for l.hasNext() && unicode.IsDigit(l.peek()) {
			l.consume()
		}
	}
	return NewToken(l.captureLexeme(), Number, l.startingPos)
}

func (l *Lexer) recognizeSymbol() *Token {
	l.consume()
	return NewToken(l.captureLexeme(), Symbol, l.startingPos)
}

func (l *Lexer) recognizeWord() *Token {
	for l.hasNext() && (unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek())) {
		l.consume()
	}

	if l.hasNext() {
		curr := l.peek()
		if curr == '_' || curr == '-' {
			l.consume()
			for l.hasNext() && (unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek())) {
				l.consume()
			}
		}
	}

	return NewToken(l.captureLexeme(), Word, l.startingPos)
}

func (l *Lexer) captureLexeme() string {
	return string(l.input[l.startingPos:l.pos])
}

func (l *Lexer) Tokenize() []*Token {
	var result []*Token
	for l.hasNext() {
		// TODO make it so that we can generate
		// lexers from the config files
		// allowing the user to put token
		// matching criteria in here. for now
		// we'll just go with a simple lexer
		// that splits strings by spaces/tabs/etc

		// skip all the layout characters
		// we dont care about these.
		for l.hasNext() && l.peek() <= ' ' {
			l.consume()
		}

		l.startingPos = l.pos

		if token := func() *Token {
			if !l.hasNext() {
				return nil
			}

			curr := l.peek()
			switch {
			case curr == '"':
				return l.recognizeString()
			case curr == '\'':
				return l.recognizeCharacter()
			case unicode.IsLetter(curr):
				return l.recognizeWord()
			case unicode.IsDigit(curr):
				return l.recognizeNumber()
			case unicode.IsGraphic(curr):
				return l.recognizeSymbol()
			case curr == ' ':
				return nil
			}

			panic(fmt.Sprintln("unhandled input! ", string(curr)))
		}(); token != nil {
			result = append(result, token)
		}

	}
	return result
}
