# nate-editor
Nate is a re-write of an old text-editor I wrote earlier this year.

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
Right now the configuration files are very much unimplemented. At the moment 
configuration files are loaded, but they do not actually modify the behaviour
of the editor, nor are they error checked.

```toml
[editor]
aliased = true
```

# license
[mit](/LICENSE)
