package gui

import (
	"log"
	"strings"

	"github.com/felixangell/fuzzysearch/fuzzy"
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
	"github.com/veandco/go-sdl2/sdl"
)

var commandSet []string

func init() {
	commandSet = make([]string, len(actions))
	for _, action := range actions {
		commandSet = append(commandSet, action.name)
	}
}

type CommandPalette struct {
	BaseComponent
	buff       *Buffer
	parentBuff *Buffer
	conf       *cfg.TomlConfig
	parent     *View

	suggestionIndex   int
	recentSuggestions *[]suggestion
}

var suggestionBoxHeight, suggestionBoxWidth = 48, 0

type suggestion struct {
	parent *CommandPalette
	name   string
}

func (s *suggestion) renderHighlighted(x, y int, ctx *strife.Renderer) {
	// wewlad
	conf := s.parent.conf.Theme.Palette

	border := 5
	ctx.SetColor(strife.HexRGB(conf.Outline))
	ctx.Rect(x-border, y-border, suggestionBoxWidth+(border*2), suggestionBoxHeight+(border*2), strife.Fill)

	ctx.SetColor(strife.HexRGB(conf.Suggestion.Selected_Background))
	ctx.Rect(x, y, suggestionBoxWidth, suggestionBoxHeight, strife.Fill)

	ctx.SetColor(strife.HexRGB(conf.Suggestion.Selected_Foreground))

	// FIXME strife library needs something to get
	// text width and heights... for now we render offscreen to measure... lol
	_, h := ctx.String("foo", -500000, -50000)

	yOffs := (suggestionBoxHeight / 2) - (h / 2)
	ctx.String(s.name, x+border, y+yOffs)
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

	// FIXME strife library needs something to get
	// text width and heights... for now we render offscreen to measure... lol
	_, h := ctx.String("foo", -500000, -50000)

	yOffs := (suggestionBoxHeight / 2) - (h / 2)
	ctx.String(s.name, x+border, y+yOffs)
}

func NewCommandPalette(conf cfg.TomlConfig, view *View) *CommandPalette {
	conf.Editor.Show_Line_Numbers = false
	conf.Editor.Highlight_Line = false

	palette := &CommandPalette{
		conf:   &conf,
		parent: view,
		buff: NewBuffer(&conf, BufferConfig{
			conf.Theme.Palette.Background,
			conf.Theme.Palette.Foreground,

			conf.Theme.Palette.Cursor,
			conf.Theme.Palette.Cursor, // TODO invert

			// we dont show line numbers
			// so these aren't necessary
			0x0, 0x0,
		}, nil, 0),
		parentBuff: nil,
	}
	palette.buff.appendLine("")

	palette.Resize(view.w/3, 48)
	palette.Translate((view.w/2)-(palette.w/2), 10)

	// the buffer is not rendered
	// relative to the palette so we have to set its position
	palette.buff.Resize(palette.w, palette.h)
	palette.buff.Translate((view.w/2)-(palette.w/2), 10)

	// this is technically a hack. this ex is an xoffset
	// for the line numbers but we're going to use it for
	// general border offsets. this is a real easy fixme
	// for general code clean but maybe another day!
	palette.buff.ex = 5
	palette.buff.ey = 0

	suggestionBoxWidth = palette.w

	return palette
}

func (b *CommandPalette) processCommand() {
	tokenizedLine := strings.Split(b.buff.contents[0].String(), " ")
	command := tokenizedLine[0]

	log.Println("command palette: ", tokenizedLine)

	action, exists := actions[command]
	if !exists {
		return
	}

	action.proc(b.parent, tokenizedLine[1:])
}

func (b *CommandPalette) calculateSuggestions() {
	tokenizedLine := strings.Split(b.buff.contents[0].String(), " ")
	command := tokenizedLine[0]

	if command == "" {
		b.recentSuggestions = nil
		return
	}

	ranks := fuzzy.RankFind(command, commandSet)
	suggestions := []suggestion{}

	for _, r := range ranks {
		cmdName := commandSet[r.Index]
		if cmdName == "" {
			continue
		}
		suggestions = append(suggestions, suggestion{b, cmdName})
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

	border := 5
	xPos := b.x - border
	yPos := b.y - border
	paletteWidth := b.w + (border * 2)
	paletteHeight := b.h + (border * 2)

	ctx.SetColor(strife.HexRGB(conf.Outline))
	ctx.Rect(xPos, yPos, paletteWidth, paletteHeight, strife.Fill)

	_, charHeight := ctx.String("foo", -5000, -5000)
	b.buff.ey = (suggestionBoxHeight / 2) - (charHeight / 2)

	b.buff.OnRender(ctx)

	if b.recentSuggestions != nil {
		for i, sugg := range *b.recentSuggestions {
			if b.suggestionIndex != i {
				sugg.render(b.x, b.y+((i+1)*(suggestionBoxHeight+border)), ctx)
			} else {
				sugg.renderHighlighted(b.x, b.y+((i+1)*(suggestionBoxHeight+border)), ctx)
			}
		}
	}

	if DEBUG_MODE {
		ctx.SetColor(strife.HexRGB(0xff00ff))
		ctx.Rect(xPos, yPos, paletteWidth, paletteHeight, strife.Line)
	}
}
