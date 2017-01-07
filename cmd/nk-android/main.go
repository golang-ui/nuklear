package main

import (
	"log"
	"runtime"
	"time"

	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/android-go/android"
	"github.com/xlab/android-go/app"
	gl "github.com/xlab/android-go/gles3"
)

func init() {
	app.SetLogTag("Nuklear Activity")
}

const (
	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	nativeWindowEvents := make(chan app.NativeWindowEvent, 1)
	inputQueueEvents := make(chan app.InputQueueEvent, 1)
	inputQueueChan := make(chan *android.InputQueue, 1)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	// sansFont := nk.NkFontAtlasAddFromFile(atlas, s("assets/FreeSans.ttf"), 16, nil)
	defaultFont := nk.NkFontAtlasAddDefault(atlas, 16, nil)
	nk.NkFontStashEnd()

	var ctx *nk.Context
	appState := &State{
		bgColor: nk.NkRgba(28, 48, 62, 255),
	}
	fpsTicker := time.NewTimer(time.Minute)
	fpsTicker.Stop()

	app.Main(func(a app.NativeActivity) {
		a.HandleNativeWindowEvents(nativeWindowEvents)
		a.HandleInputQueueEvents(inputQueueEvents)
		go app.HandleInputQueues(inputQueueChan, func() {
			a.InputQueueHandled()
		}, app.LogInputEvents)
		a.InitDone()
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
				gfxMain(ctx, appState)
				fpsTicker.Reset(time.Second / 30)
			case event := <-nativeWindowEvents:
				switch event.Kind {
				case app.NativeWindowRedrawNeeded:
					gfxMain(ctx, appState)
					a.NativeWindowRedrawDone()
				case app.NativeWindowCreated:
					ctx = nk.NkPlatformInit(event.Window, nk.PlatformInstallCallbacks)
					if ctx == nil {
						log.Fatalln("Nuklear failed to init")
					}
					if defaultFont != nil {
						nk.NkStyleSetFont(ctx, defaultFont.Handle())
					}
					fpsTicker.Reset(time.Second / 30)
				case app.NativeWindowDestroyed:
					fpsTicker.Stop()
					nk.NkPlatformShutdown()
				}
			}
		}
	})
}

// gfxMain is the main GUI code that is borrowed directly from the desktop example.
func gfxMain(ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(200, 200, 600, 600)
	update := nk.NkBegin(ctx, s("Demo"), bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update > 0 {
		nk.NkLayoutRowStatic(ctx, 30, 80, 1)
		if nk.NkButtonLabel(ctx, s("button")) > 0 {
			log.Println("[INFO] button pressed!")
		}
		nk.NkLayoutRowDynamic(ctx, 30, 2)
		if nk.NkOptionLabel(ctx, s("easy"), flag(state.opt == Easy)) > 0 {
			state.opt = Easy
		}
		if nk.NkOptionLabel(ctx, s("hard"), flag(state.opt == Hard)) > 0 {
			state.opt = Hard
		}
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		nk.NkPropertyInt(ctx, s("Compression:"), 0, &state.prop, 100, 10, 1)
		{
			nk.NkLayoutRowDynamic(ctx, 20, 1)
			nk.NkLabel(ctx, s("background:"), nk.TextLeft)
			nk.NkLayoutRowDynamic(ctx, 25, 1)
			size := nk.NkVec2(nk.NkWidgetWidth(ctx), 400)
			if nk.NkComboBeginColor(ctx, state.bgColor, size) > 0 {
				nk.NkLayoutRowDynamic(ctx, 120, 1)
				state.bgColor = nk.NkColorPicker(ctx, state.bgColor, nk.ColorFormatRGBA)
				nk.NkLayoutRowDynamic(ctx, 25, 1)
				r, g, b, a := state.bgColor.RGBAi()
				r = nk.NkPropertyi(ctx, s("#R:"), 0, r, 255, 1, 1)
				g = nk.NkPropertyi(ctx, s("#G:"), 0, g, 255, 1, 1)
				b = nk.NkPropertyi(ctx, s("#B:"), 0, b, 255, 1, 1)
				a = nk.NkPropertyi(ctx, s("#A:"), 0, a, 255, 1, 1)
				state.bgColor.SetRGBAi(r, g, b, a)
				nk.NkComboEnd(ctx)
			}
		}
	}
	nk.NkEnd(ctx)

	// Render
	bg := []float32{38, 38, 38, 255}
	nk.NkColorFv(bg, state.bgColor)

	handle := nk.NkPlatformDisplayHandle()
	width, height := handle.Width, handle.Height
	state.width, state.height = width, height
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
}

type Option uint8

const (
	Easy Option = 0
	Hard Option = 1
)

type State struct {
	width   int
	height  int
	bgColor nk.Color
	prop    int32
	opt     Option
}
