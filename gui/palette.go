package gui

import (
	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"strings"
)

type CommandPalette struct {
	BaseComponent
	HasFocus   bool
	buff       *Buffer
	parentBuff *Buffer
}

func NewCommandPalette(conf cfg.TomlConfig) *CommandPalette {
	conf.Editor.Show_Line_Numbers = false
	conf.Editor.Highlight_Line = false

	palette := &CommandPalette{
		buff:       NewBuffer(&conf, nil, 0),
		parentBuff: nil,
		HasFocus:   false,
	}
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

	action(b.parentBuff)
}

func (b *CommandPalette) OnUpdate() bool {
	if !b.HasFocus {
		return true
	}

	override := func(k int) bool {
		if k != sdl.K_RETURN {
			return false
		}

		b.processCommand()

		// regardless we close the command
		// palette and re-focus on the buffer
		// that we transferred from.
		b.parentBuff.SetInputHandler(b.inputHandler)
		b.parentBuff.HasFocus = true

		return true
	}
	return b.buff.doUpdate(override)
}

func (b *CommandPalette) OnRender(ctx *strife.Renderer) {
	if !b.HasFocus {
		return
	}

	border := 5

	ctx.SetColor(strife.White)
	ctx.Rect(b.x-border, b.y-border, b.w+(border*2), b.h+(border*2), strife.Fill)

	b.buff.OnRender(ctx)
}

func (b *CommandPalette) OnDispose() {
	log.Println("poop diddity scoop")
}
