package buff

import (
	"github.com/felixangell/phi/lex"
	"log"
	"strconv"
)

// BufferDirtyState ...
type BufferDirtyState bool

type BufferActionFn func(*BufferView, []*lex.Token) BufferDirtyState

type BufferAction struct {
	name          string
	proc          BufferActionFn
	showInPalette bool
}

func NewBufferAction(name string, proc BufferActionFn) BufferAction {
	return BufferAction{
		name:          name,
		proc:          proc,
		showInPalette: true,
	}
}

func OpenFile(v *BufferView, _ []*lex.Token) BufferDirtyState {
	panic("unimplemented")
}

func NewFile(v *BufferView, commands []*lex.Token) BufferDirtyState {
	// TODO some nice error stuff
	// have an error roll thing in the view?

	if !commands[0].IsType(lex.String) {
		return false
	}

	fileName := commands[0].Lexeme
	// strip out the quotes (1...n-1)
	fileName = fileName[1 : len(fileName)-1]

	buff := v.AddBuffer()
	buff.OpenFile(fileName)

	buff.SetFocus(true)
	v.focusedBuff = buff.index

	return false
}

func GotoLine(v *BufferView, commands []*lex.Token) BufferDirtyState {
	if len(commands) == 0 {
		return false
	}

	if !commands[0].IsType(lex.Number) {
		return false
	}

	lineNum, err := strconv.ParseInt(commands[0].Lexeme, 10, 64)
	if err != nil {
		log.Println("goto line invalid argument ", err.Error())
		return false
	}

	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	b.gotoLine(lineNum)
	return false
}

func focusLeft(v *BufferView, _ []*lex.Token) BufferDirtyState {
	if v == nil {
		return false
	}
	v.ChangeFocus(-1)
	return false
}

func focusRight(v *BufferView, _ []*lex.Token) BufferDirtyState {
	if v == nil {
		return false
	}
	v.ChangeFocus(1)
	return false
}

func pageDown(v *BufferView, _ []*lex.Token) BufferDirtyState {
	if v == nil {
		return false
	}
	buff := v.getCurrentBuff()
	if buff == nil {
		return false
	}

	buff.scrollDown(DefaultScrollAmount)
	for i := 0; i < DefaultScrollAmount; i++ {
		buff.moveDown()
	}
	return false
}

func pageUp(v *BufferView, _ []*lex.Token) BufferDirtyState {
	if v == nil {
		return false
	}
	buff := v.getCurrentBuff()
	if buff == nil {
		return false
	}

	buff.scrollUp(DefaultScrollAmount)
	for i := 0; i < DefaultScrollAmount; i++ {
		buff.moveUp()
	}
	return false
}
