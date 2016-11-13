package main

import (
	"github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/sdl_ttf"
	"github.com/felixangell/nate/gui"
	"fmt"
)

type NateEditor struct {
	window *sdl.Window
	surface *sdl.Surface
	running bool
    panels []*gui.Panel
}

func (n *NateEditor) init() {
    // setup a default panel
    testPanel := gui.NewPanel()
    testPanel.AddComponent(gui.NewBuffer())
    n.panels = append(n.panels, testPanel)
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

    for _, panel := range n.panels {
        panel.Update()
    }
}

func (n *NateEditor) render() {
    w, h := n.window.GetSize()
    n.surface.FillRect(&sdl.Rect{0, 0, int32(w), int32(h)}, 0xffffff)

    for _, panel := range n.panels {
        panel.Render(n.surface)
    }

    n.window.UpdateSurface()
}

func main() {
    sdl.Init(sdl.INIT_EVERYTHING)
    defer sdl.Quit()

    if err := ttf.Init(); err != nil {
        panic(err)
    }

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

    editor := &NateEditor{window: window, surface: surface, running: true}
    editor.init()

    timer := sdl.GetTicks()
    var tick_interval uint32 = 20
    num_frames := 0
    frame_idx := sdl.GetTicks() / tick_interval
    simulate := true

    for editor.running {
        if sdl.GetTicks() / tick_interval > frame_idx {
            frame_idx = sdl.GetTicks() / tick_interval
            simulate = true
        }

    	if simulate {
            editor.update()
        }

    	editor.render()
        num_frames += 1

        if sdl.GetTicks() - timer > 1000 {
            timer = sdl.GetTicks()
            fmt.Println("frames: ", num_frames)
            num_frames = 0
        }
    }
}
