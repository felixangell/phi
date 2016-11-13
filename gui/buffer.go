package gui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/vinzmay/go-rope"
)

type Cursor struct {
	x, y int32
	rx, ry int32
}

type Buffer struct {
	x, y int
	font *ttf.Font
	contents []*rope.Rope
	curs *Cursor
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
	buff.appendLine("Hello, World!\nHello!")
	return buff
}

func (b *Buffer) appendLine(val string) {
	b.contents = append(b.contents, rope.New(val))
}

func (b *Buffer) Update() {

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
