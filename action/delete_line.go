package action

import "github.com/felixangell/phi/buff"

func DeleteLine(v *buff.BufferView, commands []string) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	b.deleteLine()
	return true
}
