package gui

import (
	"strings"
	"unicode/utf8"

	"github.com/felixangell/nate/cfg"
	"github.com/felixangell/nate/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/vinzmay/go-rope"
)

var (
	timer        uint32 = 0
	reset_timer  uint32 = 0
	should_draw  bool   = true
	should_flash bool
)

type Buffer struct {
	BaseComponent
	font     *ttf.Font
	contents []*rope.Rope
	curs     *Cursor
	cfg      *cfg.TomlConfig
}

func NewBuffer(conf *cfg.TomlConfig) *Buffer {
	font, err := ttf.OpenFont("./res/firacode.ttf", 24)
	if err != nil {
		panic(err)
	}

	config := conf
	if config == nil {
		config = cfg.NewDefaultConfig()
	}

	buff := &Buffer{
		contents: []*rope.Rope{},
		font:     font,
		curs:     &Cursor{},
		cfg:      config,
	}
	buff.appendLine("This is a test 世界.")
	return buff
}

func (b *Buffer) Dispose() {
	for _, texture := range TEXTURE_CACHE {
		texture.Destroy()
	}
}

func (b *Buffer) Init() {}

func (b *Buffer) appendLine(val string) {
	b.contents = append(b.contents, rope.New(val))
	b.curs.move(len(val), 0)
}

func (b *Buffer) processTextInput(t *sdl.TextInputEvent) {
	// TODO: how the fuck do decode this properly?
	rawVal, size := utf8.DecodeLastRune(t.Text[:1])
	if rawVal == utf8.RuneError || size == 0 {
		return
	}

	b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string(rawVal))
	b.curs.move(1, 0)
}

func (b *Buffer) processActionKey(t *sdl.KeyDownEvent) {
	switch t.Keysym.Scancode {
	case sdl.SCANCODE_RETURN:
		initial_x := b.curs.x
		prevLineLen := b.contents[b.curs.y].Len()

		var newRope *rope.Rope
		if initial_x < prevLineLen && initial_x > 0 {
			left, right := b.contents[b.curs.y].Split(initial_x)
			newRope = right
			b.contents[b.curs.y] = left
		} else if initial_x == 0 {
			b.contents = append(b.contents, new(rope.Rope))      // grow
			copy(b.contents[b.curs.y+1:], b.contents[b.curs.y:]) // shift
			b.contents[b.curs.y] = new(rope.Rope)                // set
			b.curs.move(0, 1)
			return
		} else {
			newRope = rope.New(" ")
		}

		b.curs.move(0, 1)
		for x := 0; x < initial_x; x++ {
			// TODO(Felix): there's a bug here where
			// this doesn't account for the rendered x
			// position when we use tabs as tabs and not spaces
			b.curs.move(-1, 0)
		}
		b.contents = append(b.contents, newRope)
	case sdl.SCANCODE_BACKSPACE:
		if b.curs.x > 0 {
			offs := -1
			if !b.cfg.Editor.Tabs_Are_Spaces {
				if b.contents[b.curs.y].Index(b.curs.x) == '\t' {
					offs = int(-b.cfg.Editor.Tab_Size)
				}
			} else if b.cfg.Editor.Hungry_Backspace && b.curs.x >= int(b.cfg.Editor.Tab_Size) {
				// FIXME wtf how does Substr even work
				// cut out the last {TAB_SIZE} amount of characters
				// and check em
				tabSize := int(b.cfg.Editor.Tab_Size)
				lastTabSizeChars := b.contents[b.curs.y].Substr(b.curs.x+1-tabSize, tabSize).String()
				artificialTab := string(make([]rune, tabSize, ' '))
				if strings.Compare(lastTabSizeChars, artificialTab) == 0 {
					// delete {TAB_SIZE} amount of characters
					// from the cursors x pos
					for i := 0; i < int(b.cfg.Editor.Tab_Size); i++ {
						b.contents[b.curs.y] = b.contents[b.curs.y].Delete(b.curs.x, 1)
						b.curs.move(-1, 0)
					}
					break
				}
			}

			b.contents[b.curs.y] = b.contents[b.curs.y].Delete(b.curs.x, 1)
			b.curs.moveRender(-1, 0, offs, 0)
		} else if b.curs.x == 0 && b.curs.y > 0 {
			// start of line, wrap to previous
			// two cases here:

			// the line_len is zero, in which case
			// we delete the line and go to the end
			// of the previous line
			if b.contents[b.curs.y].Len() == 0 {
				b.curs.move(b.contents[b.curs.y-1].Len(), -1)
				// FIXME, delete from the curs.y dont pop!
				b.contents = b.contents[:len(b.contents)-1]
				return
			}

			// TODO(Felix): handle all the edge cases here...

			// or, the line has characters, so we join
			// that line with the previous line
			prevLineLen := b.contents[b.curs.y-1].Len()
			b.contents[b.curs.y-1] = b.contents[b.curs.y-1].Concat(b.contents[b.curs.y])
			b.contents = append(b.contents[:b.curs.y], b.contents[b.curs.y+1:]...)
			b.curs.move(prevLineLen, -1)
		}
	case sdl.SCANCODE_RIGHT:
		currLineLength := b.contents[b.curs.y].Len()
		if b.curs.x >= currLineLength && b.curs.y < len(b.contents)-1 {
			// we're at the end of the line and we have
			// some lines after, let's wrap around
			b.curs.move(0, 1)
			b.curs.move(-currLineLength, 0)
		} else if b.curs.x < b.contents[b.curs.y].Len() {
			// we have characters to the right, let's move along
			b.curs.move(1, 0)
		}
	case sdl.SCANCODE_LEFT:
		if b.curs.x == 0 && b.curs.y > 0 {
			b.curs.move(b.contents[b.curs.y-1].Len(), -1)

		} else if b.curs.x > 0 {
			b.curs.move(-1, 0)
		}
	case sdl.SCANCODE_TAB:
		if b.cfg.Editor.Tabs_Are_Spaces {
			// make an empty rune array of TAB_SIZE, cast to string
			// and insert it.
			b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, "    ")
			b.curs.move(int(b.cfg.Editor.Tab_Size), 0)
		} else {
			b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string('\t'))
			// the actual position is + 1, but we make it
			// move by TAB_SIZE characters on the view.
			b.curs.moveRender(1, 0, int(b.cfg.Editor.Tab_Size), 0)
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

	if b.inputHandler == nil {
		panic("help")
	}

	if b.inputHandler.Event != nil {
		switch t := b.inputHandler.Event.(type) {
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

	if !should_flash && sdl.GetTicks()-reset_timer > b.cfg.Editor.Cursor_Reset_Delay {
		should_flash = true
	}

	if sdl.GetTicks()-timer > b.cfg.Editor.Cursor_Flash_Rate && (should_flash && b.cfg.Editor.Flash_Cursor) {
		timer = sdl.GetTicks()
		should_draw = !should_draw
	}
}

var last_w, last_h int32
var TEXTURE_CACHE map[rune]*sdl.Texture = map[rune]*sdl.Texture{}

func (b *Buffer) OnRender(ctx *sdl.Renderer) {

	// render the ol' cursor
	if should_draw && b.cfg.Editor.Draw_Cursor {
		gfx.SetDrawColorHex(ctx, 0x657B83)
		ctx.FillRect(&sdl.Rect{
			b.x + (int32(b.curs.rx)+1)*last_w,
			b.y + int32(b.curs.ry)*last_h,
			last_w,
			last_h,
		})
	}

	// TODO(Felix): cull this so that
	// we dont render lines we cant see.

	var y_col int32
	for _, rope := range b.contents {
		// this is because if we had the following
		// text input:
		//
		// Foo
		// _			<-- underscore is a space!
		// Blah
		// and we delete that underscore... it causes
		// a panic because there are no characters in
		// the empty string!
		if rope.Len() == 0 {
			// even though the string is empty
			// we still need to offset it by a line
			y_col += 1
			continue
		}

		var x_col int32
		for _, char := range rope.String() {
			switch char {
			case '\n':
				x_col = 0
				y_col += 1
				continue
			case '\t':
				x_col += b.cfg.Editor.Tab_Size
				continue
			}

			char_colour := 0x7a7a7a
			if b.curs.x == int(x_col) && b.curs.y == int(y_col) {
				char_colour = 0xff00ff
			}

			x_col += 1

			text := renderString(b.font, string(char), gfx.HexColor(uint32(char_colour)), true)
			defer text.Free()

			last_w = text.W
			last_h = text.H

			texture, ok := TEXTURE_CACHE[char]
			if !ok {
				// can't find it in the cache so we
				// load and then cache it.
				texture, _ = ctx.CreateTextureFromSurface(text)
				TEXTURE_CACHE[char] = texture
			}

			// FIXME still kinda slow
			// we can also cull so that
			// we don't render things that aren't
			// visible outside of the component
			ctx.Copy(texture, nil, &sdl.Rect{
				b.x + (x_col * text.W),
				b.y + (y_col * text.H),
				text.W,
				text.H,
			})
		}

		y_col += 1
	}
}
