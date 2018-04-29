package gui

import "os"

type BufferAction struct {
	name string
	proc func(*Buffer) bool
}

func NewBufferAction(name string, proc func(*Buffer) bool) BufferAction {
	return BufferAction{
		name: name,
		proc: proc,
	}
}

var actions = map[string]BufferAction{
	"save":         NewBufferAction("save", Save),
	"delete_line":  NewBufferAction("delete_ine", DeleteLine),
	"close_buffer": NewBufferAction("close_buffer", CloseBuffer),
	"paste":        NewBufferAction("paste", Paste),
	"show_palette": NewBufferAction("show_palette", ShowPalette),
	"exit": NewBufferAction("exit", func(*Buffer) bool {
		// TODO do this properly lol
		os.Exit(0)
		return false
	}),
}
