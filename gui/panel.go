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

func (p *Panel) Dispose() {
	for _, comp := range p.components {
		comp.Dispose()
	}
}

func (p *Panel) Init() {}

func (p *Panel) Update() {
	for _, c := range p.components {
		c.Update()
	}
}

func (p *Panel) OnRender(ctx *sdl.Renderer) {
}
