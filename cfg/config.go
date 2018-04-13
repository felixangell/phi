package cfg

import (
	"log"
	"runtime"
	"strconv"
	"strings"
)

type TomlConfig struct {
	Editor   EditorConfig
	Cursor   CursorConfig
	Render   RenderConfig
	Theme    ThemeConfig
	Commands map[string]Command
}

var DEFUALT_TOML_CONFIG string = getDefaultConfig()

func getDefaultConfig() string {
	switch strings.ToLower(runtime.GOOS) {
	case "windows":
		return DEFAULT_WINDOWS_TOML_CONFIG
	case "linux":
		return DEFAULT_LINUX_TOML_CONFIG
	case "darwin":
		return DEFAULT_MAC_TOML_CONFIG
	}

	// fallback is a windows config.
	return DEFAULT_WINDOWS_TOML_CONFIG
}

type Command struct {
	Shortcut string
}

type CursorConfig struct {
	Flash_Rate  int64
	Reset_Delay int64
	Draw        bool
	Flash       bool
	Block_Width string
}

func (c CursorConfig) GetCaretWidth() int {
	if c.Block_Width == "block" {
		return -1
	}
	if c.Block_Width == "" {
		return -1
	}

	value, err := strconv.ParseInt(c.Block_Width, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(value)
}

type RenderConfig struct {
	Aliased            bool
	Accelerated        bool
	Throttle_Cpu_Usage bool
}

// todo make this more extendable...
// e.g. .phi-editor/themes with TOML
// themes in them and we can select
// the default theme in the EditorConfig
// instead.
type ThemeConfig struct {
	Background    int32
	Foreground    int32
	Cursor        int32
	Cursor_Invert int32
}

type EditorConfig struct {
	Tab_Size             int
	Hungry_Backspace     bool
	Tabs_Are_Spaces      bool
	Match_Braces         bool
	Maintain_Indentation bool
	Highlight_Line       bool
}

func NewDefaultConfig() *TomlConfig {
	log.Println("Loading default configuration... this should never happen")
	return &TomlConfig{
		Editor: EditorConfig{},
		Theme: ThemeConfig{
			Background:    0x002649,
			Foreground:    0xf2f4f6,
			Cursor:        0xf2f4f6,
			Cursor_Invert: 0xffffff,
		},
	}
}
