package cfg

import (
	"errors"
	"log"
	"regexp"
	"strconv"

	"github.com/felixangell/strife"
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
	Foreground uint32   `toml:"foreground"`
	Background uint32   `toml:"background"`
	Match      []string `toml:"match"`
	Pattern    string   `toml:"pattern"`

	CompiledPattern *regexp.Regexp
	MatchList       map[string]bool
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
	Aliased             bool
	Accelerated         bool
	Throttle_Cpu_Usage  bool
	Always_Render       bool
	Vertical_Sync       bool
	Syntax_Highlighting bool
}

// todo make this more extendable...
// e.g. .phi-editor/themes with TOML
// themes in them and we can select
// the default theme in the EditorConfig
// instead.
type ThemeConfig struct {
	Background                uint32
	Foreground                uint32
	Cursor                    uint32
	Cursor_Invert             uint32
	Palette                   PaletteConfig
	Gutter_Background         uint32
	Gutter_Foreground         uint32
	Highlight_Line_Background uint32
}

type PaletteConfig struct {
	Background    uint32
	Foreground    uint32
	Cursor        uint32
	Outline       uint32
	Render_Shadow bool
	Shadow_Color  uint32
	Suggestion    struct {
		Background          uint32
		Foreground          uint32
		Selected_Background uint32
		Selected_Foreground uint32
	}
}

type EditorConfig struct {
	Tab_Size             int
	Hungry_Backspace     bool
	Tabs_Are_Spaces      bool
	Match_Braces         bool
	Maintain_Indentation bool
	Highlight_Line       bool
	Font_Path            string
	Font_Face            string
	Font_Size            int
	Show_Line_Numbers    bool
	Loaded_Font          *strife.Font
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
