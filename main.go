package main

import (
	"fmt"
	"runtime"

	"github.com/felixangell/nate/cfg"
	"github.com/felixangell/nate/gfx"
	"github.com/felixangell/nate/gui"
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	PRINT_FPS bool = true
)

type NateEditor struct {
	gui.BaseComponent

	window   *sdl.Window
	renderer *sdl.Renderer
	running  bool
}

func (n *NateEditor) init(cfg *cfg.TomlConfig) {
	w, h := n.window.GetSize()
	n.AddComponent(gui.NewView(w, h, cfg))

	// palette := gui.NewCommandPalette()
	// palette.Translate(int32(w/2), 20)
	// n.AddComponent(palette)
}

func (n *NateEditor) dispose() {
	for _, comp := range n.GetComponents() {
		gui.Dispose(comp)
	}
}

func (n *NateEditor) update() {
	n.GetInputHandler().Event = nil
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		n.GetInputHandler().Event = event

		switch event.(type) {
		case *sdl.QuitEvent:
			n.running = false
		case *sdl.TextEditingEvent:
			n.GetInputHandler().Event = event
		case *sdl.TextInputEvent:
			n.GetInputHandler().Event = event
		}
	}

	for _, comp := range n.GetComponents() {
		gui.Update(comp)
	}
}

func (n *NateEditor) render() {
	gfx.SetDrawColorHex(n.renderer, 0xffffff)
	n.renderer.Clear()

	for _, component := range n.GetComponents() {
		gui.Render(component, n.renderer)
	}

	n.renderer.Present()
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	config := cfg.Setup()

	windowWidth, windowHeight := 800, 600
	{
		// calculate the size of the window
		// based on the resolution of the monitor
		// this is the display width
		var displayMode sdl.DisplayMode
		sdl.GetDisplayMode(0, 0, &displayMode)
		windowWidth = int(float32(displayMode.W) / 1.5)
		windowHeight = windowWidth / 16 * 9
	}

	window, err := sdl.CreateWindow("Nate Editor", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	{
		img.Init(img.INIT_PNG)
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
		icon, err := img.Load("./res/icons/icon" + size + ".png")
		if err != nil {
			panic(err)
		}
		window.SetIcon(icon)
	}

	var mode uint32 = sdl.RENDERER_SOFTWARE
	if config.Render.Accelerated {
		mode = sdl.RENDERER_ACCELERATED
	}

	renderer, err := sdl.CreateRenderer(window, -1, mode)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	editor := &NateEditor{window: window, renderer: renderer, running: true}
	editor.SetInputHandler(&gui.InputHandler{})
	editor.init(&config)

	timer := sdl.GetTicks()
	num_frames := 0

	for editor.running {
		editor.update()
		editor.render()
		num_frames += 1

		if sdl.GetTicks()-timer > 1000 {
			timer = sdl.GetTicks()
			if PRINT_FPS {
				fmt.Println("frames: ", num_frames)
			}
			num_frames = 0
		}

		sdl.Delay(2)
	}

	editor.dispose()
}
