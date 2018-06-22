package gui

import (
	"log"
	"strconv"
	"strings"
)

type BufferAction struct {
	name          string
	proc          func(*View, []string) bool
	showInPalette bool
}

func NewBufferAction(name string, proc func(*View, []string) bool) BufferAction {
	return BufferAction{
		name:          name,
		proc:          proc,
		showInPalette: true,
	}
}

func OpenFile(v *View, commands []string) bool {
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

func NewFile(v *View, commands []string) bool {
	// TODO some nice error stuff
	// have an error roll thing in the view?

	buff := v.AddBuffer()
	buff.OpenFile(commands[0])

	buff.SetFocus(true)
	v.focusedBuff = buff.index

	return false
}

func GotoLine(v *View, commands []string) bool {
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

func focusLeft(v *View, commands []string) bool {
	if v == nil {
		return false
	}
	v.ChangeFocus(-1)
	return false
}

func focusRight(v *View, commands []string) bool {
	if v == nil {
		return false
	}
	v.ChangeFocus(1)
	return false
}

func pageDown(v *View, commands []string) bool {
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

func pageUp(v *View, commands []string) bool {
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

var actions = map[string]BufferAction{
	"page_down": NewBufferAction("page_down", pageDown),
	"page_up":   NewBufferAction("page_up", pageUp),

	"undo": NewBufferAction("undo", Undo),
	"redo": NewBufferAction("redo", Redo),

	"focus_left":  NewBufferAction("focus_left", focusLeft),
	"focus_right": NewBufferAction("focus_right", focusRight),

	"goto":         NewBufferAction("goto", GotoLine),
	"new":          NewBufferAction("new", NewFile),
	"open":         NewBufferAction("open", OpenFile),
	"save":         NewBufferAction("save", Save),
	"delete_line":  NewBufferAction("delete_line", DeleteLine),
	"close_buffer": NewBufferAction("close_buffer", CloseBuffer),
	"paste":        NewBufferAction("paste", Paste),
	"show_palette": NewBufferAction("show_palette", ShowPalette),
	"exit":         NewBufferAction("exit", ExitPhi),
}
