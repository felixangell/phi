package gui

import "log"

func Save(b *Buffer) bool {
	log.Println("Save: unimplemented!")
	return false
}

// FIXME make this better!
// behaviour is a little off at the moment.
func DeleteLine(b *Buffer) bool {
	if b.curs.y == 0 {
		for b.curs.x < b.contents[b.curs.y].Len() {
			b.curs.move(1, 0)
		}
		b.deleteBeforeCursor()
		return true
	}
	if b.curs.y >= len(b.contents) {
		return false
	}
	b.contents = remove(b.contents[:], b.curs.y)
	return true
}
