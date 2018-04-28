package gui

func CloseBuffer(b *Buffer) bool {
	view := b.parent
	view.ChangeFocus(-1)
	view.removeBuffer(b.index)
	return false
}
