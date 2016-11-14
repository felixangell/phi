package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Component interface {
	Update()
	Render(*sdl.Renderer)

	GetInputHandler() *InputHandler
	SetInputHandler(h *InputHandler)
}
