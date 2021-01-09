package buff

import "github.com/felixangell/phi/lex"

func DeleteLine(v *BufferView, _ []*lex.Token) BufferDirtyState {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	b.deleteLine()
	return true
}
