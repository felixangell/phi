package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/felixangell/phi-editor/cfg"
	"github.com/felixangell/phi-editor/gui"
	"github.com/felixangell/strife"
)

const (
	PRINT_FPS bool = true
)

type PhiEditor struct {
	gui.BaseComponent
	running     bool
	defaultFont *strife.Font
}

func (n *PhiEditor) init(cfg *cfg.TomlConfig) {
	mainView := gui.NewView(1280, 720, cfg)

	args := os.Args
	if len(args) > 1 {
		// TODO check these are files
		// that actually exist here?
		for _, arg := range args[1:] {
			mainView.AddBuffer().OpenFile(arg)
		}
	} else {
		// we have no args, open up a scratch file
		tempFile, err := ioutil.TempFile("/var/tmp/", "phi-editor-")
		if err != nil {
			log.Println("Failed to create temp file", err.Error())
			os.Exit(1)
		}

		mainView.AddBuffer().OpenFile(tempFile.Name())
	}

	n.AddComponent(mainView)

	font, err := strife.LoadFont("./res/firacode.ttf", 20)
	if err != nil {
		panic(err)
	}
	n.defaultFont = font
}

func (n *PhiEditor) dispose() {
	for _, comp := range n.GetComponents() {
		gui.Dispose(comp)
	}
}

func (n *PhiEditor) update() bool {
	needsRender := false
	for _, comp := range n.GetComponents() {
		dirty := gui.Update(comp)
		if dirty {
			needsRender = true
		}
	}
	return needsRender
}

func (n *PhiEditor) render(ctx *strife.Renderer) {
	ctx.SetFont(n.defaultFont)

	for _, child := range n.GetComponents() {
		gui.Render(child, ctx)
	}
}

func main() {
	config := cfg.Setup()

	ww, wh := 1280, 720
	window := strife.SetupRenderWindow(ww, wh, strife.DefaultConfig())
	window.SetTitle("Hello world!")
	window.SetResizable(true)

	editor := &PhiEditor{running: true}
	window.HandleEvents(func(evt strife.StrifeEvent) {
		switch evt.(type) {
		case *strife.CloseEvent:
			window.Close()
		}
	})

	window.Create()

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
			log.Println("unrecognized runtime ", runtime.GOOS)
		}

		icon, err := strife.LoadImage("./res/icons/icon" + size + ".png")
		if err != nil {
			panic(err)
		}
		window.SetIconImage(icon)
		defer icon.Destroy()
	}

	editor.init(&config)

	timer := strife.CurrentTimeMillis()
	frames, updates := 0, 0
	fps, ups := frames, updates

	ctx := window.GetRenderContext()

	ctx.Clear()
	editor.render(ctx)
	ctx.Display()

	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		if editor.update() {
			ctx.Clear()
			editor.render(ctx)

			// this is only printed on each
			// render...
			ctx.SetColor(strife.White)
			ctx.String(fmt.Sprintf("fps: %d, ups %d", fps, ups), ww-256, wh-128)

			ctx.Display()
			frames += 1
		}
		updates += 1

		if strife.CurrentTimeMillis()-timer > 1000 {
			timer = strife.CurrentTimeMillis()
			fps, ups = frames, updates
			frames, updates = 0, 0
		}

		if config.Render.Throttle_Cpu_Usage {
			// todo put in the config how long
			// we sleep for!
			time.Sleep(16)
		}
	}

	editor.dispose()
}
