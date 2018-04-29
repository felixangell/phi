package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/felixangell/go-rope"
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/phi/lex"
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

type BufferConfig struct {
	background        int32
	foreground        int32
	cursor            int32
	cursorInvert      int32
	lineNumBackground int32
	lineNumForeground int32
}

type Buffer struct {
	BaseComponent
	index        int
	parent       *View
	font         *strife.Font
	contents     []*rope.Rope
	curs         *Cursor
	cfg          *cfg.TomlConfig
	buffOpts     BufferConfig
	cam          *camera
	filePath     string
	languageInfo *cfg.LanguageSyntaxConfig
}

func NewBuffer(conf *cfg.TomlConfig, buffOpts BufferConfig, parent *View, index int) *Buffer {
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
		buffOpts: buffOpts,
		filePath: "",
		cam:      &camera{0, 0},
	}
	return buff
}

// sx, sy -> the starting x, y pos of the cursor
// ex, ey -> the end x, y of the selection
type selection struct {
	parent *Buffer
	sx, sy int
	ex, ey int
}

func (s *selection) renderAt(ctx *strife.Renderer, xOff int, yOff int) {
	ctx.SetColor(strife.Red)

	xd := (s.ex - s.sx) + 1
	yd := (s.ey - s.sy) + 1

	for y := 0; y < yd; y++ {
		lineLen := s.parent.contents[s.ey+y].Len() - s.sx
		if y == yd-1 {
			lineLen = xd
		}

		width := lineLen * last_w

		ctx.Rect(xOff+(s.sx*last_w), yOff+((s.sy+y)*last_h), width, last_h, strife.Fill)
	}
}

func (b *Buffer) OpenFile(filePath string) {
	b.filePath = filePath

	log.Println("Opening file ", filePath)

	ext := path.Ext(filePath)

	var err error
	b.languageInfo, err = b.cfg.GetSyntaxConfig(ext)
	if err != nil {
		log.Println(err.Error())
	}

	// if the file doesn't exist, try to create it before reading it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		} else {
			f.Close()
		}
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

func (b *Buffer) setLine(idx int, val string) {
	b.contents[idx] = rope.New(val)
	if b.curs.y == idx {
		b.moveToEndOfLine()
	}
}

func (b *Buffer) appendLine(val string) {
	b.contents = append(b.contents, rope.New(val))
	// because we've added a new line
	// we have to set the x to the start
	b.curs.x = 0
}

// inserts a string, handling all of the newlines etc
func (b *Buffer) insertString(idx int, val string) {
	lines := strings.Split(val, "\n")

	for _, l := range lines {
		b.contents = append(b.contents, new(rope.Rope))
		copy(b.contents[b.curs.y+idx+1:], b.contents[b.curs.y+idx:])
		b.contents[b.curs.y+idx] = rope.New(l)
		b.moveDown()
	}

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
			if action, ok := actions[actionName]; ok {
				log.Println("Executing action '" + actionName + "'")
				return action.proc(b)
			}
		} else {
			log.Println("warning, unimplemented shortcut ctrl +", string(unicode.ToLower(r)), actionName)
		}
	}

	if SUPER_DOWN {
		actionName, actionExists := cfg.Shortcuts.Supers[string(unicode.ToLower(r))]
		if actionExists {
			if action, ok := actions[actionName]; ok {
				return action.proc(b)
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

func (b *Buffer) moveToEndOfLine() {
	lineLen := b.contents[b.curs.y].Len()

	if b.curs.x > lineLen {
		distToMove := b.curs.x - lineLen
		for i := 0; i < distToMove; i++ {
			b.moveLeft()
		}
	} else if b.curs.x < lineLen {
		distToMove := lineLen - b.curs.x
		for i := 0; i < distToMove; i++ {
			b.moveRight()
		}
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

var lastSelection *selection

func (b *Buffer) processSelection(key int) bool {
	if lastSelection == nil {
		lastSelection = &selection{
			b,
			b.curs.x, b.curs.y,
			b.curs.x, b.curs.y,
		}
	}

	switch key {
	case sdl.K_LEFT:
		lastSelection.ex--
		b.moveLeft()
		break
	case sdl.K_RIGHT:
		lastSelection.ex++
		b.moveRight()
		break
	case sdl.K_UP:
		lastSelection.ey--
		b.moveUp()
		break
	case sdl.K_DOWN:
		lastSelection.ey++
		b.moveDown()
		break
	}

	return true
}

// processes a key press. returns if a key was processed
// or not, for example the letter 'a' could run through this
// which is not an action key, therefore we return false
// because it was not processed.
func (b *Buffer) processActionKey(key int) bool {
	if SHIFT_DOWN {
		switch key {
		case sdl.K_LEFT:
			fallthrough
		case sdl.K_RIGHT:
			fallthrough
		case sdl.K_DOWN:
			fallthrough
		case sdl.K_UP:
			return b.processSelection(key)
		}
	}

	switch key {
	case sdl.K_CAPSLOCK:
		CAPS_LOCK = !CAPS_LOCK

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

	case sdl.K_BACKSPACE:
		if SUPER_DOWN {
			b.deleteBeforeCursor()
		} else {
			b.deletePrev()
		}

	case sdl.K_RIGHT:
		currLineLength := b.contents[b.curs.y].Len()

		if CONTROL_DOWN && b.parent != nil {
			b.parent.ChangeFocus(1)
			break
		}

		if SUPER_DOWN {
			for b.curs.x < currLineLength {
				b.curs.move(1, 0)
			}
			break
		}

		// FIXME this is weird!
		// this will move to the next blank or underscore
		// character
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
			break
		}

		b.moveRight()
	case sdl.K_LEFT:
		if CONTROL_DOWN && b.parent != nil {
			b.parent.ChangeFocus(-1)
			break
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
			break
		}

		b.moveLeft()
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

	case sdl.K_END:
		currLine := b.contents[b.curs.y]
		if b.curs.x < currLine.Len() {
			distToMove := currLine.Len() - b.curs.x
			b.curs.move(distToMove, 0)
		}

	case sdl.K_HOME:
		if b.curs.x > 0 {
			b.curs.move(-b.curs.x, 0)
		}

	case sdl.K_PAGEUP:
		b.scrollUp()

	case sdl.K_PAGEDOWN:
		b.scrollDown()

	case sdl.K_DELETE:
		b.deleteNext()

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

	default:
		return false
	}

	return true
}

var (
	SHIFT_DOWN   bool = false
	SUPER_DOWN        = false // cmd on mac, ctrl on windows
	CONTROL_DOWN      = false // what is this on windows?
	ALT_DOWN          = false // option on mac
	CAPS_LOCK         = false
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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

var lastCursorDraw = time.Now()
var renderFlashingCursor = true

func (b *Buffer) OnUpdate() bool {
	return b.doUpdate(nil)
}

func (b *Buffer) doUpdate(pred func(r int) bool) bool {
	if !b.HasFocus() {
		return false
	}

	SHIFT_DOWN = strife.KeyPressed(sdl.K_LSHIFT) || strife.KeyPressed(sdl.K_RSHIFT)
	SUPER_DOWN = strife.KeyPressed(sdl.K_LGUI) || strife.KeyPressed(sdl.K_RGUI)
	ALT_DOWN = strife.KeyPressed(sdl.K_LALT) || strife.KeyPressed(sdl.K_RALT)
	CONTROL_DOWN = strife.KeyPressed(sdl.K_LCTRL) || strife.KeyPressed(sdl.K_RCTRL)

	if strife.PollKeys() {
		keyCode := strife.PopKey()

		if pred != nil {
			if val := pred(keyCode); val {
				return val
			}
		}

		// try process this key input as an
		// action first
		actionPerformed := b.processActionKey(keyCode)
		if actionPerformed {

			// check if we need to remove the selection
			// or not, this is if we moved the cursor somehow
			if lastSelection != nil {
				if b.curs.x != lastSelection.ex {
					lastSelection = nil
				} else if b.curs.y != lastSelection.ey {
					lastSelection = nil
				}
			}

			return true
		}

		textEntered := b.processTextInput(rune(keyCode))
		if textEntered {
			return true
		}
	}

	// handle cursor flash
	if b.cfg.Cursor.Flash && runtime.GOOS == "spookyghost" {
		if time.Now().Sub(lastCursorDraw) >= time.Duration(b.cfg.Cursor.Flash_Rate)*time.Millisecond {
			renderFlashingCursor = !renderFlashingCursor
			lastCursorDraw = time.Now()
		}
	}

	if !b.HasFocus() {
		return false
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

// editor x and y offsets
var ex, ey = 0, 0

// runs up a lexer instance
func lexFindMatches(matches *map[int]syntaxRuneInfo, currLine string, toMatch map[string]bool, bg int, fg int) {
	// start up a lexer instance and
	// lex the line.
	lexer := lex.New(currLine)

	tokenStream := lexer.Tokenize()

	for _, tok := range tokenStream {
		if _, ok := toMatch[tok.Lexeme]; ok {
			(*matches)[tok.Start] = syntaxRuneInfo{bg, -1, len(tok.Lexeme)}
		}
	}
}

func (b *Buffer) syntaxHighlightLine(currLine string) map[int]syntaxRuneInfo {
	matches := map[int]syntaxRuneInfo{}

	subjects := make([]*cfg.SyntaxCriteria, len(b.languageInfo.Syntax))
	colours := make([]int, len(b.languageInfo.Syntax))

	idx := 0
	for _, criteria := range b.languageInfo.Syntax {
		colours[idx] = criteria.Colour
		subjects[idx] = criteria
		idx++
	}

	// HOLY SLOW BATMAN
	for syntaxIndex, syntax := range subjects {
		if syntax.Pattern != "" {
			for charIndex := 0; charIndex < len(currLine); charIndex++ {
				a := string(currLine[charIndex:])

				matched := syntax.CompiledPattern.FindStringIndex(a)
				if matched != nil {
					if _, ok := matches[charIndex]; !ok {
						matchedStrLen := (matched[1] - matched[0])
						matches[charIndex+matched[0]] = syntaxRuneInfo{colours[syntaxIndex], -1, matchedStrLen}
						charIndex += matchedStrLen
						continue
					}
				}
			}
		} else {
			background := colours[syntaxIndex]
			foreground := 0
			lexFindMatches(&matches, currLine, syntax.MatchList, background, foreground)
		}
	}

	return matches
}

func (b *Buffer) renderAt(ctx *strife.Renderer, rx int, ry int) {
	// BACKGROUND
	ctx.SetColor(strife.HexRGB(b.buffOpts.background))
	ctx.Rect(b.x, b.y, b.w, b.h, strife.Fill)

	if b.cfg.Editor.Highlight_Line && b.HasFocus() {
		ctx.SetColor(strife.Black) // highlight_line_col?

		highlightLinePosY := ey + (ry + b.curs.ry*last_h) - (b.cam.y * last_h)
		highlightLinePosX := ex + rx

		ctx.Rect(highlightLinePosX, highlightLinePosY, b.w-ex, last_h, strife.Fill)
	}

	var visibleLines, visibleChars int = 50, -1

	// HACK
	// force a render off screen
	// so we can calculate the size of characters
	{
		if int(last_h) == 0 || int(last_w) == 0 {
			last_w, last_h = ctx.String("_", -50, -50)
		}
	}

	// last_h > 0 means we have done
	// a render.
	if int(last_h) > 0 && int(b.h) != 0 {
		// render an extra three lines just
		// so we dont cut anything off if its
		// not evenly divisible
		visibleLines = (int(b.h) / int(last_h)) + 3
	}

	// calculate how many chars we can fit
	// on the screen.
	if int(last_w) > 0 && int(b.w) != 0 {
		visibleChars = (int(b.w-ex) / int(last_w)) - 3
	}

	start := b.cam.y
	upper := b.cam.y + visibleLines
	if start > len(b.contents) {
		start = len(b.contents)
	}
	if upper > len(b.contents) {
		upper = len(b.contents)
	}

	// render the selection if any
	if lastSelection != nil {
		lastSelection.renderAt(ctx, b.x+ex, b.y+ey)
	}

	numLines := len(b.contents)

	var y_col int
	for lineNum, rope := range b.contents[start:upper] {
		currLine := []rune(rope.String())

		// slice the visible characters only.
		currLine = currLine[:min(visibleChars, len(currLine))]

		// char index => colour
		var matches map[int]syntaxRuneInfo
		if b.languageInfo != nil && len(currLine) > 0 {
			matches = b.syntaxHighlightLine(string(currLine))
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

			ctx.SetColor(strife.HexRGB(b.buffOpts.foreground))

			// if we're currently over a character then set
			// the font colour to something else
			// ONLY SET THE COLOUR IF WE HAVE FOCUS ALSO!
			if b.HasFocus() && b.curs.x+1 == x_col && b.curs.y == y_col && should_draw {
				ctx.SetColor(strife.HexRGB(b.buffOpts.cursorInvert))
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
			last_w, last_h = ctx.String(string(char), ex+(rx+((x_col-1)*last_w)), (ry + (y_col * last_h)))
		}

		if b.cfg.Editor.Show_Line_Numbers {
			gutterPadPx := 10
			numLinesWidth := len(string(numLines)) + 1
			gutterWidth := last_w*numLinesWidth + (gutterPadPx * 2)

			// render the line numbers
			ctx.SetColor(strife.HexRGB(b.buffOpts.lineNumBackground))
			ctx.Rect(rx, (ry + (y_col * last_h)), gutterWidth, b.h, strife.Fill)

			ctx.SetColor(strife.HexRGB(b.buffOpts.lineNumForeground))
			ctx.String(fmt.Sprintf("%*d", numLinesWidth, start+lineNum), rx+gutterPadPx, (ry + (y_col * last_h)))

			ex = gutterWidth
		}

		y_col += 1
	}

	// render the ol' cursor
	if b.HasFocus() && renderFlashingCursor && b.cfg.Cursor.Draw {
		cursorWidth := b.cfg.Cursor.GetCaretWidth()
		if cursorWidth == -1 {
			cursorWidth = last_w
		}

		ctx.SetColor(strife.HexRGB(b.buffOpts.cursor)) // caret colour
		ctx.Rect(ex+(rx+b.curs.rx*last_w)-(b.cam.x*last_w), (ry+b.curs.ry*last_h)-(b.cam.y*last_h), cursorWidth, last_h, strife.Fill)
	}
}

func (b *Buffer) OnRender(ctx *strife.Renderer) {
	b.renderAt(ctx, b.x, b.y)
}
