package gui

import (
	"github.com/felixangell/phi-editor/cfg"
	"github.com/felixangell/strife"
)

type CommandPalette struct {
	BaseComponent
	buff *Buffer
}

func NewCommandPalette(conf *cfg.TomlConfig) *CommandPalette {
	palette := &CommandPalette{
		buff: NewBuffer(conf, nil, 0),
	}
	palette.buff.HasFocus = true
	return palette
}

func (b *CommandPalette) OnInit() {

}

func (b *CommandPalette) OnUpdate() bool {
	return b.buff.OnUpdate()
}

func (b *CommandPalette) OnRender(ctx *strife.Renderer) {
	ctx.SetColor(strife.Red)
	ctx.Rect(b.x, b.y, b.w, b.h, strife.Fill)

	b.buff.OnRender(ctx)
}

func (b *CommandPalette) OnDispose() {

}
