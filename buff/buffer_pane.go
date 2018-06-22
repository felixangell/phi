package buff

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/phi/gui"
	"github.com/felixangell/strife"
)

var metaPanelHeight = 32

type BufferPane struct {
	gui.BaseComponent
	Buff *Buffer
	font *strife.Font
}

func NewBufferPane(buff *Buffer) *BufferPane {
	fontPath := filepath.Join(cfg.FONT_FOLDER, buff.cfg.Editor.Font_Face+".ttf")
	metaPanelFont, err := strife.LoadFont(fontPath, 14)
	if err != nil {
		log.Println("Note: failed to load meta panel font ", fontPath)
		metaPanelFont = buff.buffOpts.font
	}

	return &BufferPane{
		gui.BaseComponent{},
		buff,
		metaPanelFont,
	}
}

var lastWidth int

func (b *BufferPane) SetFocus(focus bool) {
	b.Buff.SetFocus(focus)
	b.BaseComponent.SetFocus(focus)
}

func (b *BufferPane) renderMetaPanel(ctx *strife.Renderer) {
	conf := b.Buff.cfg.Theme.Palette

	x, y := b.GetPos()
	w, h := b.GetSize()

	pad := 6
	mpY := (y + h) - (metaPanelHeight)

	focused := b.Buff.index == b.Buff.parent.focusedBuff

	colour := strife.HexRGB(conf.Suggestion.Background)

	if focused {
		nr := int(colour.R) + 10
		ng := int(colour.G) + 10
		nb := int(colour.B) + 10
		colour = strife.RGB(nr, ng, nb)
	}

	// panel backdrop
	ctx.SetColor(colour)
	ctx.Rect(x, mpY, w, metaPanelHeight, strife.Fill)

	// tab info etc. on right hand side
	{
		tabSize := b.Buff.cfg.Editor.Tab_Size

		// TODO
		syntaxName := "Undefined"

		infoLine := fmt.Sprintf("Tab Size: %d    Syntax: %s", tabSize, syntaxName)
		ctx.SetColor(strife.HexRGB(conf.Suggestion.Foreground))

		ctx.SetFont(b.font)
		lastWidth, _ = ctx.Text(infoLine, ((x + w) - (lastWidth + (pad))), mpY+(pad/2))
	}

	{
		modified := ' '
		if b.Buff.modified {
			modified = '*'
		}

		infoLine := fmt.Sprintf("%s%c Line %d, Column %d", b.Buff.filePath, modified, b.Buff.curs.y+1, b.Buff.curs.x)

		if cfg.DebugMode {
			infoLine = fmt.Sprintf("%s, BuffIndex: %d", infoLine, b.Buff.index)
		}

		ctx.SetColor(strife.HexRGB(conf.Suggestion.Foreground))

		ctx.SetFont(b.font)
		_, strHeight := ctx.Text(infoLine, x+pad, mpY+(pad/2)+1)
		metaPanelHeight = strHeight + pad
	}

	// resize to match new height if any
	b.Buff.Resize(w, h-metaPanelHeight)
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
