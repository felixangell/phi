package gui

func CloseBuffer(v *View, commands []string) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	if b.modified {
		// do command palette thing!
		return false
	}

	v.removeBuffer(b.index)
	return false
}
