package cfg

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

// TODO we only had double key combos
// e.g. cmd+s. we want to handle things
// like cmd+alt+s
type shortcutRegister struct {
	Supers   map[string]string
	Controls map[string]string
}

var Shortcuts = &shortcutRegister{
	Supers:   map[string]string{},
	Controls: map[string]string{},
}

func configureAndValidate(conf *TomlConfig) {
	// config & validate the keyboard shortcuts
	log.Println("Configuring keyboard shortcuts")
	{
		// keyboard commands
		for commandName, cmd := range conf.Commands {
			shortcut := cmd.Shortcut
			vals := strings.Split(shortcut, "+")

			// TODO handle conflicts

			switch vals[0] {
			case "super":
				Shortcuts.Supers[vals[1]] = commandName
			case "ctrl":
				Shortcuts.Controls[vals[1]] = commandName
			}
		}
	}

	log.Println("Syntax Highlighting")
	{
		conf.associations = map[string]string{}

		for lang, extSet := range conf.Associations {
			log.Println(lang, "=>", extSet.Extensions)

			for _, ext := range extSet.Extensions {
				log.Println("registering", ext, "as", lang)
				conf.associations[ext] = lang
			}
		}

		for name, conf := range conf.Syntax {
			log.Println(name + ":")
			for name, val := range conf {
				log.Println(name, val)
			}
		}
	}
}

func Setup() TomlConfig {
	log.Println("Setting up Phi Editor")

	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}
	CONFIG_DIR := filepath.Join(home, CONFIG_DIR_PATH)
	CONFIG_PATH := filepath.Join(CONFIG_DIR, CONFIG_TOML_FILE)
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

	configureAndValidate(&conf)
	return conf
}
