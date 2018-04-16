package gui

import rope "github.com/felixangell/go-rope"

func DeleteLine(b *Buffer) bool {
	if b.curs.y == 0 {
		b.contents[b.curs.y] = new(rope.Rope)
		b.moveToEndOfLine()
		return false
	}

	if b.curs.y >= len(b.contents) {
		return false
	}

	prevLineLen := b.contents[b.curs.y].Len()
	b.contents = remove(b.contents, b.curs.y)

	if b.curs.y >= len(b.contents) {
		b.moveUp()
		b.moveToEndOfLine()
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
