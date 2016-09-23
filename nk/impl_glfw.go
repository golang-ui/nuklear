package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT -DNK_INCLUDE_VERTEX_BUFFER_OUTPUT
#cgo linux LDFLAGS: -lglfw -lGL -lm
#cgo darwin LDFLAGS: -lglfw3 -framework OpenGL -lm
#cgo windows LDFLAGS: -lglfw3 -lopengl32 -lm

#include <string.h>

#include <GLFW/glfw3.h>

#define NK_IMPLEMENTATION
#define NK_GLFW_GL3_IMPLEMENTATION
#include "nuklear.h"
#include "nuklear_glfw_gl3.h"
*/
import "C"

type GLFWwindow C.GLFWwindow
