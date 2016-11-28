package cfg

import "strconv"

type TomlConfig struct {
	Editor EditorConfig
	Cursor CursorConfig
	Render RenderConfig
	Theme  ThemeConfig
}

var DEFUALT_TOML_CONFIG string = `[editor]
tab_size = 2
hungry_backspace = true
tabs_are_spaces = true
match_braces = false
maintain_indentation = true
highlight_line = true

[render]
aliased = true

[theme]
background = "0xfdf6e3"
foreground = "0x7a7a7a"
cursor = "0x657B83"
cursor_invert = "0xffffff"

[cursor]
flash_rate = 400
reset_delay = 400
draw = true
flash = true
`

type CursorConfig struct {
	Flash_Rate  uint32
	Reset_Delay uint32
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
	Aliased bool
}

// todo make this more extendable...
// e.g. .nate-editor/themes with TOML
// themes in them and we can select
// the default theme in the EditorConfig
// instead.
type ThemeConfig struct {
	Background    string
	Foreground    string
	Cursor        string
	Cursor_Invert string
}

type EditorConfig struct {
	Tab_Size             int32
	Hungry_Backspace     bool
	Tabs_Are_Spaces      bool
	Match_Braces         bool
	Maintain_Indentation bool
	Highlight_Line       bool
}

func NewDefaultConfig() *TomlConfig {
	return &TomlConfig{
		Editor: EditorConfig{},
		Theme: ThemeConfig{
			Background:    "0xfdf6e3",
			Foreground:    "0x7a7a7a",
			Cursor:        "0x657B83",
			Cursor_Invert: "0xffffff",
		},
	}
}
