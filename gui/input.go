package gui

import (
	"github.com/felixangell/strife"
	"github.com/veandco/go-sdl2/sdl"
)

type InputHandler struct {
	Event sdl.Event
}

func HandleEvent(comp Component, evt strife.StrifeEvent) {
	comp.HandleEvent(evt)
	for _, child := range comp.GetComponents() {
		if child == nil {
			continue
		}
		child.HandleEvent(evt)
	}
}
