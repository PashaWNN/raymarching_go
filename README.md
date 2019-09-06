# Ray marching (Go)

I've implemented this algorithm again, now in Go, because I just started to learn this language.
And yeah, this version supports multithreading and will use all of your CPU cores.

## Changes in version 2

* now generates iterative GIF instead of staic PNG
* fixed terrible bug: program was 150 times slower because all algorithm was inside for loop with number of iterations

## Planned in version 3

* port to WebAssembly to allow play around directly in browser

![demo][image]

[image]: https://github.com/pashawnn/raymarching_go/blob/master/picture.png

Feel free to contribute if you want!