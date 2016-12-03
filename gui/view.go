package gui

import (
	"github.com/felixangell/nate/cfg"
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
	n.addBuffer()
}

func (n *View) OnUpdate() {}

func (n *View) OnRender(ctx *strife.Renderer) {}

func (n *View) OnDispose() {}

func (n *View) addBuffer() {
	c := NewBuffer(n.conf)

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	bufferWidth := n.w / (len(n.components) + 1)
	c.Resize(bufferWidth, n.h)

	// setup and add the panel for the buffer
	panel := NewPanel(n.inputHandler)
	c.SetInputHandler(n.inputHandler)
	panel.AddComponent(c)
	n.components = append(n.components, panel)

	// translate all the components accordingly.
	for i, p := range n.components {
		p.Translate(bufferWidth*i, 0)
	}
}
