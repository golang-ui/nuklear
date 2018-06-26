// +build !android
// +build !sdl2

package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT -DNK_INCLUDE_VERTEX_BUFFER_OUTPUT -Wno-implicit-function-declaration
#cgo windows LDFLAGS: -Wl,--allow-multiple-definition
#include <string.h>

#include "nuklear.h"
*/
import "C"
import (
	"unsafe"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type PlatformInitOption int

const (
	PlatformDefault PlatformInitOption = iota
	PlatformInstallCallbacks
)

func NkPlatformInit(win *glfw.Window, opt PlatformInitOption) *Context {
	state.win = win
	if opt == PlatformInstallCallbacks {
		win.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
			state.scroll.SetX(state.scroll.X() + float32(xoff))
			state.scroll.SetY(state.scroll.Y() + float32(yoff))
		})
		win.SetCharCallback(func(w *glfw.Window, char rune) {
			if len(state.text) < 256 { // NK_GLFW_TEXT_MAX
				state.text += string(char)
			}
		})
	}

	state.ctx = NewContext()
	NkInitDefault(state.ctx, nil)
	deviceCreate()
	return state.ctx

	// TODO(xlab): clipboard
	// state.ctx.clip.copy = nk_glfw3_clipbard_copy;
	// state.ctx.clip.paste = nk_glfw3_clipbard_paste;
	// state.ctx.clip.userdata = nk_handle_ptr(0);
}

func NkPlatformShutdown() {
	NkFontAtlasClear(state.atlas)
	NkFree(state.ctx)
	deviceDestroy()
	state = nil
}

func NkFontStashBegin(atlas **FontAtlas) {
	state.atlas = NewFontAtlas()
	NkFontAtlasInitDefault(state.atlas)
	NkFontAtlasBegin(state.atlas)
	*atlas = state.atlas
}

func NkFontStashEnd() {
	var width, height int32
	image := NkFontAtlasBake(state.atlas, &width, &height, FontAtlasRgba32)
	deviceUploadAtlas(image, width, height)
	NkFontAtlasEnd(state.atlas, NkHandleId(int32(state.ogl.font_tex)), &state.ogl.null)
	if font := state.atlas.DefaultFont(); font != nil {
		NkStyleSetFont(state.ctx, font.Handle())
	}
}

func NkPlatformNewFrame() {
	win := state.win
	ctx := state.ctx
	state.width, state.height = win.GetSize()
	state.display_width, state.display_height = win.GetFramebufferSize()
	state.fbScaleX = float32(state.display_width) / float32(state.width)
	state.fbScaleY = float32(state.display_height) / float32(state.height)

	NkInputBegin(ctx)
	for _, r := range state.text {
		NkInputUnicode(ctx, Rune(r))
	}

	// optional grabbing behavior
	m := ctx.Input().Mouse()
	if m.Grab() {
		win.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	} else if m.Ungrab() {
		win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}

	NkInputKey(ctx, KeyDel, keyPressed(win, glfw.KeyDelete))
	NkInputKey(ctx, KeyEnter, keyPressed(win, glfw.KeyEnter))
	NkInputKey(ctx, KeyEnter, keyPressed(win, glfw.KeyKPEnter))
	NkInputKey(ctx, KeyTab, keyPressed(win, glfw.KeyTab))
	NkInputKey(ctx, KeyBackspace, keyPressed(win, glfw.KeyBackspace))
	NkInputKey(ctx, KeyUp, keyPressed(win, glfw.KeyUp))
	NkInputKey(ctx, KeyDown, keyPressed(win, glfw.KeyDown))
	NkInputKey(ctx, KeyTextStart, keyPressed(win, glfw.KeyHome))
	NkInputKey(ctx, KeyTextEnd, keyPressed(win, glfw.KeyEnd))
	NkInputKey(ctx, KeyScrollStart, keyPressed(win, glfw.KeyHome))
	NkInputKey(ctx, KeyScrollEnd, keyPressed(win, glfw.KeyEnd))
	NkInputKey(ctx, KeyScrollUp, keyPressed(win, glfw.KeyPageUp))
	NkInputKey(ctx, KeyScrollDown, keyPressed(win, glfw.KeyPageDown))
	NkInputKey(ctx, KeyShift, keysPressed(win, glfw.KeyLeftShift, glfw.KeyRightShift))
	if keysPressed(win, glfw.KeyLeftControl, glfw.KeyRightControl) > 0 {
		NkInputKey(ctx, KeyCopy, keyPressed(win, glfw.KeyC))
		NkInputKey(ctx, KeyPaste, keyPressed(win, glfw.KeyV))
		NkInputKey(ctx, KeyCut, keyPressed(win, glfw.KeyX))
		NkInputKey(ctx, KeyTextUndo, keyPressed(win, glfw.KeyZ))
		NkInputKey(ctx, KeyTextRedo, keyPressed(win, glfw.KeyR))
		NkInputKey(ctx, KeyTextWordLeft, keyPressed(win, glfw.KeyLeft))
		NkInputKey(ctx, KeyTextWordRight, keyPressed(win, glfw.KeyRight))
		NkInputKey(ctx, KeyTextLineStart, keyPressed(win, glfw.KeyB))
		NkInputKey(ctx, KeyTextLineEnd, keyPressed(win, glfw.KeyE))
	} else {
		NkInputKey(ctx, KeyLeft, keyPressed(win, glfw.KeyLeft))
		NkInputKey(ctx, KeyRight, keyPressed(win, glfw.KeyRight))
		NkInputKey(ctx, KeyCopy, 0)
		NkInputKey(ctx, KeyPaste, 0)
		NkInputKey(ctx, KeyCut, 0)
		NkInputKey(ctx, KeyShift, 0)
	}
	x, y := win.GetCursorPos()
	NkInputMotion(ctx, int32(x), int32(y))
	if m := ctx.Input().Mouse(); m.Grabbed() {
		prevX, prevY := m.Prev()
		win.SetCursorPos(float64(prevX), float64(prevY))
		m.SetPos(prevX, prevY)
	}

	NkInputButton(ctx, ButtonLeft, int32(x), int32(y), buttonPressed(win, glfw.MouseButtonLeft))
	NkInputButton(ctx, ButtonMiddle, int32(x), int32(y), buttonPressed(win, glfw.MouseButtonMiddle))
	NkInputButton(ctx, ButtonRight, int32(x), int32(y), buttonPressed(win, glfw.MouseButtonRight))
	NkInputScroll(ctx, state.scroll)
	NkInputEnd(ctx)
	state.text = ""
	state.scroll.Reset()
}

var (
	sizeofDrawIndex = unsafe.Sizeof(DrawIndex(0))
	emptyVertex     = platformVertex{}
)

type platformVertex struct {
	position [2]float32
	uv       [2]float32
	col      [4]Byte
}

const (
	platformVertexSize  = unsafe.Sizeof(platformVertex{})
	platformVertexAlign = unsafe.Alignof(platformVertex{})
)

type platformState struct {
	win *glfw.Window

	width          int
	height         int
	display_width  int
	display_height int

	ogl   *platformDevice
	ctx   *Context
	atlas *FontAtlas

	fbScaleX float32
	fbScaleY float32

	text   string
	scroll Vec2
}

func NkPlatformDisplayHandle() *glfw.Window {
	if state != nil {
		return state.win
	}
	return nil
}

func keyPressed(win *glfw.Window, key glfw.Key) int32 {
	if win.GetKey(key) == glfw.Press {
		return 1
	}
	return 0
}

func buttonPressed(win *glfw.Window, button glfw.MouseButton) int32 {
	if win.GetMouseButton(button) == glfw.Press {
		return 1
	}
	return 0
}

func keysPressed(win *glfw.Window, keys ...glfw.Key) int32 {
	for i := range keys {
		if win.GetKey(keys[i]) == glfw.Press {
			return 1
		}
	}
	return 0
}
