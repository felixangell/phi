package gui

import (
	"github.com/felixangell/nate/cfg"
	"github.com/felixangell/nate/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/vinzmay/go-rope"
	"unicode/utf8"
)

type Cursor struct {
	x, y   int
	rx, ry int
}

func (c *Cursor) move(x, y int) {
	c.move_render(x, y, x, y)
}

// moves the cursors position, and the
// rendered coordinates by the given amount
func (c *Cursor) move_render(x, y, rx, ry int) {
	c.x += x
	c.y += y

	c.rx += rx
	c.ry += ry
}

const (
	cursor_flash_ms uint32 = 400
	reset_delay_ms  uint32 = 400
)

var (
	should_draw  bool   = false
	should_flash bool   = true
	timer        uint32 = 0
	reset_timer  uint32 = 0
)

type Buffer struct {
	x, y          int
	font          *ttf.Font
	contents      []*rope.Rope
	curs          *Cursor
	input_handler *InputHandler
	cfg           *cfg.TomlConfig
}

func NewBuffer(conf *cfg.TomlConfig) *Buffer {
	font, err := ttf.OpenFont("./res/firacode.ttf", 24)
	if err != nil {
		panic(err)
	}

	buff := &Buffer{
		contents: []*rope.Rope{},
		font:     font,
		curs:     &Cursor{},
		cfg:      conf,
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
	b.curs.move(len(val), 0)
}

func (b *Buffer) processTextInput(t *sdl.TextInputEvent) {
	// how the fuck do we decode a 32 byte array of junk?

	raw_val, size := utf8.DecodeLastRune(t.Text[1:5])
	if raw_val == utf8.RuneError || size == 0 {
		return
	}

	b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string(raw_val))
	b.curs.move(1, 0)
}

func (b *Buffer) processActionKey(t *sdl.KeyDownEvent) {
	switch t.Keysym.Scancode {
	case sdl.SCANCODE_RETURN:
		line_len := -b.contents[b.curs.y].Len()
		b.curs.move(line_len, 1)
		b.contents = append(b.contents, rope.New(" "))
	case sdl.SCANCODE_BACKSPACE:
		if b.curs.x > 0 {
			b.contents[b.curs.y] = b.contents[b.curs.y].Delete(b.curs.x, 1)
			b.curs.move(-1, 0)
		}
	}
}

func renderString(font *ttf.Font, val string, col sdl.Color, smooth bool) *sdl.Surface {
	if smooth {
		text, err := font.RenderUTF8_Blended(val, col)
		if err != nil {
			panic(err)
		}
		return text
	} else {
		text, err := font.RenderUTF8_Solid(val, col)
		if err != nil {
			panic(err)
		}
		return text
	}
	return nil
}

func (b *Buffer) Update() {
	prev_x := b.curs.x
	prev_y := b.curs.y

	if b.input_handler.Event != nil {
		switch t := b.input_handler.Event.(type) {
		case *sdl.TextInputEvent:
			b.processTextInput(t)
		case *sdl.KeyDownEvent:
			b.processActionKey(t)
		}
	}

	if b.curs.x != prev_x || b.curs.y != prev_y {
		should_draw = true
		should_flash = false
		reset_timer = sdl.GetTicks()
	}

	if !should_flash && sdl.GetTicks()-reset_timer > reset_delay_ms {
		should_flash = true
	}

	if sdl.GetTicks()-timer > cursor_flash_ms && should_flash {
		timer = sdl.GetTicks()
		should_draw = !should_draw
	}
}

var last_w, last_h int32

func (b *Buffer) Render(ctx *sdl.Renderer) {

	// render the ol' cursor
	if should_draw {
		gfx.SetDrawColorHex(ctx, 0x657B83)
		ctx.FillRect(&sdl.Rect{
			(int32(b.curs.rx) + 1) * last_w,
			int32(b.curs.ry) * last_h,
			last_w,
			last_h,
		})
	}

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

			text := renderString(b.font, string(char), gfx.HexColor(0x7a7a7a), b.cfg.Editor.Aliased)
			defer text.Free()

			last_w = text.W
			last_h = text.H

			// FIXME very slow
			texture, _ := ctx.CreateTextureFromSurface(text)
			defer texture.Destroy()

			ctx.Copy(texture, nil, &sdl.Rect{
				(x_col * text.W),
				(y_col * text.H),
				text.W,
				text.H,
			})
		}

		y_col += 1
	}
}
