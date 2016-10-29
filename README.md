## Nuklear [![GoDoc](https://godoc.org/github.com/golang-ui/nuklear/nk?status.svg)](https://godoc.org/github.com/golang-ui/nuklear/nk)

Package nk provides Go bindings for nuklear.h — a small ANSI C gui library. See [github.com/vurtun/nuklear](https://github.com/vurtun/nuklear).<br />
All the binding code has automatically been generated with rules defined in [nk.yml](/nk.yml).

An idiomatic wrapper package isn't coming I guess, because it will require strong interest in further development and I have no time for that now. The `nk` package is fine for the start, then we'll figure out something better that just a wrapper.

### Features (plain C version)

* Immediate mode graphical user interface toolkit
* Single header library
* Written in C89 (ANSI C)
* Small codebase (~15kLOC)
* Focus on portability, efficiency and simplicity
* No dependencies (not even the standard library if not wanted)
* Fully skinnable and customizable
* Low memory footprint with total memory control if needed or wanted
* UTF-8 support
* No global or hidden state
* Customizable library modules (you can compile and use only what you need)
* Optional font baker and vertex buffer output

#### About nuklear.h

This is a minimal state immediate mode graphical user interface single header
toolkit written in ANSI C and licensed under public domain.
It was designed as a simple embeddable user interface for application and does
not have any dependencies, a default renderbackend or OS window and input handling
but instead provides a very modular library approach by using simple input state
for input and draw commands describing primitive shapes as output.
So instead of providing a layered library that tries to abstract over a number
of platform and render backends it only focuses on the actual UI.

### Installation of nk

#### OS X

```bash
$ brew install glfw3 # must be >= 3.2
$ go get github.com/golang-ui/nuklear/nk

# consult your distro package archives for GLFW if you are under Linux
```

Both OpenGL 2.1 and OpenGL 3.3 contexts are working fine, but by default OpenGL 2.1 is used.

#### Windows

1. Get MinGW compiler toolchain and MSYS via [MinGW installer](https://sourceforge.net/projects/mingw/files/latest/download);
2. Get GLFW 3.2.1 pre-built distro from the official site http://www.glfw.org/;
3. Unpack GLFW in some simple location on `C:\`;
4. Open MSYS shell (usually `C:\MinGW\msys\1.0\msys.bat`);

Then everything should go smooth:
```
$ go version
go version go1.6.2 windows/386

$ gcc -v
COLLECT_GCC=C:\MinGW\bin\gcc.exe
Thread model: posix
gcc version 5.3.0 (GCC)

$ CGO_CFLAGS="-I/c/dev/glfw-3.2.1/include" CGO_LDFLAGS="-L/c/dev/glfw-3.2.1/lib-mingw" go install github.com/golang-ui/nuklear/nk
```

See a [screenshot](https://cl.ly/1r0j2Y3D2I2W/Screen%20Shot%202016-10-29%20at%2015.11.24.png) of running [cmd/nk-example](/cmd/nk-example). The OpenGL 2.1 context backend works fine under Windows 7.

### Demo

There is an example app [nk-example](https://github.com/golang-ui/nuklear/blob/master/cmd/nk-example/main.go) that shows the usage of Nuklear GUI library, based on the official demos.

```bash
$ go get github.com/golang-ui/nuklear/cmd/nk-example

$ nk-example
2016/09/23 23:13:09 glfw: created window 400x500
2016/09/23 23:13:10 [INFO] button pressed!
2016/09/23 23:13:10 [INFO] button pressed!
2016/09/23 23:13:10 [INFO] button pressed!
```

<img alt="demo screenshot nuklear" src="assets/demo.png" width="500"/>

Another more realistic Golang app that uses Nuklear to do its GUI, [a simple WebM player](https://github.com/xlab/libvpx-go):

<a href="https://www.youtube.com/watch?v=5kj5ApnhPAE"><img alt="nuklear screenshot webm" src="assets/demo2.png" width="800"/></a>

### Rebuilding the package

You will need to get the [cgogen](https://git.io/cgogen) tool installed first.

```
$ git clone https://github.com/golang-ui/nuklear && cd nuklear
$ make clean
$ make
```

### License

All the code except when stated otherwise is licensed under the [MIT license](https://xlab.mit-license.org).
