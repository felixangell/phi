<p align="center"><img src="https://raw.githubusercontent.com/felixangell/phi/gh-pages/images/icon96.png"></p>

<h1>phi</h1>

[![Build Status](https://travis-ci.org/felixangell/phi.svg?branch=master)](https://travis-ci.org/felixangell/phi)

Phi is a minimal code editor designed to look pretty, run fast, and be easy
to configure and use. It's primary function is for editing code.

The editor is still a work in progress. There is a chance it will **eat up your battery**, **run quite slowly**, 
and probably **crash frequently**.

<p align="center">**Do not edit your precious files with this editor!**</p>

Here's a screenshot of Phi in action:

<p align="center"><img src="https://raw.githubusercontent.com/felixangell/phi/gh-pages/images/screenshot.png"></p>

# goals
The editor must:

* run fast;
* load and edit large files with ease;
* look pretty; and finally
* be easy to use

## non-goals
The editor probably wont:

* have any plugin support;
* be very customizable in terms of layout;
* support many non utf8 encodings;
* support non true-type-fonts;
* support right-to-left languages;

Avoiding most of these is to avoid complexity in the code-base
and general architecture of the editor and is beyond the scope of this project currently.

# why?
The editor does not exist as a serious replacement to Sublime Text/VSCode/Emacs/[editor name here]. 

Though one of my big goals for the project is to possibly replace sublime text for my own personal use. Thus the editor is somewhat optimized for my own work-flow.

The code is up purely for people to look at and maybe use or contribute or whatever. Sharing is caring!

# reporting bugs/troubleshooting
Note the editor is still unstable. Please report any bugs you find so I can
squash them! It is appreciated if you skim the issue (or search!) handler to make sure
you aren't reporting duplicate bugs.

## before filing an issue
Just to make sure it's an issue with the editor currently and not due to a 
broken change - please can you:

* make sure the repository is up to date
* make sure all the dependencies are updated, especially "github.com/felixangell/strife"
* try removing the ~/.phi-config folder manually and letting the editor re-load it

# building
See the [BUILDING](/BUILDING.md) file.

# license
[MIT License](/LICENSE)
