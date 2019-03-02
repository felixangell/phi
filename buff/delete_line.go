package buff

import "github.com/felixangell/phi/lex"

func DeleteLine(v *BufferView, commands []*lex.Token) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	b.deleteLine()
	return true
}
