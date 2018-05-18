package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/felixangell/fuzzysearch/fuzzy"
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/phi/lex"
	"github.com/felixangell/piecetable"
	"github.com/felixangell/strife"
	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	timer        int64 = 0
	reset_timer  int64 = 0
	should_draw  bool  = true
	should_flash bool
)

const (
	DEFAULT_SCROLL_AMOUNT = 10
)

// TODO move into config
// line pad:
var pad = 6
var halfPad = pad / 2

type camera struct {
	x  int
	y  int
	dx int
	dy int
}

// TODO maybe have a thread that finds
// words in the current file
// and adds them to the vocabulary.
type AutoCompleteBox struct {
	vocabularyRegister map[string]bool
	vocabulary         []string
	suggestions        []string
	lastRunes          []rune
}

func (a *AutoCompleteBox) hasSuggestions() bool {
	return len(a.suggestions) > 0
}

func (a *AutoCompleteBox) process(r rune) {
	// space or non letter clears last word
	if r == ' ' || !unicode.IsLetter(r) {
		// we completed the word
		// so we add it to the vocabulary
		if r == ' ' {
			word := string(a.lastRunes)
			if _, ok := a.vocabularyRegister[word]; !ok {
				a.vocabulary = append(a.vocabulary, word)
				a.vocabularyRegister[word] = true
			}
		}

		a.lastRunes = []rune{}
		return
	}

	a.lastRunes = append(a.lastRunes, r)

	// don't bother unless its 4 or more letters.
	if len(a.lastRunes) <= 3 {
		return
	}

	word := string(a.lastRunes)
	println("looking up word: ", word)

	results := fuzzy.RankFind(word, a.vocabulary)
	a.suggestions = []string{}
	for _, res := range results {
		a.suggestions = append(a.suggestions, res.Target)
	}
}

func (a *AutoCompleteBox) renderAt(x, y int, ctx *strife.Renderer) {
	height := last_h + pad
	itemCount := len(a.suggestions)

	ctx.SetColor(strife.HexRGB(0x000000))
	ctx.Rect(x, y, 120, height*itemCount, strife.Fill)

	ctx.SetColor(strife.HexRGB(0xffffff))
	for idx, sugg := range a.suggestions {
		ctx.String(sugg, x, y+(idx*height))
	}
}

func newAutoCompleteBox() *AutoCompleteBox {
	return &AutoCompleteBox{
		map[string]bool{},
		[]string{},
		[]string{},
		[]rune{},
	}
}

type BufferConfig struct {
	background        int32
	foreground        int32
	cursor            int32
	cursorInvert      int32
	lineNumBackground int32
	lineNumForeground int32
	font              *strife.Font
}

type Buffer struct {
	BaseComponent
	index        int
	parent       *View
	curs         *Cursor
	cfg          *cfg.TomlConfig
	buffOpts     BufferConfig
	cam          *camera
	table        *piecetable.PieceTable
	filePath     string
	languageInfo *cfg.LanguageSyntaxConfig
	ex, ey       int
	modified     bool
	autoComplete *AutoCompleteBox
}

func NewBuffer(conf *cfg.TomlConfig, buffOpts BufferConfig, parent *View, index int) *Buffer {
	config := conf
	if config == nil {
		config = cfg.NewDefaultConfig()
	}

	curs := sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM)

	// TODO do this only if we are hovering over the buffer
	sdl.SetCursor(curs)

	buff := &Buffer{
		index:        index,
		parent:       parent,
		curs:         &Cursor{},
		cfg:          config,
		table:        piecetable.MakePieceTable(""),
		buffOpts:     buffOpts,
		filePath:     "",
		cam:          &camera{0, 0, 0, 0},
		autoComplete: newAutoCompleteBox(),
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
		lineLen := s.parent.table.Lines[s.ey+y].Len() - s.sx
		if y == yd-1 {
			lineLen = xd
		}

		width := lineLen * last_w

		ctx.Rect(xOff+(s.sx*last_w), yOff+((s.sy+y)*last_h), width, last_h, strife.Fill)
	}
}

func (b *Buffer) reload() {
	// if the file doesn't exist, try to create it before reading it
	if _, err := os.Stat(b.filePath); os.IsNotExist(err) {
		// this shouldn't really happen, for some
		// reason the file no longer exists?
		log.Println("File does not exist when reloading?! " + b.filePath)
		return
	}

	// if the file has modifications made to it
	// ask if the user wants to reload the file or not
	// otherwise re-load it anyway.
	if b.modified {
		ok := dialog.Message("This file has been modified, would you like to reload?").YesNo()
		if !ok {
			return
		}
	}

	contents, err := ioutil.ReadFile(b.filePath)
	if err != nil {
		panic(err)
	}

	b.table = piecetable.MakePieceTable(string(contents))

	// TODO perhaps when we reload the current line might not exist or something
	// try and set the cursor to what it was before but maybe make sure its not out
	// of bounds, etc.

	b.modified = false
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

	// add the file to the watcher.
	b.parent.registerFile(filePath, b)

	b.table = piecetable.MakePieceTable(string(contents))

	// because appendLine sets modified to true
	// we should reset this to false since weve
	// only loaded a file.
	b.modified = false
}

func (b *Buffer) setLine(idx int, val string) {
	b.modified = true

	b.table.Lines[idx] = piecetable.NewLine(val, b.table)
	if b.curs.y == idx {
		b.moveToEndOfLine()
	}
}

func (b *Buffer) appendLineAt(val *piecetable.Line, idx int) {
	b.modified = true

	b.table.Lines = append(b.table.Lines, val)
	copy(b.table.Lines[idx+1:], b.table.Lines[idx:])
	b.table.Lines[idx] = val
}

func (b *Buffer) appendStringAt(val string, idx int) {
	b.appendLineAt(piecetable.NewLine(val, b.table), idx)
}

// appendLine adds a string to the end of
// the buffer.
func (b *Buffer) appendLine(val string) {
	b.modified = true

	b.table.Lines = append(b.table.Lines, piecetable.NewLine(val, b.table))
	// because we've added a new line
	// we have to set the x to the start
	b.curs.x = 0
}

// inserts a string, handling all of the newlines etc
func (b *Buffer) insertString(idx int, val string) {
	b.modified = true

	lines := strings.Split(val, "\n")

	for _, l := range lines {
		b.table.Lines = append(b.table.Lines, piecetable.NewLine("", b.table))
		copy(b.table.Lines[b.curs.y+idx+1:], b.table.Lines[b.curs.y+idx:])
		b.table.Lines[b.curs.y+idx] = piecetable.NewLine(l, b.table)
		b.moveDown()
	}

}

func (b *Buffer) insertRune(r rune) {
	b.modified = true

	log.Println("Inserting rune ", r, " into current line at ", b.curs.x, ":", b.curs.y)
	log.Println("Line before insert> ", b.table.Lines[b.curs.y].String())

	b.table.Insert(string(r), b.curs.y, b.curs.x)
	b.moveRight()
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

func (b *Buffer) deleteLine() {
	// HACK FIXME
	b.modified = true

	if len(b.table.Lines) > 1 {
		b.table.Lines = remove(b.table.Lines, b.curs.y)
	} else {
		// we are on the first line
		// and there is nothing else to delete
		// so we just clear the line
		b.table.Lines[b.curs.y] = piecetable.NewLine("", b.table)
	}

	if b.curs.y >= len(b.table.Lines) {
		if b.curs.y > 0 {
			b.moveUp()
		}
	}

	b.moveToEndOfLine()
}

func (b *Buffer) processTextInput(r rune) bool {
	if ALT_DOWN && r == '\t' {
		// nop, we dont want to
		// insert tabs when we
		// alt tab out of view of this app
		return true
	}

	b.autoComplete.process(r)

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

	mainSuper, shortcutName := CONTROL_DOWN, "ctrl"
	source := cfg.Shortcuts.Controls
	if runtime.GOOS == "darwin" {
		mainSuper, shortcutName = SUPER_DOWN, "super"
		source = cfg.Shortcuts.Supers
	}

	if mainSuper {
		left := 1073741903
		right := 1073741904

		// map to left/right/etc.
		// FIXME
		var key string
		switch int(r) {
		case left:
			key = "left"
		case right:
			key = "right"
		default:
			key = string(unicode.ToLower(r))
		}

		actionName, actionExists := source[key]
		if actionExists {
			if action, ok := actions[actionName]; ok {
				return action.proc(b.parent, []string{})
			}
		} else {
			log.Println("warning, unimplemented shortcut", shortcutName, "+", unicode.ToLower(r), "#", int(r), actionName)
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
			currLine := b.table.Lines[b.curs.y]
			if b.curs.x < currLine.Len() {
				curr := b.table.Index(b.curs.y, b.curs.x+1)
				if curr == r {
					b.moveRight()
					return true
				} else {
					log.Print("no it's ", curr)
				}
			}
		}
	}

	// HACK FIXME
	b.modified = true

	b.table.Insert(string(r), b.curs.y, b.curs.x)
	b.moveRight()

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
		b.table.Insert(string(rune(matchingPair)), b.curs.y, b.curs.x)
	}

	return true
}

func remove(slice []*piecetable.Line, s int) []*piecetable.Line {
	return append(slice[:s], slice[s+1:]...)
}

func (b *Buffer) deleteNext() {
	b.moveRight()
	b.deletePrev()
}

// FIXME clean this up!
func (b *Buffer) deletePrev() {
	if b.curs.x > 0 {
		offs := -1
		if !b.cfg.Editor.Tabs_Are_Spaces {
			if b.table.Index(b.curs.y, b.curs.x) == '\t' {
				offs = int(-b.cfg.Editor.Tab_Size)
			}
		} else if b.cfg.Editor.Hungry_Backspace && b.curs.x >= int(b.cfg.Editor.Tab_Size) {
			// cut out the last {TAB_SIZE} amount of characters
			// and check em
			tabSize := int(b.cfg.Editor.Tab_Size)

			// render the line...
			currLine := b.table.Lines[b.curs.y].String()
			before := currLine[b.curs.x-tabSize:]

			if strings.HasPrefix(before, b.makeTab()) {
				// delete {TAB_SIZE} amount of characters
				// from the cursors x pos
				for i := 0; i < int(b.cfg.Editor.Tab_Size); i++ {

					b.table.Delete(b.curs.y, b.curs.x)
					b.curs.move(-1, 0)
				}
				return
			}
		}

		b.table.Delete(b.curs.y, b.curs.x)

		if len(b.autoComplete.lastRunes) > 0 {
			// pop the last word off in the auto complete
			b.autoComplete.lastRunes = b.autoComplete.lastRunes[:len(b.autoComplete.lastRunes)-1]
		}

		b.curs.moveRender(-1, 0, offs, 0)
	} else if b.curs.x == 0 && b.curs.y > 0 {
		// TODO this should set the previous word
		// in the auto complete to the word before the cursor after
		// wrapping back the line?

		// start of line, wrap to previous
		prevLineLen := b.table.Lines[b.curs.y-1].Len()

		val := b.table.Lines[b.curs.y].String()
		b.table.Insert(val, b.curs.y-1, b.table.Lines[b.curs.y-1].Len())

		b.table.Lines = append(b.table.Lines[:b.curs.y], b.table.Lines[b.curs.y+1:]...)
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
		b.curs.move(b.table.Lines[b.curs.y-1].Len(), -1)

	} else if b.curs.x > 0 {
		b.curs.move(-1, 0)
	}
}

func (b *Buffer) moveRight() {
	currLineLength := b.table.Lines[b.curs.y].Len()

	if b.curs.x >= currLineLength && b.curs.y < len(b.table.Lines)-1 {
		// we're at the end of the line and we have
		// some lines after, let's wrap around
		b.curs.move(-currLineLength, 0)
		b.moveDown()
	} else if b.curs.x < b.table.Lines[b.curs.y].Len() {
		// we have characters to the right, let's move along
		b.curs.move(1, 0)
	}
}

func (b *Buffer) moveToStartOfLine() {
	for b.curs.x > 0 {
		b.moveLeft()
	}
}

func (b *Buffer) moveToEndOfLine() {
	lineLen := b.table.Lines[b.curs.y].Len()

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

// TODO make this scroll auto magically.
func (b *Buffer) gotoLine(num int64) {
	distToMove := float64((num - 1) - int64(b.curs.y))
	for i := int64(0); i < int64(math.Abs(distToMove)); i++ {
		if distToMove < 0 {
			b.moveUp()
		} else {
			b.moveDown()
		}
	}
}

func (b *Buffer) moveUp() {
	if b.curs.y > 0 {
		offs := 0
		prevLineLen := b.table.Lines[b.curs.y-1].Len()
		if b.curs.x > prevLineLen {
			offs = prevLineLen - b.curs.x
		}

		if b.cam.y > 0 {
			if b.curs.y == b.cam.y {
				b.scrollUp(1)
			}
		}

		// TODO: offset should account for tabs
		b.curs.move(offs, -1)
	}
}

func (b *Buffer) moveDown() {
	if b.curs.y < len(b.table.Lines)-1 {
		offs := 0
		nextLineLen := b.table.Lines[b.curs.y+1].Len()
		if b.curs.x > nextLineLen {
			offs = nextLineLen - b.curs.x
		}

		// TODO: offset should account for tabs

		b.curs.move(offs, 1)

		visibleLines := int(math.Ceil(float64(b.h-b.ey) / float64(last_h+pad)))

		if b.curs.y >= visibleLines && b.curs.y-b.cam.y == visibleLines {
			b.scrollDown(1)
		}
	}
}

func (b *Buffer) swapLineUp() bool {
	if b.curs.y > 0 {
		currLine := b.table.Lines[b.curs.y]
		prevLine := b.table.Lines[b.curs.y-1]
		b.table.Lines[b.curs.y-1] = currLine
		b.table.Lines[b.curs.y] = prevLine
		b.moveUp()
	}
	return true
}

func (b *Buffer) swapLineDown() bool {
	if b.curs.y < len(b.table.Lines) {
		currLine := b.table.Lines[b.curs.y]
		nextLine := b.table.Lines[b.curs.y+1]
		b.table.Lines[b.curs.y+1] = currLine
		b.table.Lines[b.curs.y] = nextLine
		b.moveDown()
	}
	return true
}

func (b *Buffer) scrollUp(lineScrollAmount int) {
	if b.cam.y > 0 {
		b.cam.dy -= lineScrollAmount
	}

}

func (b *Buffer) scrollDown(lineScrollAmount int) {
	if b.cam.y < len(b.table.Lines) {
		b.cam.dy += lineScrollAmount
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

	case sdl.K_ESCAPE:
		return true

	case sdl.K_RETURN:
		if SUPER_DOWN {
			// in sublime this goes
			// into the next block
			// nicely indented!
		}

		// HACK FIXME
		b.modified = true

		// FIXME
		// clear the last runes
		b.autoComplete.lastRunes = []rune{}

		// START OF LINE:
		if b.curs.x == 0 {
			// we're at the start of a line, so we want to
			// shift the line down and insert an empty line
			// above it!
			b.table.Lines = append(b.table.Lines, piecetable.NewLine("", b.table)) // grow
			copy(b.table.Lines[b.curs.y+1:], b.table.Lines[b.curs.y:])             // shift
			b.table.Lines[b.curs.y] = piecetable.NewLine("", b.table)              // set
			b.moveDown()
			return true
		}

		initialX := b.curs.x
		prevLineLen := b.table.Lines[b.curs.y].Len()

		// END OF LINE:
		if initialX == prevLineLen {
			b.appendStringAt("", b.curs.y+1)
			b.moveDown()
			b.moveToStartOfLine()
			return true
		}

		// we're not at the end of the line, but we're not at
		// the start, i.e. we're SPLITTING the line
		left := b.table.Lines[b.curs.y].String()
		rightPart := left[initialX:]

		for i := 0; i < len(rightPart); i++ {
			// TODO POP in piecetable?
			b.table.Delete(b.curs.y, len(left)-i)
		}

		b.appendStringAt(rightPart, b.curs.y+1)
		b.moveDown()
		b.moveToStartOfLine()

	case sdl.K_BACKSPACE:
		// HACK FIXME
		b.modified = true

		if SUPER_DOWN {
			b.deleteBeforeCursor()
		} else {
			b.deletePrev()
		}

	case sdl.K_RIGHT:
		currLineLength := b.table.Lines[b.curs.y].Len()

		if SUPER_DOWN {
			for b.curs.x < currLineLength {
				b.moveLeft()
			}
			break
		}

		// FIXME this is weird!
		// this will move to the next blank or underscore
		// character
		if ALT_DOWN {
			line := b.table.Lines[b.curs.y]

			i := b.curs.x + 1 // ?

			for i < len(line.String())-1 {
				curr := b.table.Index(b.curs.y, i)

				switch curr {
				case ' ':
					fallthrough
				case '\n':
					fallthrough
				case '_':
					goto rightWordOuter
				}

				i = i + 1
			}

		rightWordOuter:

			dist := i - b.curs.x
			for j := 0; j < dist; j++ {
				b.moveRight()
			}
			break
		}

		b.moveRight()
	case sdl.K_LEFT:
		if SUPER_DOWN {
			// TODO go to the nearest \t
			// if no \t (i.e. start of line) go to
			// the start of the line!
			b.curs.gotoStart()
		} else if ALT_DOWN {
			i := b.curs.x
			for i > 0 {
				currChar := b.table.Index(b.curs.y, i)

				switch currChar {
				case ' ':
					fallthrough
				case '\n':
					fallthrough
				case '_':
					i = i - 1
					goto outer
				default:
					break
				}
				i = i - 1
			}

		outer:

			start := b.curs.x
			for j := 0; j < start-i; j++ {
				b.moveLeft()
			}

			if start == 0 {
				b.moveUp()
				b.moveToEndOfLine()
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

		b.moveUp()

	case sdl.K_DOWN:
		if ALT_DOWN {
			return b.swapLineDown()
		}

		if SUPER_DOWN {
			// go to the end of the file
		}

		b.moveDown()

	case sdl.K_TAB:
		// HACK FIXME
		b.modified = true

		if b.cfg.Editor.Tabs_Are_Spaces {
			// make an empty rune array of TAB_SIZE, cast to string
			// and insert it.
			b.table.Insert(b.makeTab(), b.curs.y, b.curs.x)
			b.curs.move(int(b.cfg.Editor.Tab_Size), 0)
		} else {
			b.table.Insert(string('\t'), b.curs.y, b.curs.x)
			// the actual position is + 1, but we make it
			// move by TAB_SIZE characters on the view.
			b.curs.moveRender(1, 0, int(b.cfg.Editor.Tab_Size), 0)
		}

	case sdl.K_END:
		currLine := b.table.Lines[b.curs.y]
		if b.curs.x < currLine.Len() {
			distToMove := currLine.Len() - b.curs.x
			b.curs.move(distToMove, 0)
		}

	case sdl.K_HOME:
		if b.curs.x > 0 {
			b.curs.move(-b.curs.x, 0)
		}

		// TODO remove since this is handled in the keymap!
	case sdl.K_PAGEUP:
		b.scrollUp(DEFAULT_SCROLL_AMOUNT)
		for i := 0; i < DEFAULT_SCROLL_AMOUNT; i++ {
			b.moveUp()
		}
	case sdl.K_PAGEDOWN:
		b.scrollDown(DEFAULT_SCROLL_AMOUNT)
		for i := 0; i < DEFAULT_SCROLL_AMOUNT; i++ {
			b.moveDown()
		}

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
			b.scrollDown(DEFAULT_SCROLL_AMOUNT)
		}
		if event.Y < 0 {
			b.scrollUp(DEFAULT_SCROLL_AMOUNT)
		}
	}
}

var lastCursorDraw = time.Now()
var renderFlashingCursor = true

var lastTimer = time.Now()
var ldx, ldy = 0, 0

var ySpeed = 1
var last = time.Now()

func (b *Buffer) processLeftClick() {
	// here we set the cursor y position
	// based off the click location
	yPos := strife.MouseCoords()[1]

	yPosToLine := ((yPos / (last_h + pad)) + 1) + b.cam.y

	fmt.Println(yPos, " line ", yPosToLine)
	b.gotoLine(int64(yPosToLine))
}

func (b *Buffer) OnUpdate() bool {
	if "animations on" == "true" {
		if b.cam.y < b.cam.dy {
			b.cam.y += ySpeed
		}
		if b.cam.y > b.cam.dy {
			b.cam.y -= ySpeed
		}
	} else {
		b.cam.x = b.cam.dx
		b.cam.y = b.cam.dy
	}

	// clamp camera
	if b.cam.y <= 0 {
		b.cam.y = 0
	}
	if b.cam.x <= 0 {
		b.cam.x = 0
	}

	switch strife.MouseButtonsState() {
	case strife.LeftMouseButton:
		b.processLeftClick()
	}

	return false
}

func (b *Buffer) processInput(pred func(r int) bool) bool {
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
	if b.cfg.Cursor.Flash && 15 == 12 {
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
	// TODO load this from config files!

	// BACKGROUND
	ctx.SetColor(strife.HexRGB(b.buffOpts.background))
	ctx.Rect(b.x, b.y, b.w, b.h, strife.Fill)

	if b.cfg.Editor.Highlight_Line && b.HasFocus() {
		ctx.SetColor(strife.Black) // highlight_line_col?

		highlightLinePosY := b.ey + (ry + b.curs.ry*(last_h+pad)) - (b.cam.y * (last_h + pad))
		highlightLinePosX := b.ex + rx

		ctx.Rect(highlightLinePosX, highlightLinePosY, b.w-b.ex, (last_h+pad)-b.ey, strife.Fill)
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
		visibleLines = (int(b.h-b.ey) / int(last_h)) + 3
	}

	// calculate how many chars we can fit
	// on the screen.
	if int(last_w) > 0 && int(b.w) != 0 {
		visibleChars = (int(b.w-b.ex) / int(last_w)) - 3
	}

	start := b.cam.y
	upper := b.cam.y + visibleLines
	if start > len(b.table.Lines) {
		start = len(b.table.Lines)
	}

	if upper > len(b.table.Lines) {
		upper = len(b.table.Lines)
	}

	// render the selection if any
	if lastSelection != nil {
		lastSelection.renderAt(ctx, b.x+b.ex, b.y+b.ey)
	}

	numLines := len(b.table.Lines)

	var y_col int
	for lineNum, rope := range b.table.Lines[start:upper] {
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

			// 13 is a carriage return
			case 13:
				continue

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

			lineHeight := last_h + pad
			xPos := b.ex + (rx + ((x_col - 1) * last_w))
			yPos := b.ey + (ry + (y_col * lineHeight)) + halfPad

			last_w, last_h = ctx.String(string(char), xPos, yPos)

			if DEBUG_MODE {
				ctx.SetColor(strife.HexRGB(0xff00ff))
				ctx.Rect(xPos, yPos, last_w, last_h, strife.Line)
			}
		}

		if b.cfg.Editor.Show_Line_Numbers {
			gutterPadPx := 10

			// how many chars we need
			numLinesCharWidth := len(string(numLines)) + 2

			gutterWidth := last_w*numLinesCharWidth + (gutterPadPx * 2)

			lineHeight := last_h + pad
			yPos := ((ry + y_col) * lineHeight) + halfPad

			// render the line numbers
			ctx.SetColor(strife.HexRGB(b.buffOpts.lineNumBackground))
			ctx.Rect(rx, yPos, gutterWidth, b.h, strife.Fill)

			if DEBUG_MODE {
				ctx.SetColor(strife.HexRGB(0xff00ff))
				ctx.Rect(rx, yPos, gutterWidth, b.h, strife.Line)
			}

			ctx.SetColor(strife.HexRGB(b.buffOpts.lineNumForeground))
			ctx.String(fmt.Sprintf("%*d", numLinesCharWidth, (start+lineNum)+1), rx+gutterPadPx, yPos)

			b.ex = gutterWidth
		}

		y_col += 1
	}

	cursorHeight := last_h + pad

	// render the ol' cursor
	if b.HasFocus() && (renderFlashingCursor || b.curs.moving) && b.cfg.Cursor.Draw {
		cursorWidth := b.cfg.Cursor.GetCaretWidth()
		if cursorWidth == -1 {
			cursorWidth = last_w
		}

		xPos := b.ex + (rx + b.curs.rx*last_w) - (b.cam.x * last_w)
		yPos := b.ey + (ry + b.curs.ry*cursorHeight) - (b.cam.y * cursorHeight)

		ctx.SetColor(strife.HexRGB(b.buffOpts.cursor))
		ctx.Rect(xPos, yPos, cursorWidth, cursorHeight, strife.Fill)

		if DEBUG_MODE {
			ctx.SetColor(strife.HexRGB(0xff00ff))
			ctx.Rect(xPos, yPos, cursorWidth, cursorHeight, strife.Line)
		}
	}

	if b.autoComplete.hasSuggestions() {

		xPos := b.ex + (rx + b.curs.rx*last_w) - (b.cam.x * last_w)
		yPos := b.ey + (ry + b.curs.ry*cursorHeight) - (b.cam.y * cursorHeight)

		autoCompleteBoxHeight := len(b.autoComplete.suggestions) * cursorHeight
		yPos = yPos - autoCompleteBoxHeight

		b.autoComplete.renderAt(xPos, yPos, ctx)
	}

	if DEBUG_MODE {
		ctx.SetColor(strife.HexRGB(0xff00ff))
		ctx.Rect(b.ex+rx, b.ey+ry, b.w-b.ex, b.h-b.ey, strife.Line)
	}
}

func (b *Buffer) OnRender(ctx *strife.Renderer) {
	ctx.SetFont(b.buffOpts.font)
	b.renderAt(ctx, b.x, b.y)
}
