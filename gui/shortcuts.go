package gui

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
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

func genFileName(dir, prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(dir, prefix+hex.EncodeToString(randBytes)+suffix)
}

// NOTE: all shortcuts return a bool
// this is whether or not they have
// modified the buffer
// if the buffer is modified it will be
// re-rendered.

func Save(v *View, commands []string) bool {
	// TODO Config option for this.
	atomicFileSave := true

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

	filePath := b.filePath
	if atomicFileSave {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			log.Println("cant get abs dir of path ", b.filePath, " can't save!")
			return false
		}

		ext := filepath.Ext(absPath)
		dir := filepath.Dir(absPath)

		filePath = genFileName(dir, "", ext)
	}

	err := ioutil.WriteFile(filePath, buffer.Bytes(), 0775)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Wrote file '" + b.filePath + "' to disk")
	if atomicFileSave {
		log.Println("- Wrote atomically to file:", filePath)
	}

	if atomicFileSave {
		log.Println("- Over-writing atomic file save")

		err := os.Rename(filePath, b.filePath)
		if err != nil {
			log.Println("Failed to save!", err.Error())
			return false
		}
	}

	b.modified = false
	return false
}
