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

var lastWidth int

func (b *BufferPane) renderMetaPanel(ctx *strife.Renderer) {
	conf := b.Buff.cfg.Theme.Palette

	pad := 6
	mpY := (b.y + b.h) - (metaPanelHeight)

	// panel backdrop
	ctx.SetColor(strife.HexRGB(conf.Suggestion.Background))
	ctx.Rect(b.x, mpY, b.w, metaPanelHeight, strife.Fill)

	// tab info etc. on right hand side
	{
		tabSize := b.Buff.cfg.Editor.Tab_Size
		syntaxName := "Undefined"

		infoLine := fmt.Sprintf("Tab Size: %d    Syntax: %s", tabSize, syntaxName)
		ctx.SetColor(strife.HexRGB(conf.Suggestion.Foreground))

		lastWidth, _ = ctx.String(infoLine, ((b.x + b.w) - (lastWidth + (pad))), mpY+(pad/2)+1)
	}

	{
		modified := ' '
		if b.Buff.modified {
			modified = '*'
		}

		infoLine := fmt.Sprintf("%s%c Line %d, Column %d", b.Buff.filePath, modified, b.Buff.curs.y, b.Buff.curs.x)
		ctx.SetColor(strife.HexRGB(conf.Suggestion.Foreground))
		_, strHeight := ctx.String(infoLine, b.x+pad, mpY+(pad/2)+1)
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
