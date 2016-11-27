package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type CommandPalette struct {
	BaseComponent
	buff *Buffer
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

func (c *CommandPalette) Update() {
	for _, c := range c.components {
		c.Update()
	}
}

func (c *CommandPalette) OnRender(ctx *sdl.Renderer) {}
