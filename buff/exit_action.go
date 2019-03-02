package buff

import (
	"log"
	"os"

	"github.com/felixangell/phi/lex"
)

func ExitPhi(v *BufferView, commands []*lex.Token) bool {
	// todo this probably wont work...
	// would also be nice to have a thing
	// that asks if we want to save all buffers
	// rather than going thru each one specifically?
	for i := range v.buffers {
		CloseBuffer(v, []*lex.Token{})
		log.Println("Closing buffer ", i)
	}

	os.Exit(0)
	return false
}
