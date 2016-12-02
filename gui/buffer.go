package gui

import (
	"io/ioutil"
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

var TEXTURE_CACHE map[rune]*sdl.Texture = map[rune]*sdl.Texture{}

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

	buffContents := []*rope.Rope{}
	buff := &Buffer{
		contents: buffContents,
		font:     font,
		curs:     &Cursor{},
		cfg:      config,
	}

	{
		contents, err := ioutil.ReadFile(cfg.CONFIG_FULL_PATH)
		if err != nil {
			panic(err)
		}

		lines := strings.Split(string(contents), "\n")
		for _, line := range lines {
			buff.appendLine(line)
		}
	}

	return buff
}

func (b *Buffer) OnDispose() {
	for _, texture := range TEXTURE_CACHE {
		texture.Destroy()
	}
}

func (b *Buffer) OnInit() {}

func (b *Buffer) appendLine(val string) {
	b.contents = append(b.contents, rope.New(val))
	b.curs.move(len(val), 0)
}

func (b *Buffer) processTextInput(t *sdl.TextInputEvent) {
	firstRune := t.Text[:4]
	r, _ := utf8.DecodeRune(firstRune)
	if r == utf8.RuneError {
		panic("oh dear!")
	}

	b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string(r))
	b.curs.move(1, 0)

	// we don't need to match braces
	// let's not continue any further
	if !b.cfg.Editor.Match_Braces {
		return
	}

	matchingPair := int(r)

	// the offset in the ASCII Table is +2 for { and for [
	// but its +1 for parenthesis (
	offset := 2

	switch r {
	case '(':
		offset = 1
		fallthrough
	case '{':
		fallthrough
	case '[':
		matchingPair += offset
		b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string(rune(matchingPair)))
	}
}

// TODO(Felix): refactor me!)
func (b *Buffer) processActionKey(t *sdl.KeyDownEvent) {
	switch t.Keysym.Scancode {
	case sdl.SCANCODE_RETURN:
		initial_x := b.curs.x
		prevLineLen := b.contents[b.curs.y].Len()

		var newRope *rope.Rope
		if initial_x < prevLineLen && initial_x > 0 {
			// we're not at the end of the line, but we're not at
			// the start, i.e. we're SPLITTING the line
			left, right := b.contents[b.curs.y].Split(initial_x)
			newRope = right
			b.contents[b.curs.y] = left
		} else if initial_x == 0 {
			// we're at the start of a line, so we want to
			// shift the line down and insert an empty line
			// above it!
			b.contents = append(b.contents, new(rope.Rope))      // grow
			copy(b.contents[b.curs.y+1:], b.contents[b.curs.y:]) // shift
			b.contents[b.curs.y] = new(rope.Rope)                // set
			b.curs.move(0, 1)
			return
		} else {
			// we're at the end of a line
			newRope = new(rope.Rope)
		}

		b.curs.move(0, 1)
		for x := 0; x < initial_x; x++ {
			// TODO(Felix): there's a bug here where
			// this doesn't account for the rendered x
			// position when we use tabs as tabs and not spaces
			b.curs.move(-1, 0)
		}

		b.contents = append(b.contents, nil)
		copy(b.contents[b.curs.y+1:], b.contents[b.curs.y:])
		b.contents[b.curs.y] = newRope
	case sdl.SCANCODE_BACKSPACE:
		if b.curs.x > 0 {
			offs := -1
			if !b.cfg.Editor.Tabs_Are_Spaces {
				if b.contents[b.curs.y].Index(b.curs.x) == '\t' {
					offs = int(-b.cfg.Editor.Tab_Size)
				}
			} else if b.cfg.Editor.Hungry_Backspace && b.curs.x >= int(b.cfg.Editor.Tab_Size) {
				// cut out the last {TAB_SIZE} amount of characters
				// and check em
				tabSize := int(b.cfg.Editor.Tab_Size)
				lastTabSizeChars := b.contents[b.curs.y].Substr(b.curs.x+1-tabSize, tabSize).String()
				if strings.Compare(lastTabSizeChars, b.makeTab()) == 0 {
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
	case sdl.SCANCODE_UP:
		if b.curs.y > 0 {
			offs := 0
			prevLineLen := b.contents[b.curs.y-1].Len()
			if b.curs.x > prevLineLen {
				offs = prevLineLen - b.curs.x
			}
			// TODO: offset should account for tabs
			b.curs.move(offs, -1)
		}
	case sdl.SCANCODE_DOWN:
		if b.curs.y < len(b.contents)-1 {
			offs := 0
			nextLineLen := b.contents[b.curs.y+1].Len()
			if b.curs.x > nextLineLen {
				offs = nextLineLen - b.curs.x
			}
			// TODO: offset should account for tabs
			b.curs.move(offs, 1)
		}
	case sdl.SCANCODE_TAB:
		if b.cfg.Editor.Tabs_Are_Spaces {
			// make an empty rune array of TAB_SIZE, cast to string
			// and insert it.
			b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, b.makeTab())
			b.curs.move(int(b.cfg.Editor.Tab_Size), 0)
		} else {
			b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string('\t'))
			// the actual position is + 1, but we make it
			// move by TAB_SIZE characters on the view.
			b.curs.moveRender(1, 0, int(b.cfg.Editor.Tab_Size), 0)
		}
	}
}

// TODO(Felix) this is really stupid
func (b *Buffer) makeTab() string {
	blah := []rune{}
	for i := 0; i < int(b.cfg.Editor.Tab_Size); i++ {
		blah = append(blah, ' ')
	}
	return string(blah)
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

func (b *Buffer) OnUpdate() {
	prev_x := b.curs.x
	prev_y := b.curs.y

	// FIXME handle focus properly
	if b.inputHandler == nil {
		return
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

	if !should_flash && sdl.GetTicks()-reset_timer > b.cfg.Cursor.Reset_Delay {
		should_flash = true
	}

	if sdl.GetTicks()-timer > b.cfg.Cursor.Flash_Rate && (should_flash && b.cfg.Cursor.Flash) {
		timer = sdl.GetTicks()
		should_draw = !should_draw
	}
}

// dimensions of the last character we rendered
var last_w, last_h int32
var lineIndex int = 0

func (b *Buffer) OnRender(ctx *sdl.Renderer) {
	gfx.SetDrawColorHexString(ctx, b.cfg.Theme.Background)
	ctx.FillRect(&sdl.Rect{b.x, b.y, b.w, b.h})

	if b.cfg.Editor.Highlight_Line {
		gfx.SetDrawColorHexString(ctx, "0x001629")
		ctx.FillRect(&sdl.Rect{
			b.x,
			b.y + int32(b.curs.ry)*last_h,
			b.w,
			last_h,
		})
	}

	// render the ol' cursor
	if should_draw && b.cfg.Cursor.Draw {
		cursorWidth := int32(b.cfg.Cursor.GetCaretWidth())
		if cursorWidth == -1 {
			cursorWidth = last_w
		}

		gfx.SetDrawColorHexString(ctx, b.cfg.Theme.Cursor)
		ctx.FillRect(&sdl.Rect{
			b.x + (int32(b.curs.rx))*last_w,
			b.y + int32(b.curs.ry)*last_h,
			cursorWidth,
			last_h,
		})
	}

	source := b.contents
	if int(last_h) > 0 && int(b.h) != 0 {
		// work out how many lines can fit into
		// the buffer, and set the source to
		// slice the line buffer accordingly
		visibleLines := int(b.h) / int(last_h)
		if len(b.contents) > visibleLines {
			if lineIndex+visibleLines >= len(b.contents) {
				lineIndex = len(b.contents) - visibleLines
			}
			source = b.contents[lineIndex : lineIndex+visibleLines]
		}
	}

	var y_col int32
	for _, rope := range source {
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

			x_col += 1

			texture, ok := TEXTURE_CACHE[char]
			if !ok {
				text := renderString(b.font, string(char), gfx.HexColorString(b.cfg.Theme.Foreground), b.cfg.Render.Aliased)
				last_w = text.W
				last_h = text.H

				// can't find it in the cache so we
				// load and then cache it.
				texture, _ = ctx.CreateTextureFromSurface(text)
				TEXTURE_CACHE[char] = texture
				text.Free()
			}

			// set the colour of the currently selected
			// character to white IF the cursor is being
			// drawn
			source, allocated := texture, false
			if b.curs.x+1 == int(x_col) && b.curs.y == int(y_col) && should_draw {
				text := renderString(b.font, string(char), gfx.HexColorString(b.cfg.Theme.Cursor_Invert), b.cfg.Render.Aliased)
				last_w = text.W
				last_h = text.H

				texture, _ = ctx.CreateTextureFromSurface(text)
				text.Free()
				source, allocated = texture, true
			}

			ctx.Copy(source, nil, &sdl.Rect{
				b.x + ((x_col - 1) * last_w),
				b.y + (y_col * last_h),
				last_w,
				last_h,
			})

			// it's not cached so we have to free it
			// ourselves
			if allocated {
				source.Destroy()
			}
		}

		y_col += 1
	}

}
