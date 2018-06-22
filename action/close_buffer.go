package action

import (
	"fmt"

	"github.com/felixangell/phi/buff"
)

func CloseBuffer(v *buff.BufferView, commands []string) bool {
	b := v.getCurrentBuff()
	if b == nil {
		return false
	}

	// TODO eventually we should have our own
	// little dialog IO message thingies.
	if b.modified {
		// TODO basename?
		text := fmt.Sprintf("Do you want to save the changes you made to %s?", b.filePath)

		// TODO
		panic(text)

		// dontSave := dialog.Message("%s", text).YesNo()
		// if !dontSave {
		// 	return false
		// }

		// save the buffer!
		// Save(v, []string{})
	}

	v.removeBuffer(b.index)
	return false
}
