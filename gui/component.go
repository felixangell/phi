package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Component interface {
	Translate(x, y int32)

	Init()
	Update()
	Render(*sdl.Renderer)

	// gross refactor me pls

	AddComponent(c Component)
	GetComponents() []Component

	GetInputHandler() *InputHandler
	SetInputHandler(h *InputHandler)
}

type ComponentLocation struct {
	x, y int32
}
