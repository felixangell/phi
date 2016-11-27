package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type CommandBuffer struct {
	*Buffer
}

func (b *CommandBuffer) processActionKey(t *sdl.KeyDownEvent) {
	switch t.Keysym.Scancode {
	case sdl.SCANCODE_RETURN:
		println("dope!")
		return
	}
	b.Buffer.processActionKey(t)
}

type CommandPalette struct {
	BaseComponent
	buffer *CommandBuffer
}

func NewCommandPalette() *CommandPalette {
	palette := &CommandPalette{}
	palette.buffer = &CommandBuffer{
		Buffer: NewBuffer(nil),
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
