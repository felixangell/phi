package cfg

import (
	"io/ioutil"
	"log"
	"os"

	// fork of BurntSushi with hexadecimal support.
	"github.com/felixangell/toml"
)

// TODO:
// - make the $HOME/.phi-editor folder if it doesn't exist
// - make the $HOME/.phi-editor/config.toml file if it doesn't exist
// - write a default toml file
//

const (
	CONFIG_DIR_PATH  = "/.phi-editor/"
	CONFIG_TOML_FILE = "config.toml"
)

var CONFIG_FULL_PATH string = ""

func Setup() TomlConfig {
	log.Println("Setting up Phi Editor")

	CONFIG_DIR := os.Getenv("HOME") + CONFIG_DIR_PATH
	CONFIG_PATH := CONFIG_DIR + CONFIG_TOML_FILE
	CONFIG_FULL_PATH = CONFIG_PATH

	// if the user doesn't have a /.phi-editor
	// directory we create it for them.
	if _, err := os.Stat(CONFIG_DIR); os.IsNotExist(err) {
		if err := os.Mkdir(CONFIG_DIR, 0775); err != nil {
			panic(err)
		}
	}

	// make sure a config.toml file exists in the
	// phi-editor directory.
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		configFile, fileCreateErr := os.Create(CONFIG_PATH)
		if fileCreateErr != nil {
			panic(fileCreateErr)
		}
		defer configFile.Close()

		// write some stuff
		_, writeErr := configFile.Write([]byte(DEFUALT_TOML_CONFIG))
		if writeErr != nil {
			panic(writeErr)
		}
		configFile.Sync()
	}

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

	return conf
}
