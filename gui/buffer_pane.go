package gui

import (
	"github.com/felixangell/strife"
)

type BufferPane struct {
	BaseComponent
	Buff *Buffer
}

func NewBufferPane(buff *Buffer) *BufferPane {
	return &BufferPane{
		BaseComponent{},
		buff,
	}
}

func (b *BufferPane) OnUpdate() bool {
	b.Buff.processInput(nil)
	return b.Buff.OnUpdate()
}

func (b *BufferPane) OnRender(ctx *strife.Renderer) {
	b.Buff.OnRender(ctx)
}
