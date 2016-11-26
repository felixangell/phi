package gui

import (
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

var (
	cursor_flash_ms uint32 = 400
	reset_delay_ms  uint32 = 400

	// buffer
	TAB_SIZE int32 = 4
	HUNGRY_BACKSPACE bool = true
	TABS_ARE_SPACES bool = true	

	// cursor 
	should_draw  bool   = false
	should_flash bool   = true
	timer        uint32 = 0
	reset_timer  uint32 = 0
)

type Buffer struct {
	ComponentLocation
	font          *ttf.Font
	contents      []*rope.Rope
	curs          *Cursor
	input_handler *InputHandler
}

func NewBuffer() *Buffer {
	font, err := ttf.OpenFont("./res/firacode.ttf", 24)
	if err != nil {
		panic(err)
	}

	buff := &Buffer{
		contents: []*rope.Rope{},
		font:     font,
		curs:     &Cursor{},
	}
	buff.appendLine("This is a test.")
	return buff
}

func (b *Buffer) Dispose() {
	for _, texture := range TEXTURE_CACHE {
		texture.Destroy()
	}
}

func (b *Buffer) Init() {}

func (b *Buffer) GetComponents() []Component {
	return []Component{}
}

func (b *Buffer) AddComponent(c Component) {}

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
	// TODO: how the fuck do decode this properly?
	raw_val, size := utf8.DecodeLastRune(t.Text[:1])
	if raw_val == utf8.RuneError || size == 0 {
		return
	}

	b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string(raw_val))
	b.curs.move(1, 0)
}

func (b *Buffer) processActionKey(t *sdl.KeyDownEvent) {
	switch t.Keysym.Scancode {
	case sdl.SCANCODE_RETURN:
		initial_x := b.curs.x
		prev_line_len := b.contents[b.curs.y].Len()

		var new_rope *rope.Rope
		if initial_x < prev_line_len && initial_x > 0 {
			left, right := b.contents[b.curs.y].Split(initial_x)
			new_rope = right
			b.contents[b.curs.y] = left
		} else if initial_x == 0 {
			b.contents = append(b.contents, new(rope.Rope))      // grow
			copy(b.contents[b.curs.y+1:], b.contents[b.curs.y:]) // shift
			b.contents[b.curs.y] = new(rope.Rope)                // set
			b.curs.move(0, 1)
			return
		} else {
			new_rope = rope.New(" ")
		}

		b.curs.move(0, 1)
		for x := 0; x < initial_x; x++ {
			// TODO(Felix): there's a bug here where
			// this doesn't account for the rendered x
			// position when we use tabs as tabs and not spaces
			b.curs.move(-1, 0)
		}
		b.contents = append(b.contents, new_rope)
	case sdl.SCANCODE_BACKSPACE:
		if b.curs.x > 0 {
			offs := -1
			if !TABS_ARE_SPACES {
				if b.contents[b.curs.y].Index(b.curs.x) == '\t' {
					offs = int(-TAB_SIZE)
				}
			} else if HUNGRY_BACKSPACE && b.curs.x >= int(TAB_SIZE) && TABS_ARE_SPACES {
				// why x + 1 here? wtf
				if b.contents[b.curs.y].Substr((b.curs.x + 1) - int(TAB_SIZE), int(TAB_SIZE)).String() == "    " {
					// delete {TAB_SIZE} amount of characters
					// from the cursors x pos
					for i := 0; i < int(TAB_SIZE); i++ {
						b.contents[b.curs.y] = b.contents[b.curs.y].Delete(b.curs.x, 1)
						b.curs.move(-1, 0)
					}
					break
				}
			} 

			b.contents[b.curs.y] = b.contents[b.curs.y].Delete(b.curs.x, 1)
			b.curs.move_render(-1, 0, offs, 0)				
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
			prev_line_len := b.contents[b.curs.y-1].Len()
			b.contents[b.curs.y-1] = b.contents[b.curs.y-1].Concat(b.contents[b.curs.y])
			b.contents = append(b.contents[:b.curs.y], b.contents[b.curs.y+1:]...)
			b.curs.move(prev_line_len, -1)
		}
	case sdl.SCANCODE_RIGHT:
		curr_line_length := b.contents[b.curs.y].Len()
		if b.curs.x >= curr_line_length && b.curs.y < len(b.contents)-1 {
			// we're at the end of the line and we have
			// some lines after, let's wrap around
			b.curs.move(0, 1)
			b.curs.move(-curr_line_length, 0)
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
		if TABS_ARE_SPACES {
			// make an empty rune array of TAB_SIZE, cast to string
			// and insert it.
			b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, "    ")
			b.curs.move(int(TAB_SIZE), 0)
		} else {
			b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string('\t'))
			// the actual position is + 1, but we make it
			// move by TAB_SIZE characters on the view.
			b.curs.move_render(1, 0, int(TAB_SIZE), 0)
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

func (b *Buffer) Translate(x, y int32) {
	b.x += x
	b.y += y
}

func (b *Buffer) Update() {
	prev_x := b.curs.x
	prev_y := b.curs.y

	if b.input_handler == nil {
		panic("help")
	}

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
var TEXTURE_CACHE map[rune]*sdl.Texture = map[rune]*sdl.Texture{}

func (b *Buffer) Render(ctx *sdl.Renderer) {

	// render the ol' cursor
	if should_draw {
		gfx.SetDrawColorHex(ctx, 0x657B83)
		ctx.FillRect(&sdl.Rect{
			b.x + (int32(b.curs.rx)+1)*last_w,
			b.y + int32(b.curs.ry)*last_h,
			last_w,
			last_h,
		})
	}

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
				x_col += TAB_SIZE
				continue
			}

			x_col += 1

			text := renderString(b.font, string(char), gfx.HexColor(0x7a7a7a), true)
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
