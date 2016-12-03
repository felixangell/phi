package gui

import (
	"github.com/felixangell/strife"
)

type Component interface {
	Translate(x, y int)
	Resize(w, h int)

	OnInit()
	OnUpdate()
	OnRender(*strife.Renderer)
	OnDispose()

	AddComponent(c Component)
	GetComponents() []Component

	GetInputHandler() *InputHandler
	SetInputHandler(h *InputHandler)
}

type BaseComponent struct {
	x, y         int
	w, h         int
	components   []Component
	inputHandler *InputHandler
}

func (b *BaseComponent) Translate(x, y int) {
	b.x += x
	b.y += y
	for _, c := range b.components {
		c.Translate(x, y)
	}
}

func (b *BaseComponent) Resize(w, h int) {
	b.w = w
	b.h = h
}

func (b *BaseComponent) GetComponents() []Component {
	return b.components
}

func (b *BaseComponent) AddComponent(c Component) {
	b.components = append(b.components, c)
	c.SetInputHandler(b.inputHandler)
	Init(c)
}

func (b *BaseComponent) SetInputHandler(i *InputHandler) {
	b.inputHandler = i
}

func (b *BaseComponent) GetInputHandler() *InputHandler {
	return b.inputHandler
}

func Update(c Component) {
	c.OnUpdate()
	for _, child := range c.GetComponents() {
		Update(child)
	}
}

func Render(c Component, ctx *strife.Renderer) {
	c.OnRender(ctx)
	for _, child := range c.GetComponents() {
		Render(child, ctx)
	}
}

func Init(c Component) {
	c.OnInit()
	for _, child := range c.GetComponents() {
		Init(child)
	}
}

func Dispose(c Component) {
	c.OnDispose()
	for _, child := range c.GetComponents() {
		Dispose(child)
	}
}
