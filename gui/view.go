package gui

import (
	"log"
	"runtime"
	"unicode"

	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
	"github.com/veandco/go-sdl2/sdl"
)

// View is an array of buffers basically.
type View struct {
	BaseComponent

	conf           *cfg.TomlConfig
	buffers        map[int]*BufferPane
	focusedBuff    int
	commandPalette *CommandPalette
}

func NewView(width, height int, conf *cfg.TomlConfig) *View {
	view := &View{
		conf:    conf,
		buffers: map[int]*BufferPane{},
	}

	view.Translate(width, height)
	view.Resize(width, height)

	view.commandPalette = NewCommandPalette(*conf, view)
	view.UnfocusBuffers()

	return view
}

func (n *View) hidePalette() {
	p := n.commandPalette
	p.clearInput()
	p.SetFocus(false)

	// set focus to the buffer
	// that invoked the cmd palette
	if p.parentBuff != nil {
		p.parentBuff.SetFocus(true)
		n.focusedBuff = p.parentBuff.index
	}

	// remove focus from palette
	p.buff.SetFocus(false)
}

func (n *View) focusPalette(buff *Buffer) {
	p := n.commandPalette
	p.SetFocus(true)

	// focus the command palette
	p.buff.SetFocus(true)

	// remove focus from the buffer
	// that invoked the command palette
	p.parentBuff = buff
}

func (n *View) UnfocusBuffers() {
	// clear focus from buffers
	for _, buffPane := range n.buffers {
		buffPane.Buff.SetFocus(false)
	}
}

func sign(dir int) int {
	if dir > 0 {
		return 1
	} else if dir < 0 {
		return -1
	}
	return 0
}

func (n *View) removeBuffer(index int) {
	log.Println("Removing buffer index:", index)
	delete(n.buffers, index)

	// only resize the buffers if we have
	// some remaining in the window
	if len(n.buffers) > 0 {
		bufferWidth := n.w / len(n.buffers)

		// translate all the components accordingly.
		for i, buffPane := range n.buffers {
			buffPane.Buff.Resize(bufferWidth, n.h)
			buffPane.Buff.SetPosition(bufferWidth*i, 0)
		}
	}

}

func (n *View) ChangeFocus(dir int) {
	prevBuff, _ := n.buffers[n.focusedBuff]

	if dir == -1 {
		n.focusedBuff--
	} else if dir == 1 {
		n.focusedBuff++
	}

	if n.focusedBuff < 0 {
		n.focusedBuff = len(n.buffers) - 1
	} else if n.focusedBuff >= len(n.buffers) {
		n.focusedBuff = 0
	}

	if prevBuff != nil {
		prevBuff.Buff.SetFocus(false)
	}

	if buffPane, ok := n.buffers[n.focusedBuff]; ok {
		buffPane.Buff.SetFocus(true)
	}
}

func (n *View) getCurrentBuff() *Buffer {
	if buffPane, ok := n.buffers[n.focusedBuff]; ok {
		return buffPane.Buff
	}
	return nil
}

func (n *View) OnInit() {
}

func (n *View) OnUpdate() bool {
	dirty := false

	CONTROL_DOWN = strife.KeyPressed(sdl.K_LCTRL) || strife.KeyPressed(sdl.K_RCTRL)
	SUPER_DOWN = strife.KeyPressed(sdl.K_LGUI) || strife.KeyPressed(sdl.K_RGUI)

	shortcutName := "ctrl"
	source := cfg.Shortcuts.Controls

	if strife.PollKeys() && (SUPER_DOWN || CONTROL_DOWN) {
		if runtime.GOOS == "darwin" {
			if SUPER_DOWN {
				source = cfg.Shortcuts.Supers
				shortcutName = "super"
			} else if CONTROL_DOWN {
				source = cfg.Shortcuts.Controls
				shortcutName = "control"
			}
		} else {
			source = cfg.Shortcuts.Supers
		}

		r := rune(strife.PopKey())

		if r == sdl.K_F12 {
			DEBUG_MODE = !DEBUG_MODE
		}

		left := sdl.K_LEFT
		right := sdl.K_RIGHT
		up := sdl.K_UP
		down := sdl.K_DOWN

		// map to left/right/etc.
		// FIXME
		var key string
		switch int(r) {
		case left:
			key = "left"
		case right:
			key = "right"
		case up:
			key = "up"
		case down:
			key = "down"
		default:
			key = string(unicode.ToLower(r))
		}

		actionName, actionExists := source[key]
		if actionExists {
			if action, ok := actions[actionName]; ok {
				log.Println("Executing action '" + actionName + "'")
				return action.proc(n, []string{})
			}
		} else {
			log.Println("view: unimplemented shortcut", shortcutName, "+", string(unicode.ToLower(r)), "#", int(r), actionName, key)
		}
	}

	if buffPane, ok := n.buffers[n.focusedBuff]; ok {
		buffPane.OnUpdate()
	}

	n.commandPalette.OnUpdate()

	return dirty
}

func (n *View) OnRender(ctx *strife.Renderer) {
	for _, buffPane := range n.buffers {
		buffPane.OnRender(ctx)
	}

	n.commandPalette.OnRender(ctx)
}

func (n *View) OnDispose() {}

func (n *View) AddBuffer() *Buffer {
	n.UnfocusBuffers()

	cfg := n.conf
	c := NewBuffer(cfg, BufferConfig{
		cfg.Theme.Background,
		cfg.Theme.Foreground,
		cfg.Theme.Cursor,
		cfg.Theme.Cursor_Invert,
		cfg.Theme.Gutter_Background,
		cfg.Theme.Gutter_Foreground,
		cfg.Editor.Loaded_Font,
	}, n, len(n.buffers))

	c.SetFocus(true)

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	var bufferWidth int
	bufferWidth = n.w / (c.index + 1)

	n.buffers[c.index] = NewBufferPane(c)
	n.focusedBuff = c.index

	// translate all the buffers accordingly.
	for i, buffPane := range n.buffers {
		buffPane.Resize(bufferWidth, n.h)
		buffPane.SetPosition(bufferWidth*i, 0)
	}

	return c
}
