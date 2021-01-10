package buff

import (
	"log"
	"strings"

	"github.com/felixangell/phi/internal/cfg"
	"github.com/felixangell/phi/internal/gui"
	"github.com/felixangell/phi/internal/lex"
	"github.com/felixangell/strife"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/veandco/go-sdl2/sdl"
)

var commandSet []string

type CommandPalette struct {
	gui.BaseComponent
	buff       *Buffer
	parentBuff *Buffer
	conf       *cfg.PhiEditorConfig
	parent     *BufferView

	pathToIndex map[string]int

	suggestionIndex   int
	recentSuggestions *[]suggestion
}

func (p *CommandPalette) SetFocus(focus bool) {
	p.buff.SetFocus(focus)
	p.BaseComponent.SetFocus(focus)
}

var suggestionBoxHeight, suggestionBoxWidth = 48, 0

type suggestion struct {
	parent *CommandPalette
	name   string
}

func (s *suggestion) renderHighlighted(x, y int, ctx *strife.Renderer) {
	conf := s.parent.conf.Theme.Palette

	border := 5
	ctx.SetColor(strife.HexRGB(conf.Outline))
	ctx.Rect(x-border, y-border, suggestionBoxWidth+(border*2), suggestionBoxHeight+(border*2), strife.Fill)

	ctx.SetColor(strife.HexRGB(conf.Suggestion.SelectedBackground))
	ctx.Rect(x, y, suggestionBoxWidth, suggestionBoxHeight, strife.Fill)

	ctx.SetColor(strife.HexRGB(conf.Suggestion.SelectedForeground))

	// FIXME strife library needs something to get
	// text width and heights... for now we render offscreen to measure... lol
	_, h := ctx.Text("foo", -500000, -50000)

	yOffs := (suggestionBoxHeight / 2) - (h / 2)
	ctx.Text(s.name, x+border, y+yOffs)
}

func (s *suggestion) render(x, y int, ctx *strife.Renderer) {
	// wewlad
	conf := s.parent.conf.Theme.Palette

	border := 5
	ctx.SetColor(strife.HexRGB(conf.Outline))
	ctx.Rect(x-border, y-border, suggestionBoxWidth+(border*2), suggestionBoxHeight+(border*2), strife.Fill)

	ctx.SetColor(strife.HexRGB(conf.Suggestion.Background))
	ctx.Rect(x, y, suggestionBoxWidth, suggestionBoxHeight, strife.Fill)

	ctx.SetColor(strife.HexRGB(conf.Suggestion.Foreground))

	_, textHeight := ctx.GetStringDimension("foo")

	yOffs := (suggestionBoxHeight / 2) - (textHeight / 2)
	ctx.Text(s.name, x+border, y+yOffs)
}

func NewCommandPalette(conf cfg.PhiEditorConfig, view *BufferView) *CommandPalette {
	conf.Editor.ShowLineNumbers = false
	conf.Editor.HighlightLine = false

	newSize := int(float64(conf.Editor.FontSize) * cfg.ScaleFactor)
	paletteFont, err := gui.GetDefaultFont().DeriveFont(newSize)
	if err != nil {
		panic(err)
	}

	palette := &CommandPalette{
		conf:   &conf,
		parent: view,
		buff: NewBuffer(&conf, BufferConfig{
			conf.Theme.Palette.Background,
			conf.Theme.Palette.Foreground,

			conf.Theme.Palette.Cursor,
			conf.Theme.Palette.Cursor, // TODO invert

			conf.Theme.HighlightLineBackground,

			// we dont show line numbers
			// so these aren't necessary
			0x0, 0x0,

			paletteFont,
		}, nil, 0),
		parentBuff: nil,
	}
	palette.buff.appendLine("")

	vW, _ := view.GetSize()

	// set the palette to be 1/3rd of the size of the view
	palette.Resize(vW/3, int(float64(newSize)*1.5))

	// translate it to the centre of the view.
	pW, pH := palette.GetSize()

	palette.Translate((vW/2)-(pW/2), 10)

	// the buffer is not rendered
	// relative to the palette so we have to set its position
	palette.buff.Resize(pW, pH)
	palette.buff.Translate((vW/2)-(pW/2), 10)

	// this is technically a hack. this ex is an xoffset
	// for the line numbers but we're going to use it for
	// general border offsets. this is a real easy fixme
	// for general code clean but maybe another day!
	palette.buff.ex = 5
	palette.buff.ey = 0

	suggestionBoxWidth = pW

	return palette
}

func (b *CommandPalette) processCommand() {
	input := b.buff.table.Lines[0].String()
	tokens := lex.New(input).Tokenize()

	if len(tokens) <= 1 {
		return
	}

	if tokens[0].Equals("!") && tokens[1].IsType(lex.Word) {
		commandName := tokens[1].Lexeme
		args := tokens[2:]
		log.Println("executing action", commandName, "with arguments", args)
		ExecuteCommandIfExist(commandName, b.parent, args...)
	}

	if index, ok := b.pathToIndex[input]; ok {
		b.parent.setFocusTo(index)
	}
}

func (b *CommandPalette) calculateCommandSuggestions() {
	input := b.buff.table.Lines[0].String()
	input = strings.TrimSpace(input)

	tokenizedLine := strings.Split(input, " ")

	if strings.Compare(tokenizedLine[0], "!") != 0 {
		return
	}

	if len(tokenizedLine) == 1 {
		suggestions := make([]suggestion, len(commandSet))

		for idx, cmd := range commandSet {
			suggestions[idx] = suggestion{b, "! " + cmd}
		}

		b.recentSuggestions = &suggestions
		return
	}

	// slice off the command thingy
	tokenizedLine = tokenizedLine[1:]

	command := tokenizedLine[0]

	if command == "" {
		b.recentSuggestions = nil
		return
	}

	ranks := fuzzy.RankFind(command, commandSet)
	suggestions := []suggestion{}

	for _, r := range ranks {
		cmdName := commandSet[r.OriginalIndex]
		if cmdName == "" {
			continue
		}
		cmdName = "! " + cmdName
		suggestions = append(suggestions, suggestion{b, cmdName})
	}

	b.recentSuggestions = &suggestions
}

func (b *CommandPalette) calculateSuggestions() {
	input := b.buff.table.Lines[0].String()
	input = strings.TrimSpace(input)

	if len(input) == 0 {
		return
	}

	if input[0] == '!' {
		b.calculateCommandSuggestions()
		return
	}

	// FIXME
	// fill it with the currently open files!

	openFiles := make([]string, len(b.parent.buffers))

	b.pathToIndex = map[string]int{}

	for i, pane := range b.parent.buffers {
		path := pane.Buff.filePath
		openFiles[i] = path
		b.pathToIndex[path] = i
	}

	ranks := fuzzy.RankFind(input, openFiles)
	var suggestions []suggestion
	for _, r := range ranks {
		pane := b.parent.buffers[r.OriginalIndex]
		if pane != nil {
			suggestions = append(suggestions, suggestion{
				b,
				pane.Buff.filePath,
			})
		}
	}

	b.recentSuggestions = &suggestions
}

func (b *CommandPalette) scrollSuggestion(dir int) {
	if b.recentSuggestions != nil {
		b.suggestionIndex += dir

		if b.suggestionIndex < 0 {
			b.suggestionIndex = len(*b.recentSuggestions) - 1
		} else if b.suggestionIndex >= len(*b.recentSuggestions) {
			b.suggestionIndex = 0
		}
	}
}

func (b *CommandPalette) clearInput() {
	b.buff.deleteLine()
}

func (b *CommandPalette) setToSuggested() {
	if b.recentSuggestions == nil {
		return
	}

	// set the buffer
	suggestions := *b.recentSuggestions
	sugg := suggestions[b.suggestionIndex]
	b.buff.setLine(0, sugg.name)

	// remove all suggestions
	b.recentSuggestions = nil
	b.suggestionIndex = -1
}

func (b *CommandPalette) OnUpdate() bool {
	if !b.HasFocus() {
		return false
	}

	override := func(key int) bool {
		switch key {

		case sdl.K_UP:
			b.scrollSuggestion(-1)
			return false
		case sdl.K_DOWN:
			b.scrollSuggestion(1)
			return false

		// any other key we calculate
		// the suggested commands
		default:
			b.suggestionIndex = -1
			b.calculateSuggestions()
			return false

		case sdl.K_RETURN:
			// we have a suggestion so let's
			// fill the buffer with that instead!
			if b.suggestionIndex != -1 {
				b.setToSuggested()
				return true
			}

			b.processCommand()
			break

		case sdl.K_ESCAPE:
			break
		}

		b.parent.hidePalette()
		return true
	}
	return b.buff.processInput(override)
}

func (b *CommandPalette) OnRender(ctx *strife.Renderer) {
	if !b.HasFocus() {
		return
	}

	conf := b.conf.Theme.Palette

	x, y := b.GetPos()
	w, h := b.GetSize()

	border := 5
	xPos := x - border
	yPos := y - border
	paletteWidth := w + (border * 2)
	paletteHeight := h + (border * 2)

	ctx.SetColor(strife.HexRGB(conf.Outline))
	ctx.Rect(xPos, yPos, paletteWidth, paletteHeight, strife.Fill)

	_, charHeight := ctx.Text("foo", -5000, -5000)
	b.buff.ey = (suggestionBoxHeight / 2) - (charHeight / 2)

	b.buff.OnRender(ctx)

	if b.recentSuggestions != nil {
		for i, suggestion := range *b.recentSuggestions {
			if b.suggestionIndex != i {
				suggestion.render(x, y+((i+1)*(suggestionBoxHeight+border)), ctx)
			} else {
				suggestion.renderHighlighted(x, y+((i+1)*(suggestionBoxHeight+border)), ctx)
			}
		}
	}

	if cfg.DebugMode {
		ctx.SetColor(strife.HexRGB(cfg.DebugModeRenderColour))
		ctx.Rect(xPos, yPos, paletteWidth, paletteHeight, strife.Line)
	}
}
