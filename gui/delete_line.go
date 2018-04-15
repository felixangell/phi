package gui

func DeleteLine(b *Buffer) bool {
	if b.curs.y >= len(b.contents) {
		return false
	}

	prevLineLen := b.contents[b.curs.y].Len()
	b.contents = remove(b.contents[:], b.curs.y)

	currLineLen := b.contents[b.curs.y].Len()

	if b.curs.x > currLineLen {
		amountToMove := prevLineLen - currLineLen
		for i := 0; i < amountToMove; i++ {
			b.moveLeft()
		}
	}
	return true
}
