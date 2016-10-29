package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT -DNK_INCLUDE_VERTEX_BUFFER_OUTPUT
#cgo windows CFLAGS: -D_GLFW_WIN32
#cgo darwin CFLAGS: -D_GLFW_COCOA -D_GLFW_USE_CHDIR -D_GLFW_USE_MENUBAR -D_GLFW_USE_RETINA -Wno-deprecated-declarations
#cgo linux,!wayland CFLAGS: -D_GLFW_X11
#cgo linux,wayland CFLAGS: -D_GLFW_WAYLAND
#cgo freebsd,!wayland CFLAGS: -D_GLFW_X11 -D_GLFW_HAS_GLXGETPROCADDRESSARB -D_GLFW_HAS_DLOPEN
#cgo freebsd,wayland CFLAGS: -D_GLFW_WAYLAND -D_GLFW_HAS_DLOPEN

#cgo linux,!wayland LDFLAGS: -lglfw -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt
#cgo linux,wayland LDFLAGS: -lglfw -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama -ldl -lrt
#cgo darwin LDFLAGS: -lglfw3 -framework Cocoa -framework OpenGL -framework IOKit -framework CoreVideo -lm
#cgo windows LDFLAGS: -lglfw3 -lopengl32 -lgdi32 -lm
#cgo freebsd,!wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama
#cgo freebsd,wayland LDFLAGS: -lGL -lX11 -lXrandr -lXxf86vm -lXi -lXcursor -lm -lXinerama

#include <string.h>

#include <GLFW/glfw3.h>

#define NK_IMPLEMENTATION
#define NK_GLFW_GL2_IMPLEMENTATION
#include "nuklear.h"
#include "nuklear_glfw_gl2.h"
*/
import "C"

type GLFWwindow C.GLFWwindow
