# nate-editor
Nate is a re-write of an old text-editor I wrote earlier this year. It's
very buggy, may or may not work on macOS, and is still a work in progress!

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
