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
	"github.com/pelletier/go-toml"
)

// TODO:
// - make the $HOME/.phi-editor folder if it doesn't exist
// - make the $HOME/.phi-editor/config.toml file if it doesn't exist
// - write a default toml file
//

const (
	ConfigDirPath  = "/.phi-editor/"
	ConfigTomlFile = "config.toml"
)

var FontFolder = ""

// this is the absolute path to the
// config.toml file. todo rename/refactor
var ConfigFullPath = ""

// the absolute path to the config directory
// rename/refactor due here too!
var configDirAbsPath = ""

var IconDirPath = ""

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

	conf := &LanguageSyntaxConfig{}
	if err := toml.Unmarshal(syntaxTomlData, conf); err != nil {
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
		if len(conf.Editor.FontPath) == 0 {
			switch runtime.GOOS {
			case "windows":
				FontFolder = filepath.Join(os.Getenv("WINDIR"), "fonts")
			case "darwin":
				FontFolder = "/Library/Fonts/"
			case "linux":
				FontFolder = findFontFolder()
			}
			// and set it accordingly.
			conf.Editor.FontPath = FontFolder
		}

		// we only support ttf at the moment.
		fontPath := filepath.Join(conf.Editor.FontPath, conf.Editor.FontFace) + ".ttf"
		if _, err := os.Stat(fontPath); os.IsNotExist(err) {
			log.Fatal("No such font '" + fontPath + "'")
		}

		// load the font!
		font, err := strife.LoadFont(fontPath, conf.Editor.FontSize)
		if err != nil {
			panic(err)
		}
		conf.Editor.LoadedFont = font

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
		var syntaxSet []*LanguageSyntaxConfig
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

	ConfigDir := filepath.Join(home, ConfigDirPath)
	configDirAbsPath = ConfigDir

	ConfigPath := filepath.Join(ConfigDir, ConfigTomlFile)

	// this folder is where we store all of the language syntax
	SyntaxConfigDir := filepath.Join(ConfigDir, "syntax")

	ConfigFullPath = ConfigPath

	// if the user doesn't have a /.phi-editor
	// directory we create it for them.
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(ConfigDir, 0775); err != nil {
			panic(err)
		}
	}

	// ----
	// downloads the icon from github
	// and puts it into the phi-editor config folder.
	IconDirPath = filepath.Join(ConfigDir, "icons")
	if _, err := os.Stat(IconDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(IconDirPath, 0775); err != nil {
			panic(err)
		}

		log.Println("setting up the icons folder")

		// https://raw.githubusercontent.com/felixangell/phi/gh-pages/images/icon128.png
		downloadIcon := func(iconSize int) {
			log.Println("downloading the phi icon ", iconSize, "x", iconSize, " png image.")

			file, err := os.Create(filepath.Join(IconDirPath, fmt.Sprintf("icon%d.png", iconSize)))
			if err != nil {
				log.Println(err.Error())
				return
			}
			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()

			// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
			resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/felixangell/phi/gh-pages/images/icon%d.png", iconSize))
			if err != nil {
				log.Println("Failed to download icon", iconSize, "!", err.Error())
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				log.Println(err.Error())
			}
		}

		size := getIconSize()

		// download the icon and
		// write it to the phi-editor folder.
		downloadIcon(size)
	}

	// try make the syntax config folder.
	if _, err := os.Stat(SyntaxConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(SyntaxConfigDir, 0775); err != nil {
			panic(err)
		}

		// load all of the default language syntax
		for name, syntaxDef := range DefaultSyntaxSet {
			languagePath := filepath.Join(SyntaxConfigDir, name+".toml")
			if _, err := os.Stat(languagePath); os.IsNotExist(err) {
				file, err := os.Create(languagePath)
				if err != nil {
					panic(err)
				}
				if _, err := file.Write([]byte(syntaxDef)); err != nil {
					panic(err)
				}
				log.Println("Wrote syntax for language '" + name + "'")
				if err := file.Close(); err != nil {
					panic(err)
				}
			}
		}
	}

	// make sure a config.toml file exists in the
	// phi-editor directory.
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		configFile, fileCreateErr := os.Create(ConfigPath)
		if fileCreateErr != nil {
			panic(fileCreateErr)
		}
		defer func() {
			if err := configFile.Close(); err != nil {
				panic(err)
			}
		}()

		_, writeErr := configFile.Write([]byte(DEFUALT_TOML_CONFIG))
		if writeErr != nil {
			panic(writeErr)
		}
		if err := configFile.Sync(); err != nil {
			panic(err)
		}
	}

	if _, err := os.Open(ConfigPath); err != nil {
		panic(err)
	}

	configTomlData, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		panic(err)
	}

	conf := TomlConfig{}
	if err := toml.Unmarshal(configTomlData, &conf); err != nil {
		panic(err)
	}

	configureAndValidate(&conf)
	return conf
}

func getIconSize() int {
	size := 16
	switch runtime.GOOS {
	case "windows":
		size = 64
	case "darwin":
		size = 512
	case "linux":
		size = 96
	}
	return size
}
