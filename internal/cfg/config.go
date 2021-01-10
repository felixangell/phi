package cfg

import (
	"errors"
	"log"
	"regexp"
	"strconv"
)

type PhiEditorConfig struct {
	Editor       *EditorConfig               `toml:"editor"`
	Cursor       *CursorConfig               `toml:"cursor"`
	Render       *RenderConfig               `toml:"render"`
	Theme        *ThemeConfig                `toml:"theme"`
	Associations map[string]FileAssociations `toml:"file_associations"`
	Commands     map[string]Command          `toml:"commands"`

	associations map[string]*LanguageSyntaxConfig
}

// GetSyntaxConfig returns a pointer to the parsed
// syntax language file for the given file extension
// e.g. what syntax def we need for a .cpp file or a .h file
func (p *PhiEditorConfig) GetSyntaxConfig(ext string) (*LanguageSyntaxConfig, error) {
	if val, ok := p.associations[ext]; ok {
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
	FlashRate  int64  `toml:"flash_rate"`
	ResetDelay int64  `toml:"reset_delay"`
	Draw       bool   `toml:"draw"`
	Flash      bool   `toml:"flash"`
	BlockWidth string `toml:"block_width"`
}

func (c CursorConfig) GetCaretWidth() int {
	if c.BlockWidth == "block" {
		return -1
	}
	if c.BlockWidth == "" {
		return -1
	}

	value, err := strconv.ParseInt(c.BlockWidth, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(value)
}

type RenderConfig struct {
	Aliased            bool   `toml:"aliased"`
	Accelerated        bool   `toml:"accelerated"`
	ThrottleCpuUsage   bool   `toml:"throttle_cpu_usage"`
	AlwaysRender       bool   `toml:"always_render"`
	VerticalSync       bool   `toml:"vertical_sync"`
	SyntaxHighlighting bool   `toml:"syntax_highlighting"`
	FrameSleepInterval uint32 `toml:"frame_sleep_interval"`
}

// todo make this more extendable...
// e.g. .phi-editor/themes with TOML
// themes in them and we can select
// the default theme in the EditorConfig
// instead.
type ThemeConfig struct {
	Background              uint32
	Foreground              uint32
	Cursor                  uint32
	CursorInvert            uint32 `toml:"cursor_invert"`
	Palette                 PaletteConfig
	GutterBackground        uint32 `toml:"gutter_background"`
	GutterForeground        uint32 `toml:"gutter_foreground"`
	HighlightLineBackground uint32 `toml:"highlight_line_background"`
}

type PaletteConfig struct {
	Background   uint32 `toml:"background"`
	Foreground   uint32 `toml:"foreground"`
	Cursor       uint32 `toml:"cursor"`
	Outline      uint32 `toml:"outline"`
	RenderShadow bool   `toml:"render_shadow"`
	ShadowColor  uint32 `toml:"shadow_color"`
	Suggestion   struct {
		Background         uint32 `toml:"background"`
		Foreground         uint32 `toml:"foreground"`
		SelectedBackground uint32 `toml:"selected_background"`
		SelectedForeground uint32 `toml:"selected_foreground"`
	}
}

type EditorConfig struct {
	TabSize             int    `toml:"tab_size"`
	HungryBackspace     bool   `toml:"hungry_backspace"`
	TabsAreSpaces       bool   `toml:"tabs_are_spaces"`
	MatchBraces         bool   `toml:"match_braces"`
	MaintainIndentation bool   `toml:"maintain_indentation"`
	HighlightLine       bool   `toml:"highlight_line"`
	FontPath            string `toml:"font_path"`
	FontFace            string `toml:"font_face"`
	FontSize            int    `toml:"font_size"`
	ShowLineNumbers     bool   `toml:"show_line_numbers"`
}

func NewDefaultConfig() *PhiEditorConfig {
	log.Println("Loading default configuration")

	return &PhiEditorConfig{
		Render: &RenderConfig{
			Aliased:            true,
			Accelerated:        true,
			ThrottleCpuUsage:   true,
			AlwaysRender:       true,
			VerticalSync:       true,
			SyntaxHighlighting: true,
		},
		Editor: &EditorConfig{
			TabSize:             4,
			HungryBackspace:     false,
			TabsAreSpaces:       true,
			MatchBraces:         true,
			MaintainIndentation: true,
			HighlightLine:       true,
			FontPath:            "/Library/Fonts",
			FontFace:            "Go-Mono",
			FontSize:            24,
			ShowLineNumbers:     true,
		},
		Theme: &ThemeConfig{
			Background:   0x002649,
			Foreground:   0xf2f4f6,
			Cursor:       0xf2f4f6,
			CursorInvert: 0xffffff,
		},
		Cursor: &CursorConfig{
			FlashRate:  0,
			ResetDelay: 0,
			Draw:       true,
			Flash:      false,
			BlockWidth: "block",
		},
		Commands: map[string]Command{
			"undo":         {"super+z"},
			"redo":         {"super+y"},
			"exit":         {"super+q"},
			"save":         {"super+s"},
			"page_down":    {"ctrl+down"},
			"page_up":      {"ctrl+up"},
			"show_palette": {"super+p"},
			"focus_left":   {"super+left"},
			"focus_right":  {"super+right"},
			"paste":        {"super+v"},
			"close_buffer": {"super+w"},
			"delete_line":  {"super+d"},
		},

		// TODO syntax defaults
		associations: map[string]*LanguageSyntaxConfig{},
	}
}
