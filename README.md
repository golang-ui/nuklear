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

### Installation of nk

```bash
$ brew install glfw3 # must be >= 3.2
$ go get github.com/golang-ui/nuklear/nk

# consult your distro package archives for GLFW if you are under Linux
```

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
