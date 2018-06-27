// +build !android
// +build sdl2

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

	"github.com/veandco/go-sdl2/sdl"
)

type PlatformInitOption int

const (
	PlatformDefault PlatformInitOption = iota
	PlatformInstallCallbacks
)

func textScrollCallback(e sdl.Event, userdata interface{}) bool {
	state := userdata.(*platformState)
	switch t := e.(type) {
	case *sdl.MouseWheelEvent:
		state.scroll.SetX(state.scroll.X() + float32(t.X))
		state.scroll.SetY(state.scroll.Y() + float32(t.Y))
	case *sdl.KeyboardEvent:
		state.text += sdl.GetKeyName(sdl.GetKeyFromScancode(t.Keysym.Scancode))
	}
	return true
}

func NkPlatformInit(win *sdl.Window, context sdl.GLContext, opt PlatformInitOption) *Context {
	state.win = win
	if opt == PlatformInstallCallbacks {
		sdl.AddEventWatchFunc(textScrollCallback, state)
	}
	state.ctx = NewContext()
	NkInitDefault(state.ctx, nil)
	deviceCreate()
	return state.ctx
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
	state.display_width, state.display_height = win.GetSize()
	state.fbScaleX = float32(state.display_width) / float32(state.width)
	state.fbScaleY = float32(state.display_height) / float32(state.height)

	NkInputBegin(ctx)
	for _, r := range state.text {
		NkInputUnicode(ctx, Rune(r))
	}

	// optional grabbing behavior
	m := ctx.Input().Mouse()
	if m.Grab() {
		sdl.SetRelativeMouseMode(true)
	} else if m.Ungrab() {
		sdl.SetRelativeMouseMode(false)
	}

	keys := sdl.GetKeyboardState()

	NkInputKey(ctx, KeyDel, int32(keys[sdl.SCANCODE_DELETE]))
	NkInputKey(ctx, KeyEnter, int32(keys[sdl.SCANCODE_RETURN]))
	NkInputKey(ctx, KeyEnter, int32(keys[sdl.SCANCODE_KP_ENTER]))
	NkInputKey(ctx, KeyTab, int32(keys[sdl.SCANCODE_TAB]))
	NkInputKey(ctx, KeyBackspace, int32(keys[sdl.SCANCODE_BACKSPACE]))
	NkInputKey(ctx, KeyUp, int32(keys[sdl.SCANCODE_UP]))
	NkInputKey(ctx, KeyDown, int32(keys[sdl.SCANCODE_DOWN]))
	NkInputKey(ctx, KeyTextStart, int32(keys[sdl.SCANCODE_HOME]))
	NkInputKey(ctx, KeyTextEnd, int32(keys[sdl.SCANCODE_END]))
	NkInputKey(ctx, KeyScrollStart, int32(keys[sdl.SCANCODE_HOME]))
	NkInputKey(ctx, KeyScrollEnd, int32(keys[sdl.SCANCODE_END]))
	NkInputKey(ctx, KeyScrollUp, int32(keys[sdl.SCANCODE_PAGEUP]))
	NkInputKey(ctx, KeyScrollDown, int32(keys[sdl.SCANCODE_PAGEDOWN]))

	shiftHeld := int32(0)
	if keys[sdl.KMOD_LSHIFT] == 1 || keys[sdl.KMOD_RSHIFT] == 1 {
		shiftHeld = int32(1)
	}
	NkInputKey(ctx, KeyShift, shiftHeld)

	controlHeld := false
	if keys[sdl.KMOD_LCTRL] == 1 || keys[sdl.KMOD_RCTRL] == 1 {
		controlHeld = true
	}

	if controlHeld {
		NkInputKey(ctx, KeyCopy, int32(keys[sdl.SCANCODE_C]))
		NkInputKey(ctx, KeyPaste, int32(keys[sdl.SCANCODE_V]))
		NkInputKey(ctx, KeyCut, int32(keys[sdl.SCANCODE_X]))
		NkInputKey(ctx, KeyTextUndo, int32(keys[sdl.SCANCODE_Z]))
		NkInputKey(ctx, KeyTextRedo, int32(keys[sdl.SCANCODE_R]))
		NkInputKey(ctx, KeyTextWordLeft, int32(keys[sdl.SCANCODE_LEFT]))
		NkInputKey(ctx, KeyTextWordRight, int32(keys[sdl.SCANCODE_RIGHT]))
		NkInputKey(ctx, KeyTextLineStart, int32(keys[sdl.SCANCODE_B]))
		NkInputKey(ctx, KeyTextLineEnd, int32(keys[sdl.SCANCODE_E]))
	} else {
		NkInputKey(ctx, KeyLeft, int32(keys[sdl.SCANCODE_LEFT]))
		NkInputKey(ctx, KeyRight, int32(keys[sdl.SCANCODE_RIGHT]))
		NkInputKey(ctx, KeyCopy, 0)
		NkInputKey(ctx, KeyPaste, 0)
		NkInputKey(ctx, KeyCut, 0)
		NkInputKey(ctx, KeyShift, 0)
	}

	x, y, mouseState := sdl.GetMouseState()
	NkInputMotion(ctx, int32(x), int32(y))
	if m := ctx.Input().Mouse(); m.Grabbed() {
		prevX, prevY := m.Prev()
		win.WarpMouseInWindow(int32(prevX), int32(prevY))
		m.SetPos(prevX, prevY)
	}

	NkInputButton(ctx, ButtonLeft, int32(x), int32(y), int32(mouseState&sdl.ButtonLMask()))
	NkInputButton(ctx, ButtonMiddle, int32(x), int32(y), int32(mouseState&sdl.ButtonMMask()))
	NkInputButton(ctx, ButtonRight, int32(x), int32(y), int32(mouseState&sdl.ButtonRMask()))

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
	win *sdl.Window

	width          int32
	height         int32
	display_width  int32
	display_height int32

	ogl   *platformDevice
	ctx   *Context
	atlas *FontAtlas

	fbScaleX float32
	fbScaleY float32

	text   string
	scroll Vec2
}

func NkPlatformDisplayHandle() *sdl.Window {
	if state != nil {
		return state.win
	}
	return nil
}
