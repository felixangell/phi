package lex

type Lexer struct {
	pos int
	input []rune
}

func New(input string) *Lexer {
	return &Lexer {
		pos: 0,
		input: []rune(input),
	}
}

func (l *Lexer) consume() rune {
	consumed := l.peek()
	l.pos++
	return consumed
}

func (l *Lexer) next(offs int) rune {
	return l.input[l.pos + offs]
}

func (l *Lexer) peek() rune {
	return l.input[l.pos]
}

func (l *Lexer) hasNext() bool {
	return l.pos < len(l.input)
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

		startPos := l.pos
		for l.hasNext() {
			// we run into a layout character 
			if l.peek() <= ' ' {
				break
			}

			l.consume()
		}

		// this should be a recognized
		// token i think?

		lexeme := string(l.input[startPos:l.pos])
		tok := NewToken(lexeme, Word, startPos)
		result = append(result, tok)
	}
	return result
}