package cfg

import (
	"os"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"fmt"
)

// TODO:
// - make the $HOME/.nate-editor folder if it doesn't exist
// - make the $HOME/.nate-editor/config.toml file if it doesn't exist
// - write a default toml file
// 

const (
	CONFIG_DIR_PATH = "/.nate-editor/"
	CONFIG_TOML_FILE = "config.toml"
)

func Setup() TomlConfig {
	CONFIG_PATH := os.Getenv("HOME") + CONFIG_DIR_PATH + CONFIG_TOML_FILE

	if _, err := os.Open(CONFIG_PATH); err != nil {
		panic(err)
	}

	configTomlData, err := ioutil.ReadFile(CONFIG_PATH)
    if err != nil {
        panic(err)
    }

	var conf TomlConfig
	if _, err := toml.Decode(string(configTomlData), &conf); err != nil {
		panic(err)
	}

	fmt.Println("Loaded '" + CONFIG_PATH + "'.")
	fmt.Println(conf)
	return conf
}
