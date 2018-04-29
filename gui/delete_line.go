package gui

func DeleteLine(v *View, commands []string) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	b.deleteLine()
	return true
}
