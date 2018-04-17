package cfg

import (
	"errors"
	"log"
	"strconv"
)

type TomlConfig struct {
	Editor       EditorConfig                `toml:"editor"`
	Cursor       CursorConfig                `toml:"cursor"`
	Render       RenderConfig                `toml:"render"`
	Theme        ThemeConfig                 `toml:"theme"`
	Associations map[string]FileAssociations `toml:"file_associations"`
	Commands     map[string]Command          `toml:"commands"`

	associations map[string]*LanguageSyntaxConfig
}

// GetSyntaxConfig returns a pointer to the parsed
// syntax language file for the given file extension
// e.g. what syntax def we need for a .cpp file or a .h file
func (t *TomlConfig) GetSyntaxConfig(ext string) (*LanguageSyntaxConfig, error) {
	if val, ok := t.associations[ext]; ok {
		return val, nil
	}
	return nil, errors.New("no language for extension '" + ext + "'")
}

type FileAssociations struct {
	Extensions []string
}

type SyntaxCriteria struct {
	Colour  int      `toml:"colouring"`
	Match   []string `toml:"match"`
	Pattern string   `toml:"pattern"`
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
	Always_Render      bool
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
	Font_Face            string
	Font_Size            int
	Show_Line_Numbers    bool
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
