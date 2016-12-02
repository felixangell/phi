<h1>nate-editor</h1>

<img align="left" src="res/icons/icon64.png">
Nate is a minimal text editor designed to look pretty, run fast, and be easy
to configure and use. It's primary function is for editing code. Note that this
is a re-write of the initial editor that I wrote last year. It's still a work in
progress and is very buggy!

# goals
The editor must:

* run at 60 fps;
* load and edit large files with ease;
* look pretty; and finally
* be easy to use 

# building
You'll need `veandco/sdl2` and `veandco/SDL2_ttf`, as well as `BurntSushi/toml` and `vinzmay/go-rope`.

```bash
$ go get github.com/felixangell/nate
```

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
