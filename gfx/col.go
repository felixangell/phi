package gfx

import "github.com/veandco/go-sdl2/sdl"

func SetDrawColorHex(ctx *sdl.Renderer, col uint32) {
	a := uint8(255)
	r := uint8(col&0xff0000 >> 16)
	g := uint8(col&0xff00 >> 8)
	b := uint8(col&0xff)
	ctx.SetDrawColor(r, g, b, a)
}

func HexColor(col uint32) sdl.Color {
	a := uint8(255)
	r := uint8(col&0xff0000 >> 16)
	g := uint8(col&0xff00 >> 8)
	b := uint8(col&0xff)
	return sdl.Color{r, g, b, a}
}
