Nuklear Activity
================

This android app is built using [android-go](http://github.com/xlab/android-go) framework,
and the GUI is done via [nuklear](http://github.com/golang-ui/nuklear) that uses Android NDK as a backend.
It manages an OpenGL ES 2.0 or 3.0 context via EGL to draw stuff on the screen and gets touch input from the sensors.

### Prerequisites

There is no additional prerequisites, the project fully inherits the same
structure as the Android [example] app provides, make sure you were able to run it
first. Just to make sure that everything works smoothly for your OS, environment
setup and the device itself.

[example]: https://github.com/xlab/android-go/tree/master/example

### Structure

```
$ tree .
.
├── Makefile
├── README.md
├── android
│   ├── AndroidManifest.xml
│   ├── Makefile
│   └── jni
│       ├── Android.mk
│       └── Application.mk
├── assets
│   └── DroidSans.ttf
├── bindata.go
├── main.go
└── util.go

4 directories, 9 files
```

Droid Sans font is packed using go-bindata and being linked statically on the compile time.

### Screenshots

<img src="https://cl.ly/2c0s3R3Q2g3V/Screenshot_20170108-045948.png" width="300"/>&nbsp;&nbsp;<img src="https://cl.ly/2X133z0Z3S1j/Screenshot_20170108-050051.png" width="300"/>

[Click for video](https://www.youtube.com/watch?v=3-MiceegZlM)

### Running

```
$ make
$ make install
$ make listen
```
