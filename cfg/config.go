package cfg

type TomlConfig struct {
	Editor EditorConfig `toml:"editor"`
	Cursor CursorConfig `toml:"cursor"`
	Render RenderConfig `toml:"render"`
	Theme  ThemeConfig  `toml:"theme"`
}

type CursorConfig struct {
	Flash_Rate  uint32
	Reset_Delay uint32
	Draw        bool
	Flash       bool
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
	Tab_Size         int32
	Hungry_Backspace bool
	Tabs_Are_Spaces  bool
	Match_Braces     bool
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
