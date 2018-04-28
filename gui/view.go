package gui

import (
	"fmt"
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
	"log"
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
	p.SetFocus(false)

	// set focus to the buffer
	// that invoked the cmd palette
	p.parentBuff.SetFocus(true)

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
	for _, buff := range n.buffers {
		buff.SetFocus(false)
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
		for i, buff := range n.buffers {
			buff.Resize(bufferWidth, n.h)
			buff.SetPosition(bufferWidth*i, 0)
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

	prevBuff.SetFocus(false)
	if buff, ok := n.buffers[n.focusedBuff]; ok {
		buff.SetFocus(true)
	}
}

func (n *View) OnInit() {
}

func (n *View) OnUpdate() bool {
	dirty := false

	if buff, ok := n.buffers[n.focusedBuff]; ok {
		buff.OnUpdate()
	}

	n.commandPalette.OnUpdate()

	return dirty
}

func (n *View) OnRender(ctx *strife.Renderer) {
	for idx, buffer := range n.buffers {
		buffer.OnRender(ctx)

		ctx.String(fmt.Sprintf("idx %d", idx), (buffer.x+buffer.w)-150, (buffer.y+buffer.h)-150)
	}

	n.commandPalette.OnRender(ctx)
}

func (n *View) OnDispose() {}

func (n *View) AddBuffer() *Buffer {
	if buf, ok := n.buffers[n.focusedBuff]; ok {
		buf.SetFocus(false)
	}

	c := NewBuffer(n.conf, n, len(n.buffers))
	c.SetFocus(true)

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	var bufferWidth int
	bufferWidth = n.w / (c.index + 1)

	n.buffers[c.index] = c
	n.focusedBuff = c.index

	// translate all the components accordingly.
	for i, buff := range n.buffers {
		buff.Resize(bufferWidth, n.h)
		buff.SetPosition(bufferWidth*i, 0)
	}

	return c
}
