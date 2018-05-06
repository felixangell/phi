package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/phi/gui"
	"github.com/felixangell/strife"
)

const (
	PRINT_FPS bool = true
)

type PhiEditor struct {
	running     bool
	defaultFont *strife.Font
	mainView    *gui.View
}

func (n *PhiEditor) handleEvent(evt strife.StrifeEvent) {

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
		tempFile, err := ioutil.TempFile("", "phi-editor-")
		if err != nil {
			log.Println("Failed to create temp file", err.Error())
			os.Exit(1)
		}

		mainView.AddBuffer().OpenFile(tempFile.Name())
	}

	n.mainView = mainView
	n.defaultFont = cfg.Editor.Loaded_Font
}

func (n *PhiEditor) dispose() {

}

func (n *PhiEditor) update() bool {
	return n.mainView.OnUpdate()
}

func (n *PhiEditor) render(ctx *strife.Renderer) {
	ctx.SetFont(n.defaultFont)
	n.mainView.OnRender(ctx)
}

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

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
		default:
			editor.handleEvent(evt)
		}
	})

	window.Create()

	{
		size := 16
		switch runtime.GOOS {
		case "windows":
			size = 64
		case "darwin":
			size = 512
		case "linux":
			size = 96
		default:
			log.Println("unrecognized runtime ", runtime.GOOS)
		}

		iconFile := fmt.Sprintf("icon%d.png", size)
		icon, err := strife.LoadImage(filepath.Join(cfg.ICON_DIR_PATH, iconFile))
		if err != nil {
			log.Println("Failed to load icon ", err.Error())
		} else {
			window.SetIconImage(icon)
			defer icon.Destroy()
		}
	}

	editor.init(&config)

	lastDebugRender := time.Now()
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

		shouldRender := editor.update()

		if shouldRender || config.Render.Always_Render {
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

		if time.Now().Sub(lastDebugRender) >= time.Second {
			lastDebugRender = time.Now()
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
