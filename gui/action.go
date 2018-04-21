package gui

type BufferAction func(*Buffer) bool

var actions = map[string]BufferAction{
	"save":         Save,
	"delete_line":  DeleteLine,
	"close_buffer": CloseBuffer,
	"paste": Paste,
}
