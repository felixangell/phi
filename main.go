package main

import (
	"fmt"
	"github.com/felixangell/nate/cfg"
	"github.com/felixangell/nate/gfx"
	"github.com/felixangell/nate/gui"
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"runtime"
)

const (
	PRINT_FPS bool = true
)

type NateEditor struct {
	window        *sdl.Window
	renderer      *sdl.Renderer
	running       bool
	bufferPanels  []*gui.Panel
	input_handler *gui.InputHandler
}

func (n *NateEditor) addBuffer(c gui.Component) {
	panel := gui.NewPanel(n.input_handler)
	c.SetInputHandler(n.input_handler)
	panel.AddComponent(c)
	n.bufferPanels = append(n.bufferPanels, panel)

	w, _ := n.window.GetSize()
	bufferWidth := w / len(n.bufferPanels)
	for i, p := range n.bufferPanels {
		p.Translate(int32(bufferWidth)*int32(i), 0)
	}
}

func (n *NateEditor) init(cfg *cfg.TomlConfig) {
	n.addBuffer(gui.NewBuffer())
	n.addBuffer(gui.NewBuffer())

	/*
		bufferPanel := gui.NewPanel(n.input_handler)
		palette := gui.NewCommandPalette()
		palette.SetInputHandler(n.input_handler)
		bufferPanel.AddComponent(palette)
		n.panels = append(n.panels, bufferPanel)
	*/
}

func (n *NateEditor) update() {
	for _, panel := range n.bufferPanels {
		panel.Update()
	}

	n.input_handler.Event = nil
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		n.input_handler.Event = event

		switch event.(type) {
		case *sdl.QuitEvent:
			n.running = false
		case *sdl.TextInputEvent:
			n.input_handler.Event = event
		}
	}

}

func (n *NateEditor) render() {
	gfx.SetDrawColorHex(n.renderer, 0xfdf6e3)
	n.renderer.Clear()

	for _, panel := range n.bufferPanels {
		panel.Render(n.renderer)
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

	window, err := sdl.CreateWindow("Nate Editor",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		1280, 720,
		sdl.WINDOW_SHOWN)
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

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	editor := &NateEditor{window: window, renderer: renderer, running: true, input_handler: &gui.InputHandler{}}
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
}
