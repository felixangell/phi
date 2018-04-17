package gui

import rope "github.com/felixangell/go-rope"

func DeleteLine(b *Buffer) bool {
	var prevLineLen = 0
	if len(b.contents) > 1 {
		prevLineLen = b.contents[b.curs.y].Len()
		b.contents = remove(b.contents, b.curs.y)
	} else {
		// we are on the first line
		// and there is nothing else to delete
		// so we just clear the line
		b.contents[b.curs.y] = new(rope.Rope)
		b.moveToEndOfLine()
		return true
	}

	if b.curs.y >= len(b.contents) {
		if b.curs.y > 0 {
			b.moveUp()
		}
		return false
	}

	currLineLen := b.contents[b.curs.y].Len()

	if b.curs.x > currLineLen {
		amountToMove := prevLineLen - currLineLen
		for i := 0; i < amountToMove; i++ {
			b.moveLeft()
		}
	}
	return true
}
