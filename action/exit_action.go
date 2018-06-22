package action

import (
	"log"
	"os"

	"github.com/felixangell/phi/buff"
)

func ExitPhi(v *buff.BufferView, commands []string) bool {
	// todo this probably wont work...
	// would also be nice to have a thing
	// that asks if we want to save all buffers
	// rather than going thru each one specifically?
	for i, _ := range v.buffers {
		CloseBuffer(v, []string{})
		log.Println("Closing buffer ", i)
	}

	os.Exit(0)
	return false
}
