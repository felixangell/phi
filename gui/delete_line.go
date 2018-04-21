package gui

import rope "github.com/felixangell/go-rope"

func DeleteLine(b *Buffer) bool {
	if len(b.contents) > 1 {
		b.contents = remove(b.contents, b.curs.y)
	} else {
		// we are on the first line
		// and there is nothing else to delete
		// so we just clear the line
		b.contents[b.curs.y] = new(rope.Rope)
	}

	if b.curs.y >= len(b.contents) {
		if b.curs.y > 0 {
			b.moveUp()
		}
	}

	b.moveToEndOfLine()
	return true
}
