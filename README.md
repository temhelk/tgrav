# tgrav
A simple terminal gravity simulator written in Go

## Controls
* `+`/`-` changes the speed of the simulation
* `f` enables the rendering of the acceleration magnitude
* `mouse wheel` scrolling zooms in and out
* `mouse panning` (while holding left mouse button) moves the view around

## Notes
* For better rendering resolution all bodies are drawn as dots using Braille symbols
* The camera automatically moves with the center of mass of the system
* I'm running it in [Kitty](https://github.com/kovidgoyal/kitty) on Linux, it also works on Windows, but can be slow (tested with Windows Terminal) 

## Demo
### A system of 4 bodies showing the acceleration amplitude surrounding them
<a href="https://asciinema.org/a/656354" ><img src="https://asciinema.org/a/656354.png" width="400"></a>

### A system of 3 bodies with the same mass
<a href="https://asciinema.org/a/656358" ><img src="https://asciinema.org/a/656358.png" width="400"></a>

### Gravitational slingshot
<a href="https://asciinema.org/a/656397" ><img src="https://asciinema.org/a/656397.png" width="400"></a>
