package cfg

type TomlConfig struct {
	Editor EditorConfig `toml:"editor"`
}

type EditorConfig struct {
	Aliased bool
	Tab_Size int32
}

func NewDefaultConfig() *TomlConfig {
	return &TomlConfig{
		Editor: EditorConfig {
			Aliased: false,
		},
	}
}
