package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type CommandBuffer struct {
	*Buffer
}

func (b *CommandBuffer) processActionKey() {
	// TODO
	// b.Buffer.processActionKey()
}

type CommandPalette struct {
	BaseComponent
	buffer *CommandBuffer
}

func NewCommandPalette() *CommandPalette {
	palette := &CommandPalette{}
	palette.buffer = &CommandBuffer{
		Buffer: NewBuffer(nil, nil, 0),
	}
	palette.AddComponent(palette.buffer)
	return palette
}

func (p *CommandPalette) OnDispose() {}

func (p *CommandPalette) OnInit() {
}

func (c *CommandPalette) OnUpdate() {

}

func (c *CommandPalette) OnRender(ctx *sdl.Renderer) {

}
