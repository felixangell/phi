package action

import (
	"log"
	"strconv"
	"strings"

	"github.com/felixangell/phi/buff"
)

type BufferAction struct {
	name          string
	proc          func(*buff.BufferView, []string) bool
	showInPalette bool
}

func NewBufferAction(name string, proc func(*buff.BufferView, []string) bool) BufferAction {
	return BufferAction{
		name:          name,
		proc:          proc,
		showInPalette: true,
	}
}

func OpenFile(v *buff.BufferView, commands []string) bool {
	path := ""
	if path == "" {
		panic("unimplemented")
		// ive removed this since the cross platform
		// thing causes too much hassle on diff. platforms
		// going to wriet my own file open viewer thing built
		// into the editor instead.
	}

	buff := v.AddBuffer()
	if len(strings.TrimSpace(path)) == 0 {
		return false
	}

	buff.OpenFile(path)
	buff.SetFocus(true)
	v.focusedBuff = buff.index

	return false
}

func NewFile(v *buff.BufferView, commands []string) bool {
	// TODO some nice error stuff
	// have an error roll thing in the view?

	buff := v.AddBuffer()
	buff.OpenFile(commands[0])

	buff.SetFocus(true)
	v.focusedBuff = buff.index

	return false
}

func GotoLine(v *buff.BufferView, commands []string) bool {
	if len(commands) == 0 {
		return false
	}

	lineNum, err := strconv.ParseInt(commands[0], 10, 64)
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

func focusLeft(v *buff.BufferView, commands []string) bool {
	if v == nil {
		return false
	}
	v.ChangeFocus(-1)
	return false
}

func focusRight(v *buff.BufferView, commands []string) bool {
	if v == nil {
		return false
	}
	v.ChangeFocus(1)
	return false
}

func pageDown(v *buff.BufferView, commands []string) bool {
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

func pageUp(v *buff.BufferView, commands []string) bool {
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
