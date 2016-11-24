package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Panel struct {
	ComponentLocation

	components    []Component
	input_handler *InputHandler
}

func (p *Panel) Translate(x, y int32) {
	p.x += x
	p.y += y
	for _, c := range p.components {
		c.Translate(x, y)
	}
}

func NewPanel(input *InputHandler) *Panel {
	return &Panel{
		components:    []Component{},
		input_handler: input,
	}
}

func (p *Panel) Init() {}

func (p *Panel) GetComponents() []Component {
	return p.components
}

func (p *Panel) AddComponent(c Component) {
	c.Init()
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
