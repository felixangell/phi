package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type CommandPalette struct {
	*ComponentLocation

	components    []Component
	input_handler *InputHandler
	buff          *Buffer
}

func NewCommandPalette() *CommandPalette {
	palette := &CommandPalette{
		buff: NewBuffer(nil),
	}
	return palette
}

func (p *CommandPalette) Dispose() {
	for _, comp := range p.components {
		comp.Dispose()
	}
}

func (p *CommandPalette) Init() {
	p.buff.SetInputHandler(p.input_handler)
	p.AddComponent(p.buff)
}

func (p *CommandPalette) AddComponent(c Component) {
	p.components = append(p.components, c)
	c.SetInputHandler(p.input_handler)
	c.Init()
}

func (p *CommandPalette) GetComponents() []Component {
	return p.components
}

func (c *CommandPalette) Update() {
	for _, c := range c.components {
		c.Update()
	}
}

func (c *CommandPalette) Render(ctx *sdl.Renderer) {
	for _, c := range c.components {
		c.Render(ctx)
	}
}

func (c *CommandPalette) GetInputHandler() *InputHandler {
	return c.input_handler
}

func (c *CommandPalette) SetInputHandler(h *InputHandler) {
	c.input_handler = h
}
