package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT
#cgo linux LDFLAGS: -lSDL -lSDL_gfx -lm
#cgo darwin LDFLAGS: -lSDL -lSDL_gfx -lm
#cgo windows LDFLAGS: -lmingw32 -lSDLmain -lSDL -lm

#define NK_IMPLEMENTATION
#define NK_SDL_IMPLEMENTATION
#include "nuklear.h"
#include "nuklear_sdl.h"

#include <SDL/SDL.h>
*/
import "C"

type SDLSurface C.SDL_Surface

type SDLEvent C.SDL_Event
