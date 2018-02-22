package gui

import (
	"github.com/felixangell/strife"
)

type Panel struct {
	BaseComponent
}

func NewPanel(input *InputHandler) *Panel {
	panel := &Panel{}
	panel.SetInputHandler(input)
	return panel
}

func (p *Panel) OnDispose() {}

func (p *Panel) OnInit() {}

func (p *Panel) OnUpdate() bool {
	return false
}

func (p *Panel) OnRender(ctx *strife.Renderer) {}
