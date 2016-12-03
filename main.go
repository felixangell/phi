package main

import (
	"fmt"
	"runtime"

	"github.com/felixangell/nate/cfg"
	"github.com/felixangell/nate/gui"
	"github.com/felixangell/strife"
)

const (
	PRINT_FPS bool = true
)

type NateEditor struct {
	gui.BaseComponent
	running     bool
	defaultFont *strife.Font
}

func (n *NateEditor) init(cfg *cfg.TomlConfig) {
	n.AddComponent(gui.NewView(800, 600, cfg))

	font, err := strife.LoadFont("./res/firacode.ttf")
	if err != nil {
		panic(err)
	}
	n.defaultFont = font
}

func (n *NateEditor) dispose() {
	for _, comp := range n.GetComponents() {
		gui.Dispose(comp)
	}
}

func (n *NateEditor) update() {
	for _, comp := range n.GetComponents() {
		gui.Update(comp)
	}
}

func (n *NateEditor) render(ctx *strife.Renderer) {
	ctx.Clear()

	ctx.SetFont(n.defaultFont)

	for _, child := range n.GetComponents() {
		gui.Render(child, ctx)
	}

	ctx.Display()
}

func main() {
	config := cfg.Setup()

	windowWidth, windowHeight := 800, 600
	window, err := strife.CreateRenderWindow(windowWidth, windowHeight, &strife.RenderConfig{
		Alias:        true,
		Accelerated:  false,
		VerticalSync: false,
	})
	if err != nil {
		panic(err)
	}

	{
		size := "16"
		switch runtime.GOOS {
		case "windows":
			size = "256"
		case "darwin":
			size = "512"
		case "linux":
			size = "96"
		default:
			panic("you runtime is " + runtime.GOOS)
		}

		icon, err := strife.LoadImage("./res/icons/icon" + size + ".png")
		if err != nil {
			panic(err)
		}
		window.SetIcon(icon)
	}

	editor := &NateEditor{running: true}
	editor.init(&config)

	timer := strife.CurrentTimeMillis()
	num_frames := 0

	for !window.CloseRequested() {
		editor.update()
		editor.render(window.GetRenderContext())
		num_frames += 1

		if strife.CurrentTimeMillis()-timer > 1000 {
			timer = strife.CurrentTimeMillis()
			if PRINT_FPS {
				fmt.Println("frames: ", num_frames)
			}
			num_frames = 0
		}
	}

	editor.dispose()
}
