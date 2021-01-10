package buff

import "github.com/felixangell/phi/internal/lex"

var register = map[string]BufferAction{
	"page_down":    NewBufferAction("page_down", pageDown),
	"page_up":      NewBufferAction("page_up", pageUp),
	"undo":         NewBufferAction("undo", Undo),
	"redo":         NewBufferAction("redo", Redo),
	"focus_left":   NewBufferAction("focus_left", focusLeft),
	"focus_right":  NewBufferAction("focus_right", focusRight),
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

func ExecuteCommandIfExist(command string, view *BufferView, tokens ...*lex.Token) BufferDirtyState {
	if cmd, ok := register[command]; ok {
		return cmd.proc(view, tokens)
	}
	return false
}
