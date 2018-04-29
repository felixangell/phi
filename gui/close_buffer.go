package gui

func CloseBuffer(v *View) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	if len(v.buffers) > 1 {
		v.ChangeFocus(-1)
	}

	v.removeBuffer(b.index)

	return false
}
