<p align="center"><img src="res/icons/icon96.png"></p>

<h1>nate-editor</h1>
Nate is a minimal text editor designed to look pretty, run fast, and be easy
to configure and use. It's primary function is for editing code. Note that this
is a re-write of the initial editor that I wrote last year. It's still a work in
progress and is very buggy! In addition to this, the editor is written as if it's a game,
so it will probably eat up your battery, run quite slow on a laptop, and probably crash
quite frequently.

<br>

Here's a screenshot of the editor right now with the default (temporary!) colour schemes
editing the config file for the editor itself.

<p align="center"><img src="screenshot.png"></p>

# goals
The editor must:

* run at 60 fps;
* load and edit large files with ease;
* look pretty; and finally
* be easy to use 

# building
You'll need Go with the GOPATH, GOBIN, etc. setup, as well as SDL2, SDL2\_image, and SDL2\_ttf. Here's
an example for Ubuntu:

```bash
$ sudo apt-get install libsdl2-dev libsdl2-image-dev libsdl2-ttf-dev
$ go get github.com/felixangell/nate
$ cd $GOPATH/src/github.com/felixangell/nate
$ go build
$ ./nate
```

If you're on macOS, you can get these dependencies via. homebrew. If you're on windows; you have my condolences.

## configuration
Configuration files are stored in `$HOME/.nate-editor/config.toml`, here's
an example, which just so happens to be the defualt configuration:

```toml
[editor]
tab_size = 2
hungry_backspace = true
tabs_are_spaces = true
match_braces = false

[render]
aliased = true

[theme]
background = "0xfdf6e3"
foreground = "0x7a7a7a"
cursor = "0x657B83"
cursor_invert = "0xffffff"

[cursor]
flash_rate = 400
reset_delay = 400
draw = true
flash = true
```

# license
[MIT License](/LICENSE)
