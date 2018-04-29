package gui

import "os"

type BufferAction struct {
	name string
	proc func(*View, []string) bool
}

func NewBufferAction(name string, proc func(*View, []string) bool) BufferAction {
	return BufferAction{
		name: name,
		proc: proc,
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

var actions = map[string]BufferAction{
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
