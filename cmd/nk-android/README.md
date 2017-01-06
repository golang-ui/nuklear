Nuklear Activity
================

This android app is built using [android-go](http://github.com/xlab/android-go) framework,
and the GUI is done via [nuklear](http://github.com/golang-ui/nuklear) that uses Android NDK as a backend.
It manages an OpenGL ES 3.0 context via EGL to draw stuff on the screen and gets touch input from the sensors.

### Prerequisites

There is no additional prerequisites, the project fully inherits the same
structure as the [example] app provides, make sure you were able to run it
first. Just to make sure that everything works smoothly for your OS, environment
setup and the device itself.

[example]: https://github.com/xlab/android-go/tree/master/example

### Structure

```
$ tree .
.
├── Makefile
├── android
│   ├── AndroidManifest.xml
│   ├── Makefile
│   ├── jni
│   │   ├── Android.mk
│   │   └── Application.mk
│   └── res
├── util.go
└── main.go

3 directories, 7 files
```

### Running

```
$ make
$ make install
$ make listen
```
