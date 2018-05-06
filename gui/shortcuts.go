package gui

import (
	"bytes"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"log"
)

func ShowPalette(v *View, commands []string) bool {
	b := v.getCurrentBuff()
	v.UnfocusBuffers()
	v.focusPalette(b)
	return true
}

func Paste(v *View, commands []string) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

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

func Save(v *View, commands []string) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

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

	// ALWAYS use atomic save possibly?
	// start writing the file on a gothread
	// and then move the file overwriting when the
	// file has been written?

	err := ioutil.WriteFile(b.filePath, buffer.Bytes(), 0775)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Wrote file '" + b.filePath + "' to disk")

	b.modified = false

	return false
}
