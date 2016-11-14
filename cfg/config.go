package cfg

type Config struct {
	Aliased bool
}

func NewDefaultConfig() *Config {
	return &Config{
		Aliased: false,
	}
}
