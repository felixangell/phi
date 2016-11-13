package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Component interface {
	Update()
	Render(*sdl.Surface)
	
	GetInputHandler() *InputHandler
	SetInputHandler(h *InputHandler)
}
