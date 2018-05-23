package cfg

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/felixangell/strife"

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

var FONT_FOLDER string = ""

// this is the absolute path to the
// config.toml file. todo rename/refactor
var CONFIG_FULL_PATH string = ""

// the absolute path to the config directory
// rename/refactor due here too!
var configDirAbsPath string = ""

var ICON_DIR_PATH string = ""

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

func loadSyntaxDef(lang string) *LanguageSyntaxConfig {
	languagePath := filepath.Join(configDirAbsPath, "syntax", lang+".toml")
	log.Println("Loading lang from ", languagePath)

	syntaxTomlData, err := ioutil.ReadFile(languagePath)
	if err != nil {
		log.Println("Failed to load highlighting for language '"+lang+"' from path: ", languagePath)
		return nil
	}

	var conf = &LanguageSyntaxConfig{}
	if _, err := toml.Decode(string(syntaxTomlData), conf); err != nil {
		panic(err)
	}

	log.Println("Loaded syntax definition for language", lang)
	return conf
}

func findFontFolder() string {
	// TODO
	return "/usr/share/fonts/"
}

func configureAndValidate(conf *TomlConfig) {
	// fonts
	log.Println("Configuring fonts")
	{
		// the font path has not been set
		// so we have to figure out what it is.
		if len(conf.Editor.Font_Path) == 0 {
			switch runtime.GOOS {
			case "windows":
				FONT_FOLDER = filepath.Join(os.Getenv("WINDIR"), "fonts")
			case "darwin":
				FONT_FOLDER = "/Library/Fonts/"
			case "linux":
				FONT_FOLDER = findFontFolder()
			}
			// and set it accordingly.
			conf.Editor.Font_Path = FONT_FOLDER
		}

		// we only support ttf at the moment.
		fontPath := filepath.Join(conf.Editor.Font_Path, conf.Editor.Font_Face) + ".ttf"
		if _, err := os.Stat(fontPath); os.IsNotExist(err) {
			log.Fatal("No such font '" + fontPath + "'")
			// TODO cool error messages for the toml format?
			os.Exit(0)
		}

		// load the font!
		font, err := strife.LoadFont(fontPath, conf.Editor.Font_Size)
		if err != nil {
			panic(err)
		}
		conf.Editor.Loaded_Font = font
	}

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
		syntaxSet := []*LanguageSyntaxConfig{}

		conf.associations = map[string]*LanguageSyntaxConfig{}

		for lang, extSet := range conf.Associations {
			log.Println(lang, "=>", extSet.Extensions)
			languageConfig := loadSyntaxDef(lang)
			// check for errors here

			syntaxSet = append(syntaxSet, languageConfig)

			for _, ext := range extSet.Extensions {
				log.Println("registering", ext, "as", lang)
				conf.associations[ext] = languageConfig
			}
		}

		// go through each language
		// and store the matches keywords
		// as a hashmap for faster lookup
		// in addition to this we compile any
		// regular expressions if necessary.
		for _, language := range syntaxSet {
			for _, syn := range language.Syntax {
				syn.MatchList = map[string]bool{}

				if syn.Pattern != "" {
					regex, err := regexp.Compile(syn.Pattern)
					if err != nil {
						log.Println(err.Error())
						continue
					}
					syn.CompiledPattern = regex
				} else {
					for _, item := range syn.Match {
						if _, ok := syn.MatchList[item]; ok {
							log.Println("Warning duplicate match item '" + item + "'")
							continue
						}

						syn.MatchList[item] = true
					}
				}

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
	configDirAbsPath = CONFIG_DIR

	CONFIG_PATH := filepath.Join(CONFIG_DIR, CONFIG_TOML_FILE)

	// this folder is where we store all of the language syntax
	SYNTAX_CONFIG_DIR := filepath.Join(CONFIG_DIR, "syntax")

	CONFIG_FULL_PATH = CONFIG_PATH

	// if the user doesn't have a /.phi-editor
	// directory we create it for them.
	if _, err := os.Stat(CONFIG_DIR); os.IsNotExist(err) {
		if err := os.Mkdir(CONFIG_DIR, 0775); err != nil {
			panic(err)
		}
	}

	// ----
	// downloads the icon from github
	// and puts it into the phi-editor config folder.
	ICON_DIR_PATH = filepath.Join(CONFIG_DIR, "icons")
	if _, err := os.Stat(ICON_DIR_PATH); os.IsNotExist(err) {
		if err := os.Mkdir(ICON_DIR_PATH, 0775); err != nil {
			panic(err)
		}

		log.Println("setting up the icons folder")

		// https://raw.githubusercontent.com/felixangell/phi/gh-pages/images/icon128.png
		downloadIcon := func(iconSize int) {
			log.Println("downloading the phi icon ", iconSize, "x", iconSize, " png image.")

			file, err := os.Create(filepath.Join(ICON_DIR_PATH, fmt.Sprintf("icon%d.png", iconSize)))
			defer file.Close()
			if err != nil {
				log.Println(err.Error())
				return
			}

			// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
			resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/felixangell/phi/gh-pages/images/icon%d.png", iconSize))
			if err != nil {
				log.Println("Failed to download icon", iconSize, "!", err.Error())
				return
			}
			defer resp.Body.Close()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				log.Println(err.Error())
			}
		}

		size := 16
		switch runtime.GOOS {
		case "windows":
			size = 64
		case "darwin":
			size = 512
		case "linux":
			size = 96
		}

		// download the icon and
		// write it to the phi-editor folder.
		downloadIcon(size)
	}

	// try make the syntax config folder.
	if _, err := os.Stat(SYNTAX_CONFIG_DIR); os.IsNotExist(err) {
		if err := os.Mkdir(SYNTAX_CONFIG_DIR, 0775); err != nil {
			panic(err)
		}

		// load all of the default language syntax
		for name, syntaxDef := range DefaultSyntaxSet {
			languagePath := filepath.Join(SYNTAX_CONFIG_DIR, name+".toml")
			if _, err := os.Stat(languagePath); os.IsNotExist(err) {
				file, err := os.Create(languagePath)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				if _, err := file.Write([]byte(syntaxDef)); err != nil {
					panic(err)
				}
				log.Println("Wrote syntax for language '" + name + "'")
			}
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
