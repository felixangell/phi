package gui

import "os"

type BufferAction func(*Buffer) bool

var actions = map[string]BufferAction{
	"save":         Save,
	"delete_line":  DeleteLine,
	"close_buffer": CloseBuffer,
	"paste":        Paste,
	"show_palette": ShowPalette,
	"exit": func(*Buffer) bool {
		// TODO do this properly lol
		os.Exit(0)
		return false
	},
}
