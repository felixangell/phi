package cfg

type TomlConfig struct {
	Editor EditorConfig `toml:"editor"`
}

type EditorConfig struct {
	Aliased bool
}

func NewDefaultConfig() *TomlConfig {
	return &TomlConfig{
		Editor: EditorConfig {
			Aliased: false,
		},
	}
}
