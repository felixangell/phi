package gui

import (
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/felixangell/go-rope"
	"github.com/felixangell/phi-editor/cfg"
	"github.com/felixangell/strife"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	timer        int64 = 0
	reset_timer  int64 = 0
	should_draw  bool  = true
	should_flash bool
)

// TODO: allow font setting or whatever

type camera struct {
	x int
	y int
}

type Buffer struct {
	BaseComponent
	HasFocus bool
	index    int
	parent   *View
	font     *strife.Font
	contents []*rope.Rope
	curs     *Cursor
	cfg      *cfg.TomlConfig
	cam      *camera
	filePath string

	languageInfo string
}

func NewBuffer(conf *cfg.TomlConfig, parent *View, index int) *Buffer {
	config := conf
	if config == nil {
		config = cfg.NewDefaultConfig()
	}

	buffContents := []*rope.Rope{}
	buff := &Buffer{
		index:    index,
		parent:   parent,
		contents: buffContents,
		curs:     &Cursor{},
		cfg:      config,
		filePath: "/tmp/phi_file_" + time.Now().String(), // TODO make this a randomly chosen temp file
		cam:      &camera{0, 0},
	}

	return buff
}

func (b *Buffer) OpenFile(filePath string) {
	b.filePath = filePath

	log.Println("Opening file ", filePath)

	ext := path.Ext(filePath)
	lang, err := b.cfg.GetLanguageFromExt(ext)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("- this file is a ", lang, " language program")
		b.languageInfo = lang
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		b.appendLine(line)
	}
}

func (b *Buffer) OnDispose() {
	// hm!
	// os.Remove(b.fileHandle)
}

func (b *Buffer) OnInit() {}

func (b *Buffer) appendLine(val string) {
	b.contents = append(b.contents, rope.New(val))

	// because we've added a new line
	// we have to set the x to the start
	b.curs.x = 0
}

func (b *Buffer) insertRune(r rune) {
	log.Println("Inserting rune ", r, " into current line at ", b.curs.x, ":", b.curs.y)
	log.Println("Line before insert> ", b.contents[b.curs.y])

	b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string(r))
	b.curs.move(1, 0)
}

// TODO handle EVERYTHING but for now im handling
// my UK macbook key layout.
var shiftAlternative = map[rune]rune{
	'1':  '!',
	'2':  '@',
	'3':  '£',
	'4':  '$',
	'5':  '%',
	'6':  '^',
	'7':  '&',
	'8':  '*',
	'9':  '(',
	'0':  ')',
	'-':  '_',
	'=':  '+',
	'`':  '~',
	'/':  '?',
	'.':  '>',
	',':  '<',
	'[':  '{',
	']':  '}',
	';':  ':',
	'\'': '"',
	'\\': '|',
	'§':  '±',
}

var altAlternative = map[rune]rune{
	'1':  '¡',
	'2':  '€',
	'3':  '#',
	'4':  '¢',
	'5':  '∞',
	'6':  '§',
	'7':  '¶',
	'8':  '•',
	'9':  'ª',
	'0':  'º',
	'-':  '–',
	'=':  '≠',
	'`':  '`',
	'/':  '÷',
	'.':  '≥',
	',':  '≤',
	'[':  '“',
	']':  '‘',
	';':  '…',
	'\'': 'æ',
	'\\': '«',
}

func (b *Buffer) processTextInput(r rune) bool {
	if ALT_DOWN && r == '\t' {
		// nop, we dont want to
		// insert tabs when we
		// alt tab out of view of this app
		return true
	}

	// only do the alt alternatives on mac osx
	// todo change this so it's not checking on every
	// input
	if runtime.GOOS == "darwin" && ALT_DOWN {
		if val, ok := altAlternative[r]; ok {
			r = val
		}
	}

	if CAPS_LOCK {
		if unicode.IsLetter(r) {
			r = unicode.ToUpper(r)
		}
	}

	if CONTROL_DOWN {
		actionName, actionExists := cfg.Shortcuts.Controls[string(unicode.ToLower(r))]
		if actionExists {
			if proc, ok := actions[actionName]; ok {
				return proc(b)
			}
		} else {
			log.Println("warning, unimplemented shortcut ctrl+", unicode.ToLower(r), actionName)
		}
	}

	if SUPER_DOWN {
		actionName, actionExists := cfg.Shortcuts.Supers[string(unicode.ToLower(r))]
		if actionExists {
			if proc, ok := actions[actionName]; ok {
				return proc(b)
			}
		} else {
			log.Println("warning, unimplemented shortcut ctrl+", unicode.ToLower(r), actionName)
		}
	}

	if SHIFT_DOWN {
		// if it's a letter convert to uppercase
		if unicode.IsLetter(r) {
			r = unicode.ToUpper(r)
		} else {

			// otherwise we have to look in our trusy
			// shift mapping thing.
			if val, ok := shiftAlternative[r]; ok {
				r = val
			}

		}
	}

	// NOTE: we have to do this AFTER we map the
	// shift combo for the value!
	// this will not insert a ), }, or ] if there
	// is one to the right of us... basically
	// this escapes out of a closing bracket
	// rather than inserting a new one IF we are inside
	// brackets.
	if b.cfg.Editor.Match_Braces {
		if r == ')' || r == '}' || r == ']' {
			currLine := b.contents[b.curs.y]
			if b.curs.x < currLine.Len() {
				curr := currLine.Index(b.curs.x + 1)
				if curr == r {
					b.curs.move(1, 0)
					return true
				} else {
					log.Print("no it's ", curr)
				}
			}
		}
	}

	b.contents[b.curs.y] = b.contents[b.curs.y].Insert(b.curs.x, string(r))
	b.curs.move(1, 0)

	// we don't need to match braces
	// let's not continue any further
	if !b.cfg.Editor.Match_Braces {
		return true
	}

	// TODO: shall we match single quotes and double quotes too?

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

	return true
}

func remove(slice []*rope.Rope, s int) []*rope.Rope {
	return append(slice[:s], slice[s+1:]...)
}

func (b *Buffer) deleteNext() {
	b.moveRight()
	b.deletePrev()
}

func (b *Buffer) deletePrev() {
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
				return
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
}

func (b *Buffer) deleteBeforeCursor() {
	// delete so we're at the end
	// of the previous line
	if b.curs.x == 0 {
		b.deletePrev()
		return
	}

	for b.curs.x > 0 {
		b.deletePrev()
	}
}

func (b *Buffer) moveLeft() {
	if b.curs.x == 0 && b.curs.y > 0 {
		b.curs.move(b.contents[b.curs.y-1].Len(), -1)

	} else if b.curs.x > 0 {
		b.curs.move(-1, 0)
	}
}

func (b *Buffer) moveRight() {
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
}

func (b *Buffer) moveUp() {
	if b.curs.y > 0 {
		b.curs.move(0, -1)
	}
}

func (b *Buffer) moveDown() {
	if b.curs.y < len(b.contents) {
		b.curs.move(0, 1)
	}
}

func (b *Buffer) swapLineUp() bool {
	if b.curs.y > 0 {
		currLine := b.contents[b.curs.y]
		prevLine := b.contents[b.curs.y-1]
		b.contents[b.curs.y-1] = currLine
		b.contents[b.curs.y] = prevLine
		b.moveUp()
	}
	return true
}

func (b *Buffer) swapLineDown() bool {
	if b.curs.y < len(b.contents) {
		currLine := b.contents[b.curs.y]
		nextLine := b.contents[b.curs.y+1]
		b.contents[b.curs.y+1] = currLine
		b.contents[b.curs.y] = nextLine
		b.moveDown()
	}
	return true
}

func (b *Buffer) scrollUp() {
	if b.cam.y > 0 {
		// TODO move the cursor down 45 lines
		// IF the buffer exceeds the window size.
		lineScrollAmount := 10
		b.cam.y -= lineScrollAmount
		for i := 0; i < lineScrollAmount; i++ {
			b.moveUp()
		}
	}

}

func (b *Buffer) scrollDown() {
	if b.cam.y < len(b.contents) {
		// TODO move the cursor down 45 lines
		// IF the buffer exceeds the window size.
		lineScrollAmount := 10

		b.cam.y += lineScrollAmount
		for i := 0; i < lineScrollAmount; i++ {
			b.moveDown()
		}
	}
}

// processes a key press. returns if there
// was a key that MODIFIED the buffer.
func (b *Buffer) processActionKey(key int) bool {
	switch key {
	case sdl.K_CAPSLOCK:
		CAPS_LOCK = !CAPS_LOCK
		return true
	case sdl.K_RETURN:
		if SUPER_DOWN {
			// in sublime this goes
			// into the next block
			// nicely indented!
		}

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
			return true
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
		return true
	case sdl.K_BACKSPACE:
		if SUPER_DOWN {
			b.deleteBeforeCursor()
		} else {
			b.deletePrev()
		}
		return true
	case sdl.K_RIGHT:
		currLineLength := b.contents[b.curs.y].Len()

		if CONTROL_DOWN && b.parent != nil {
			b.parent.ChangeFocus(1)
			return true
		}

		if SUPER_DOWN {
			for b.curs.x < currLineLength {
				b.curs.move(1, 0)
			}
			return true
		}

		// FIXME this is weird!
		if ALT_DOWN {
			currLine := b.contents[b.curs.y]

			var i int
			for i = b.curs.x + 1; i < currLine.Len(); i++ {
				curr := currLine.Index(i)
				if curr <= ' ' || curr == '_' {
					break
				}
			}

			for j := 0; j < i; j++ {
				b.moveRight()
			}
			return true
		}

		b.moveRight()
		return true
	case sdl.K_LEFT:
		if CONTROL_DOWN && b.parent != nil {
			b.parent.ChangeFocus(-1)
			return true
		}

		if SUPER_DOWN {
			// TODO go to the nearest \t
			// if no \t (i.e. start of line) go to
			// the start of the line!
			b.curs.gotoStart()
		}

		if ALT_DOWN {
			currLine := b.contents[b.curs.y]

			i := b.curs.x
			for i > 0 {
				currChar := currLine.Index(i)
				// TODO is a seperator thing?
				if currChar <= ' ' || currChar == '_' {
					// move over one more?
					i = i - 1
					break
				}
				i = i - 1
			}

			start := b.curs.x
			for j := 0; j < start-i; j++ {
				b.moveLeft()
			}
			return true
		}

		b.moveLeft()
		return true
	case sdl.K_UP:
		if ALT_DOWN {
			return b.swapLineUp()
		}

		if SUPER_DOWN {
			// go to the start of the file
		}

		if b.curs.y > 0 {
			offs := 0
			prevLineLen := b.contents[b.curs.y-1].Len()
			if b.curs.x > prevLineLen {
				offs = prevLineLen - b.curs.x
			}
			// TODO: offset should account for tabs
			b.curs.move(offs, -1)
		}
		return true
	case sdl.K_DOWN:
		if ALT_DOWN {
			return b.swapLineDown()
		}

		if SUPER_DOWN {
			// go to the end of the file
		}

		if b.curs.y < len(b.contents)-1 {
			offs := 0
			nextLineLen := b.contents[b.curs.y+1].Len()
			if b.curs.x > nextLineLen {
				offs = nextLineLen - b.curs.x
			}
			// TODO: offset should account for tabs
			b.curs.move(offs, 1)
		}
		return true
	case sdl.K_TAB:
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
		return true

	case sdl.K_END:
		currLine := b.contents[b.curs.y]
		if b.curs.x < currLine.Len() {
			distToMove := currLine.Len() - b.curs.x
			b.curs.move(distToMove, 0)
		}
		return true

	case sdl.K_HOME:
		if b.curs.x > 0 {
			b.curs.move(-b.curs.x, 0)
		}
		return true

	case sdl.K_PAGEUP:
		b.scrollUp()
		return true

	case sdl.K_PAGEDOWN:
		b.scrollDown()
		return true

	case sdl.K_DELETE:
		b.deleteNext()
		return true

	case sdl.K_LGUI:
		fallthrough
	case sdl.K_RGUI:
		fallthrough

	case sdl.K_LALT:
		fallthrough
	case sdl.K_RALT:
		fallthrough

	case sdl.K_LCTRL:
		fallthrough
	case sdl.K_RCTRL:
		fallthrough

	case sdl.K_LSHIFT:
		fallthrough
	case sdl.K_RSHIFT:
		return true
	}

	return false
}

var (
	SHIFT_DOWN   bool = false
	SUPER_DOWN        = false // cmd on mac, ctrl on windows
	CONTROL_DOWN      = false // what is this on windows?
	ALT_DOWN          = false // option on mac
	CAPS_LOCK         = false
)

// TODO(Felix) this is really stupid
func (b *Buffer) makeTab() string {
	blah := []rune{}
	for i := 0; i < int(b.cfg.Editor.Tab_Size); i++ {
		blah = append(blah, ' ')
	}
	return string(blah)
}

func (b *Buffer) HandleEvent(evt strife.StrifeEvent) {
	switch event := evt.(type) {
	case *strife.MouseWheelEvent:
		if event.Y > 0 {
			b.scrollDown()
		}
		if event.Y < 0 {
			b.scrollUp()
		}
	}
}

func (b *Buffer) OnUpdate() bool {
	b.HasFocus = true

	prev_x := b.curs.x
	prev_y := b.curs.y

	SHIFT_DOWN = strife.KeyPressed(sdl.K_LSHIFT) || strife.KeyPressed(sdl.K_RSHIFT)
	SUPER_DOWN = strife.KeyPressed(sdl.K_LGUI) || strife.KeyPressed(sdl.K_RGUI)
	ALT_DOWN = strife.KeyPressed(sdl.K_LALT) || strife.KeyPressed(sdl.K_RALT)
	CONTROL_DOWN = strife.KeyPressed(sdl.K_LCTRL) || strife.KeyPressed(sdl.K_RCTRL)

	if strife.PollKeys() {
		keyCode := strife.PopKey()

		// try process this key input as an
		// action first
		actionPerformed := b.processActionKey(keyCode)
		if actionPerformed {
			return true
		}

		textEntered := b.processTextInput(rune(keyCode))
		if textEntered {
			return true
		}
	}

	// FIXME handle focus properly
	if b.inputHandler == nil {
		return false
	}

	if b.curs.x != prev_x || b.curs.y != prev_y {
		should_draw = true
		should_flash = false
		reset_timer = strife.CurrentTimeMillis()
	}

	// fixme to not use CurrentTimeMillis
	if !should_flash && strife.CurrentTimeMillis()-reset_timer > b.cfg.Cursor.Reset_Delay {
		should_flash = true
	}

	if strife.CurrentTimeMillis()-timer > b.cfg.Cursor.Flash_Rate && (should_flash && b.cfg.Cursor.Flash) {
		timer = strife.CurrentTimeMillis()
		should_draw = !should_draw
	}

	return false
}

type syntaxRuneInfo struct {
	background int
	foreground int
	length     int
}

// dimensions of the last character we rendered
var last_w, last_h int

func (b *Buffer) renderAt(ctx *strife.Renderer, rx int, ry int) {
	// BACKGROUND
	ctx.SetColor(strife.HexRGB(b.cfg.Theme.Background))
	ctx.Rect(b.x, b.y, b.w, b.h, strife.Fill)

	if b.cfg.Editor.Highlight_Line && b.HasFocus {
		ctx.SetColor(strife.Black) // highlight_line_col?
		ctx.Rect(rx, (ry+b.curs.ry*last_h)-(b.cam.y*last_h), b.w, last_h, strife.Fill)
	}

	// render the ol' cursor
	if should_draw && b.cfg.Cursor.Draw && b.HasFocus {
		cursorWidth := b.cfg.Cursor.GetCaretWidth()
		if cursorWidth == -1 {
			cursorWidth = last_w
		}

		ctx.SetColor(strife.HexRGB(b.cfg.Theme.Cursor)) // caret colour
		ctx.Rect((rx+b.curs.rx*last_w)-(b.cam.x*last_w), (ry+b.curs.ry*last_h)-(b.cam.y*last_h), cursorWidth, last_h, strife.Fill)
	}

	var visibleLines int = 50

	// last_h > 0 means we have done
	// a render.
	if int(last_h) > 0 && int(b.h) != 0 {
		// render an extra three lines just
		// so we dont cut anything off if its
		// not evenly divisible
		visibleLines = (int(b.h) / int(last_h)) + 3
	}

	start := b.cam.y
	upper := b.cam.y + visibleLines
	if upper > len(b.contents) {
		upper = len(b.contents)
	}

	var y_col int
	for _, rope := range b.contents[start:upper] {
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

		currLine := []rune(rope.String())

		// char index => colour
		matches := map[int]syntaxRuneInfo{}

		subjects := []cfg.SyntaxCriteria{}
		colours := []int{}

		stuff := b.cfg.Syntax[b.languageInfo]
		for _, criteria := range stuff {
			colours = append(colours, criteria.Colour)
			subjects = append(subjects, criteria)
		}

		// HOLY SLOW BATMAN
		for idx, _ := range currLine {
			for syntaxIndex, syntax := range subjects {

				// we have a regex pattern
				if syntax.Pattern != "" {

					// FIXME this is also very slow!
					// we could easily compile all of these
					// regular expressions when we load the
					// syntax highlighter.
					a := string(currLine[idx:])
					match, _ := regexp.MatchString(syntax.Pattern, a)
					if match {
						// for some reason this affects the whole line
						if _, ok := matches[idx]; !ok {
							matches[idx] = syntaxRuneInfo{colours[syntaxIndex], -1, len(a)}
							continue
						}
					}

				} else {

					for _, subject := range syntax.Match {
						if idx+len(subject)+1 > len(currLine) {
							continue
						}
						a := currLine[idx : idx+len(subject)+1]

						// we only want to match words. so we check that it has a space
						// before or after the subject word.
						if strings.Compare(string(a), subject+" ") == 0 || strings.Compare(string(a), " "+subject) == 0 {
							if _, ok := matches[idx]; !ok {
								matches[idx] = syntaxRuneInfo{colours[syntaxIndex], -1, len(string(a))}
								break
							}
							idx += len(subject)
						}
					}

				}
			}
		}

		colorStack := []int{}

		var x_col int
		for idx, char := range currLine {
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

			ctx.SetColor(strife.HexRGB(b.cfg.Theme.Foreground))

			// if we're currently over a character then set
			// the font colour to something else
			// ONLY SET THE COLOUR IF WE HAVE FOCUS ALSO!
			if b.HasFocus && b.curs.x+1 == x_col && b.curs.y == y_col && should_draw {
				ctx.SetColor(strife.HexRGB(b.cfg.Theme.Cursor_Invert))
			}

			if info, ok := matches[idx]; ok {
				for i := 0; i < info.length; i++ {
					colorStack = append(colorStack, info.background)
				}
			}

			if len(colorStack) > 0 {
				var a int32
				a, colorStack = int32(colorStack[len(colorStack)-1]), colorStack[:len(colorStack)-1]
				ctx.SetColor(strife.HexRGB(a))
			}
			last_w, last_h = ctx.String(string(char), (rx + ((x_col - 1) * last_w)), (ry + (y_col * last_h)))
		}

		y_col += 1
	}
}

func (b *Buffer) OnRender(ctx *strife.Renderer) {
	b.renderAt(ctx, b.x, b.y)
}
