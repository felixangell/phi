package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Panel struct {
	components []Component
	input_handler *InputHandler
}

func NewPanel(input *InputHandler) *Panel {
	return &Panel{
		components: []Component{},
		input_handler: input,
	}
}

func (p *Panel) AddComponent(c Component) {
	p.components = append(p.components, c)
	c.SetInputHandler(p.input_handler)
}

func (p *Panel) SetInputHandler(i *InputHandler) {
	p.input_handler = i
}

func (p *Panel) GetInputHandler() *InputHandler {
	return p.input_handler
}

func (p *Panel) Update() {
	for _, c := range p.components {
		c.Update()
	}
}

func (p *Panel) Render(ctx *sdl.Renderer) {
	for _, c := range p.components {
		c.Render(ctx)
	}
}
