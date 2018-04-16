package gui

import (
	"github.com/felixangell/phi-editor/cfg"
	"github.com/felixangell/strife"
)

type View struct {
	BaseComponent
	conf        *cfg.TomlConfig
	buffers     map[int]*Buffer
	focusedBuff int
}

func NewView(width, height int, conf *cfg.TomlConfig) *View {
	view := &View{
		conf:    conf,
		buffers: map[int]*Buffer{},
	}
	view.Translate(width, height)
	view.Resize(width, height)
	view.focusPalette()
	return view
}

func (n *View) focusPalette() {
	// clear focus from buffers
	for _, buff := range n.buffers {
		buff.HasFocus = false
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
	// remove focus from the curr buffer.
	if buf, ok := n.buffers[n.focusedBuff]; ok {
		buf.HasFocus = false
	}

	newIndex := n.focusedBuff + sign(dir)
	if newIndex >= n.NumComponents() {
		newIndex = 0
	} else if newIndex < 0 {
		newIndex = n.NumComponents() - 1
	}

	if buff := n.components[newIndex]; buff != nil {
		n.focusedBuff = newIndex
	}
}

func (n *View) OnInit() {
}

func (n *View) OnUpdate() bool {
	dirty := false
	for _, comp := range n.components {
		if comp == nil {
			continue
		}

		if Update(comp) {
			dirty = true
		}
	}
	return dirty
}

func (n *View) OnRender(ctx *strife.Renderer) {}

func (n *View) OnDispose() {}

func (n *View) AddBuffer() *Buffer {
	if buf, ok := n.buffers[n.focusedBuff]; ok {
		buf.HasFocus = false
	}

	c := NewBuffer(n.conf, n, n.NumComponents())
	c.HasFocus = true

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	var bufferWidth int

	// NOTE: because we're ADDING a component
	// here we add 1 to the components since
	// we want to calculate the sizes _after_
	// we've added this component.
	numComponents := n.NumComponents() + 1
	if numComponents > 0 {
		bufferWidth = n.w / numComponents
	} else {
		bufferWidth = n.w
	}

	n.AddComponent(c)
	n.buffers[c.index] = c
	n.focusedBuff = c.index

	// translate all the components accordingly.
	for i, p := range n.components {
		if p == nil {
			continue
		}

		p.Resize(bufferWidth, n.h)
		p.SetPosition(bufferWidth*i, 0)
	}

	return c
}
