package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/android-go/android"
	"github.com/xlab/android-go/app"
	gl "github.com/xlab/android-go/gles3"
)

func init() {
	app.SetLogTag("NuklearActivity")
}

const (
	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

var pt float32

func initPt() {
	// we live in a world of 600x600
	const squareSide = 600
	display := nk.NkPlatformDisplayHandle()
	if display.Width < display.Height {
		if newPt := float32(display.Width) / squareSide; newPt > 0 {
			pt = newPt
		}
	} else if newPt := float32(display.Height) / squareSide; newPt > 0 {
		pt = newPt
	}
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	nativeWindowEvents := make(chan app.NativeWindowEvent, 1)
	inputQueueEvents := make(chan app.InputQueueEvent, 1)
	inputQueueChan := make(chan *android.InputQueue, 1)

	var ctx *nk.Context
	appState := &State{
		buffer:  make([]byte, 256*1024), // for text box
		bgColor: nk.NkRgba(0, 145, 118, 255),
	}
	copy(appState.buffer, []byte("Write text here...\x00"))
	fpsTicker := time.NewTimer(time.Minute)
	fpsTicker.Stop()

	app.Main(func(a app.NativeActivity) {
		a.HandleNativeWindowEvents(nativeWindowEvents)
		a.HandleInputQueueEvents(inputQueueEvents)
		go func() {
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			app.HandleInputQueues(inputQueueChan, func() {
				a.InputQueueHandled()
			}, func(ev *android.InputEvent) {
				switch android.InputEventGetType(ev) {
				case android.InputEventTypeKey:
					key := android.KeyEventGetKeyCode(ev)
					action := android.KeyEventGetAction(ev)
					meta := android.KeyEventGetMetaState(ev)
					switch action {
					case android.KeyEventActionDown,
						android.KeyEventActionUp:
						nk.NkPlatformInput(nil, &nk.PlatformKeyEvent{
							Activity:  a.NativeActivity(),
							Action:    action,
							KeyCode:   key,
							MetaState: meta,
						})
					}
				case android.InputEventTypeMotion:
					action := android.MotionEventGetAction(ev)
					switch action {
					case android.MotionEventActionDown,
						android.MotionEventActionMove,
						android.MotionEventActionUp:
						x := android.MotionEventGetX(ev, 0)
						y := android.MotionEventGetY(ev, 0)
						nk.NkPlatformInput(&nk.PlatformTouchEvent{
							Action: action,
							X:      int32(x),
							Y:      int32(y),
						}, nil)
					}
				}
			})
		}()
		a.InitDone()
		const fpsTime = time.Second / 60
		for {
			select {
			case <-a.LifecycleEvents():
			case event := <-inputQueueEvents:
				switch event.Kind {
				case app.QueueCreated:
					inputQueueChan <- event.Queue
				case app.QueueDestroyed:
					inputQueueChan <- nil
				}
			case <-fpsTicker.C:
				gfxMain(a.NativeActivity(), ctx, appState)
				fpsTicker.Reset(fpsTime)
			case event := <-nativeWindowEvents:
				switch event.Kind {
				case app.NativeWindowRedrawNeeded:
					initPt()
					gfxMain(a.NativeActivity(), ctx, appState)
					a.NativeWindowRedrawDone()
					fpsTicker.Reset(fpsTime)
				case app.NativeWindowCreated:
					ctx = nk.NkPlatformInit(event.Window, nk.PlatformInstallCallbacks)
					if ctx == nil {
						log.Fatalln("Nuklear failed to init")
					}
					initPt()
					atlas := nk.NewFontAtlas()
					nk.NkFontStashBegin(&atlas)
					sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/DroidSans.ttf"), 20*pt, nil)
					// defaultFont := nk.NkFontAtlasAddDefault(atlas, 16*pt, nil)
					nk.NkFontStashEnd()
					if sansFont != nil {
						nk.NkStyleSetFont(ctx, sansFont.Handle())
					}
				case app.NativeWindowDestroyed:
					fpsTicker.Stop()
					nk.NkPlatformShutdown()
				}
			}
		}
	})
}

func filler(ctx *nk.Context, height float32) {
	nk.NkLayoutRowStatic(ctx, height, 0, 0)
}

// gfxMain is the main GUI code that is borrowed directly from the desktop example.
func gfxMain(activity *android.NativeActivity, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(50*pt, 300*pt, 500*pt, 500*pt)
	update := nk.NkBegin(ctx, "Demo", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update > 0 {
		filler(ctx, 10*pt)
		nk.NkLayoutRowStatic(ctx, 40*pt, int32(140*pt), 3)
		{
			if nk.NkButtonLabel(ctx, "button") > 0 {
				log.Println("[INFO] button pressed!")
				state.times++
			}
			nk.NkSpacing(ctx, 1)
			if nk.NkButtonLabel(ctx, "toggle keyboard") > 0 {
				toggleKeyboard(activity, state)
			}
		}
		nk.NkLayoutRowDynamic(ctx, 20*pt, 1)
		{
			nk.NkLabel(ctx, fmt.Sprintf("button pressed %d times", state.times), nk.TextAlignLeft)
		}
		filler(ctx, 10*pt)
		nk.NkLayoutRowDynamic(ctx, 100*pt, 1)
		{
			flags := nk.NkEditStringZeroTerminated(ctx,
				nk.EditBox, state.buffer, int32(len(state.buffer)), nk.NkFilterDefault)
			if flags&nk.EditActivated == nk.EditActivated && !state.kbShown {
				toggleKeyboard(activity, state)
			} else if flags&nk.EditDeactivated == nk.EditDeactivated && state.kbShown {
				toggleKeyboard(activity, state)
			}
		}
		filler(ctx, 10*pt)
		nk.NkLayoutRowDynamic(ctx, 30*pt, 2)
		{
			if nk.NkOptionLabel(ctx, "easy", flag(state.opt == Easy)) > 0 {
				state.opt = Easy
			}
			if nk.NkOptionLabel(ctx, "hard", flag(state.opt == Hard)) > 0 {
				state.opt = Hard
			}
		}
		filler(ctx, 10*pt)
		nk.NkLayoutRowDynamic(ctx, 25*pt, 1)
		{
			nk.NkPropertyInt(ctx, "Compression:", 0, &state.prop, 100, 10, 1)
		}
		filler(ctx, 10*pt)
		nk.NkLayoutRowDynamic(ctx, 20*pt, 1)
		{
			nk.NkLabel(ctx, "background:", nk.TextLeft)
		}
		filler(ctx, 10*pt)
		nk.NkLayoutRowDynamic(ctx, 25*pt, 1)
		{
			size := nk.NkVec2(nk.NkWidgetWidth(ctx), 400*pt)
			if nk.NkComboBeginColor(ctx, state.bgColor, size) > 0 {
				nk.NkLayoutRowDynamic(ctx, 120*pt, 1)
				state.bgColor = nk.NkColorPicker(ctx, state.bgColor, nk.ColorFormatRGBA)
				nk.NkLayoutRowDynamic(ctx, 25*pt, 1)
				r, g, b, a := state.bgColor.RGBAi()
				r = nk.NkPropertyi(ctx, "#R:", 0, r, 255, 1, 1)
				g = nk.NkPropertyi(ctx, "#G:", 0, g, 255, 1, 1)
				b = nk.NkPropertyi(ctx, "#B:", 0, b, 255, 1, 1)
				a = nk.NkPropertyi(ctx, "#A:", 0, a, 255, 1, 1)
				state.bgColor.SetRGBAi(r, g, b, a)
				nk.NkComboEnd(ctx)
			}
		}
	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)

	handle := nk.NkPlatformDisplayHandle()
	width, height := handle.Width, handle.Height
	state.width, state.height = width, height
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
}

func toggleKeyboard(activity *android.NativeActivity, state *State) {
	if activity == nil {
		return
	}
	if state.kbShown {
		err := activity.SetSoftKeyboardState(android.SoftKeyboardHidden)
		if err != nil {
			log.Println("[WARN] hide keyboard error:", err)
			if err := activity.SetSoftKeyboardState(android.SoftKeyboardVisible); err != nil {
				log.Println("[WARN] show keyboard error:", err)
			} else {
				state.kbShown = true
				return
			}
		}
		state.kbShown = false
		return
	}
	err := activity.SetSoftKeyboardState(android.SoftKeyboardVisible)
	if err != nil {
		log.Println("[WARN] show keyboard error:", err)
		state.kbShown = false
		return
	}
	state.kbShown = true
}

type Option uint8

const (
	Easy Option = 0
	Hard Option = 1
)

type State struct {
	buffer  []byte
	width   int
	height  int
	kbShown bool
	bgColor nk.Color
	prop    int32
	opt     Option
	times   int
}
