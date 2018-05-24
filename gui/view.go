package gui

import (
	"fmt"
	"log"
	"runtime"
	"unicode"

	"github.com/felixangell/phi/cfg"
	"github.com/felixangell/strife"
	"github.com/fsnotify/fsnotify"
	"github.com/veandco/go-sdl2/sdl"
)

type bufferEvent interface {
	Process(view *View)
	String() string
}

type ReloadBufferEvent struct {
	buff *Buffer
}

func (r *ReloadBufferEvent) Process(view *View) {
	log.Println("reloading buffer", r.buff.filePath)
	r.buff.reload()
}

func (r *ReloadBufferEvent) String() string {
	return "reload-buffer-event"
}

// View is an array of buffers basically.
type View struct {
	BaseComponent

	conf           *cfg.TomlConfig
	buffers        []*BufferPane
	focusedBuff    int
	commandPalette *CommandPalette

	watcher      *fsnotify.Watcher
	bufferMap    map[string]*Buffer
	bufferEvents chan bufferEvent
}

func NewView(width, height int, conf *cfg.TomlConfig) *View {
	view := &View{
		conf:         conf,
		buffers:      []*BufferPane{},
		bufferMap:    map[string]*Buffer{},
		bufferEvents: make(chan bufferEvent),
	}

	view.Translate(width, height)
	view.Resize(width, height)

	view.commandPalette = NewCommandPalette(*conf, view)
	view.UnfocusBuffers()

	var err error
	view.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		// ?
	}

	// goroutine to handle all of the fsnotify events
	// converts them into events phi can handle cleanly.
	go func() {
		for {
			select {
			case event := <-view.watcher.Events:
				log.Println("evt: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {

					// modified so we specify a reload event
					buff, ok := view.bufferMap[event.Name]
					if !ok {
						panic(fmt.Sprintf("no such buffer for file '%s'", event.Name))
						break
					}

					view.bufferEvents <- &ReloadBufferEvent{buff}
					log.Println("modified file:", event.Name)
				}
			case err := <-view.watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// handles all of the phi events
	go func() {
		for {
			event := <-view.bufferEvents
			event.Process(view)
		}
	}()

	return view
}

func (n *View) registerFile(path string, buff *Buffer) {
	log.Println("Registering file ", path)

	err := n.watcher.Add(path)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to register file '%s'", path), "to buffer ", buff.index)
		return
	}

	n.bufferMap[path] = buff
}

func (n *View) Close() {
	n.watcher.Close()
}

func (n *View) hidePalette() {
	p := n.commandPalette
	p.clearInput()
	p.SetFocus(false)

	// remove focus from palette
	p.buff.SetFocus(false)
}

func (n *View) focusPalette(buff *Buffer) {
	p := n.commandPalette
	p.SetFocus(true)

	// remove focus from the buffer
	// that invoked the command palette
	p.parentBuff = buff
}

func (n *View) UnfocusBuffers() {
	// clear focus from buffers
	for _, buffPane := range n.buffers {
		buffPane.SetFocus(false)
	}
}

func sign(dir int) int {
	if dir > 0 {
		return 1
	} else if dir < 0 {
		return -1
	}
	return 0
}

func (n *View) removeBuffer(index int) {
	log.Println("Removing buffer index:", index)
	log.Println("num buffs before delete: ", len(n.buffers))

	n.buffers = append(n.buffers[:index], n.buffers[index+1:]...)

	// only resize the buffers if we have
	// some remaining in the window
	if len(n.buffers) > 0 {
		bufferWidth := n.w / len(n.buffers)

		// translate all the components accordingly.
		for idx, buffPane := range n.buffers {
			// re-write the index.
			buffPane.Buff.index = idx

			buffPane.Resize(bufferWidth, n.h)
			buffPane.SetPosition(bufferWidth*idx, 0)
		}
	}

	dir := -1
	if n.focusedBuff == 0 {
		dir = 1
	}

	n.ChangeFocus(dir)
}

func (n *View) setFocusTo(index int) {
	log.Println("Moving focus from ", n.focusedBuff, " to ", index)

	n.UnfocusBuffers()
	n.focusedBuff = index
	buff := n.getCurrentBuffPane()
	buff.SetFocus(true)
}

// FIXME
func (n *View) ChangeFocus(dir int) {
	// we cant change focus if there are no
	// buffers to focus to
	if len(n.buffers) == 0 {
		return
	}

	newFocus := n.focusedBuff

	if dir == -1 {
		newFocus--
	} else if dir == 1 {
		newFocus++
	}

	if newFocus < 0 {
		newFocus = len(n.buffers) - 1
	} else if newFocus >= len(n.buffers) {
		newFocus = 0
	}

	n.UnfocusBuffers()
	n.setFocusTo(newFocus)
}

func (n *View) getCurrentBuffPane() *BufferPane {
	if len(n.buffers) == 0 {
		return nil
	}

	if buffPane := n.buffers[n.focusedBuff]; buffPane != nil {
		return buffPane
	}
	return nil
}

func (n *View) getCurrentBuff() *Buffer {
	buff := n.getCurrentBuffPane()
	if buff != nil {
		return buff.Buff
	}
	return nil
}

func (n *View) OnInit() {
}

func (n *View) OnUpdate() bool {
	dirty := false

	CONTROL_DOWN = strife.KeyPressed(sdl.K_LCTRL) || strife.KeyPressed(sdl.K_RCTRL)
	SUPER_DOWN = strife.KeyPressed(sdl.K_LGUI) || strife.KeyPressed(sdl.K_RGUI)

	shortcutName := "ctrl"
	source := cfg.Shortcuts.Controls

	if strife.PollKeys() && (SUPER_DOWN || CONTROL_DOWN) {
		if runtime.GOOS == "darwin" {
			if SUPER_DOWN {
				source = cfg.Shortcuts.Supers
				shortcutName = "super"
			} else if CONTROL_DOWN {
				source = cfg.Shortcuts.Controls
				shortcutName = "control"
			}
		} else {
			source = cfg.Shortcuts.Controls
		}

		r := rune(strife.PopKey())

		if r == 'l' {
			DEBUG_MODE = !DEBUG_MODE
		}

		left := sdl.K_LEFT
		right := sdl.K_RIGHT
		up := sdl.K_UP
		down := sdl.K_DOWN

		// map to left/right/etc.
		// FIXME
		var key string
		switch int(r) {
		case left:
			key = "left"
		case right:
			key = "right"
		case up:
			key = "up"
		case down:
			key = "down"
		default:
			key = string(unicode.ToLower(r))
		}

		actionName, actionExists := source[key]
		if actionExists {
			if action, ok := actions[actionName]; ok {
				log.Println("Executing action '" + actionName + "'")
				return action.proc(n, []string{})
			}
		} else {
			log.Println("view: unimplemented shortcut", shortcutName, "+", string(unicode.ToLower(r)), "#", int(r), actionName, key)
		}
	}

	buff := n.getCurrentBuffPane()
	if buff != nil {
		buff.OnUpdate()
	}

	n.commandPalette.OnUpdate()

	return dirty
}

func (n *View) OnRender(ctx *strife.Renderer) {
	for _, buffPane := range n.buffers {
		buffPane.OnRender(ctx)
	}

	n.commandPalette.OnRender(ctx)

	if DEBUG_MODE {
		ctx.SetColor(strife.HexRGB(0xff00ff))
		mPos := strife.MouseCoords()
		ctx.Rect(mPos[0], mPos[1], 16, 16, strife.Line)

		renderDebugPane(ctx, 10, 10)
	}
}

func (n *View) OnDispose() {}

func (n *View) AddBuffer() *Buffer {
	n.UnfocusBuffers()

	cfg := n.conf
	c := NewBuffer(cfg, BufferConfig{
		cfg.Theme.Background,
		cfg.Theme.Foreground,
		cfg.Theme.Cursor,
		cfg.Theme.Cursor_Invert,
		cfg.Theme.Gutter_Background,
		cfg.Theme.Gutter_Foreground,
		cfg.Editor.Loaded_Font,
	}, n, len(n.buffers))

	c.SetFocus(true)

	// work out the size of the buffer and set it
	// note that we +1 the components because
	// we haven't yet added the panel
	var bufferWidth int
	bufferWidth = n.w / (c.index + 1)

	n.focusedBuff = c.index

	n.buffers = append(n.buffers, NewBufferPane(c))

	// translate all the buffers accordingly.
	idx := 0
	for _, buffPane := range n.buffers {
		buffPane.Resize(bufferWidth, n.h)
		buffPane.SetPosition(bufferWidth*idx, 0)
		idx++
	}

	return c
}
