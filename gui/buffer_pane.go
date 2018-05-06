package gui

import (
	"fmt"
	"github.com/felixangell/strife"
)

var metaPanelHeight = 32

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

func (b *BufferPane) renderMetaPanel(ctx *strife.Renderer) {
	pad := 10
	mpY := (b.y + b.h) - (metaPanelHeight)

	// panel backdrop
	ctx.SetColor(strife.Black)
	ctx.Rect(b.x, mpY, b.w, metaPanelHeight, strife.Fill)

	{
		infoLine := fmt.Sprintf("Line %d, Column %d", b.Buff.curs.y, b.Buff.curs.x)
		ctx.SetColor(strife.White)
		_, strHeight := ctx.String(infoLine, b.x+(pad/2), mpY+(pad/2))
		metaPanelHeight = strHeight + pad
	}

	// resize to match new height if any
	b.Buff.Resize(b.w, b.h-metaPanelHeight)
}

func (b *BufferPane) Resize(w, h int) {
	b.BaseComponent.Resize(w, h)
	b.Buff.Resize(w, h)
}

func (b *BufferPane) SetPosition(x, y int) {
	b.BaseComponent.SetPosition(x, y)
	b.Buff.SetPosition(x, y)
}

func (b *BufferPane) OnUpdate() bool {
	b.Buff.processInput(nil)
	return b.Buff.OnUpdate()
}

func (b *BufferPane) OnRender(ctx *strife.Renderer) {
	b.Buff.OnRender(ctx)
	b.renderMetaPanel(ctx)
}
