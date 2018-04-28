package gui

import (
	"bytes"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"log"
)

func ShowPalette(b *Buffer) bool {
	b.parent.UnfocusBuffers()
	b.parent.focusPalette(b)
	return true
}

func Paste(b *Buffer) bool {
	str, err := clipboard.ReadAll()

	if err == nil {
		b.insertString(b.curs.x, str)
		b.moveToEndOfLine()
		return true
	}

	log.Println("Failed to paste from clipboard: ", err.Error())
	return false
}

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

	// TODO:
	// - files probably dont have to be entirely
	//   re-saved all the time!
	// - we can probably stream this somehow?
	// - multi threaded?
	// - lots of checks to do here: does the file exist/not exist
	//   handle the errors... etc.

	err := ioutil.WriteFile(b.filePath, buffer.Bytes(), 0775)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Wrote file '" + b.filePath + "' to disk")
	return false
}
