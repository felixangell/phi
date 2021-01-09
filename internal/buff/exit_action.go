package buff

import (
	"log"
	"os"

	"github.com/felixangell/phi/internal/lex"
)

func ExitPhi(v *BufferView, _ []*lex.Token) BufferDirtyState {
	for i := range v.buffers {
		log.Println("Closing buffer ", i)
		CloseBuffer(v, []*lex.Token{})
	}

	log.Println("Exiting!")
	os.Exit(0)
	return false
}
