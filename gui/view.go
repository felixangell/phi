package gui

import (
	"github.com/felixangell/phi-editor/cfg"
	"github.com/felixangell/strife"
)

type View struct {
	BaseComponent
	conf *cfg.TomlConfig
}

func NewView(width, height int, conf *cfg.TomlConfig) *View {
	view := &View{conf: conf}
	view.Translate(width, height)
	view.Resize(width, height)
	return view
}

func (n *View) OnInit() {
}

func (n *View) OnUpdate() bool {
	return false
}

func (n *View) OnRender(ctx *strife.Renderer) {}

func (n *View) OnDispose() {}

func (n *View) AddBuffer() *Buffer {
	c := NewBuffer(n.conf, n, n.NumComponents())

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

	// setup and add the panel for the buffer
	panel := NewPanel(n.inputHandler)
	c.SetInputHandler(n.inputHandler)

	panel.AddComponent(c)
	n.AddComponent(panel)

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
