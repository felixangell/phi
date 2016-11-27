package cfg

type TomlConfig struct {
	Editor EditorConfig `toml:"editor"`
}

type EditorConfig struct {
	Aliased            bool   `toml:"aliased"`
	Tab_Size           int32  `toml:"tab_size"`
	Hungry_Backspace   bool   `toml:"hungry_backspace"`
	Tabs_Are_Spaces    bool   `toml:"tabs_are_spaces"`
	Draw_Cursor        bool   `toml:"draw_cursor"`
	Flash_Cursor       bool   `toml:"flash_cursor"`
	Cursor_Flash_Rate  uint32 `toml:"cursor_flash_rate"`
	Cursor_Reset_Delay uint32 `toml:"cursor_reset_delay"`
	Match_Braces       bool   `toml:"match_braces"`
}

func NewDefaultConfig() *TomlConfig {
	return &TomlConfig{
		Editor: EditorConfig{},
	}
}
