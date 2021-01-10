package main

import (
	"fmt"
	"github.com/felixangell/phi/internal/cfg"
	"github.com/felixangell/phi/internal/editor"
	"github.com/felixangell/strife"
	"log"
	"runtime"
	"time"
)

func main() {
	runtime.LockOSThread()

	config := cfg.NewDefaultConfig()

	windowConfig := strife.DefaultConfig()
	windowConfig.Accelerated = config.Render.Accelerated
	windowConfig.Alias = config.Render.Aliased
	windowConfig.VerticalSync = config.Render.VerticalSync

	scaledWidth, scaledHeight := calcScaledWindowDimension(800, 600)
	window := strife.SetupRenderWindow(scaledWidth, scaledHeight, windowConfig)
	window.AllowHighDPI()
	window.SetTitle("Hello world!")
	window.SetResizable(true)

	editorInst := editor.NewPhiEditor()
	window.HandleEvents(func(evt strife.StrifeEvent) {
		switch event := evt.(type) {
		case *strife.CloseEvent:
			window.Close()
		case *strife.WindowResizeEvent:
			log.Println("window resize is unimplemented: size", event.Width, event.Height)
		default:
			editorInst.HandleEvent(evt)
		}
	})

	if err := window.Create(); err != nil {
		panic(err)
	}

	editorInst.ApplyConfig(config)

	ctx := window.GetRenderContext()

	// singular render before we enter into the event loop.
	ctx.Clear()
	editorInst.Render(ctx)
	ctx.Display()

	lastDebugRender := time.Now()
	frames, updates := 0, 0
	fps, ups := frames, updates
	for {
		window.PollEvents()
		if window.CloseRequested() {
			break
		}

		shouldRender := editorInst.Update()

		if shouldRender || config.Render.AlwaysRender {
			ctx.Clear()
			editorInst.Render(ctx)

			// this is only printed on each
			// render...
			ctx.SetColor(strife.White)
			ctx.Text(fmt.Sprintf("fps: %d, ups %d", fps, ups), int(scaledWidth-256), int(scaledHeight-128))

			ctx.Display()
			frames++
		}
		updates++

		if time.Now().Sub(lastDebugRender) >= time.Second {
			lastDebugRender = time.Now()
			fps, ups = frames, updates
			frames, updates = 0, 0
		}

		if config.Render.ThrottleCpuUsage {
			time.Sleep(time.Duration(config.Render.FrameSleepInterval) * time.Millisecond)
		}
	}
}

func calcScaledWindowDimension(width, height float32) (int, int) {
	dpi, defDpi := strife.GetDisplayDPI(0)

	cfg.ScaleFactor = float64(dpi / defDpi)

	scaledWidth := int((width * dpi) / defDpi)
	scaledHeight := int((height * dpi) / defDpi)
	return scaledWidth, scaledHeight
}

func getIconSizeForOS() int {
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
	return size
}
