package buff

import (
	"fmt"
	"github.com/felixangell/phi/internal/cfg"
	"github.com/felixangell/phi/internal/gui"
	"github.com/felixangell/phi/internal/lex"
	"github.com/felixangell/phi/pkg/piecetable"
	"github.com/felixangell/strife"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"unicode"
)

const (
	DefaultScrollAmount = 10
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

// AutoCompleteBox ...
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
	height := lastCharH + pad
	itemCount := len(a.suggestions)

	ctx.SetColor(strife.HexRGB(0x000000))
	ctx.Rect(x, y, 120, height*itemCount, strife.Fill)

	ctx.SetColor(strife.HexRGB(0xffffff))
	for idx, sugg := range a.suggestions {
		ctx.Text(sugg, x, y+(idx*height))
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
	background        uint32
	foreground        uint32
	cursor            uint32
	cursorInvert      uint32
	highlightLine     uint32
	lineNumBackground uint32
	lineNumForeground uint32
	font              *strife.Font
}

// Buffer is a structure representing
// a buffer of text.
type Buffer struct {
	gui.BaseComponent
	index        int
	parent       *BufferView
	curs         *Cursor
	cfg          *cfg.PhiEditorConfig
	buffOpts     BufferConfig
	cam          *camera
	table        *piecetable.PieceTable
	filePath     string
	languageInfo *cfg.LanguageSyntaxConfig
	ex, ey       int
	modified     bool
	autoComplete *AutoCompleteBox
}

// NewBuffer creates a new buffer with the given configurations
func NewBuffer(conf *cfg.PhiEditorConfig, buffOpts BufferConfig, parent *BufferView, index int) *Buffer {
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
		curs:         nil,
		cfg:          config,
		table:        piecetable.MakePieceTable(""),
		buffOpts:     buffOpts,
		filePath:     "",
		cam:          &camera{0, 0, 0, 0},
		autoComplete: newAutoCompleteBox(),
	}
	buff.curs = newCursor(buff)
	return buff
}

// sx, sy -> the starting x, y pos of the cursor
// ex, ey -> the end x, y of the selection
type selection struct {
	parent *Buffer
	sx, sy int
	ex, ey int
}

func (s *selection) renderAt(ctx *strife.Renderer, _ int, _ int) {
	ctx.SetColor(strife.Blue)

	yd := (s.ey - s.sy) + 1

	b := s.parent

	// renders the highlighting for a line.
	for y := 0; y < yd; y++ {
		lineLen := s.parent.table.Lines[s.sy+y].Len()

		// purely for the aesthetics?
		// feel like this will create bugs.
		if lineLen == 0 {
			lineLen = 1
		}

		height := (lastCharH + pad) - b.ey

		// width of box should be entire line
		width := lineLen * lastCharW

		// UNLESS we are on the current line
		if y == yd-1 {
			width = s.ex * lastCharW
		}

		xPos := s.sx * lastCharW
		yPos := (s.sy + y) * (lastCharH + pad)

		ctx.SetColor(strife.Blue)
		ctx.Rect(b.ex+xPos, b.ey+yPos, width, height, strife.Fill)
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
		// ok := dialog.Message("This file has been modified, would you like to reload?").YesNo()
		panic("this should show a file modified thing but this has been removed for now!")
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

// OpenFile will open the given file path into this buffer.
// This also handles loading of syntax stuff for syntax highlighting.
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
			if err := f.Close(); err != nil {
				panic(err)
			}
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
	if altDown && r == '\t' {
		// nop, we dont want to
		// insert tabs when we
		// alt tab out of view of this app
		return true
	}

	b.autoComplete.process(r)

	// FIXME
	// only do the alt alternatives on mac osx
	// todo change this so it's not checking on every
	// input
	if runtime.GOOS == "darwin" && altDown {
		if val, ok := altAlternative[r]; ok {
			r = val
		}
	}

	if capsLockDown {
		if unicode.IsLetter(r) {
			r = unicode.ToUpper(r)
		}
	}

	mainSuper, shortcutName := controlDown, "ctrl"
	source := cfg.Shortcuts.Controls
	if runtime.GOOS == "darwin" {
		mainSuper, shortcutName = superDown, "super"
		source = cfg.Shortcuts.Supers
	}

	if mainSuper {
		// FIXME magic numbers!
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
			if action, ok := register[actionName]; ok {
				return bool(action.proc(b.parent, []*lex.Token{}))
			}
		} else {
			log.Println("warning, unimplemented shortcut", shortcutName, "+", unicode.ToLower(r), "#", int(r), actionName)
		}
	}

	if shiftDown {
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
	if b.cfg.Editor.MatchBraces {
		if r == ')' || r == '}' || r == ']' {
			currLine := b.table.Lines[b.curs.y]
			if b.curs.x < currLine.Len() {
				curr := b.table.Index(b.curs.y, b.curs.x+1)
				if curr == r {
					b.moveRight()
					return true
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
	if !b.cfg.Editor.MatchBraces {
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
		if b.cfg.Editor.HungryBackspace && b.curs.x >= int(b.cfg.Editor.TabSize) {
			// cut out the last {TAB_SIZE} amount of characters
			// and check em
			tabSize := int(b.cfg.Editor.TabSize)

			// render the line...
			currLine := b.table.Lines[b.curs.y].String()
			before := currLine[b.curs.x-tabSize:]

			if strings.HasPrefix(before, b.makeTab()) {
				// delete {TAB_SIZE} amount of characters
				// from the cursors x pos
				for i := 0; i < int(b.cfg.Editor.TabSize); i++ {

					b.table.Delete(b.curs.y, b.curs.x)

					// FIXME this was -1, maybe we dont
					// want to handle tabs?
					b.moveLeft()
				}
				return
			}
		}

		b.table.Delete(b.curs.y, b.curs.x)

		if len(b.autoComplete.lastRunes) > 0 {
			// pop the last word off in the auto complete
			b.autoComplete.lastRunes = b.autoComplete.lastRunes[:len(b.autoComplete.lastRunes)-1]
		}

		b.moveLeft()
	} else if b.curs.x == 0 && b.curs.y > 0 {
		// FIXME

		// TODO this should set the previous word
		// in the auto complete to the word before the cursor after
		// wrapping back the line?

		// start of line, wrap to previous
		val := b.table.Lines[b.curs.y].String()
		b.table.Insert(val, b.curs.y-1, b.table.Lines[b.curs.y-1].Len())

		b.table.Lines = append(b.table.Lines[:b.curs.y], b.table.Lines[b.curs.y+1:]...)

		b.moveUp()
		b.moveToEndOfLine()
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
		str := b.table.Lines[b.curs.y].String()
		inBounds := (b.curs.x-1 >= 0 && b.curs.x-1 < len(str))

		charWidth := 1
		if inBounds && str[b.curs.x-1] == '\t' {
			charWidth = 4
		}

		b.curs.moveRender(-1, 0, -charWidth, 0)
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

		charWidth := 1
		str := b.table.Lines[b.curs.y].String()[b.curs.x]
		if str == '\t' {
			charWidth = 4
		}

		b.curs.moveRender(1, 0, charWidth, 0)
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

		_, h := b.GetSize()
		visibleLines := int(math.Ceil(float64(h-b.ey) / float64(lastCharH+pad)))

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
		if lastSelection.ex == 0 {
			lastSelection.ey--
			lineLen := b.table.Lines[lastSelection.ey].Len()
			lastSelection.ex = lineLen
		} else {
			lastSelection.ex--
		}

		b.moveLeft()
		break
	case sdl.K_RIGHT:
		lineLen := b.table.Lines[lastSelection.ey].Len()
		if lastSelection.ex == lineLen {
			lastSelection.ey++
			lastSelection.ex = 0
		} else {
			lastSelection.ex++
		}

		b.moveRight()
		break
	case sdl.K_UP:
		lastSelection.ey--
		b.moveUp()
		lastSelection.ex = b.curs.x
		break
	case sdl.K_DOWN:
		lastSelection.ey++
		b.moveDown()
		lastSelection.ex = b.curs.x
		break
	}

	return true
}

// processes a key press. returns if a key was processed
// or not, for example the letter 'a' could run through this
// which is not an action key, therefore we return false
// because it was not processed.
func (b *Buffer) processActionKey(key int) bool {
	if shiftDown {
		switch key {
		case strife.KEY_LEFT:
			fallthrough
		case strife.KEY_RIGHT:
			fallthrough
		case strife.KEY_DOWN:
			fallthrough
		case strife.KEY_UP:
			return b.processSelection(key)
		}
	}

	switch key {
	case strife.KEY_CAPSLOCK:
		capsLockDown = !capsLockDown

	case strife.KEY_ESCAPE:
		return true

	case strife.KEY_RETURN:
		if superDown {
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

	case strife.KEY_BACKSPACE:
		// HACK FIXME
		b.modified = true

		if superDown {
			b.deleteBeforeCursor()
		} else {
			b.deletePrev()
		}

	case strife.KEY_RIGHT:
		currLineLength := b.table.Lines[b.curs.y].Len()

		if superDown {
			for b.curs.x < currLineLength {
				b.moveLeft()
			}
			break
		}

		// FIXME this is weird!
		// this will move to the next blank or underscore
		// character
		if altDown {
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
	case strife.KEY_LEFT:
		if superDown {
			// TODO go to the nearest \t
			// if no \t (i.e. start of line) go to
			// the start of the line!
			b.curs.gotoStart()
		} else if altDown {
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
	case strife.KEY_UP:
		if altDown {
			return b.swapLineUp()
		}

		if superDown {
			// go to the start of the file
		}

		b.moveUp()

	case strife.KEY_DOWN:
		if altDown {
			return b.swapLineDown()
		}

		if superDown {
			// go to the end of the file
		}

		b.moveDown()

	case strife.KEY_TAB:
		// HACK FIXME
		b.modified = true

		if b.cfg.Editor.TabsAreSpaces {
			// make an empty rune array of TAB_SIZE, cast to string
			// and insert it.
			tab := b.makeTab()
			b.table.Insert(tab, b.curs.y, b.curs.x)
			for i := 0; i < len(tab); i++ {
				b.moveRight()
			}
		} else {
			b.table.Insert(string('\t'), b.curs.y, b.curs.x)
			b.moveRight()
		}

	case strife.KEY_END:
		currLine := b.table.Lines[b.curs.y]
		if b.curs.x < currLine.Len() {
			distToMove := currLine.Len() - b.curs.x
			for i := 0; i < distToMove; i++ {
				b.moveRight()
			}
		}

	case strife.KEY_HOME:
		if b.curs.x > 0 {
			b.moveToStartOfLine()
		}

		// TODO remove since this is handled in the keymap!
	case strife.KEY_PAGEUP:
		b.scrollUp(DefaultScrollAmount)
		for i := 0; i < DefaultScrollAmount; i++ {
			b.moveUp()
		}
	case strife.KEY_PAGEDOWN:
		b.scrollDown(DefaultScrollAmount)
		for i := 0; i < DefaultScrollAmount; i++ {
			b.moveDown()
		}

	case strife.KEY_DELETE:
		b.deleteNext()

	case strife.KEY_LGUI:
		fallthrough
	case strife.KEY_RGUI:
		fallthrough

	case strife.KEY_LALT:
		fallthrough
	case strife.KEY_RALT:
		fallthrough

	case strife.KEY_LCTRL:
		fallthrough
	case strife.KEY_RCTRL:
		fallthrough

	case strife.KEY_LSHIFT:
		fallthrough
	case strife.KEY_RSHIFT:
		return true

	default:
		return false
	}

	return true
}

var (
	shiftDown    = false
	superDown    = false // cmd on mac, ctrl on windows
	controlDown  = false // what is this on windows?
	altDown      = false // option on mac
	capsLockDown = false
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
	for i := 0; i < int(b.cfg.Editor.TabSize); i++ {
		blah = append(blah, ' ')
	}
	return string(blah)
}

func (b *Buffer) HandleEvent(evt strife.StrifeEvent) {
	switch event := evt.(type) {
	case *strife.MouseWheelEvent:
		if event.Y > 0 {
			b.scrollDown(DefaultScrollAmount)
		}
		if event.Y < 0 {
			b.scrollUp(DefaultScrollAmount)
		}
	}
}

var lastCursorDraw = time.Now()
var renderFlashingCursor = true

func (b *Buffer) processLeftClick() {
	// here we set the cursor y position
	// based off the click location
	xPos, yPos := strife.MouseCoords()

	yPosToLine := (((yPos) / (lastCharH + pad)) + 1) + b.cam.y
	xPosToLine := ((xPos - b.ex) / (lastCharW)) + b.cam.x

	fmt.Println(yPos, " line ", yPosToLine, " - ", xPos, " char ", xPosToLine)

	b.gotoLine(int64(yPosToLine))

	// we should be at the start of the line but lets
	// move there anyways just to make sure

	b.moveToStartOfLine()

	for i := 0; i < xPosToLine; i++ {
		b.moveRight()
	}
}

const ySpeed = 1

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

	shiftDown = strife.KeyPressed(sdl.K_LSHIFT) || strife.KeyPressed(sdl.K_RSHIFT)
	superDown = strife.KeyPressed(sdl.K_LGUI) || strife.KeyPressed(sdl.K_RGUI)
	altDown = strife.KeyPressed(sdl.K_LALT) || strife.KeyPressed(sdl.K_RALT)
	controlDown = strife.KeyPressed(sdl.K_LCTRL) || strife.KeyPressed(sdl.K_RCTRL)

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
				// TODO invalidate selection when necessary

			}

			return true
		}

		textEntered := b.processTextInput(rune(keyCode))
		if textEntered {
			return true
		}
	}

	// FIXME for now this is only enabled in debug mode.
	if b.cfg.Cursor.Flash && cfg.DebugMode {
		if time.Now().Sub(lastCursorDraw) >= time.Duration(b.cfg.Cursor.FlashRate)*time.Millisecond {
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
	background uint32
	foreground uint32
	length     int
}

// dimensions of the last character we rendered
var lastCharW, lastCharH int

func lexFindMatches(matches *map[int]syntaxRuneInfo, currLine string, toMatch map[string]bool, bg uint32, fg uint32) {
	lexer := lex.New(currLine)
	tokenStream := lexer.Tokenize()
	for _, tok := range tokenStream {
		if _, ok := toMatch[tok.Lexeme]; ok {
			(*matches)[tok.Start] = syntaxRuneInfo{bg, fg, len(tok.Lexeme)}
		}
	}
}

type charColouring struct {
	bg uint32
	fg uint32
}

const syntaxCacheEvictionTime = time.Second * 30

type syntaxInfo struct {
	data      map[int]syntaxRuneInfo
	cacheTime time.Time
}

// hacky. FIXME with a proper solution.
var syntaxCache = map[string]syntaxInfo{}
var syntaxHighlights = 0

// syntaxHighlightLine will highlight the given string
// it returns a map of column positions -> bg/fg colours + the length of the colouring information.
//
// this could do with a lot of optimisation. this is executed pretty much constantly
// and will run a regex on the line
//
// just as I wrote this comment though I added some very naive caching that will probably
// break in a few edge cases!
func (b *Buffer) syntaxHighlightLine(currLine string) map[int]syntaxRuneInfo {
	if info, ok := syntaxCache[currLine]; ok {
		if time.Now().Sub(info.cacheTime) < syntaxCacheEvictionTime {
			return info.data
		}
	}

	syntaxHighlights++

	matches := map[int]syntaxRuneInfo{}

	subjects := make([]*cfg.SyntaxCriteria, len(b.languageInfo.Syntax))
	colours := make([]charColouring, len(b.languageInfo.Syntax))

	idx := 0
	for _, criteria := range b.languageInfo.Syntax {
		colours[idx] = charColouring{
			criteria.Background,
			criteria.Foreground,
		}
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

						colouring := colours[syntaxIndex]
						matches[charIndex+matched[0]] = syntaxRuneInfo{
							colouring.bg,
							colouring.fg,
							matchedStrLen,
						}

						charIndex += matchedStrLen
						continue
					}
				}
			}
		} else {
			colouring := colours[syntaxIndex]
			lexFindMatches(&matches, currLine, syntax.MatchList, colouring.bg, colouring.fg)
		}
	}

	syntaxCache[currLine] = syntaxInfo{
		data:      matches,
		cacheTime: time.Now(),
	}

	return matches
}

func (b *Buffer) renderAt(ctx *strife.Renderer, rx int, ry int) {
	// TODO load this from config files!

	x, y := b.GetPos()
	w, h := b.GetSize()

	// BACKGROUND
	ctx.SetColor(strife.HexRGB(b.buffOpts.background))
	ctx.Rect(x, y, w, h, strife.Fill)

	if b.cfg.Editor.HighlightLine && b.HasFocus() {
		ctx.SetColor(strife.HexRGB(b.buffOpts.highlightLine)) // highlight_line_col?

		highlightLinePosY := b.ey + (ry + b.curs.ry*(lastCharH+pad)) - (b.cam.y * (lastCharH + pad))
		highlightLinePosX := b.ex + rx

		width := w - b.ex
		height := (lastCharH + pad) - b.ey
		ctx.Rect(highlightLinePosX, highlightLinePosY, width, height, strife.Fill)
	}

	var visibleLines, visibleChars int = 50, -1

	// HACK
	// force a render off screen
	// so we can calculate the size of characters
	{
		if int(lastCharH) == 0 || int(lastCharW) == 0 {
			lastCharW, lastCharH = ctx.Text("_", -50, -50)
		}
	}

	// lastCharH > 0 means we have done
	// a render.
	if int(lastCharH) > 0 && int(h) != 0 {
		// render an extra three lines just
		// so we dont cut anything off if its
		// not evenly divisible
		visibleLines = (int(h-b.ey) / int(lastCharH)) + 3
	}

	// calculate how many chars we can fit
	// on the screen.
	if int(lastCharW) > 0 && int(w) != 0 {
		visibleChars = (int(w-b.ex) / int(lastCharW)) - 3
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
		lastSelection.renderAt(ctx, x+b.ex, y+b.ey)
	}

	// calculate cursor sizes... does
	// this have to be done every frame?
	{
		cursorWidth := b.cfg.Cursor.GetCaretWidth()
		if cursorWidth == -1 {
			cursorWidth = lastCharW
		}
		cursorHeight := lastCharH + pad

		b.curs.SetSize(cursorWidth, cursorHeight)
	}

	// render the ol' cursor
	if b.HasFocus() && (renderFlashingCursor || b.curs.moving) && b.cfg.Cursor.Draw {
		b.curs.Render(ctx, rx, ry)
	}

	numLines := len(b.table.Lines)

	var yCol int
	for lineNum, rope := range b.table.Lines[start:upper] {
		currLine := []rune(rope.String())

		if visibleChars >= 0 {
			// slice the visible characters only.
			currLine = currLine[:min(visibleChars, len(currLine))]
		}

		// char index => colour
		var matches map[int]syntaxRuneInfo
		if b.languageInfo != nil && len(currLine) > 0 {
			matches = b.syntaxHighlightLine(string(currLine))
		}

		var colorStack []charColouring

		// TODO move this into a struct
		// or something.

		// the x position of the _character_
		var xCol int

		for idx, char := range currLine {
			switch char {

			// 13 is a carriage return
			case 13:
				continue

			case '\n':
				xCol = 0
				yCol++
				continue
			case '\t':
				xCol += b.cfg.Editor.TabSize
				continue
			}

			xCol++

			if info, ok := matches[idx]; ok {
				if colorStack == nil || len(colorStack) == 0 {
					colorStack = make([]charColouring, info.length)
					for i := 0; i < info.length; i++ {
						colorStack[i] = charColouring{info.background, info.foreground}
					}
				}
			}

			characterColor := charColouring{0, b.buffOpts.foreground}

			if len(colorStack) > 0 {
				a := colorStack[len(colorStack)-1]
				colorStack = colorStack[:len(colorStack)-1]
				characterColor = a
			}

			if b.HasFocus() && (b.curs.x-b.cam.x) == (xCol-1) && (b.curs.y-b.cam.y) == yCol {
				characterColor = charColouring{0, b.buffOpts.cursorInvert}
			}

			lineHeight := lastCharH + pad
			xPos := b.ex + (rx + ((xCol - 1) * lastCharW))
			yPos := b.ey + (ry + (yCol * lineHeight)) + halfPad

			// todo render background

			ctx.SetColor(strife.HexRGB(characterColor.fg))
			lastCharW, lastCharH = ctx.Text(string(char), xPos, yPos)

			if cfg.DebugMode {
				ctx.SetColor(strife.HexRGB(0xff00ff))
				ctx.Rect(xPos, yPos, lastCharW, lastCharH, strife.Line)
			}
		}

		if b.cfg.Editor.ShowLineNumbers {
			gutterPadPx := 10

			// how many chars we need
			numLinesCharWidth := len(string(numLines)) + 2

			gutterWidth := lastCharW*numLinesCharWidth + (gutterPadPx * 2)

			lineHeight := lastCharH + pad
			yPos := ((ry + yCol) * lineHeight) + halfPad

			// render the line numbers
			ctx.SetColor(strife.HexRGB(b.buffOpts.lineNumBackground))
			ctx.Rect(rx, yPos, gutterWidth, h, strife.Fill)

			if cfg.DebugMode {
				ctx.SetColor(strife.HexRGB(0xff00ff))
				ctx.Rect(rx, yPos, gutterWidth, h, strife.Line)
			}

			ctx.SetColor(strife.HexRGB(b.buffOpts.lineNumForeground))
			ctx.Text(fmt.Sprintf("%*d", numLinesCharWidth, (start+lineNum)+1), rx+gutterPadPx, yPos)

			b.ex = gutterWidth
		}

		yCol++
	}

	if b.autoComplete.hasSuggestions() {
		xPos := b.ex + (rx + b.curs.rx*lastCharW) - (b.cam.x * lastCharW)
		yPos := b.ey + (ry + b.curs.ry*b.curs.height) - (b.cam.y * b.curs.height)

		autoCompleteBoxHeight := len(b.autoComplete.suggestions) * b.curs.height
		yPos = yPos - autoCompleteBoxHeight

		b.autoComplete.renderAt(xPos, yPos, ctx)
	}

	if cfg.DebugMode {
		ctx.SetColor(strife.HexRGB(0xff00ff))
		ctx.Rect(b.ex+rx, b.ey+ry, w-b.ex, h-b.ey, strife.Line)
	}
}

func (b *Buffer) OnRender(ctx *strife.Renderer) {
	x, y := b.GetPos()

	ctx.SetFont(b.buffOpts.font)
	b.renderAt(ctx, x, y)
}
