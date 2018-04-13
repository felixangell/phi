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

	GetInputHandler() *InputHandler
	SetInputHandler(h *InputHandler)
}

type BaseComponent struct {
	x, y          int
	w, h          int
	components    []Component
	numComponents int
	inputHandler  *InputHandler
}

func (b *BaseComponent) DeleteComponent(index int) {
	b.components[index] = nil
	b.numComponents--
}

func (b *BaseComponent) SetPosition(x, y int) {
	b.x = x
	b.y = y
	for _, c := range b.components {
		if c == nil {
			continue
		}

		c.SetPosition(x, y)
	}
}

func (b *BaseComponent) NumComponents() int {
	return b.numComponents
}

func (b *BaseComponent) Translate(x, y int) {
	b.x += x
	b.y += y
	for _, c := range b.components {
		if c == nil {
			continue
		}

		c.Translate(x, y)
	}
}

func (b *BaseComponent) Resize(w, h int) {
	b.w = w
	b.h = h
	for _, c := range b.components {
		if c == nil {
			continue
		}

		c.Resize(w, h)
	}
}

func (b *BaseComponent) GetComponents() []Component {
	return b.components
}

func (b *BaseComponent) AddComponent(c Component) {
	b.components = append(b.components, c)
	b.numComponents++
	c.SetInputHandler(b.inputHandler)
	Init(c)
}

func (b *BaseComponent) SetInputHandler(i *InputHandler) {
	b.inputHandler = i
}

func (b *BaseComponent) GetInputHandler() *InputHandler {
	return b.inputHandler
}

func Update(c Component) bool {
	needsRender := c.OnUpdate()
	for _, child := range c.GetComponents() {
		if child == nil {
			continue
		}
		dirty := Update(child)
		if dirty {
			needsRender = true
		}
	}
	return needsRender
}

func Render(c Component, ctx *strife.Renderer) {
	c.OnRender(ctx)
	for _, child := range c.GetComponents() {
		if child == nil {
			continue
		}

		Render(child, ctx)
	}
}

func Init(c Component) {
	c.OnInit()
	for _, child := range c.GetComponents() {
		if child == nil {
			continue
		}

		Init(child)
	}
}

func Dispose(c Component) {
	c.OnDispose()
	for _, child := range c.GetComponents() {
		if child == nil {
			continue
		}

		Dispose(child)
	}
}
