// +build android
// +build !sdl2

package nk

import (
	"sync"
	"time"
	"unsafe"

	"github.com/xlab/android-go/android"
	"github.com/xlab/android-go/egl"
)

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
	dev := state.ogl
	var width, height int32
	image := NkFontAtlasBake(state.atlas, &width, &height, FontAtlasRgba32)
	deviceUploadAtlas(image, width, height)
	NkFontAtlasEnd(state.atlas, NkHandleId(int32(state.ogl.font_tex[0])), &dev.null)
	if font := state.atlas.DefaultFont(); font != nil {
		NkStyleSetFont(state.ctx, font.Handle())
	}
}

type PlatformKeyEvent struct {
	Activity  *android.NativeActivity
	Action    int32
	KeyCode   int32
	MetaState int32
}

type PlatformTouchEvent struct {
	Action int32
	X, Y   int32
}

func NkPlatformInput(touch *PlatformTouchEvent, key *PlatformKeyEvent) {
	if touch != nil && state != nil {
		state.touch.Add(*touch)
	}
	if key != nil {
		state.keys.Add(*key)
		if key.Action == android.KeyEventActionUp && key.Activity != nil {
			switch key.KeyCode {
			case android.KeycodeUnknown:
			case android.KeycodeEnter, android.KeycodeTab:
			default:
				r, err := key.Activity.KeyEventGetUnicodeChar(key.Action, key.KeyCode, key.MetaState)
				if err == nil && r > 0 {
					state.text += string(r)
				}
			}
		}
	}
}

func NkPlatformNewFrame() {
	display := state.display
	ctx := state.ctx

	// for Android scale ratio can be tricky stuff
	state.width, state.height = display.Width, display.Height
	state.display_width, state.display_height = display.Width, display.Height
	state.fbScaleX = float32(state.display_width) / float32(state.width)
	state.fbScaleY = float32(state.display_height) / float32(state.height)

	NkInputBegin(ctx)

	for _, r := range state.text {
		NkInputUnicode(ctx, Rune(r))
	}
	state.text = ""

	if state.touch.CurrentAction() == android.MotionEventActionUp {
		ctx.Input().Mouse().SetPos(0, 0)
	}
	state.touch.Observe(func(action, x, y int32) {
		switch action {
		case android.MotionEventActionDown:
			ctx.Input().Mouse().SetPos(x, y)
			NkInputButton(ctx, ButtonLeft, x, y, 1)
		case android.MotionEventActionMove:
			NkInputMotion(ctx, x, y)
		case android.MotionEventActionUp:
			ctx.Input().Mouse().SetPos(x, y)
			NkInputButton(ctx, ButtonLeft, x, y, 0)
		}
	})

	toDeliver := map[int32][]int32{
		android.KeycodeForwardDel: {},
		android.KeycodeEnter:      {},
		android.KeycodeTab:        {},
		android.KeycodeDel:        {},
		android.KeycodeDpadUp:     {},
		android.KeycodeDpadDown:   {},
		android.KeycodeDpadLeft:   {},
		android.KeycodeDpadRight:  {},
		android.KeycodeShiftLeft:  {},
		android.KeycodeShiftRight: {},
		android.KeycodeCtrlLeft:   {},
		android.KeycodeCtrlRight:  {},
	}
	state.keys.Observe(func(action, key, metaState int32) {
		toDeliver[key] = append(toDeliver[key], action)
	})
	for key, actions := range toDeliver {
		deliverKey := func(key, isDown int32) {
			switch key {
			case android.KeycodeForwardDel:
				NkInputKey(ctx, KeyDel, isDown)
			case android.KeycodeEnter:
				NkInputKey(ctx, KeyEnter, isDown)
			case android.KeycodeTab:
				NkInputKey(ctx, KeyTab, isDown)
			case android.KeycodeDel:
				NkInputKey(ctx, KeyBackspace, isDown)
			case android.KeycodeDpadUp:
				NkInputKey(ctx, KeyUp, isDown)
			case android.KeycodeDpadDown:
				NkInputKey(ctx, KeyDown, isDown)
			case android.KeycodeDpadLeft:
				NkInputKey(ctx, KeyLeft, isDown)
			case android.KeycodeDpadRight:
				NkInputKey(ctx, KeyRight, isDown)
			case android.KeycodeShiftLeft, android.KeycodeShiftRight:
				NkInputKey(ctx, KeyShift, isDown)
			case android.KeycodeCtrlLeft, android.KeycodeCtrlRight:
				NkInputKey(ctx, KeyCtrl, isDown)
			}
		}
		switch {
		case len(actions) == 0:
			deliverKey(key, 0)
		case len(actions) == 0:
			deliverKey(key, isDown(actions[0]))
		case len(actions) >= 2:
			if actions[0] == actions[1] {
				deliverKey(key, isDown(actions[0]))
			} else {
				deliverKey(key, isDown(actions[0]))
				deliverKey(key, isDown(actions[1]))
			}
		}
	}
	NkInputEnd(ctx)
}

func isDown(action int32) int32 {
	if action == android.KeyEventActionDown {
		return 1
	}
	return 0
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
	display *egl.DisplayHandle
	touch   *touchHandler
	keys    *keyHandler

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

func newPlatformState() *platformState {
	return &platformState{
		ogl:   &platformDevice{},
		touch: newTouchHandler(),
		keys:  newKeyHandler(),
	}
}

func NkPlatformDisplayHandle() *egl.DisplayHandle {
	if state != nil {
		return state.display
	}
	return nil
}

var state *platformState

type platformDevice struct {
	cmds *Buffer
	null DrawNullTexture

	vbo, vao, ebo []uint32
	prog          uint32
	vert_shdr     uint32
	frag_shdr     uint32

	attrib_pos   uint32
	attrib_uv    uint32
	attrib_col   uint32
	uniform_tex  int32
	uniform_proj int32

	font_tex []uint32
}

const touchDecayTime = 500 * time.Millisecond

type touchHandler struct {
	current *PlatformTouchEvent
	queue   []PlatformTouchEvent
	mux     *sync.RWMutex
	decay   *time.Timer
}

func newTouchHandler() *touchHandler {
	h := &touchHandler{
		queue: make([]PlatformTouchEvent, 0, 1024),
		mux:   new(sync.RWMutex),
	}
	h.decay = time.NewTimer(time.Minute)
	h.decay.Stop()
	go func() {
		for range h.decay.C {
			h.decay.Reset(touchDecayTime)
		}
	}()
	return h
}

func (t *touchHandler) Add(ev PlatformTouchEvent) {
	t.mux.Lock()
	t.queue = append(t.queue, ev)
	t.current = &ev
	t.decay.Reset(touchDecayTime)
	t.mux.Unlock()
}

func (t *touchHandler) Reset() {
	t.mux.Lock()
	if ql := len(t.queue); ql > 0 {
		t.queue = t.queue[:0]
	}
	t.mux.Unlock()
}

func (t *touchHandler) Observe(fn func(action, x, y int32)) {
	t.mux.Lock()
	if len(t.queue) > 0 {
		for i := range t.queue {
			fn(t.queue[i].Action, t.queue[i].X, t.queue[i].Y)
		}
		t.queue = t.queue[:0]
	}
	t.mux.Unlock()
}

func (t *touchHandler) CurrentPos() (x int32, y int32) {
	t.mux.RLock()
	if t.current != nil {
		x, y = t.current.X, t.current.Y
	}
	t.mux.RUnlock()
	return x, y
}

func (t *touchHandler) CurrentAction() (a int32) {
	t.mux.RLock()
	if t.current != nil {
		a = t.current.Action
	} else {
		a = -1
	}
	t.mux.RUnlock()
	return a
}

const keyDecayTime = 500 * time.Millisecond

type keyHandler struct {
	current *PlatformKeyEvent
	queue   []PlatformKeyEvent
	mux     *sync.RWMutex
	decay   *time.Timer
}

func newKeyHandler() *keyHandler {
	h := &keyHandler{
		queue: make([]PlatformKeyEvent, 0, 1024),
		mux:   new(sync.RWMutex),
	}
	h.decay = time.NewTimer(time.Minute)
	h.decay.Stop()
	go func() {
		for range h.decay.C {
			h.decay.Reset(keyDecayTime)
		}
	}()
	return h
}

func (h *keyHandler) Add(ev PlatformKeyEvent) {
	h.mux.Lock()
	h.queue = append(h.queue, ev)
	h.current = &ev
	h.decay.Reset(keyDecayTime)
	h.mux.Unlock()
}

func (h *keyHandler) Reset() {
	h.mux.Lock()
	if ql := len(h.queue); ql > 0 {
		h.queue = h.queue[:0]
	}
	h.mux.Unlock()
}

func (h *keyHandler) Observe(fn func(action, key, metaState int32)) {
	h.mux.Lock()
	if len(h.queue) > 0 {
		for i := range h.queue {
			ev := &h.queue[i]
			fn(ev.Action, ev.KeyCode, ev.MetaState)
		}
		h.queue = h.queue[:0]
	}
	h.mux.Unlock()
}

func (h *keyHandler) Current() (a, k int32) {
	h.mux.RLock()
	if h.current != nil {
		a = h.current.Action
		k = h.current.KeyCode
	}
	h.mux.RUnlock()
	return a, k
}
