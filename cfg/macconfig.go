package cfg

var DEFAULT_MAC_TOML_CONFIG string = `[editor]
tab_size = 4
hungry_backspace = true
tabs_are_spaces = true
match_braces = false
maintain_indentation = true
highlight_line = true
font_face = "Courier New"
font_size = 20
show_line_numbers = true

[render]
aliased = true
accelerated = true
throttle_cpu_usage = true
always_render = true

[file_associations]
[file_associations.c]
extensions = [".c", ".h", ".cc"]

[file_associations.go]
extensions = [".go"]

[theme]
background = 0x002649
foreground = 0xf2f4f6
cursor = 0xf2f4f6
cursor_invert = 0x000000

[cursor]
flash_rate = 400
reset_delay = 400
draw = true
flash = true

[commands]
[commands.save]
shortcut = "super+s"

[commands.close_buffer]
shortcut = "super+w"

[commands.delete_line]
shortcut = "super+d"
`
