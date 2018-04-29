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
	buff              *Buffer
	parentBuff        *Buffer
	recentSuggestions *[]suggestion
}

const suggestionBoxHeight = 128

type suggestion struct {
	name string
}

func (s *suggestion) render(x, y int, ctx *strife.Renderer) {
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

	return palette
}

func (b *CommandPalette) OnInit() {
	b.buff.Translate(b.x, b.y)
	b.buff.Resize(b.w, b.h)
}

func (b *CommandPalette) processCommand() {
	tokenizedLine := strings.Split(b.buff.contents[0].String(), " ")
	command := tokenizedLine[0]

	action, exists := actions[command]
	if !exists {
		return
	}

	action.proc(b.parentBuff)
}

func (b *CommandPalette) calculateSuggestions() {
	tokenizedLine := strings.Split(b.buff.contents[0].String(), " ")
	command := tokenizedLine[0]

	ranks := fuzzy.RankFind(command, commandSet)
	suggestions := make([]suggestion, len(ranks))

	for _, r := range ranks {
		cmdName := commandSet[r.Index]
		suggestions = append(suggestions, suggestion{cmdName})
	}

	b.recentSuggestions = &suggestions
}

func (b *CommandPalette) clearInput() {
	actions["delete_line"].proc(b.buff)
}

func (b *CommandPalette) OnUpdate() bool {
	if !b.HasFocus() {
		return false
	}

	override := func(k int) bool {
		if k != sdl.K_RETURN && k != sdl.K_ESCAPE {
			b.calculateSuggestions()
			return false
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
			sugg.render(b.x, b.y+(i*suggestionBoxHeight), ctx)
		}
	}
}

func (b *CommandPalette) OnDispose() {
	log.Println("poop diddity scoop")
}
