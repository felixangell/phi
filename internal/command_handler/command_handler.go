package command_handler

import (
	"fmt"
	"github.com/felixangell/phi/internal/cfg"
	"github.com/felixangell/strife"
	"strings"
)

var commandHandlerInst *CommandHandler

func DeduceCommand(superDown bool, controlDown bool, keys ...int) (string, bool) {
	if commandHandlerInst == nil {
		panic("CommandHandler not initialised yet!")
	}
	return commandHandlerInst.DeduceCommandName(superDown, controlDown, keys...)
}

func SetupCommandHandler(config *cfg.PhiEditorConfig) {
	commandHandlerInst = newCommandHandler(config)
}

type CommandHandler struct {
	triggers map[intHashSetHashKey]string
}

// TODO build this map up further.
var shortcutMap = map[string]int{
	// for the simplicity of making the API work
	// we assume that control and super map to the left
	// of their keys, even if the right is pressed. it is upto
	// the caller to consolidate the two (for now).
	"ctrl": strife.KEY_LCTRL,
	"super":   strife.KEY_LGUI,

	"a": strife.KEY_A,
	"b": strife.KEY_B,
	"c": strife.KEY_C,
	"d": strife.KEY_D,
	"e": strife.KEY_E,
	"f": strife.KEY_F,
	"g": strife.KEY_G,
	"h": strife.KEY_H,
	"i": strife.KEY_I,
	"j": strife.KEY_J,
	"k": strife.KEY_K,
	"l": strife.KEY_L,
	"m": strife.KEY_M,
	"n": strife.KEY_N,
	"o": strife.KEY_O,
	"p": strife.KEY_P,
	"q": strife.KEY_Q,
	"r": strife.KEY_R,
	"s": strife.KEY_S,
	"t": strife.KEY_T,
	"u": strife.KEY_U,
	"v": strife.KEY_V,
	"w": strife.KEY_W,
	"x": strife.KEY_X,
	"y": strife.KEY_Y,
	"z": strife.KEY_Z,

	"left":  strife.KEY_LEFT,
	"right": strife.KEY_RIGHT,
	"up":    strife.KEY_UP,
	"down":  strife.KEY_DOWN,
}

func mapShortcutKeyword(keyword string) int {
	if val, ok := shortcutMap[keyword]; ok {
		return val
	}
	panic(fmt.Sprintf("unsupported keyword %s", keyword))
}

func parseHashCombo(combo string) intHashSetHashKey {
	hs := newIntHashSet()
	parts := strings.Split(combo, "+")
	for _, part := range parts {
		hs.Store(mapShortcutKeyword(part))
	}
	return hs.Hash()
}

func (c *CommandHandler) DeduceCommandName(superDown bool, controlDown bool, keys ...int) (string, bool) {
	hs := newIntHashSet(keys...)
	if superDown {
		hs.Store(mapShortcutKeyword("super"))
	}
	if controlDown {
		hs.Store(mapShortcutKeyword("ctrl"))
	}

	cmdName, ok := c.triggers[hs.Hash()]
	return cmdName, ok
}

func populateTriggers(config *cfg.PhiEditorConfig, handler *CommandHandler) {
	for cmdName, cmd := range config.Commands {
		hash := parseHashCombo(cmd.Shortcut)
		handler.triggers[hash] = cmdName
	}
}

func newCommandHandler(config *cfg.PhiEditorConfig) *CommandHandler {
	handler := &CommandHandler{
		triggers: map[intHashSetHashKey]string{},
	}
	populateTriggers(config, handler)
	return handler
}
