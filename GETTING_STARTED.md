# Getting Started
This document is a quick getting started guide for users who are interested
in trying Phi.

## Installing
Notes for installing and building are detailed in the README file. You will probably
have to build from source. This shouldn't be too much of a hassle with Linux 
and MacOS - Windows might be a bit of a hassle.

## Configuration Files
Configuration files for Phi are located in `~/.phi-editor`. This is in your 
$HOME directory.

## Fonts
The font loading in Phi is still very much a work in progress. Loading fonts
in a platform independent way is a particular struggle when dealing with Linux.
Below are some notes for users to help get started with the editor.

### Linux/Ubuntu
Ubuntu's fonts are in the `/usr/share/fonts` folder, though they are categorized
by font type. Here is a quick example to get you started on Ubuntu:

```toml
font_path = "/usr/share/fonts/truetype/dejavu"
font_face = "DejaVuSansMono"
```

For Linux in general, you may have to hunt down your fonts folder yourself
and set the `font_path` variable accordingly.