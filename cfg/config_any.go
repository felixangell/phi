// +build !linux,!darwin

package cfg

var DEFUALT_TOML_CONFIG = `[editor]
tab_size = 4
hungry_backspace = true
tabs_are_spaces = true
match_braces = false
maintain_indentation = true
highlight_line = true
font_face = "Consola"
font_size = 20
show_line_numbers = true

[render]
aliased = true
accelerated = true
throttle_cpu_usage = true
always_render = true
syntax_highlighting = true

[file_associations]
[file_associations.toml]
extensions = [".toml"]

[file_associations.c]
extensions = [".c", ".h", ".cc"]

[file_associations.go]
extensions = [".go"]

[file_associations.md]
extensions = [".md"]

[theme]
background = 0x002649
foreground = 0xf2f4f6
cursor = 0xf2f4f6
cursor_invert = 0x000000
gutter_background = 0x002649
gutter_foreground = 0xf2f4f6

[theme.palette]
outline = 0xebedef
background = 0xffffff
foreground = 0x000000
cursor = 0xf2f4f6
render_shadow = true
shadow_color = 0x000000

[theme.palette.suggestion]
background = 0xebedef
foreground = 0x3a3839
selected_background = 0xc7cbd1
selected_foreground = 0x3a3839

[cursor]
flash_rate = 400
reset_delay = 400
draw = true
flash = true

[commands]
[commands.save]
shortcut = "ctrl+s"

[commands.focus_left]
shortcut = "ctrl+left"

[commands.focus_right]
shortcut = "ctrl+right"

[commands.show_palette]
shortcut = "ctrl+p"

[commands.paste]
shortcut = "ctrl+v"

[commands.exit]
shortcut = "ctrl+q"

[commands.close_buffer]
shortcut = "ctrl+w"

[commands.delete_line]
shortcut = "ctrl+d"
`
