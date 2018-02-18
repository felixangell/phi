package main

import (
	"fmt"
	"runtime"

	"github.com/ark-lang/ark/src/util/log"
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

	font, err := strife.LoadFont("./res/firacode.ttf", 14)
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

	ww, wh := 800, 600
	window := strife.SetupRenderWindow(ww, wh, strife.DefaultConfig())
	window.SetTitle("Hello world!")
	window.SetResizable(true)
	window.Create()

	window.HandleEvents(func(evt strife.StrifeEvent) {
		switch evt.(type) {
		case *strife.CloseEvent:
			window.Close()
		}
	})

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
			log.Error("unrecognized runtime ", runtime.GOOS)
		}

		icon, err := strife.LoadImage("./res/icons/icon" + size + ".png")
		if err != nil {
			panic(err)
		}
		window.SetIconImage(icon)
		defer icon.Destroy()
	}

	editor := &NateEditor{running: true}
	editor.init(&config)

	timer := strife.CurrentTimeMillis()
	num_frames := 0

	ctx := window.GetRenderContext()

	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		editor.update()

		// TODO: we dont have to constantly
		// render, we can render when we need to
		// i.e. the cursor moves or something
		editor.render(ctx)

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
