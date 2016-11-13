package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
)

type NateEditor struct {
	window *sdl.Window
	surface *sdl.Surface
	running bool
}

func (n *NateEditor) update() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			n.running = false
		case *sdl.TextInputEvent:
			fmt.Println("todo, text input", t)
		}
	}
}

func (n *NateEditor) render() {

}

func main() {
    sdl.Init(sdl.INIT_EVERYTHING)

    window, err := sdl.CreateWindow("Nate Editor", 
    	sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 
    	1280, 720, 
    	sdl.WINDOW_SHOWN)
    if err != nil {
        panic(err)
    }
    defer window.Destroy()

    surface, err := window.GetSurface()
    if err != nil {
        panic(err)
    }

    editor := NateEditor {
    	window: window,
    	surface: surface,
    	running: true,
    }

    for editor.running {
    	editor.update()
    	editor.render()
    }

    sdl.Quit()
}
