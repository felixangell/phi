package gui

import (
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
)

// View is an array of buffers basically.
type View struct {
	BaseComponent

	conf           *cfg.TomlConfig
	buffers        map[int]*Buffer
	focusedBuff    int
	commandPalette *CommandPalette
}

func NewView(width, height int, conf *cfg.TomlConfig) *View {
	view := &View{
		conf:    conf,
		buffers: map[int]*Buffer{},
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
	p.HasFocus = false

	// set focus to the buffer
	// that invoked the cmd palette
	p.parentBuff.inputHandler = p.buff.inputHandler
	p.parentBuff.HasFocus = true

	// remove focus from palette
	p.buff.HasFocus = false
	p.buff.SetInputHandler(nil)
}

func (n *View) focusPalette(buff *Buffer) {
	p := n.commandPalette
	p.HasFocus = true

	// focus the command palette
	p.buff.HasFocus = true
	p.buff.SetInputHandler(buff.inputHandler)

	// remove focus from the buffer
	// that invoked the command palette
	buff.inputHandler = nil
	p.parentBuff = buff
}

func (n *View) UnfocusBuffers() {
	// clear focus from buffers
	for _, buff := range n.buffers {
		buff.HasFocus = false
		buff.inputHandler = nil
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

func (n *View) ChangeFocus(dir int) {
	println("implement me! ", dir)
}

func (n *View) OnInit() {
}

func (n *View) OnUpdate() bool {
	dirty := false

	for _, buffer := range n.buffers {
		if buffer.OnUpdate() {
			dirty = true
		}
	}
	n.commandPalette.OnUpdate()

	return dirty
}

func (n *View) OnRender(ctx *strife.Renderer) {
	for _, buffer := range n.buffers {
		buffer.OnRender(ctx)
	}
	n.commandPalette.OnRender(ctx)
}

func (n *View) OnDispose() {}

func (n *View) AddBuffer() *Buffer {
	if buf, ok := n.buffers[n.focusedBuff]; ok {
		buf.HasFocus = false
	}

	c := NewBuffer(n.conf, n, len(n.buffers))
	c.HasFocus = true

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	var bufferWidth int
	bufferWidth = n.w / (len(n.buffers) + 1)

	n.buffers[c.index] = c
	n.focusedBuff = c.index

	// translate all the components accordingly.
	for i, buff := range n.buffers {
		if buff == nil {
			continue
		}

		buff.Resize(bufferWidth, n.h)
		buff.SetPosition(bufferWidth*i, 0)
	}

	return c
}
