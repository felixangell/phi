package gui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/vinzmay/go-rope"
)

type Buffer struct {
	x, y int
	font *ttf.Font
	contents *rope.Rope
}

func NewBuffer() *Buffer {
	font, err := ttf.OpenFont("./res/firacode.ttf", 14)
	if err != nil {
		panic(err)
	}

	return &Buffer{
		contents: rope.New("Hello, World!"),
		font: font,
	}
}

func (b *Buffer) Update() {

}

func (b *Buffer) Render(ctx *sdl.Surface) {
	text, _ := b.font.RenderUTF8_Solid(b.contents.String(), sdl.Color{255, 0, 255, 255})
	defer text.Free()

	text.Blit(nil, ctx, nil)
}
