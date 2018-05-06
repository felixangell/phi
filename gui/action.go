package gui

import (
	"log"
	"os"
	"strconv"
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

var actions = map[string]BufferAction{
	"focus_left":   NewBufferAction("focus_left", focusLeft),
	"focus_right":  NewBufferAction("focus_right", focusRight),
	"goto":         NewBufferAction("goto", GotoLine),
	"new":          NewBufferAction("new", NewFile),
	"save":         NewBufferAction("save", Save),
	"delete_line":  NewBufferAction("delete_line", DeleteLine),
	"close_buffer": NewBufferAction("close_buffer", CloseBuffer),
	"paste":        NewBufferAction("paste", Paste),
	"show_palette": NewBufferAction("show_palette", ShowPalette),
	"exit": NewBufferAction("exit", func(*View, []string) bool {
		// TODO do this properly lol
		os.Exit(0)
		return false
	}),
}
