package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Component interface {
	Translate(x, y int32)

	Init()
	Update()
	OnRender(*sdl.Renderer)
	Dispose()

	AddComponent(c Component)
	GetComponents() []Component

	GetInputHandler() *InputHandler
	SetInputHandler(h *InputHandler)
}

type BaseComponent struct {
	x, y          int32
	components    []Component
	input_handler *InputHandler
}

func Render(c Component, ctx *sdl.Renderer) {
	c.OnRender(ctx)
	for _, child := range c.GetComponents() {
		Render(child, ctx)
	}
}

func (b *BaseComponent) Translate(x, y int32) {
	b.x += x
	b.y += y
	for _, c := range b.components {
		c.Translate(x, y)
	}
}

func (b *BaseComponent) GetComponents() []Component {
	return b.components
}

func (b *BaseComponent) AddComponent(c Component) {
	b.components = append(b.components, c)
}

func (b *BaseComponent) SetInputHandler(i *InputHandler) {
	b.input_handler = i
}

func (b *BaseComponent) GetInputHandler() *InputHandler {
	return b.input_handler
}
