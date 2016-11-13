package gui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/vinzmay/go-rope"
	"unicode/utf8"
)

type Cursor struct {
	x, y int32
	rx, ry int32
}

func (c *Cursor) move(x, y int32) {
	c.move_render(x, y, x, y)
}

// moves the cursors position, and the
// rendered coordinates by the given amount
func (c *Cursor) move_render(x, y, rx, ry int32) {
	c.x += x
	c.y += y

	c.rx += rx
	c.ry += ry
}

type Buffer struct {
	x, y int
	font *ttf.Font
	contents []*rope.Rope
	curs *Cursor
	input_handler *InputHandler
}

func NewBuffer() *Buffer {
	font, err := ttf.OpenFont("./res/firacode.ttf", 24)
	if err != nil {
		panic(err)
	}

	buff := &Buffer{
		contents: []*rope.Rope{},
		font: font,
		curs: &Cursor{},
	}
	buff.appendLine("This is a test.")
	return buff
}

func (b *Buffer) SetInputHandler(i *InputHandler) {
	b.input_handler = i
}

func (b *Buffer) GetInputHandler() *InputHandler {
	return b.input_handler
}

func (b *Buffer) appendLine(val string) {
	b.contents = append(b.contents, rope.New(val))
	b.curs.move(int32(len(val)), 0)
}

func (b *Buffer) processTextInput(t *sdl.TextInputEvent) {
	raw_val, _ := utf8.DecodeLastRune(t.Text[0:1])
	if raw_val == utf8.RuneError {
		return
	}

	b.contents[b.curs.y] = b.contents[b.curs.y].Concat(rope.New(string(raw_val)))
	b.curs.move(1, 0)
}

func (b *Buffer) Update() {
	if b.input_handler.Event != nil {
		switch t := b.input_handler.Event.(type) {
		case *sdl.TextInputEvent:
			b.processTextInput(t)
		}
	}
}

var last_w, last_h int32
func (b *Buffer) Render(ctx *sdl.Surface) {

	// render the ol' cursor
	ctx.FillRect(&sdl.Rect{
		(b.curs.rx + 1) * last_w, 
		b.curs.ry * last_h, 
		last_w, 
		last_h,
	}, 0xff00ff)

	var y_col int32
	for _, rope := range b.contents {
		
		var x_col int32
		for _, char := range rope.String() {
			switch char {
			case '\n':
				x_col = 0
				y_col += 1
				continue
			}

			x_col += 1

			text, _ := b.font.RenderUTF8_Solid(string(char), sdl.Color{0, 0, 0, 255})
			last_w = text.W
			last_h = text.H

			defer text.Free()
			text.Blit(nil, ctx, &sdl.Rect{
				(x_col * text.W), 
				(y_col * text.H), 
				text.W, 
				text.H,
			})
		}
	}
}
