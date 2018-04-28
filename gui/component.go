package gui

import (
	"github.com/felixangell/strife"
)

type Component interface {
	SetPosition(x, y int)
	Translate(x, y int)
	Resize(w, h int)

	OnInit()
	OnUpdate() bool
	OnRender(*strife.Renderer)
	OnDispose()

	NumComponents() int
	AddComponent(c Component)
	GetComponents() []Component

	HandleEvent(evt strife.StrifeEvent)

	GetInputHandler() *InputHandler
	SetInputHandler(h *InputHandler)
}

type BaseComponent struct {
	x, y         int
	w, h         int
	inputHandler *InputHandler
}

func (b *BaseComponent) HandleEvent(evt strife.StrifeEvent) {
	// NOP
}

func (b *BaseComponent) SetPosition(x, y int) {
	b.x = x
	b.y = y
}

func (b *BaseComponent) Translate(x, y int) {
	b.x += x
	b.y += y
}

func (b *BaseComponent) Resize(w, h int) {
	b.w = w
	b.h = h
}

func (b *BaseComponent) SetInputHandler(i *InputHandler) {
	b.inputHandler = i
}

func (b *BaseComponent) GetInputHandler() *InputHandler {
	return b.inputHandler
}
