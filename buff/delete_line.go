package buff

func DeleteLine(v *BufferView, commands []string) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	b.deleteLine()
	return true
}
