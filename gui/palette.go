package gui

import (
	"github.com/felixangell/fuzzysearch/fuzzy"
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"strings"
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

	suggestionIndex   int
	recentSuggestions *[]suggestion
}

var suggestionBoxHeight, suggestionBoxWidth = 48, 0

type suggestion struct {
	name string
}

func (s *suggestion) renderHighlighted(x, y int, ctx *strife.Renderer) {
	ctx.SetColor(strife.Blue)
	ctx.Rect(x, y, suggestionBoxWidth, suggestionBoxHeight, strife.Fill)

	ctx.SetColor(strife.White)
	ctx.String(s.name, x, y)
}

func (s *suggestion) render(x, y int, ctx *strife.Renderer) {
	ctx.SetColor(strife.Red)
	ctx.Rect(x, y, suggestionBoxWidth, suggestionBoxHeight, strife.Fill)

	ctx.SetColor(strife.White)
	ctx.String(s.name, x, y)
}

func NewCommandPalette(conf cfg.TomlConfig, view *View) *CommandPalette {
	conf.Editor.Show_Line_Numbers = false
	conf.Editor.Highlight_Line = false

	palette := &CommandPalette{
		buff:       NewBuffer(&conf, nil, 0),
		parentBuff: nil,
	}
	palette.buff.appendLine("")

	palette.Resize(view.w/3, 48)
	palette.Translate((view.w/2)-(palette.w/2), 10)
	palette.buff.Resize(palette.w, palette.h)
	palette.buff.Translate((view.w/2)-(palette.w/2), 10)

	suggestionBoxWidth = palette.w

	return palette
}

func (b *CommandPalette) OnInit() {
	b.buff.Translate(b.x, b.y)
	b.buff.Resize(b.w, b.h)
}

func (b *CommandPalette) processCommand() {
	tokenizedLine := strings.Split(b.buff.contents[0].String(), " ")
	command := tokenizedLine[0]

	log.Println(tokenizedLine)

	action, exists := actions[command]
	if !exists {
		return
	}

	action.proc(b.parentBuff)
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
		suggestions = append(suggestions, suggestion{cmdName})
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
	actions["delete_line"].proc(b.buff)
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

			fallthrough
		case sdl.K_ESCAPE:
			break
		}

		b.processCommand()
		b.parentBuff.parent.hidePalette()
		return true
	}
	return b.buff.doUpdate(override)
}

func (b *CommandPalette) OnRender(ctx *strife.Renderer) {
	if !b.HasFocus() {
		return
	}

	border := 5

	ctx.SetColor(strife.White)
	ctx.Rect(b.x-border, b.y-border, b.w+(border*2), b.h+(border*2), strife.Fill)

	b.buff.OnRender(ctx)

	if b.recentSuggestions != nil {
		for i, sugg := range *b.recentSuggestions {
			if b.suggestionIndex != i {
				sugg.render(b.x, b.y+((i+1)*suggestionBoxHeight), ctx)
			} else {
				sugg.renderHighlighted(b.x, b.y+((i+1)*suggestionBoxHeight), ctx)
			}
		}
	}
}

func (b *CommandPalette) OnDispose() {
	log.Println("poop diddity scoop")
}
