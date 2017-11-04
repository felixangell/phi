package gui

import (
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"github.com/felixangell/nate/cfg"
	"github.com/felixangell/strife"
	"github.com/vinzmay/go-rope"
)

var (
	timer        int64 = 0
	reset_timer  int64 = 0
	should_draw  bool  = true
	should_flash bool
)

// TODO: allow font setting or whatever

type Buffer struct {
	BaseComponent
	font     *strife.Font
	contents []*rope.Rope
	curs     *Cursor
	cfg      *cfg.TomlConfig
}

func NewBuffer(conf *cfg.TomlConfig) *Buffer {
	config := conf
	if config == nil {
		config = cfg.NewDefaultConfig()
	}

	buffContents := []*rope.Rope{}
	buff := &Buffer{
		contents: buffContents,
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

}

func (b *Buffer) OnInit() {}

func (b *Buffer) appendLine(val string) {
	b.contents = append(b.contents, rope.New(val))
	b.curs.move(len(val), 0)
}

func (b *Buffer) processTextInput() {
	firstRune := []byte{'a'}
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

func (b *Buffer) processActionKey() {
	// we dont process key input so
	// these are a bunch of stupid dummy case statements but
	// the logic for them should work! we just need to handle
	// key events which hasn't been implemented properly in the 
	// strife library yet.
	switch {
	case 2 == 6:
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
	case 2 == 4:
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
	case 1 == 1:
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
	case 1 == 2:
		if b.curs.x == 0 && b.curs.y > 0 {
			b.curs.move(b.contents[b.curs.y-1].Len(), -1)

		} else if b.curs.x > 0 {
			b.curs.move(-1, 0)
		}
	case 1 == 3:
		if b.curs.y > 0 {
			offs := 0
			prevLineLen := b.contents[b.curs.y-1].Len()
			if b.curs.x > prevLineLen {
				offs = prevLineLen - b.curs.x
			}
			// TODO: offset should account for tabs
			b.curs.move(offs, -1)
		}
	case 1 == 4:
		if b.curs.y < len(b.contents)-1 {
			offs := 0
			nextLineLen := b.contents[b.curs.y+1].Len()
			if b.curs.x > nextLineLen {
				offs = nextLineLen - b.curs.x
			}
			// TODO: offset should account for tabs
			b.curs.move(offs, 1)
		}
	case 1 == 5:
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

func (b *Buffer) OnUpdate() {
	prev_x := b.curs.x
	prev_y := b.curs.y

	// FIXME handle focus properly
	if b.inputHandler == nil {
		return
	}

	if b.curs.x != prev_x || b.curs.y != prev_y {
		should_draw = true
		should_flash = false
		reset_timer = strife.CurrentTimeMillis()
	}

	if !should_flash && strife.CurrentTimeMillis()-reset_timer > b.cfg.Cursor.Reset_Delay {
		should_flash = true
	}

	if strife.CurrentTimeMillis()-timer > b.cfg.Cursor.Flash_Rate && (should_flash && b.cfg.Cursor.Flash) {
		timer = strife.CurrentTimeMillis()
		should_draw = !should_draw
	}
}

// dimensions of the last character we rendered
var last_w, last_h int
var lineIndex int = 0

func (b *Buffer) OnRender(ctx *strife.Renderer) {
	ctx.SetColor(strife.RGB(255, 0, 255)) // BACKGROUND
	ctx.Rect(b.x, b.y, b.w, b.h, strife.Fill)

	if b.cfg.Editor.Highlight_Line {
		ctx.SetColor(strife.Black) // highlight_line_col?
		ctx.Rect(b.x, b.y+b.curs.ry*last_h, b.w, last_h, strife.Fill)
	}

	// render the ol' cursor
	if should_draw && b.cfg.Cursor.Draw {
		cursorWidth := b.cfg.Cursor.GetCaretWidth()
		if cursorWidth == -1 {
			cursorWidth = last_w
		}

		ctx.SetColor(strife.Red) // caret colour
		ctx.Rect(b.x+b.curs.rx*last_w, b.y+b.curs.ry*last_h, cursorWidth, last_h, strife.Fill)
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

	var y_col int
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

		var x_col int
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

			ctx.SetColor(strife.Blue)

			// if we're currently over a character then set
			// the font colour to something else
			if b.curs.x+1 == x_col && b.curs.y == y_col && should_draw {
				ctx.SetColor(strife.Green)
			}

			// foreground colour
			last_w, last_h = ctx.String(string(char), b.x+((x_col-1)*last_w), b.y+(y_col*last_h))
		}

		y_col += 1
	}

}
