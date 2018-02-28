package gui

import (
	"bytes"
	"io/ioutil"
	"log"
)

// NOTE: all shortcuts return a bool
// this is whether or not they have
// modified the buffer
// if the buffer is modified it will be
// re-rendered.

func Save(b *Buffer) bool {

	var buffer bytes.Buffer
	for idx, line := range b.contents {
		if idx > 0 {
			// TODO: this avoids a trailing newline
			// if we handle it like this? but if we have
			// say enforce_newline_at_eof or something we
			// might want to do this all the time
			buffer.WriteRune('\n')
		}
		buffer.WriteString(line.String())
	}

	err := ioutil.WriteFile(b.filePath, buffer.Bytes(), 0775)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Wrote file '", b.filePath, "' to disk")
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
