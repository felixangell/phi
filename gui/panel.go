package gui

import (
	"github.com/veandco/go-sdl2/sdl"
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

func (p *Panel) OnUpdate() {}

func (p *Panel) OnRender(ctx *sdl.Renderer) {}
