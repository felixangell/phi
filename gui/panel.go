package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Panel struct {
	components []Component
}

func NewPanel() *Panel {
	return &Panel{
		components: []Component{},
	}
}

func (p *Panel) AddComponent(c Component) {
	p.components = append(p.components, c)
}

func (p *Panel) Update() {
	for _, c := range p.components {
		c.Update()
	}
}

func (p *Panel) Render(ctx *sdl.Surface) {
	for _, c := range p.components {
		c.Render(ctx)
	}
}
