package cfg

import (
	"errors"
	"log"
	"strconv"
)

type TomlConfig struct {
	Editor       EditorConfig                         `toml:"editor"`
	Cursor       CursorConfig                         `toml:"cursor"`
	Render       RenderConfig                         `toml:"render"`
	Theme        ThemeConfig                          `toml:"theme"`
	Associations map[string]FileAssociations          `toml:"file_associations"`
	Commands     map[string]Command                   `toml:"commands"`
	Syntax       map[string]map[string]SyntaxCriteria `toml:"syntax"`

	// this maps ext => language
	// when we have file associations from
	// the Associations field we take
	// each extension and put them here
	// pointing it to the language.
	// basically the reverse/opposite
	associations map[string]string
}

func (t *TomlConfig) GetLanguageFromExt(ext string) (string, error) {
	if val, ok := t.associations[ext]; ok {
		return val, nil
	}
	return "", errors.New("no language for extension '" + ext + "'")
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
